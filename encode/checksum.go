package encode

import (
	"crypto/md5"
	"encoding/hex"
)

func Md5(payload interface{}, salt string) string {
	serialisedRequest := StringifyGraceful(payload, false)
	serialisedCase :=  salt + "." + serialisedRequest

	hasher := md5.New()
	hasher.Write([]byte(serialisedCase))

	return hex.EncodeToString(hasher.Sum(nil))
}
