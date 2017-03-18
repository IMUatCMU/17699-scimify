package service

import (
	"github.com/go-scim/scimify/processor"
	"github.com/go-scim/scimify/resource"
	"net/http"
)

type response struct {
	statusCode int
	headers    map[string]string
	body       []byte
}

var nil_response response

type service func(*http.Request) (response, error)

// TODO use recover() to handle panics
func endpoint(srv service) http.HandlerFunc {
	return http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		var (
			c int
			h map[string]string
			b []byte
		)
		r, e := srv(req)
		if nil != e {
			c, h, b = handleError(e)
		} else {
			c, h, b = r.statusCode, r.headers, r.body
		}

		if len(h) == 0 {
			h = make(map[string]string)
		}
		if len(b) != 0 {
			if _, ok := h["Content-Type"]; !ok {
				h["Content-Type"] = "application/json+scim"
			}
		}

		writeResponse(rw, c, h, b)
	})
}

func handleError(err error) (int, map[string]string, []byte) {
	var scimErr resource.Error
	if e, ok := err.(resource.Error); !ok {
		scimErr = resource.CreateError(resource.ServerError, err.Error())
	} else {
		scimErr = e
	}

	serializer := processor.SimpleJsonSerializationProcessor()
	ctx := &processor.ProcessorContext{
		SerializationTargetFunc: func() interface{} {
			return scimErr
		},
	}
	err = serializer.Process(ctx)
	if nil != err {
		return http.StatusInternalServerError,
			map[string]string{"Content-Type": "text/plain"},
			[]byte(scimErr.Detail)
	} else {
		return scimErr.StatusCode,
			map[string]string{"Content-Type": "application/json+scim"},
			ctx.ResponseBody
	}
}

func writeResponse(rw http.ResponseWriter, statusCode int, headers map[string]string, body []byte) {
	for k, v := range headers {
		rw.Header().Set(k, v)
	}
	rw.WriteHeader(statusCode)
	rw.Write(body)
}
