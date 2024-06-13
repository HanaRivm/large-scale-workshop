package TestServiceServant

import (
	"fmt"
	"math/rand"
	"time"
)

var cacheMap map[string]string

func HelloWorld() string {
	return "Hello World"
}

func HelloToUser(username string) string {
	return fmt.Sprintf("Hello %s", username)
}

func init() {
	cacheMap = make(map[string]string)
}

func Store(key, value string) {
	cacheMap[key] = value
}

func Get(key string) (string, bool) {
	value, exists := cacheMap[key]
	return value, exists
}

func WaitAndRand(seconds int32, sendToClient func(x int32) error) error {
	time.Sleep(time.Duration(seconds) * time.Second)
	return sendToClient(int32(rand.Intn(10)))
}
