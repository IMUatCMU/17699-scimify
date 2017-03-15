package processor

import (
	"bytes"
	"fmt"
	"github.com/go-scim/scimify/resource"
	"math"
	"reflect"
	"runtime"
	"strconv"
	"strings"
	"sync"
	"unicode/utf8"
)

type assistedJsonSerializationProcessor struct {
	argSlot RName
}

func (ajp *assistedJsonSerializationProcessor) Process(ctx *ProcessorContext) error {
	target := getR(ctx, ajp.argSlot, true, nil)
	schema := getSchema(ctx, true)
	include := getStringArray(ctx, ArgIncludePaths, false, []string{})
	exclude := getStringArray(ctx, ArgExcludePaths, false, []string{})

	switch target.(type) {
	case resource.ScimObject:
		bytes, err := ajp.serialize(target.(resource.ScimObject), include, exclude, schema)
		ctx.Results[RBodyBytes] = bytes
		return err

	case []resource.ScimObject:
		bytes, err := ajp.serializeArray(target.([]resource.ScimObject), include, exclude, schema)
		ctx.Results[RBodyBytes] = bytes
		return err

	default:
		return &PrerequisiteFailedError{reporter: "assisted json serializer", requirement: "single or array of scim object"}
	}
}

var hex = "0123456789abcdef"
var numberType = reflect.TypeOf(Number(""))

func (ajp *assistedJsonSerializationProcessor) serialize(target resource.ScimObject, inclusionPaths, exclusionPaths []string, schema *resource.Schema) ([]byte, error) {
	e := newEncodeState()
	err := e.marshal(target.Data(), encOpts{escapeHTML: true, inclusionPaths: inclusionPaths, exclusionPaths: exclusionPaths}, schema.AsAttribute())
	if nil != err {
		return nil, err
	}
	return e.Bytes(), nil
}

func (ajp *assistedJsonSerializationProcessor) serializeArray(target []resource.ScimObject, inclusionPaths, exclusionPaths []string, schema *resource.Schema) ([]byte, error) {
	e := newEncodeState()
	e.WriteByte('[')
	for i, targetElem := range target {
		bytes, err := ajp.serialize(targetElem, inclusionPaths, exclusionPaths, schema)
		if err != nil {
			return nil, err
		}
		if i > 0 {
			e.WriteByte(',')
		}
		e.Write(bytes)
	}
	e.WriteByte(']')
	return e.Bytes(), nil
}

// encode options
type encOpts struct {
	// quoted causes primitive fields to be encoded inside JSON strings.
	quoted bool
	// escapeHTML causes '<', '>', and '&' to be escaped in JSON strings.
	escapeHTML bool
	// inclusion paths are the path requested to be included in the return.
	inclusionPaths []string
	// exclusion paths are the path requested to be excluded in the return.
	exclusionPaths []string
}

// Decide whether an attribute should be included in the JSON response:
// - if "returned=always", the attribute is returned
// - if "returned=never", the attribute is not returned
// - if "returned=request", the attribute is only returned when it's assigned and it's path appears in the inclusionPaths list
// - if "returned=default", the attribute is returned except when it's path appears in the exclusionPaths list or it is not assigned
func (opt encOpts) shouldEncode(v reflect.Value, attr *resource.Attribute) bool {
	trueIfAssigned := func(v reflect.Value) bool {
		//fmt.Println(attr.Name, v)
		return v.IsValid() && attr.IsValueAssigned(v)
	}

	switch attr.Returned {
	case resource.Always:
		return true
	case resource.Never:
		return false
	case resource.Request:
		for _, p := range opt.inclusionPaths {
			if strings.ToLower(attr.Assist.Path) == strings.ToLower(p) {
				return trueIfAssigned(v)
			}
		}
		return false
	case resource.Default:
		for _, p := range opt.exclusionPaths {
			if strings.ToLower(attr.Assist.Path) == strings.ToLower(p) {
				return false
			}
		}
		return trueIfAssigned(v)
	default:
		return trueIfAssigned(v)
	}
}

// encode function
type encoderFunc func(e *encodeState, v reflect.Value, opts encOpts, attr *resource.Attribute)

// encode State
var encodeStatePool sync.Pool

func newEncodeState() *encodeState {
	if v := encodeStatePool.Get(); v != nil {
		e := v.(*encodeState)
		e.Reset()
		return e
	}
	return new(encodeState)
}

type encodeState struct {
	bytes.Buffer
	scratch [64]byte
}

func (e *encodeState) error(err error) {
	fmt.Println(string(e.Bytes()))
	panic(err)
}

func (e *encodeState) marshal(v interface{}, opts encOpts, attr *resource.Attribute) (err error) {
	defer func() {
		if r := recover(); r != nil {
			if _, ok := r.(runtime.Error); ok {
				panic(r)
			}
			if s, ok := r.(string); ok {
				panic(s)
			}
			err = r.(error)
		}
	}()
	e.reflectValue(reflect.ValueOf(v), opts, attr)
	return nil
}

func (e *encodeState) reflectValue(v reflect.Value, opts encOpts, attr *resource.Attribute) {
	valueEncoder(v, attr)(e, v, opts, attr)
}

// NOTE: keep in sync with stringBytes below.
func (e *encodeState) string(s string, escapeHTML bool) int {
	len0 := e.Len()
	e.WriteByte('"')
	start := 0
	for i := 0; i < len(s); {
		if b := s[i]; b < utf8.RuneSelf {
			if 0x20 <= b && b != '\\' && b != '"' &&
				(!escapeHTML || b != '<' && b != '>' && b != '&') {
				i++
				continue
			}
			if start < i {
				e.WriteString(s[start:i])
			}
			switch b {
			case '\\', '"':
				e.WriteByte('\\')
				e.WriteByte(b)
			case '\n':
				e.WriteByte('\\')
				e.WriteByte('n')
			case '\r':
				e.WriteByte('\\')
				e.WriteByte('r')
			case '\t':
				e.WriteByte('\\')
				e.WriteByte('t')
			default:
				// This encodes bytes < 0x20 except for \t, \n and \r.
				// If escapeHTML is set, it also escapes <, >, and &
				// because they can lead to security holes when
				// user-controlled strings are rendered into JSON
				// and served to some browsers.
				e.WriteString(`\u00`)
				e.WriteByte(hex[b>>4])
				e.WriteByte(hex[b&0xF])
			}
			i++
			start = i
			continue
		}
		c, size := utf8.DecodeRuneInString(s[i:])
		if c == utf8.RuneError && size == 1 {
			if start < i {
				e.WriteString(s[start:i])
			}
			e.WriteString(`\ufffd`)
			i += size
			start = i
			continue
		}
		// U+2028 is LINE SEPARATOR.
		// U+2029 is PARAGRAPH SEPARATOR.
		// They are both technically valid characters in JSON strings,
		// but don't work in JSONP, which has to be evaluated as JavaScript,
		// and can lead to security holes there. It is valid JSON to
		// escape them, so we do so unconditionally.
		// See http://timelessrepo.com/json-isnt-a-javascript-subset for discussion.
		if c == '\u2028' || c == '\u2029' {
			if start < i {
				e.WriteString(s[start:i])
			}
			e.WriteString(`\u202`)
			e.WriteByte(hex[c&0xF])
			i += size
			start = i
			continue
		}
		i += size
	}
	if start < len(s) {
		e.WriteString(s[start:])
	}
	e.WriteByte('"')
	return e.Len() - len0
}

// NOTE: keep in sync with string above.
func (e *encodeState) stringBytes(s []byte, escapeHTML bool) int {
	len0 := e.Len()
	e.WriteByte('"')
	start := 0
	for i := 0; i < len(s); {
		if b := s[i]; b < utf8.RuneSelf {
			if 0x20 <= b && b != '\\' && b != '"' &&
				(!escapeHTML || b != '<' && b != '>' && b != '&') {
				i++
				continue
			}
			if start < i {
				e.Write(s[start:i])
			}
			switch b {
			case '\\', '"':
				e.WriteByte('\\')
				e.WriteByte(b)
			case '\n':
				e.WriteByte('\\')
				e.WriteByte('n')
			case '\r':
				e.WriteByte('\\')
				e.WriteByte('r')
			case '\t':
				e.WriteByte('\\')
				e.WriteByte('t')
			default:
				// This encodes bytes < 0x20 except for \t, \n and \r.
				// If escapeHTML is set, it also escapes <, >, and &
				// because they can lead to security holes when
				// user-controlled strings are rendered into JSON
				// and served to some browsers.
				e.WriteString(`\u00`)
				e.WriteByte(hex[b>>4])
				e.WriteByte(hex[b&0xF])
			}
			i++
			start = i
			continue
		}
		c, size := utf8.DecodeRune(s[i:])
		if c == utf8.RuneError && size == 1 {
			if start < i {
				e.Write(s[start:i])
			}
			e.WriteString(`\ufffd`)
			i += size
			start = i
			continue
		}
		// U+2028 is LINE SEPARATOR.
		// U+2029 is PARAGRAPH SEPARATOR.
		// They are both technically valid characters in JSON strings,
		// but don't work in JSONP, which has to be evaluated as JavaScript,
		// and can lead to security holes there. It is valid JSON to
		// escape them, so we do so unconditionally.
		// See http://timelessrepo.com/json-isnt-a-javascript-subset for discussion.
		if c == '\u2028' || c == '\u2029' {
			if start < i {
				e.Write(s[start:i])
			}
			e.WriteString(`\u202`)
			e.WriteByte(hex[c&0xF])
			i += size
			start = i
			continue
		}
		i += size
	}
	if start < len(s) {
		e.Write(s[start:])
	}
	e.WriteByte('"')
	return e.Len() - len0
}

func valueEncoder(v reflect.Value, attr *resource.Attribute) encoderFunc {
	if !v.IsValid() {
		return invalidValueEncoder
	}
	return newTypeEncoder(v.Type(), attr)
}

func newTypeEncoder(t reflect.Type, attr *resource.Attribute) encoderFunc {
	if t.Kind() == reflect.Interface {
		return interfaceEncoder
	}

	multiValued, dataType := attr.MultiValued, attr.Type
	if multiValued {
		switch t.Kind() {
		case reflect.Slice:
			return newSliceEncoder(t, attr)
		case reflect.Array:
			return newArrayEncoder(t, attr)
		default:
			return unsupportedTypeEncoder
		}
	} else {
		switch dataType {
		case type_bool:
			switch t.Kind() {
			case reflect.Bool:
				return boolEncoder
			default:
				return unsupportedTypeEncoder
			}
		case type_int:
			switch t.Kind() {
			case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
				return intEncoder
			default:
				return unsupportedTypeEncoder
			}
		case type_float:
			switch t.Kind() {
			case reflect.Float32:
				return float32Encoder
			case reflect.Float64:
				return float64Encoder
			default:
				return unsupportedTypeEncoder
			}
		case type_string, type_ref, type_binary, type_datetime:
			switch t.Kind() {
			case reflect.String:
				return stringEncoder
			default:
				return unsupportedTypeEncoder
			}
		case type_complex:
			switch t.Kind() {
			case reflect.Map:
				return newMapEncoder(t, attr)
			default:
				return unsupportedTypeEncoder
			}
		default:
			return unsupportedTypeEncoder
		}
	}
}

// encoders
func invalidValueEncoder(e *encodeState, _ reflect.Value, _ encOpts, _ *resource.Attribute) {
	e.WriteString("null")
}

func interfaceEncoder(e *encodeState, v reflect.Value, opts encOpts, attr *resource.Attribute) {
	if v.IsNil() {
		e.WriteString("null")
		return
	}
	e.reflectValue(v.Elem(), opts, attr)
}

func boolEncoder(e *encodeState, v reflect.Value, opts encOpts, _ *resource.Attribute) {
	if opts.quoted {
		e.WriteByte('"')
	}
	if v.Bool() {
		e.WriteString("true")
	} else {
		e.WriteString("false")
	}
	if opts.quoted {
		e.WriteByte('"')
	}
}

func intEncoder(e *encodeState, v reflect.Value, opts encOpts, _ *resource.Attribute) {
	b := strconv.AppendInt(e.scratch[:0], v.Int(), 10)
	if opts.quoted {
		e.WriteByte('"')
	}
	e.Write(b)
	if opts.quoted {
		e.WriteByte('"')
	}
}

type floatEncoder int // number of bits

func (bits floatEncoder) encode(e *encodeState, v reflect.Value, opts encOpts, _ *resource.Attribute) {
	f := v.Float()
	if math.IsInf(f, 0) || math.IsNaN(f) {
		e.error(&UnsupportedValueError{v, strconv.FormatFloat(f, 'g', -1, int(bits))})
	}
	b := strconv.AppendFloat(e.scratch[:0], f, 'g', -1, int(bits))
	if opts.quoted {
		e.WriteByte('"')
	}
	e.Write(b)
	if opts.quoted {
		e.WriteByte('"')
	}
}

var (
	float32Encoder = (floatEncoder(32)).encode
	float64Encoder = (floatEncoder(64)).encode
)

func stringEncoder(e *encodeState, v reflect.Value, opts encOpts, attr *resource.Attribute) {
	if v.Type() == numberType {
		numStr := v.String()
		// In Go1.5 the empty string encodes to "0", while this is not a valid number literal
		// we keep compatibility so check validity after this.
		if numStr == "" {
			numStr = "0" // Number's zero-val
		}
		if !isValidNumber(numStr) {
			e.error(fmt.Errorf("json: invalid number literal %q", numStr))
		}
		e.WriteString(numStr)
		return
	}
	if opts.quoted {
		e0 := &encodeState{}
		err := e0.marshal(v, encOpts{escapeHTML: true}, attr)
		if err != nil {
			e.error(err)
		}
		e.string(string(e0.Bytes()), opts.escapeHTML)
	} else {
		e.string(v.String(), opts.escapeHTML)
	}
}

type mapEncoder struct{}

func (me *mapEncoder) encode(e *encodeState, v reflect.Value, opts encOpts, attr *resource.Attribute) {
	if v.IsNil() {
		e.WriteString("null")
		return
	}
	keyAttrs := make([]*resource.Attribute, 0, len(attr.SubAttributes))
	for _, subAttr := range attr.SubAttributes {
		keyAttrs = append(keyAttrs, subAttr)
	}
	e.WriteByte('{')
	isFirst := true
	for _, subAttr := range keyAttrs {
		val := v.MapIndex(reflect.ValueOf(subAttr.Assist.JSONName))
		if !opts.shouldEncode(val, subAttr) {
			continue
		}
		if !isFirst {
			e.WriteByte(',')
		}
		isFirst = false
		e.string(subAttr.Assist.JSONName, opts.escapeHTML)
		e.WriteByte(':')
		valueEncoder(val, subAttr)(e, val, opts, subAttr)
	}
	e.WriteByte('}')
}

func newMapEncoder(t reflect.Type, _ *resource.Attribute) encoderFunc {
	switch t.Key().Kind() {
	case reflect.String:
	default:
		return unsupportedTypeEncoder
	}
	me := &mapEncoder{}
	return me.encode
}

type sliceEncoder struct {
	arrayEnc encoderFunc
}

func (se *sliceEncoder) encode(e *encodeState, v reflect.Value, opts encOpts, attr *resource.Attribute) {
	if v.IsNil() {
		e.WriteString("null")
		return
	}
	se.arrayEnc(e, v, opts, attr)
}

func newSliceEncoder(t reflect.Type, attr *resource.Attribute) encoderFunc {
	enc := &sliceEncoder{arrayEnc: newArrayEncoder(t, attr)}
	return enc.encode
}

type arrayEncoder struct {
	elemEnc encoderFunc
}

func (ae *arrayEncoder) encode(e *encodeState, v reflect.Value, opts encOpts, attr *resource.Attribute) {
	elemAttr := attr.Clone()
	elemAttr.MultiValued = false

	e.WriteByte('[')
	n := v.Len()
	for i := 0; i < n; i++ {
		if i > 0 {
			e.WriteByte(',')
		}
		ae.elemEnc(e, v.Index(i), opts, elemAttr)
	}
	e.WriteByte(']')
}

func newArrayEncoder(t reflect.Type, attr *resource.Attribute) encoderFunc {
	enc := &arrayEncoder{elemEnc: newTypeEncoder(t.Elem(), attr)}
	return enc.encode
}

// number
// A Number represents a JSON number literal.
type Number string

// String returns the literal text of the number.
func (n Number) String() string { return string(n) }

// Float64 returns the number as a float64.
func (n Number) Float64() (float64, error) {
	return strconv.ParseFloat(string(n), 64)
}

// Int64 returns the number as an int64.
func (n Number) Int64() (int64, error) {
	return strconv.ParseInt(string(n), 10, 64)
}

// isValidNumber reports whether s is a valid JSON number literal.
func isValidNumber(s string) bool {
	// This function implements the JSON numbers grammar.
	// See https://tools.ietf.org/html/rfc7159#section-6
	// and http://json.org/number.gif

	if s == "" {
		return false
	}

	// Optional -
	if s[0] == '-' {
		s = s[1:]
		if s == "" {
			return false
		}
	}

	// Digits
	switch {
	default:
		return false

	case s[0] == '0':
		s = s[1:]

	case '1' <= s[0] && s[0] <= '9':
		s = s[1:]
		for len(s) > 0 && '0' <= s[0] && s[0] <= '9' {
			s = s[1:]
		}
	}

	// . followed by 1 or more digits.
	if len(s) >= 2 && s[0] == '.' && '0' <= s[1] && s[1] <= '9' {
		s = s[2:]
		for len(s) > 0 && '0' <= s[0] && s[0] <= '9' {
			s = s[1:]
		}
	}

	// e or E followed by an optional - or + and
	// 1 or more digits.
	if len(s) >= 2 && (s[0] == 'e' || s[0] == 'E') {
		s = s[1:]
		if s[0] == '+' || s[0] == '-' {
			s = s[1:]
			if s == "" {
				return false
			}
		}
		for len(s) > 0 && '0' <= s[0] && s[0] <= '9' {
			s = s[1:]
		}
	}

	// Make sure we are at the end.
	return s == ""
}

// errors

func unsupportedTypeEncoder(e *encodeState, v reflect.Value, _ encOpts, a *resource.Attribute) {
	e.error(&UnexpectedTypeError{v.Type(), a})
}
