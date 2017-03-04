package resource

type ScimObject interface {
	GetId() string
	Data() map[string]interface{}
}
