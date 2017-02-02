// Constants for the SCIM schema
package resource

// URN
const (
	UserUrn                  = "urn:ietf:params:scim:schemas:core:2.0:User"
	GroupUrn                 = "urn:ietf:params:scim:schemas:core:2.0:Group"
	ResourceTypeUrn          = "urn:ietf:params:scim:schemas:core:2.0:ResourceType"
	ServiceProviderConfigUrn = "urn:ietf:params:scim:schemas:core:2.0:ServiceProviderConfig"
	SchemaUrn                = "urn:ietf:params:scim:schemas:core:2.0:Schema"
)

// Internally Used Urn
const (
	CommonUrn = "urn:ietf:params:scim:schemas:common:2.0:common"
)

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
