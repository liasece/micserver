package math

import (
	"crypto/hmac"
	"crypto/sha1"
	"crypto/sha256"
	"encoding/hex"
)

// HmacSha256 func
func HmacSha256(content []byte, key []byte) string {
	mac := hmac.New(sha256.New, key)
	_, err := mac.Write(content)
	if err != nil {
	}
	return hex.EncodeToString(mac.Sum(nil))
}

// HmacSha256ByString func
func HmacSha256ByString(contentstr string, keystr string) string {
	content := []byte(contentstr)
	key := []byte(keystr)
	return HmacSha256(content, key)
}

// HmacSha1 func
func HmacSha1(content []byte, key []byte) string {
	//hmac ,use sha1
	mac := hmac.New(sha1.New, key)
	// mac := hmac.New(md5.New, key)
	_, err := mac.Write(content)
	if err != nil {
	}
	return string(mac.Sum(nil))
}
