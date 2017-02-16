package helper

import (
	"sync"
)

type MapEntryProcessor func(key string, value interface{}) (interface{}, error)
type SliceElementProcessor func(idx int, elem interface{}) (interface{}, error)

type Aggregator interface {
	Aggregate(interface{})
	Result() interface{}
}

type result struct {
	Value	interface{}
	Err 	error
}

func WalkSliceInParallel(target []interface{}, processFunc SliceElementProcessor, aggregator Aggregator) (interface{}, error) {
	done := make(chan struct{})
	defer close(done)

	c := processSliceElement(done, target, processFunc)
	for r := range c {
		if r.Err != nil {
			return nil, r.Err
		} else {
			aggregator.Aggregate(r.Value)
		}
	}

	return aggregator.Result(), nil
}

func processSliceElement(done <-chan struct{}, target []interface{}, processFunc SliceElementProcessor) (<-chan result) {
	c := make(chan result)

	go func() {
		var wg sync.WaitGroup

		// walk map
		for i, e := range target {
			index, element := i, e
			wg.Add(1)

			// walk entry
			go func() {
				r, err := processFunc(index, element)
				select {
				case c <- result{r, err}:
				case <-done:
				}
				wg.Done()
			}()

			// abort walk if done
			select {
			case <-done:
				return
			default:
				continue
			}
		}

		// close result channel when everything is done
		go func() {
			wg.Wait()
			close(c)
		}()
	}()

	return c
}

func WalkStringMapInParallel(target map[string]interface{}, processFunc MapEntryProcessor, aggregator Aggregator) (interface{}, error) {
	done := make(chan struct{})
	defer close(done)

	c := processStringMapEntry(done, target, processFunc)
	for r := range c {
		if r.Err != nil {
			return nil, r.Err
		} else {
			aggregator.Aggregate(r.Value)
		}
	}

	return aggregator.Result(), nil
}

func processStringMapEntry(done <-chan struct{}, target map[string]interface{}, processFunc MapEntryProcessor) (<-chan result) {
	c := make(chan result)

	go func() {
		var wg sync.WaitGroup

		// walk map
		for k, v := range target {
			key, val := k, v
			wg.Add(1)

			// walk entry
			go func() {
				r, err := processFunc(key, val)
				select {
				case c <- result{r, err}:
				case <-done:
				}
				wg.Done()
			}()

			// abort walk if done
			select {
			case <-done:
				return
			default:
				continue
			}
		}

		// close result channel when everything is done
		go func() {
			wg.Wait()
			close(c)
		}()
	}()

	return c
}