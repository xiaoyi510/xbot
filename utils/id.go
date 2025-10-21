package utils

import (
	"fmt"
	"sync/atomic"
	"time"
)

var (
	echoCounter uint64
)

// GenerateEcho 生成唯一的 echo ID，用于 API 调用追踪
func GenerateEcho() string {
	id := atomic.AddUint64(&echoCounter, 1)
	return fmt.Sprintf("%d-%d", time.Now().UnixNano(), id)
}

// GenerateSessionKey 生成会话 key
func GenerateSessionKey(userID, groupID int64) string {
	if groupID == 0 {
		return fmt.Sprintf("user:%d", userID)
	}
	return fmt.Sprintf("group:%d:user:%d", groupID, userID)
}
