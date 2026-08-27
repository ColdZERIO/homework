package storage

import (
	"fmt"
	"time"
)

const ttl = 5 * time.Minute

type QueryParams struct {
	limit   int
	offset  int
	where   string
	orderby string
}

func KeyCacheID(id string) string {
	return fmt.Sprintf("user:id:%s", id)
}

// Добавить функцию с параметрами запроса как ключ
func KeyCacheList(q QueryParams) string {
	return fmt.Sprintf("user:list:%d:%d:%s:%s", q.limit, q.offset, q.where, q.orderby)
}
