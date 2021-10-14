package fibonacci_calculator

import "errors"

var InvalidParametersValues = errors.New("invalid parameters values")

func Fibonacci(from uint64, to uint64) ([]uint64, error) {
	if from < 1 || to < 1 || to < from {
		return nil, InvalidParametersValues
	}

	res := make([]uint64, to-from+1)
	if from == 1 {
		res[0] = 1
	}

	var n2, n1 uint64 = 0, 1
	for i := uint64(1); i < to; i++ {
		n2, n1 = n1, n1+n2
		if i >= from-1 {
			res[i-from+1] = n1
		}
	}

	return res, nil
}
