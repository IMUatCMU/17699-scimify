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

func (e Error) AsResource() *Resource {
	return &Resource{
		Attributes: map[string]interface{}{
			"schemas":  e.Schemas,
			"detail":   e.Detail,
			"scimType": e.ScimType,
			"status":   e.Status,
		},
	}
}

// Error identifiers used to create the Error structure
const (
	InvalidFilter  = "invalidFilter"
	TooMany        = "tooMany"
	Uniqueness     = "uniqueness"
	Mutability     = "mutability"
	InvalidSyntax  = "invalidSyntax"
	InvalidPath    = "invalidPath"
	NoTarget       = "noTarget"
	InvalidValue   = "invalidValue"
	InvalidVers    = "invalidVers"
	Sensitive      = "sensitive"
	NotFound       = "notFound"
	NotImplemented = "notImplemented"

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
