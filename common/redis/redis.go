package redis

import (
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

/*
https://github.com/redis/go-redis 는 redis 업체에서 정식 verify 받은 repository
예전에는 "https://github.com/gomodule/redigo" 를 사용을 했었지만, redis 업체에서 직접 제대로 만들어서 공유를 해주고
예전 보다 많이 최적화가 되어 있다고 하고 추천을 많이 하길래 한번 츄라이를 해봄

redigo의 경우는 pool 을 직접 개발자가 어느 정도 명령어를 사용하여 사용을 하지만,
go-redis 의 경우는 자체적으로 내부에서 pooling 을 한다고 함

일단 테스트를 거쳐서 진행을 해볼 예정
참고 : https://redis.uptrace.dev/guide/go-redis-vs-redigo.html

go-redis 의 경우는 conn() 을 따로 일부러 부르면 오히려 장애가 생김
알아서 pooling 을 해주기 떄문에 client 에서 실행을 하고
connection 마다 이름을 정하거나 확인을 하거나 특정 작업을 위한 경우가 아니면 모듈에 맡기면 됨
*/

// var host string

var (
	client     *redis.Client
	redis_host string
)

func init() {
	redis_host = ":6379"
	client = conn()
}

func conn() *redis.Client {
	return redis.NewClient(&redis.Options{
		MaxIdleConns:    10,
		ConnMaxIdleTime: 240 * time.Second,
		ReadTimeout:     10 * time.Second,
		WriteTimeout:    10 * time.Second,
		DialTimeout:     10 * time.Second,
		PoolTimeout:     10 * time.Second,
		PoolSize:        1000,
		Addr:            redis_host,
		Password:        "", // no password set
		DB:              0,  // use default DB
		// TLSConfig: &tls.Config{ // aws elasticache 에서만 사용
		// 	MinVersion: tls.VersionTLS12,
		// 	//Certificates: []tls.Certificate{cert}
		// },
		OnConnect: func(ctx context.Context, cn *redis.Conn) error {
			_, err := cn.Ping(ctx).Result()
			if err != nil {
				panic(err) // ping 이 안된다는 것은 redis를 사용할 수 없다는 것
			}

			fmt.Println("connenction success")
			return nil
		},
	})

}

func GetClient() *redis.Client {
	if client == nil {
		conn()
	}

	return client
}

func Get(c context.Context, key string) (string, error) {
	return client.Get(c, key).Result()
}

func GetInterface(c context.Context, key string, i interface{}) error {
	r := client.Get(c, key)
	if r.Err() != nil {
		return r.Err()
	}

	err := r.Scan(i) // unmarchaling 임
	if err != nil {
		return err
	}

	return nil
}

func Set(c context.Context, key, value string) error {
	return client.Set(c, key, value, redis.KeepTTL).Err()
}

func SetInterface(c context.Context, key string, value interface{}) error {
	return client.Set(c, key, value, redis.KeepTTL).Err()
}

func SetEX(c context.Context, key string, value interface{}, ttl time.Duration) error {
	return client.SetEx(c, key, value, ttl).Err()
}

func Del(c context.Context, key string) (int64, error) {
	return client.Del(c, key).Result()
}

// TODO:
// cursor 값을 외부에서 사용하며, 다시 스캔을 돌릴 수 있는 것을 외부에서 판단을 할 수 있도록 변경을 해야 함
// ex) return []string, cursor, error
func Scan(c context.Context, keyword string) ([]string, error) {
	var cursor uint64
	foundKeys := make([]string, 0)

	for {
		var r []string
		var err error
		r, cursor, err = client.Scan(c, cursor, keyword, SCAN_MAX_COUNT).Result()
		if err != nil {
			return nil, err
		}

		foundKeys = append(foundKeys, r...)
		if cursor == 0 {
			break
		}
	}

	return foundKeys, nil
}

func Publish(c context.Context, channel string, message interface{}) error {
	return client.Publish(c, channel, message).Err()
}

func SubScribe(c context.Context, channel string) *redis.PubSub {
	return client.Subscribe(c, channel)
}

func ReceiveMessage(c context.Context, channel string) (*redis.Message, error) {
	pubsub := SubScribe(c, channel)
	return pubsub.ReceiveMessage(c)
}

func Sadd(c context.Context, key string, values ...interface{}) error {
	r := client.SAdd(c, key, values)
	if r.Err() != nil {
		return r.Err()
	}

	return nil
}

func SaddString(c context.Context, key string, values ...string) error {
	r := client.SAdd(c, key, values)
	if r.Err() != nil {
		return r.Err()
	}

	return nil
}

// Srem 는 set 데이터 타입의 멤버 값을 삭제를 하는 함수
// set 데이터 타입의 멤버를 삭제 할 경우, 하나의 key값에 멤버가 없을 때, 키 값이 모두 삭제가 됨
func Srem(c context.Context, key string, value ...interface{}) error {
	r := client.SRem(c, key, value)
	if r.Err() != nil {
		return r.Err()
	}

	return nil
}

// Smembers 는 set 데이터 타입 키 값으로 멤버들을 조회 해주는 함수
func Smembers(c context.Context, key string, valueStruct interface{}) error {
	r := client.SMembers(c, key)
	if r.Err() != nil {
		return r.Err()
	}

	err := r.ScanSlice(valueStruct)
	if err != nil {
		return err
	}

	return nil
}

func Sismember(c context.Context, key string, value interface{}) (bool, error) {
	return client.SIsMember(c, key, value).Result()
}
