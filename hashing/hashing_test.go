package hashing

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/require"
)

const (
	_testURL1     = "https://www.reddit.com/"
	_testURL2     = "http://google.com"
	_testHTTPBody = "Test HTTP body"
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
func _getHTTPBodyMock(url string) ([]byte, error) {
	return []byte(_testHTTPBody), nil
}

func _getHTTPBodyMockError(url string) ([]byte, error) {
	return []byte{}, errors.New("test error")
}

func TestStart(t *testing.T) {
	tests := []struct {
		testcase         string
		urls             []string
		httpRequestFails bool
		expectedHashes   int
	}{
		{
			testcase:         "urls 0 - http ok",
			urls:             []string{},
			httpRequestFails: false,
			expectedHashes:   0,
		},
		{
			testcase:         "urls 1 - http ok",
			urls:             []string{_testURL1},
			httpRequestFails: false,
			expectedHashes:   1,
		},
		{
			testcase:         "urls 2 - http ok",
			urls:             []string{_testURL1, _testURL2},
			httpRequestFails: false,
			expectedHashes:   2,
		},

		{
			testcase:         "urls 0 - http fails",
			urls:             []string{},
			httpRequestFails: true,
			expectedHashes:   0,
		},
		{
			testcase:         "urls 1 - http fails",
			urls:             []string{_testURL1},
			httpRequestFails: true,
			expectedHashes:   0,
		},
		{
			testcase:         "urls 2 - http fails",
			urls:             []string{_testURL1, _testURL2},
			httpRequestFails: true,
			expectedHashes:   0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.testcase, func(t *testing.T) {
			if tt.httpRequestFails {
				getHTTPBody = _getHTTPBodyMockError
			} else {
				getHTTPBody = _getHTTPBodyMock
			}

			h, _ := NewHashing(10, tt.urls)
			h.Start()
			actualHashes := 0
			for result := range h.hashCh {
				actualHashes++

				if !tt.httpRequestFails {
					require.Equal(t, getMD5Hash([]byte(_testHTTPBody)), result.hash)
				}
			}

			require.Equal(t, tt.expectedHashes, actualHashes)
		})
	}
}

// TODO: add tests for the parallel limit
