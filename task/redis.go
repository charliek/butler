package task

import (
	"encoding/json"
	"github.com/garyburd/redigo/redis"
	log "github.com/ngmoco/timber"
	"time"
)

var redisServer = "localhost:6379"
var pool = &redis.Pool{
	MaxIdle:     3,
	IdleTimeout: 240 * time.Second,
	Dial: func() (redis.Conn, error) {
		c, err := redis.Dial("tcp", redisServer)
		if err != nil {
			return nil, err
		}
		return c, err
	},
	TestOnBorrow: func(c redis.Conn, t time.Time) error {
		_, err := c.Do("PING")
		return err
	},
}

func recordStatus(status *TaskStatus) error {
	conn := pool.Get()
	defer conn.Close()
	bytes, err := json.Marshal(&status)
	if err != nil {
		log.Warn("Error marshalling json for task status %v", err)
	}
	_, err = conn.Do("SET", status.redisStatusKey(), bytes)
	if err != nil {
		log.Warn("Error saving task status to redis %v", err)
	}
	return err
}

func recordLine(status *TaskStatus, line string) error {
	conn := pool.Get()
	defer conn.Close()
	_, err := conn.Do("RPUSH", status.redisLineKey(), line)
	if err != nil {
		log.Warn("Error saving line with key %s to redis %s", status.redisLineKey(), err)
	}
	return err
}
