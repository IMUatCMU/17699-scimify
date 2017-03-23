package modify

import (
	"fmt"
)

type InvalidModificationError struct {
	reason string
}

func (ime *InvalidModificationError) Error() string {
	return fmt.Sprintf("modification payload is invalid due to: %s", ime.reason)
}

type ModificationFailedError struct {
	cause interface{}
}

func (mfe *ModificationFailedError) Error() string {
	switch mfe.cause.(type) {
	case string:
		return mfe.cause.(string)
	case error:
		return mfe.cause.(error).Error()
	default:
		return fmt.Sprintf("%+v", mfe.cause)
	}
}

type InvalidPathError struct {
	path   string
	reason string
}

func (ipe *InvalidPathError) Error() string {
	if len(ipe.path) == 0 {
		if len(ipe.reason) == 0 {
			return "invalid path"
		} else {
			return ipe.reason
		}
	} else {
		if len(ipe.reason) == 0 {
			return fmt.Sprintf("path %s is invalid", ipe.path)
		} else {
			return fmt.Sprintf("path %s is invalid due to: %s", ipe.path, ipe.reason)
		}
	}
}

type MissingAttributeForPathError struct {
	path string
}

func (mae *MissingAttributeForPathError) Error() string {
	return fmt.Sprintf("attribute not found for path (component) %s", mae.path)
}
