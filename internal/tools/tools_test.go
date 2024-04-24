package tools

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

// TestCurrentFunction will test `currentFunction`. Because of the way tests are ran, it won't actually return the
// name of this test, therefore it's hard-coded as the actual caller.
func TestCurrentFunction(t *testing.T) {
	thisFunc := CurrentFunction()
	assert.Equal(t, "testing.tRunner", thisFunc, "func name not as expected")
}
