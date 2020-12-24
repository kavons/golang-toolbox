package redis_cluster_test

import (
	"encoding/json"
	"fmt"
	"log"
	"strconv"
	"testing"
	"time"

	"github.com/astaxie/beego"
	"github.com/go-redis/redis"
)

var (
	//conn *redis.ClusterClient
	conn *redis.Client
)

func init() {
	//conn = redis.NewClusterClient(&redis.ClusterOptions{
	//	Addrs:    []string{"192.168.1.232:30370", "192.168.1.233:30371", "192.168.1.234:30372"},
	//	Password: "redis2019",
	//})
	conn = redis.NewClient(&redis.Options{
		Addr:     "127.0.0.1:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	_, err := conn.Ping().Result()
	if err != nil {
		log.Fatalf("redis cluster connection failed: %v", err)
	}
}

func RedisSet(key string, value interface{}, expire int) error {
	if expire > 0 {
		err := conn.Do("SET", key, value, "PX", expire).Err()
		if err != nil {
			beego.Error("RedisSet Error! key:", key, "Details:", err.Error())
			return err
		}
	} else {
		err := conn.Do("SET", key, value).Err()
		if err != nil {
			beego.Error("RedisSet Error! key:", key, "Details:", err.Error())
			return err
		}
	}

	return nil
}

func TestRedisSet(t *testing.T) {
	err := RedisSet("name", "Jack", 0)
	if err != nil {
		log.Fatalf("RedisSet Error!")
	}
}

func RedisGet(key string) (string, error) {
	value, err := conn.Do("GET", key).String()
	if err != nil {
		return "", nil
	}

	return value, nil
}

func TestRedisGet(t *testing.T) {
	name, err := RedisGet("name")
	if err != nil {
		log.Fatalf("RedisGet Error!")
	}

	fmt.Println(name)
}

func RedisExpire(key string, expire int) error {
	err := conn.Do("EXPIRE", key, expire).Err()
	if err != nil {
		beego.Error("RedisExpire Error!", key, "Details:", err.Error())
		return err
	}

	return nil
}

func TestRedisExpire(t *testing.T) {
	err := RedisExpire("name", 3600)
	if err != nil {
		log.Fatalf("RedisExpire Error!")
	}
}

func RedisPTTL(key string) (int, error) {
	ttl, err := conn.Do("PTTL", key).Int()
	if err != nil {
		return -1, err
	}

	return ttl, nil
}

func TestRedisPTTL(t *testing.T) {
	ttl, err := RedisPTTL("name")
	if err != nil {
		log.Fatalf("RedisPTTL Error!")
	}

	fmt.Println(ttl)
}

func RedisTTL(key string) (int, error) {
	ttl, err := conn.Do("TTL", key).Int()
	if err != nil {
		return -1, err
	}

	return ttl, nil
}

func TestRedisTTL(t *testing.T) {
	ttl, err := RedisTTL("name")
	if err != nil {
		log.Fatalf("RedisTTL Error!")
	}

	fmt.Println(ttl)
}

type Person struct {
	Name string `json:"name"`
	Age  int    `json:"age"`
}

func RedisSetJson(key string, value interface{}, expire int) error {
	jsonData, _ := json.Marshal(value)
	if expire > 0 {
		err := conn.Do("SET", key, jsonData, "PX", expire).Err()
		if err != nil {
			beego.Error("RedisSetJson Error! key:", key, "Details:", err.Error())
			return err
		}
	} else {
		err := conn.Do("SET", key, jsonData).Err()
		if err != nil {
			beego.Error("RedisSetJson Error! key:", key, "Details:", err.Error())
			return err
		}
	}

	return nil
}

func TestRedisSetJson(t *testing.T) {
	person := Person{
		"Snow", 35,
	}

	err := RedisSetJson("person", person, 0)
	if err != nil {
		log.Fatalf("RedisTTL Error!")
	}
}

func RedisGetJson(key string) ([]byte, error) {
	value, err := conn.Do("GET", key).String()
	if err != nil {
		return nil, nil
	}

	return []byte(value), nil
}

func TestRedisGetJson(t *testing.T) {
	personBytes, err := RedisGetJson("person")
	if err != nil {
		log.Fatalf("RedisTTL Error!")
	}

	person := new(Person)
	json.Unmarshal(personBytes, &person)
	fmt.Println(person.Name, person.Age)
}

func RedisDel(key string) error {
	err := conn.Do("DEL", key).Err()
	if err != nil {
		beego.Error("RedisDel Error! key:", key, "Details:", err.Error())
	}
	return err
}

func TestRedisDel(t *testing.T) {
	err := RedisDel("person")
	if err != nil {
		log.Fatalf("RedisDel Error!")
	}
}

func RedisHGet(key, field string) (string, error) {
	value, err := conn.Do("HGET", key, field).String()
	if err != nil {
		return "", nil
	}

	return value, nil
}

func TestRedisHGet(t *testing.T) {
	token, err := RedisHGet("token", "device")
	if err != nil {
		log.Fatalf("RedisHGet Error!")
	}

	fmt.Println(token)
}

func RedisHSet(key, field, value string) error {
	err := conn.Do("HSET", key, field, value).Err()
	if err != nil {
		beego.Error("RedisHSet Error!", key, "field:", field, "Details:", err.Error())
	}
	return err
}

func TestRedisHSet(t *testing.T) {
	err := RedisHSet("token", "device", "android")
	if err != nil {
		log.Fatalf("RedisHSet Error!")
	}
}

func RedisHDel(key, field string) error {
	err := conn.Do("HDEL", key, field).Err()
	if err != nil {
		beego.Error("RedisHDel Error!", key, "field:", field, "Details:", err.Error())
	}
	return err
}

func TestRedisHDel(t *testing.T) {
	err := RedisHDel("token", "device")
	if err != nil {
		log.Fatalf("RedisHDel Error!")
	}
}

func RedisZAdd(key, member, score string) error {
	err := conn.Do("ZADD", key, score, member).Err()
	if err != nil {
		beego.Error("RedisZAdd Error!", key, "member:", member, "score:", score, "Details:", err.Error())
	}
	return err
}

func TestRedisZAdd(t *testing.T) {
	timeStamp := int(time.Now().Unix())
	score := strconv.Itoa(timeStamp)
	err := RedisZAdd("close_test_attendees", "1", score)
	if err != nil {
		log.Fatalf("RedisZAdd Error!")
	}
}

func RedisZRank(key, member string) (int, error) {
	rank, err := conn.Do("ZRANK", key, member).Int()
	if err != nil {
		beego.Error("RedisZRank Error!", key, "member:", member, "Details:", err.Error())
		return -1, nil
	}

	return rank, err
}

func TestRedisZRank(t *testing.T) {
	rank, err := RedisZRank("close_test_attendees", "3")
	if err != nil {
		log.Fatalf("RedisZAdd Error!")
	}

	fmt.Println(rank)
}

func RedisZRange(key string, start, stop int) (values []string, err error) {
	values, err = conn.ZRange(key, int64(start), int64(stop)).Result()
	if err != nil {
		beego.Error("RedisZRange Error!", key, "start:", start, "stop:", stop, "Details:", err.Error())
		return
	}

	return
}

func TestRedisZRange(t *testing.T) {
	values, err := RedisZRange("close_test_attendees", 0, -1)
	if err != nil {
		log.Fatalf("RedisZRange Error!")
	}

	fmt.Println(values)
}

func RedisZRangeWithScores(key string, start, stop int) (values []redis.Z, err error) {
	values, err = conn.ZRangeWithScores(key, int64(start), int64(stop)).Result()
	if err != nil {
		beego.Error("RedisZRange Error!", key, "start:", start, "stop:", stop, "Details:", err.Error())
		return
	}

	return
}

func TestRedisZRangeWithScores(t *testing.T) {
	values, err := RedisZRangeWithScores("close_test_attendees", 0, -1)
	if err != nil {
		log.Fatalf("RedisZRange Error!")
	}

	for _, v := range values {
		timestamp := int64(v.Score)
		member := v.Member.(string)
		fmt.Println(timestamp, member)
	}
}

func RedisZRem(key, member string) error {
	err := conn.Do("ZREM", key, member).Err()
	if err != nil {
		beego.Error("RedisZRem Error!", key, "member:", member, "Details:", err.Error())
	}
	return err
}

func TestRedisZRem(t *testing.T) {
	err := RedisZRem("close_test_attendees", "3")
	if err != nil {
		log.Fatalf("RedisZRem Error!")
	}
}

func RedisRPUSH(key string, member string) (err error) {
	err = conn.Do("RPUSH", key, member).Err()
	if err != nil {
		beego.Error("RedisRPUSH Error!", key, member, "Details:", err.Error())
		return
	}

	return
}

func TestRedisRPUSH(t *testing.T) {
	err := RedisRPUSH("list", "3")
	if err != nil {
		log.Fatalf("RedisRPush Error!")
	}
}

func RedisBLPOP(key string, timeout time.Duration) (value []string, err error) {
	value, err = conn.BLPop(timeout, key).Result()
	if err != nil {
		beego.Error("RedisBLPOP Error!", key, timeout, "Details:", err.Error())
		return
	}

	return
}

func TestRedisBLPOP(t *testing.T) {
	value, err := RedisBLPOP("list", 15*time.Second)
	if err != nil {
		log.Fatalf("RedisRPush Error!")
	}

	fmt.Println(value)
}
