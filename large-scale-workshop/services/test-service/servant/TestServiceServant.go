package TestServiceServant

import "fmt"

func HelloWorld() string {
	return "Hello World"
}

func HelloToUser(username string) string {
	return fmt.Sprintf("Hello %s", username)
}
