package encode

import (
	"github.com/stretchr/testify/assert"
	"math"
	"testing"
)

var exampleStruct struct {
	Key string `json:"key"`
	Val int64  `json:"val"`
}

func init() {
	exampleStruct = struct {
		Key string `json:"key"`
		Val int64  `json:"val"`
	} {
		"SomeKey",
		22,
	}
}

func TestSimpleJsonSerialisation(t *testing.T) {
	expectedJson := `{"key":"SomeKey","val":22}`
	assert.Equal(t, expectedJson, StringifyGraceful(exampleStruct, false))
}

func TestPrettyJsonSerialisation(t *testing.T) {
	expectedJson := `{
	"key": "SomeKey",
	"val": 22
}`
	assert.Equal(t, expectedJson, StringifyGraceful(exampleStruct, true))
}

func TestScalarSerialisation(t *testing.T) {
	assert.Equal(t, "2", StringifyGraceful(2, false))
	assert.Equal(t, `"lala"`, StringifyGraceful("lala", false))
	assert.Equal(t, `null`, StringifyGraceful(nil, false))
}

func TestJsonIncompatibleSerialisation(t *testing.T) {
	a := math.Inf(1)
	assert.Equal(t, "+Inf", StringifyGraceful(a, false))
}
