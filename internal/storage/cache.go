package storage

import (
	"fmt"
	"time"
)

const fiveMinutes =  5 * time.Minute

func KeyCache(id string) string {
	return fmt.Sprintf("user:id:%s", id)
}