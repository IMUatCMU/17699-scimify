package processor

import "github.com/go-scim/scimify/resource"

type errorTranslatingProcessor struct{}

func (etp *errorTranslatingProcessor) Process(ctx *ProcessorContext) error {
	var finalError resource.Error
	err := getError(ctx, true)
	switch err.(type) {

	case *TypeMismatchError, *FormatError, *TypeUnsupportedError, *RequiredMissingError,
		*RequiredUnassignedError, *NoDefinedAttributeError, *UnexpectedTypeError, *UnsupportedValueError:
		finalError = resource.CreateError(resource.InvalidValue, err.Error())

	case *ValueChangedError:
		finalError = resource.CreateError(resource.Mutability, err.Error())

	case *MissingContextValueError, *AttributeMismatchWithKeyError, *PrerequisiteFailedError:
		finalError = resource.CreateError(resource.ServerError, err.Error())

	default:
		finalError = resource.CreateError(resource.ServerError, err.Error())
	}

	ctx.Results[RFinalError] = finalError
	return nil
}
