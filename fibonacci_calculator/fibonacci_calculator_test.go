package fibonacci_calculator

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestFibonacci(t *testing.T) {
	testCases := []struct {
		inputFrom   uint64
		inputTo     uint64
		expectedRes []uint64
		expectedErr error
	}{
		{inputFrom: 1, inputTo: 1, expectedRes: []uint64{1}, expectedErr: nil},
		{inputFrom: 1, inputTo: 2, expectedRes: []uint64{1, 1}, expectedErr: nil},
		{inputFrom: 1, inputTo: 3, expectedRes: []uint64{1, 1, 2}, expectedErr: nil},
		{inputFrom: 2, inputTo: 3, expectedRes: []uint64{1, 2}, expectedErr: nil},
		{inputFrom: 3, inputTo: 3, expectedRes: []uint64{2}, expectedErr: nil},
		{inputFrom: 2, inputTo: 4, expectedRes: []uint64{1, 2, 3}, expectedErr: nil},
		{inputFrom: 0, inputTo: 2, expectedRes: nil, expectedErr: InvalidParametersValues},
		{inputFrom: 2, inputTo: 0, expectedRes: nil, expectedErr: InvalidParametersValues},
		{inputFrom: 0, inputTo: 0, expectedRes: nil, expectedErr: InvalidParametersValues},
		{inputFrom: 4, inputTo: 3, expectedRes: nil, expectedErr: InvalidParametersValues},
	}
	for _, tc := range testCases {
		t.Run(fmt.Sprintf("FROM=%d, TO=%d", tc.inputFrom, tc.inputTo), func(t *testing.T) {
			res, err := Fibonacci(tc.inputFrom, tc.inputTo, make(<-chan struct{}))
			require.Equal(t, tc.expectedRes, res)
			require.Equal(t, tc.expectedErr, err)
		})
	}
}
