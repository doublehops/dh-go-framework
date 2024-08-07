package testbuffer

import (
	"fmt"
	"math/rand"
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestBuffer is to help with testing and not really needs to be tested itself. However
func TestRead(t *testing.T) {
	tests := []struct {
		name               string
		filename           string
		writeString        []byte
		expectedWriteError error
		expectedReadError  error
		expectedLen        int
		expectedValue      []byte
	}{
		{
			name:               "testSuccess",
			filename:           fmt.Sprintf("/tmp/testbuffer-%d", rand.Int()),
			writeString:        []byte("This is the test string"),
			expectedWriteError: nil,
			expectedReadError:  nil,
			expectedLen:        23,
			expectedValue:      []byte("This is the test string"),
		},
		{
			name:               "testFailToWriteAndRead",
			filename:           fmt.Sprintf("/no-permission-here"),
			writeString:        []byte("This is the test string"),
			expectedWriteError: fmt.Errorf("unable to write temp file: /no-permission-here. open /no-permission-here: permission denied"),
			expectedReadError:  fmt.Errorf("unable to read file: /no-permission-here. open /no-permission-here: no such file or directory"),
			expectedLen:        0,
			expectedValue:      []byte(""),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			Filename = tt.filename
			buf := TestBuffer{}
			length, err := buf.Write(tt.writeString)

			if err != nil || tt.expectedReadError != nil {
				assert.EqualError(t, err, tt.expectedWriteError.Error(), "error writing to buffer")
			}
			assert.Equal(t, tt.expectedLen, length)

			val, err := buf.Read()
			if err != nil || tt.expectedWriteError != nil {
				assert.Error(t, err, tt.expectedReadError.Error(), "error reading from buffer")
			}
			assert.Equal(t, tt.expectedValue, val)
		})
	}
}
