package TestServiceServant

import "fmt"

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
