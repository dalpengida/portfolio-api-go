package redis

import "time"

var (
	SCAN_MAX_COUNT int64 = 1000
	SESSION_EXPIRE       = time.Minute * 60
)

func SessionKey(userId string) string {
	return "session:" + userId
}
