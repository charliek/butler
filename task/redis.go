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

func initializeStatus(status *TaskStatus) error {
	conn := pool.Get()
	defer conn.Close()
	_, err := conn.Do("LPUSH", listKey, status.Id)
	if err != nil {
		return err
	}
	return updateStatus(status)
}

func updateStatus(status *TaskStatus) error {
	conn := pool.Get()
	defer conn.Close()
	bytes, err := json.Marshal(&status)
	if err != nil {
		log.Warn("Error marshalling json for task status %v", err)
	}
	return updateKey(status.redisStatusKey(), bytes)
}

func updateKey(key string, b []byte) error {
	conn := pool.Get()
	defer conn.Close()
	_, err := conn.Do("SET", key, b)
	if err != nil {
		log.Warn("Error updating redis key '%s' with error '%v'", key, err)
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

func ListTasks() ([]*TaskStatus, error) {
	conn := pool.Get()
	defer conn.Close()
	idList, err := redis.Strings(conn.Do("LRANGE", listKey, 0, 100))
	if err != nil {
		return nil, err
	}
	idLen := len(idList)
	log.Info("List length %v", idLen)
	taskList := make([]*TaskStatus, idLen, idLen)
	for idx, id := range idList {
		status, err := LookupTask(id)
		if err != nil {
			return nil, err
		}
		taskList[idx] = status
		log.Info("List index %v", idx)
	}
	return taskList, nil
}

func LookupTask(id string) (*TaskStatus, error) {
	conn := pool.Get()
	defer conn.Close()
	data, err := redis.Bytes(conn.Do("GET", statusKeyFromId(id)))
	if err != nil {
		return nil, err
	}
	var status TaskStatus
	err = json.Unmarshal(data, &status)
	if err != nil {
		return nil, err
	}
	return &status, nil
}
