package component

import (
	"crypto/sha1"
	"fmt"
	"io"
	"sort"
)

//GetSignature 获取签名信息
func GetSignature(msg, token, timestamp, nonce string) (signature string) {
	tmps := []string{token, msg, timestamp, nonce}
	sort.Strings(tmps)
	tmpStr := tmps[0] + tmps[1] + tmps[2] + tmps[3]
	t := sha1.New()
	io.WriteString(t, tmpStr)
	tmp := fmt.Sprintf("%x", t.Sum(nil))
	return tmp
}

// CheckSignature  检查签名
func (c *Component) CheckSignature(signature, msg, timestamp, nonce string) bool {
	msgSignature := GetSignature(msg, c.ComponentToken, timestamp, nonce)
	if msgSignature == signature {
		return true
	}
	return false
}
