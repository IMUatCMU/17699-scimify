package processor

import (
	"github.com/go-scim/scimify/resource"
	"sync"
)

var (
	oneError       sync.Once
	errorTranslate Processor
)

func ErrorTranslatingProcessor() Processor {
	oneError.Do(func() {
		errorTranslate = &errorTranslatingProcessor{}
	})
	return errorTranslate
}

type errorTranslatingProcessor struct{}

func (etp *errorTranslatingProcessor) Process(ctx *ProcessorContext) error {
	var translatedErr resource.Error
	err := etp.getError(ctx)

	switch err.(type) {
	case *TypeMismatchError, *FormatError, *TypeUnsupportedError, *RequiredMissingError,
		*RequiredUnassignedError, *NoDefinedAttributeError, *UnexpectedTypeError, *UnsupportedValueError:
		translatedErr = resource.CreateError(resource.InvalidValue, err.Error())

	case *ValueChangedError:
		translatedErr = resource.CreateError(resource.Mutability, err.Error())

	case *MissingContextValueError, *AttributeMismatchWithKeyError, *PrerequisiteFailedError:
		translatedErr = resource.CreateError(resource.ServerError, err.Error())

	default:
		translatedErr = resource.CreateError(resource.ServerError, err.Error())
	}
	ctx.Err = translatedErr

	return nil
}

func (etp *errorTranslatingProcessor) getError(ctx *ProcessorContext) error {
	if ctx.Err == nil {
		panic(&MissingContextValueError{"error thrown"})
	}
	return ctx.Err
}
