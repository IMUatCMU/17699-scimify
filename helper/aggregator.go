package helper

type Aggregator interface {
	Aggregate(key, value interface{})
	Result() interface{}
}

type DoNothingAggregator struct {}
func (d *DoNothingAggregator) Aggregate(key, value interface{}) {}
func (d *DoNothingAggregator) Result() interface{} { return true }
