package hashing

import (
	"testing"

	"github.com/stretchr/testify/require"
)

const (
	_testURL1 = "http://www.adjust.com"
	_testURL2 = "http://google.com"
)

func TestNewHashing(t *testing.T) {
	tests := []struct {
		testcase      string
		limit         int
		urls          []string
		expectedError bool
	}{
		{
			testcase:      "limit 0 - urls 0",
			limit:         0,
			urls:          []string{},
			expectedError: true,
		},
		{
			testcase:      "limit 1 - urls 0",
			limit:         1,
			urls:          []string{},
			expectedError: false,
		},
		{
			testcase:      "limit 2 - urls 0",
			limit:         2,
			urls:          []string{},
			expectedError: false,
		},

		{
			testcase:      "limit 0 - urls 1",
			limit:         0,
			urls:          []string{_testURL1},
			expectedError: true,
		},
		{
			testcase:      "limit 1 - urls 1",
			limit:         1,
			urls:          []string{_testURL1},
			expectedError: false,
		},
		{
			testcase:      "limit 2 - urls 1",
			limit:         2,
			urls:          []string{_testURL1},
			expectedError: false,
		},

		{
			testcase:      "limit 0 - urls 2",
			limit:         0,
			urls:          []string{_testURL1, _testURL2},
			expectedError: true,
		},
		{
			testcase:      "limit 1 - urls 2",
			limit:         1,
			urls:          []string{_testURL1, _testURL2},
			expectedError: false,
		},
		{
			testcase:      "limit 2 - urls 2",
			limit:         2,
			urls:          []string{_testURL1, _testURL2},
			expectedError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.testcase, func(t *testing.T) {
			h, err := NewHashing(tt.limit, tt.urls)

			if tt.expectedError {
				require.Error(t, err)
				require.Nil(t, h)
			} else {
				require.NoError(t, err)
				require.NotNil(t, h)
			}
		})
	}

}
