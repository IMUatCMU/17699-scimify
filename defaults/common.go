package defaults

import (
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

func Min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
