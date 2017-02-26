package resource

type ContextKey int

const (
	CK_Schema          = ContextKey(0)
	CK_Reference       = ContextKey(1)
	CK_ResourceType    = ContextKey(2)
	CK_ResourceTypeURI = ContextKey(3)
)
