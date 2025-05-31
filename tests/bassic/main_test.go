package bassic

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAddOne(t *testing.T) {
	// var (
	// 	input = 1
	// 	output = 3
	// )
	// actual := AddOne(input)
	// if actual != output {
	// 	t.Errorf("AddOne(%d), input %d, actual %d, expected %d", input, input, actual, output)
	// }
	// usage of assert from testify
	assert.Equal(t, AddOne(1), 3, "AddOne should return 3")
	assert.NotEqual(t, 2, 3, "AddOne should not return 2")
	assert.Nil(t, nil, nil)
}

func TestRequire(t *testing.T) {
	require.Equal(t, 2, 3)
	fmt.Println("not Excuted")
}
func TestAssert(t *testing.T) {
	assert.Equal(t, 2, 2, "2 should not be equal to 3")
	fmt.Println("not Excuted")
}

func TestAddOne2(t *testing.T) {
	assert.Equal(t, AddOne2(1), 2, "AddOne2 should return 2")
}