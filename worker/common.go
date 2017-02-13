package worker

type WrappedReturn struct {
	ReturnData interface{}
	Err error
}