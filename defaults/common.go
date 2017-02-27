package defaults

import (
	"crypto/sha1"
	"encoding/base64"
	"fmt"
	"strconv"
	"time"
)

func CurrentTime() string {
	return time.Now().Format("2006-01-02T15:04:05Z")
}

func BumpVersion(ver string) (string, error) {
	if i, err := strconv.ParseInt(ver[2:], 10, 64); err != nil {
		return "", err
	} else {
		return fmt.Sprintf("W\"%d", i+1), nil
	}
}

func GenerateNewVersion(id string) string {
	hash := sha1.New()
	now := time.Now().Format(time.RFC3339Nano)
	hash.Write([]byte(id))
	hash.Write([]byte(now))
	return fmt.Sprintf("W\\/\"%s\"", base64.StdEncoding.EncodeToString(hash.Sum(nil)))
}

func Min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
