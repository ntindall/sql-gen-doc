package main

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestPadRemainingWidth(t *testing.T) {
	testcases := []struct {
		desc        string
		inputString string
		inputWidth  int
		expectation string
		shouldPanic bool
	}{
		{
			desc:        "should pad width - len(s)",
			inputString: "hello world", // 11 characters
			inputWidth:  13,
			expectation: "  ",
		},
		{
			desc:        "returns nothing if width == len(s)",
			inputString: "hello world", // 11 characters
			inputWidth:  11,
			expectation: "",
		},
		{
			desc:        "panicks if width < len(s)",
			inputString: "hello world", // 11 characters
			inputWidth:  10,
			shouldPanic: true,
		},
	}

	for i, tc := range testcases {
		func() {
			defer func() {
				if p := recover(); p != nil {
					require.True(t, tc.shouldPanic, fmt.Sprintf("unexpected recovered from panic: %v", p))
					assert.Equal(t, p, "strings: negative Repeat count")
				}
			}()

			t.Logf("test case %d: %s", i, tc.desc)

			actual := padRemainingWidth(tc.inputString, tc.inputWidth)
			assert.Equal(t, tc.expectation, actual)
		}()
	}
}
