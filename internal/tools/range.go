package tools

import "errors"

type Range struct {
	value int
	end   int
	step  int
}

func NewRange(begin int, end int, step int) (Range, error) {
	if (end-begin)*step > 0 {
		return Range{begin, end, step}, nil
	} else {
		return Range{}, errors.New("invalid arguments")
	}
}

func (rng *Range) Next() bool {
	if ((rng.value+rng.step < rng.end) && rng.step > 0) ||
		((rng.value+rng.step > rng.end) && rng.step < 0) {
		rng.value += rng.step
		return true
	} else {
		return false
	}
}

func (rng *Range) Value() int {
	return rng.value
}
