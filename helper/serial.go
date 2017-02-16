package helper

func WalkSliceInSerial(target []interface{}, processFunc SliceElementProcessor, aggregator Aggregator) (interface{}, error) {
	for idx, elem := range target {
		result, err := processFunc(idx, elem)
		if err != nil {
			return nil, err
		}
		aggregator.Aggregate(idx, result)
	}
	return aggregator.Result(), nil
}
