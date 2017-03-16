package processor

import (
	"crypto/sha1"
	"encoding/base64"
	"fmt"
	"github.com/go-scim/scimify/resource"
	"time"
)

const (
	type_bool     = resource.Boolean
	type_int      = resource.Integer
	type_float    = resource.Decimal
	type_string   = resource.String
	type_ref      = resource.Reference
	type_binary   = resource.Binary
	type_datetime = resource.DateTime
	type_complex  = resource.Complex
)

func getCurrentTime() string {
	return time.Now().Format("2006-01-02T15:04:05Z")
}

func generateNewVersion(id string) string {
	hash := sha1.New()
	now := time.Now().Format(time.RFC3339Nano)
	hash.Write([]byte(id))
	hash.Write([]byte(now))
	return fmt.Sprintf("W/\"%s\"", base64.StdEncoding.EncodeToString(hash.Sum(nil)))
}
