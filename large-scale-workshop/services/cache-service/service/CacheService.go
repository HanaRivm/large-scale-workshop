package CacheService

import (
	"log"

	"github.com/TAULargeScaleWorkshop/HANA/large-scale-workshop/services/registry-service/servant/dht"
)

type CacheServiceServer struct {
	UnimplementedCacheServiceServer
	chord dht.Chord
}

func NewCacheServiceServer() *CacheServiceServer {
	return &CacheServiceServer{}
}

func NewCacheService(nodeName string, port int32, rootNodeName string) (*CacheService, error) {
	var chord *dht.Chord
	var err error
	if rootNodeName == "" {
		chord, err = dht.NewChord(nodeName, port)
	} else {
		chord, err = dht.JoinChord(nodeName, rootNodeName, port)
	}
	if err != nil {
		return nil, err
	}

	return &CacheService{
		chord: chord,
	}, nil
}

func (cs *CacheService) Set(key, value string) {
	cs.mu.Lock()
	defer cs.mu.Unlock()
	cs.chord.Set(key, value)
}

func (cs *CacheService) Get(key string) (string, error) {
	cs.mu.Lock()
	defer cs.mu.Unlock()
	return cs.chord.Get(key)
}

func (cs *CacheService) Delete(key string) {
	cs.mu.Lock()
	defer cs.mu.Unlock()
	cs.chord.Delete(key)
}

func (cs *CacheService) IsAlive() bool {
	alive, err := cs.chord.IsFirst()
	if err != nil {
		log.Println("Error checking IsAlive:", err)
		return false
	}
	return alive
}
