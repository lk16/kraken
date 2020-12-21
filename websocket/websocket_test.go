package websocket

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRound(t *testing.T) {

	type testCase struct {
		input          float64
		decimals       int
		expectedOutput string
	}

	testCases := []testCase{
		{3.1415, 0, "3"},
		{3.1415, 1, "3.1"},
		{3.1415, 2, "3.14"},
		{3.1415, 3, "3.142"},
	}

	for _, testCase := range testCases {
		assert.Equal(t, testCase.expectedOutput, Round(testCase.input, testCase.decimals))
	}
}
