// Constants for the SCIM schema
package schema

// schema attribute types
const (
	String    = "string"
	Boolean   = "boolean"
	Decimal   = "decimal"
	Integer   = "integer"
	DateTime  = "dateTime"
	Reference = "reference"
	Complex   = "complex"
)

// schema attribute mutability
const (
	ReadOnly  = "readOnly"
	ReadWrite = "readWrite"
	Immutable = "immutable"
	WriteOnly = "writeOnly"
)

// schema attribute returned
const (
	Always  = "always"
	Never   = "never"
	Default = "default"
	Request = "request"
)

// schema attribute uniqueness
const (
	None   = "none"
	Server = "server"
	Global = "global"
)

// schema attribute reference types
const (
	External = "external"
	Uri      = "uri"
)
