package redis

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/redis/go-redis/v9"
)

// TestRedisBasic 는 기본 구조 테스트 및 라이브러리 확인을 하기 위함, 기능 테스트
func TestRedisBasic(t *testing.T) {
	t.Skip()

	ctx := context.TODO()
	r := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	err := r.Set(ctx, "test", "test_value", 0).Err()
	if err != nil {
		t.Fatal(err)
	}

	rr, err := r.Set(ctx, "test", "test_value", 0).Result()
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println("set result : ", rr)

	v, err := r.Get(ctx, "test").Result()
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println("test", v)

	v, err = r.Get(ctx, "test2").Result()
	if err == redis.Nil {
		fmt.Println("test2 does not exist")
	} else if err != nil {
		t.Fatal(err)
	} else {
		fmt.Println("test2", v)
	}

	i, err := r.Del(ctx, "test").Result()
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println("test del result : ", i)

	i, err = r.Del(ctx, "alksjdkla").Result()
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println("invalid key del result : ", i)

	keys, err := Scan(ctx, "ss*")
	if err != nil {
		t.Fatal(err)
	}

	fmt.Println(keys)

}

// TestRedisNx 는 redis 함수의 setex setnx 차이를 알기 위하여
func TestRedisNx(t *testing.T) {
	t.Skip()

	ctx := context.TODO()
	r := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	v, err := r.SetEx(ctx, "ex", "test", 10*time.Second).Result()
	if err != nil {
		panic(err)
	}
	fmt.Println(v)

	ttl, err := r.TTL(ctx, "ex").Result()
	if err != nil {
		panic(err)
	}
	fmt.Println(ttl)

	// nxV, err := r.SetNX(ctx, "nx", "test", 0).Result()
	nxV, err := r.SetNX(ctx, "nx", "test", redis.KeepTTL).Result()
	if err != nil {
		panic(err)
	}
	fmt.Println(nxV)
}

// TestPool 는 redis pooling 관련 테스트를 하기 위함
func TestPool(t *testing.T) {
	t.Skip()

	ctx := context.TODO()
	p := GetClient()
	r, err := p.SetEx(ctx, "test", "vvvvvv", 10*time.Second).Result()
	if err != nil {
		panic(err)
	}
	fmt.Println("set ex result : ", r)

	rr, err := p.Get(ctx, "test").Result()
	if err != nil {
		panic(err)
	}
	fmt.Println("get test result : ", rr)

	for i := 0; i < 1000; i++ {
		rr, err := p.Get(ctx, "test").Result()
		if err != nil {
			panic(err)
		}
		fmt.Println("get test result : ", rr)
	}
}

// TestPubSub 는 pub / sub 를 확인 하기 위함
func TestPubSub(t *testing.T) {
	t.Skip() // 무한 대기할 꺼라 스킵

	ctx := context.TODO()
	channel := "testChannel"

	pubsub := SubScribe(ctx, channel)
	defer pubsub.Close()

	// 채널을 이용할 경우
	// ch := pubsub.Channel()
	// for msg := range ch {
	// 	fmt.Println(msg.Channel, msg.Payload)
	// }

	// 그냥 진행 할 경우
	for {
		msg, err := pubsub.ReceiveMessage(ctx)
		if err != nil {
			panic(err)
		}

		fmt.Println(msg.Channel, msg.Payload)
	}
}

// BenchmarkPoolNQuery 는 pool 을 이용했을 때 bench 마크 돌려 봄
func BenchmarkPoolNQuery(b *testing.B) {
	p := GetClient()
	ctx := context.TODO()

	r, err := p.SetEx(ctx, "bench_test", "vvvvvv", 60*time.Second).Result()
	if err != nil {
		panic(err)
	}
	fmt.Println("set ex result : ", r)
	for i := 0; i < 10000; i++ {
		_, err = p.Get(ctx, "bench_test").Result()
		if err != nil {
			panic(err)
		}
	}
}

// TestSxxxxFeature 는 Sadd, Srem, Smembers 등 S 시리즈 테스트
func TestSxxxxFeature(t *testing.T) {
	t.Skip()

	ctx := context.TODO()
	key := "test:test:feature"

	values := []string{"hello", "hi", "test", "hhhhh"}

	v := make([]string, 0)

	err := SaddString(ctx, key, values...)
	if err != nil {
		t.Fatal(err)
	}

	err = Smembers(ctx, key, &v)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(v)

	err = Srem(ctx, key, values[0])
	if err != nil {
		t.Fatal(err)
	}

	v2 := make([]string, 0)
	err = Smembers(ctx, key, &v2)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println("srem result : ", v)

	_, err = Del(ctx, key)
	if err != nil {
		t.Fatal(err)
	}

}
