package server

import (
	"project_yd/util"
	"time"

	"github.com/go-redis/redis"
)

var RedisManager *redis.Client

func RedisConnection() {
	println("Redis Connect!!")
	RedisManager = redis.NewClient(&redis.Options{
		Addr:     util.RedisIp + ":" + util.RedisPort, // Redis 서버 주소 및 포트
		Password: "",                                  // 패스워드 (설정된 경우)
		DB:       0,                                   // 데이터베이스 인덱스 (기본값: 0)
	})
	println("Redis Connect End!!")
}

func SetRedis(key, value string) error {
	err := RedisManager.Set(key, value, 0).Err()
	if err != nil {
		println("Set Redis Error")
		return err
	}
	return nil
}

// -- Redis에 키값이 없을때와 있을때 벨류값 다르게 처리
func SetNxRedis(key, value string, newValue string) error {
	setRedis, err := RedisManager.SetNX(key, value, 0).Result()
	if err != nil {
		println("Set Redis Error")
		return err
	}
	//-- 키값이 존재하지않아 새롭게 추가
	if setRedis {
		err := SetRedis(key, value)
		if err != nil {
			return err
		}
		return nil
	}

	{
		err := SetRedis(key, newValue)
		if err != nil {
			return err
		}
	}

	return nil
}
func GetRedis(key string) (string, error) {
	value, err := RedisManager.Get(key).Result()
	if err != nil {
		println("Get Redis Error")
		return "", err
	}
	return value, nil
}
func HasRedisKey(key string) bool {

	exists, err := RedisManager.Exists(key).Result()
	if err != nil {
		println("GetRedisKey Error")
		println(err.Error())
		return false
	}
	if exists == 1 {
		return true
	}
	return false
}

const HEARTBEAT_KEY = "HeartBeat_"

func SetHeartBeat(UUID string) (string, error) {
	key := HEARTBEAT_KEY + UUID
	currentTime := time.Now()
	layout := "2006-01-02 15:04:05"
	formattedTime := currentTime.Format(layout)
	value := HEARTBEAT_KEY + formattedTime
	err := SetRedis(key, value)
	if err != nil {
		println("Set HeartBeat Error")
		return "", err
	}
	println("SetHeartBeat")
	return value, nil
}
func GetHeartBeat(UUID string) (string, error) {
	key := HEARTBEAT_KEY + UUID
	heartBeat, err := GetRedis(key)
	if err != nil {
		println("GetHeartBeat Error!!")
		return "", err
	}
	println("GetHeartBeat")
	return heartBeat, nil
}
func HasHeartBeat(UUID string) bool {
	key := HEARTBEAT_KEY + UUID
	result := HasRedisKey(key)
	println("Check->HasHeartBeat")
	return result
}
