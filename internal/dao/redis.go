package dao

import (
	"encoding/json"
	"fmt"
	"github.com/go-redis/redis/v7"
	"github.com/wyy-go/go-web-template/internal/common/constant"
	"github.com/wyy-go/go-web-template/internal/config"
	"sync"
	"time"
)

const (
	prefixConnKey = "CONN:%s"
)

var (
	redisClient     *redis.Client
	onceRedisClient sync.Once
)

type ConnInfo struct {
	Uin            string `json:"uin"`
	ConnId         string `json:"conn_id"`
	Platform       string `json:"platform"`
	Device         string `json:"device"`
	Server         string `json:"server"`
	LoginTime      int64  `json:"login_time"`
	DisconnectTime int64  `json:"disconnect_time"`
	Status         int    `json:"status"`
}

func connKeyUin(uin string) string {
	return fmt.Sprintf(prefixConnKey, uin)
}

func GetRedisClient() *redis.Client {
	onceRedisClient.Do(func() {
		redisClient = redis.NewClient(&redis.Options{
			Addr:     config.GetConfig().RedisConfig.Addr,
			Password: config.GetConfig().RedisConfig.Password,
			DB:       config.GetConfig().RedisConfig.DB,
		})
	})

	return redisClient
}

func AddConn(uin string, info *ConnInfo) (err error) {
	client := GetRedisClient()
	b, err := json.Marshal(info)
	if err != nil {
		return
	}
	_, err = client.Pipelined(func(pipe redis.Pipeliner) error {
		key := connKeyUin(uin)
		if err := pipe.HSet(key, info.Platform, string(b)).Err(); err != nil {
			return err
		}

		if err := pipe.Expire(key, time.Duration(constant.PushOnlineKeepDays*24)*time.Hour).Err(); err != nil {
			return err
		}
		return nil
	})
	return
}

func DelConn(uin, platform string) (err error) {
	client := GetRedisClient()
	_, err = client.Pipelined(func(pipe redis.Pipeliner) error {
		if err := client.HDel(connKeyUin(uin), platform).Err(); err != nil {
			return err
		}
		return nil
	})

	return
}

func ExpireConn(uin string) (err error) {
	client := GetRedisClient()
	if err = client.Expire(connKeyUin(uin), time.Duration(constant.PushOnlineKeepDays*24)*time.Hour).Err(); err != nil {
		return
	}

	return
}

func GetConnByPlatform(uin, platform string) *ConnInfo {
	if uin == "" || platform == "" {
		return nil
	}
	client := GetRedisClient()

	key := connKeyUin(uin)

	if b, err := client.HGet(key, platform).Bytes(); err != nil {
		return nil
	} else {
		info := &ConnInfo{}
		if err := json.Unmarshal(b, info); err != nil {
			return nil
		}
		return info
	}
}

func GetConnByUin(uin string) (conns map[string][]*ConnInfo, err error) {
	conns = make(map[string][]*ConnInfo)
	client := GetRedisClient()
	r := client.HGetAll(connKeyUin(uin))
	if err = r.Err(); err != nil {
		return
	}

	for _, v := range r.Val() {
		info := ConnInfo{}
		if err := json.Unmarshal([]byte(v), &info); err != nil {
			continue
		}
		conns[info.Server] = append(conns[info.Server], &info)
	}

	return
}
