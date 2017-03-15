package resource

type Error struct {
	Schemas    []string `json:"schemas"`
	Detail     string   `json:"detail"`
	ScimType   string   `json:"scimType"`
	Status     string   `json:"status"`
	StatusCode int      `json:"-"`
}

func (e Error) Error() string {
	return e.Detail
}

func (e Error) GetId() string {
	return ""
}

func (e Error) Data() map[string]interface{} {
	return map[string]interface{}{
		"schemas":  e.Schemas,
		"detail":   e.Detail,
		"scimType": e.ScimType,
		"status":   e.Status,
	}
}

// Error identifiers used to create the Error structure
const (
	// The specified filter syntax was invalid, or the specified attribute and filter comparison combination is not supported.
	// Query by GET or POST, PATCH
	InvalidFilter = "invalidFilter"

	// The specified filter yields many more results than the server is willing to calculate or process.
	// Query by GET or POST
	TooMany = "tooMany"

	// One or more of the attribute values are already in use or are reserved.
	// POST, PUT, PATCH
	Uniqueness = "uniqueness"

	// The attempted modification is not compatible with the target attribute's mutability
	// PUT, PATCH
	Mutability = "mutability"

	// The request body message structure was invalid or did not conform to the request schema.
	// POST, PUT, or BULK
	InvalidSyntax = "invalidSyntax"

	// The "path" attribute was invalid or malformed
	// PATCH
	InvalidPath = "invalidPath"

	// The specified "path" did not  yield an attribute or attribute value that could be operated on.
	// This occurs when the specified "path" value contains a filter that yields no match.
	// PATCH
	NoTarget = "noTarget"

	// A required value was missing, or the value specified was not compatible with the operation or attribute type
	// QUERY GET, POST, PUT, PATCH
	InvalidValue = "invalidValue"

	// The specified SCIM protocol version is not supported
	// GET, POST, PUT, PATCH, DELETE
	InvalidVers = "invalidVers"

	// The specified request cannot be completed, due to the passing of sensitive (e.g., personal) information in a request URI.
	// GET
	Sensitive = "sensitive"

	// The requested resource is not found
	NotFound = "notFound"

	// The feature is not implemented
	NotImplemented = "notImplemented"

	// Internal server error
	ServerError = "serverError"
)

// Error Factory
func CreateError(identifier string, detail string) Error {
	switch identifier {
	case InvalidFilter:
		return Error{
			Schemas:    []string{ErrorUrn},
			ScimType:   InvalidFilter,
			Detail:     detail,
			Status:     "400",
			StatusCode: 400,
		}

	case TooMany:
		return Error{
			Schemas:    []string{ErrorUrn},
			ScimType:   TooMany,
			Detail:     detail,
			Status:     "400",
			StatusCode: 400,
		}

	case Uniqueness:
		return Error{
			Schemas:    []string{ErrorUrn},
			ScimType:   Uniqueness,
			Detail:     detail,
			Status:     "400",
			StatusCode: 400,
		}

	case Mutability:
		return Error{
			Schemas:    []string{ErrorUrn},
			ScimType:   Mutability,
			Detail:     detail,
			Status:     "400",
			StatusCode: 400,
		}

	case InvalidSyntax:
		return Error{
			Schemas:    []string{ErrorUrn},
			ScimType:   InvalidSyntax,
			Detail:     detail,
			Status:     "400",
			StatusCode: 400,
		}

	case InvalidPath:
		return Error{
			Schemas:    []string{ErrorUrn},
			ScimType:   InvalidPath,
			Detail:     detail,
			Status:     "400",
			StatusCode: 400,
		}

	case NoTarget:
		return Error{
			Schemas:    []string{ErrorUrn},
			ScimType:   NoTarget,
			Detail:     detail,
			Status:     "400",
			StatusCode: 400,
		}

	case InvalidValue:
		return Error{
			Schemas:    []string{ErrorUrn},
			ScimType:   InvalidValue,
			Detail:     detail,
			Status:     "400",
			StatusCode: 400,
		}

	case InvalidVers:
		return Error{
			Schemas:    []string{ErrorUrn},
			ScimType:   InvalidVers,
			Detail:     detail,
			Status:     "400",
			StatusCode: 400,
		}

	case Sensitive:
		return Error{
			Schemas:    []string{ErrorUrn},
			ScimType:   Sensitive,
			Detail:     detail,
			Status:     "400",
			StatusCode: 400,
		}

	case NotFound:
		return Error{
			Schemas:    []string{ErrorUrn},
			ScimType:   NotFound,
			Detail:     detail,
			Status:     "404",
			StatusCode: 404,
		}

	case ServerError:
		return Error{
			Schemas:    []string{ErrorUrn},
			ScimType:   ServerError,
			Detail:     detail,
			Status:     "500",
			StatusCode: 500,
		}

	case NotImplemented:
		return Error{
			Schemas:    []string{ErrorUrn},
			ScimType:   NotImplemented,
			Detail:     detail,
			Status:     "501",
			StatusCode: 501,
		}

	default:
		panic("unknown error identifier " + identifier)
	}
}
