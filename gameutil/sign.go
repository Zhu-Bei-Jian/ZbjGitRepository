package gameutil

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"sort"
)

func GetSignContent(params map[string]interface{}, exceptKey ...string) string {
	keys := make([]string, 0, len(params))
	for k := range params {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	var signContent string
	first := true
	for _, key := range keys {
		except := false
		for _, v := range exceptKey {
			if v == key {
				except = true
				break
			}
		}

		if except {
			continue
		}

		if !first {
			signContent += "&"
		}
		first = false

		signContent += fmt.Sprintf("%s=%v", key, params[key])
	}
	return signContent
}

func HmacSHA256(content string, key string) string {
	mac := hmac.New(sha256.New, []byte(key))
	mac.Write([]byte(content))
	return hex.EncodeToString(mac.Sum(nil))
}
