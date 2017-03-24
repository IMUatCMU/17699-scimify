package processor

import (
	"fmt"
	"github.com/go-scim/scimify/resource"
	"sync"
	"github.com/go-scim/scimify/modify"
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
	case resource.Error:
		translatedErr = err.(resource.Error)

	case *modify.InvalidModificationError:
		translatedErr = resource.CreateError(resource.InvalidSyntax, err.Error())

	case *modify.ModificationFailedError, *modify.InvalidPathError, *modify.MissingAttributeForPathError:
		translatedErr = resource.CreateError(resource.InvalidPath, err.Error())

	case *TypeMismatchError, *FormatError, *TypeUnsupportedError, *RequiredMissingError,
		*RequiredUnassignedError, *NoDefinedAttributeError, *UnexpectedTypeError, *UnsupportedValueError:
		translatedErr = resource.CreateError(resource.InvalidValue, err.Error())

	case *ValueChangedError:
		translatedErr = resource.CreateError(resource.Mutability, err.Error())

	case *MissingContextValueError, *AttributeMismatchWithKeyError, *PrerequisiteFailedError:
		translatedErr = resource.CreateError(resource.ServerError, err.Error())

	default:
		fmt.Printf("%T", err)
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
