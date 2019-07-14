package encode

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestMd5(t *testing.T) {
	data := struct{
		Key string `json:"key"`
		Val int64 `json:"val"`
	}{
		"one",
		1,
	}

	md5data := Md5(data, "some_namespace")

	//check md5 of "some_namespace.{"key":"one","val":1}" e.g. https://passwordsgenerator.net/md5-hash-generator/
	assert.Equal(t, "f190d3140a6affd97e6e508cd48d6cf2", md5data)
}

