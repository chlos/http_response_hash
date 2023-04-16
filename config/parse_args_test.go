package config

import (
	"os"
	"strconv"
	"testing"

	"github.com/stretchr/testify/require"
)

const (
	_testURL1          = "https://www.reddit.com/"
	_testURL2          = "http://google.com"
	_testParallelLimit = 13
)

func TestNewConfig(t *testing.T) {
	tests := []struct {
		testcase              string
		args                  []string
		expectedError         bool
		expectedParallelLimit int
		expectedURLs          []string
	}{
		{
			testcase:      "parallel flag 0 - urls 0",
			args:          []string{"cmd"},
			expectedError: true,
		},
		{
			testcase:              "parallel flag 0 - urls 1",
			args:                  []string{"cmd", _testURL1},
			expectedError:         false,
			expectedParallelLimit: _defaultParallelLimit,
			expectedURLs:          []string{_testURL1},
		},
		{
			testcase:              "parallel flag 0 - urls 2",
			args:                  []string{"cmd", _testURL1, _testURL2},
			expectedError:         false,
			expectedParallelLimit: _defaultParallelLimit,
			expectedURLs:          []string{_testURL1, _testURL2},
		},

		{
			testcase:      "parallel flag 1 - urls 0",
			args:          []string{"cmd", "-parallel", strconv.Itoa(_testParallelLimit)},
			expectedError: true,
		},
		{
			testcase:              "parallel flag 1 - urls 1",
			args:                  []string{"cmd", "-parallel", strconv.Itoa(_testParallelLimit), _testURL1},
			expectedError:         false,
			expectedParallelLimit: _testParallelLimit,
			expectedURLs:          []string{_testURL1},
		},
		{
			testcase:              "parallel flag 1 - urls 2",
			args:                  []string{"cmd", "-parallel", strconv.Itoa(_testParallelLimit), _testURL1, _testURL2},
			expectedError:         false,
			expectedParallelLimit: _testParallelLimit,
			expectedURLs:          []string{_testURL1, _testURL2},
		},
	}

	for _, tt := range tests {
		t.Run(tt.testcase, func(t *testing.T) {
			os.Args = tt.args
			config, err := NewConfig()

			if tt.expectedError {
				require.Error(t, err)
				require.Nil(t, config)
			} else {
				require.NoError(t, err)
				require.NotNil(t, config)
				require.Equal(t, tt.expectedParallelLimit, config.ParallelLimit)
				require.Equal(t, tt.expectedURLs, config.URLs)
			}
		})
	}
}
