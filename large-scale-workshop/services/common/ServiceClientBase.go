package common

import (
	"fmt"
	"log"
	"math/rand"
	"time"

	"google.golang.org/grpc"
)

type ServiceClientBase[client_t any] struct {
	RegistryAddresses []string
	CreateClient      func(grpc.ClientConnInterface) client_t
}

func (obj *ServiceClientBase[client_t]) PickNode() string {
	rand.Seed(time.Now().UnixNano())
	return obj.RegistryAddresses[rand.Intn(len(obj.RegistryAddresses))]
}

func (obj *ServiceClientBase[client_t]) Connect() (res client_t, closeFunc func(), err error) {
	registryAddress := obj.PickNode()
	log.Printf("picking node%s", registryAddress)
	conn, err := grpc.Dial(registryAddress, grpc.WithInsecure(), grpc.WithBlock(), grpc.WithTimeout(time.Minute))
	if err != nil {
		var empty client_t
		return empty, nil, fmt.Errorf("failed to connect client to %v: %v", obj.RegistryAddresses, err)
	}
	log.Printf("picking node2")

	c := obj.CreateClient(conn)
	return c, func() { conn.Close() }, nil
}
