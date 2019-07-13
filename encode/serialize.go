package encode

import (
	"encoding/json"
	"fmt"
)

func StringifyGraceful(in interface{}, pretty bool) (res string) {
	var serializedData []byte
	var err error
	if pretty {
		serializedData, err = json.MarshalIndent(in, "", "\t")
	} else {
		serializedData, err = json.Marshal(in)
	}

	if err != nil {
		res = fmt.Sprint(in)
		return
	}

	res = string(serializedData)
	return
}
