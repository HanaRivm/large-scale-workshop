package dht

import (
	"log"

	metaffi "github.com/MetaFFI/lang-plugin-go/api"
	metaffiruntime "github.com/MetaFFI/lang-plugin-go/go-runtime"
	"github.com/MetaFFI/plugin-sdk/compiler/go/IDL"
)

var openjdkRuntime *metaffi.MetaFFIRuntime
var chordModule *metaffi.MetaFFIModule
var newChord func(...interface{}) ([]interface{}, error)
var joinChord func(...interface{}) ([]interface{}, error)
var set func(...interface{}) ([]interface{}, error)
var get func(...interface{}) ([]interface{}, error)
var pdelete func(...interface{}) ([]interface{}, error)
var getAllKeys func(...interface{}) ([]interface{}, error)
var isFirst func(...interface{}) ([]interface{}, error)

func init() {
	openjdkRuntime = metaffi.NewMetaFFIRuntime("openjdk")
	chordModule, err := openjdkRuntime.LoadModule("./dht/Chord.class")
	newChord, err = chordModule.Load("class=dht.Chord,callable=<init>",
		[]IDL.MetaFFIType{IDL.STRING8, IDL.INT32},
		[]IDL.MetaFFIType{IDL.HANDLE})
	if err != nil {
		panic(err)
	}
	joinChord, err = chordModule.Load("class=dht.Chord,callable=<init>",
		[]IDL.MetaFFIType{IDL.STRING8, IDL.STRING8, IDL.INT32},
		[]IDL.MetaFFIType{IDL.HANDLE})
	if err != nil {
		panic(err)
	}
	set, err =
		chordModule.Load("class=dht.Chord,callable=set,instance_required", []IDL.MetaFFIType{IDL.HANDLE, IDL.STRING8, IDL.STRING8}, nil)
	if err != nil {
		panic(err)
	}
	getAllKeys, err =
		chordModule.LoadWithAlias("class=dht.Chord,callable=getAllKeys,instance_required",
			[]IDL.MetaFFITypeInfo{{StringType: IDL.HANDLE}},
			[]IDL.MetaFFITypeInfo{{StringType: IDL.STRING8_ARRAY, Dimensions: 1}})
	if err != nil {
		panic(err)
	}
	// Load get method
	get, err = chordModule.Load("class=dht.Chord,callable=get,instance_required",
		[]IDL.MetaFFIType{IDL.HANDLE, IDL.STRING8}, []IDL.MetaFFIType{IDL.STRING8})
	if err != nil {
		log.Fatalf("Failed to load get: %v", err)
	}
	// Load delete method
	pdelete, err = chordModule.Load("class=dht.Chord,callable=delete,instance_required",
		[]IDL.MetaFFIType{IDL.HANDLE, IDL.STRING8}, nil)
	if err != nil {
		log.Fatalf("Failed to load delete: %v", err)
	}
	chordModule.Load("class=dht.Chord,field=isFirst,getter,instance_required",
		[]IDL.MetaFFIType{IDL.HANDLE},
		[]IDL.MetaFFIType{IDL.BOOL})
	if err != nil {
		panic(err)
	}
	isFirst, err = chordModule.Load("class=dht.Chord,field=isFirst,getter,instance_required", []IDL.MetaFFIType{IDL.HANDLE}, []IDL.MetaFFIType{IDL.BOOL})
	if err != nil {
		log.Fatalf("Failed to load isFirst: %v", err)
	}
}

type Chord struct {
	handle metaffiruntime.MetaFFIHandle
}

func NewChord(name string, port int32) (*Chord, error) {
	h, err := newChord(name, port)
	if err != nil {
		return nil, err
	}
	c := &Chord{}
	c.handle = h[0].(metaffiruntime.MetaFFIHandle)
	return c, nil
}
func JoinChord(name string, rootNodeName string, port int32) (*Chord, error) {
	h, err := joinChord(name, rootNodeName, port)
	if err != nil {
		return nil, err
	}
	c := &Chord{}
	c.handle = h[0].(metaffiruntime.MetaFFIHandle)
	return c, nil
}
func (c *Chord) IsFirst() (bool, error) {
	res, err := isFirst(c.handle)
	if err != nil {
		return false, err
	}
	return res[0].(bool), nil
}

func (c *Chord) Set(key string, val string) error {
	_, err := set(c.handle, key, val)
	return err
}

func (c *Chord) Get(key string) (string, error) {
	res, err := get(c.handle, key)
	if err != nil {
		return "", err
	}
	return res[0].(string), nil
}

func (c *Chord) Delete(key string) error {
	_, err := pdelete(c.handle, key)
	return err
}

func (c *Chord) GetAllKeys() ([]string, error) {
	res, err := getAllKeys(c.handle)
	if err != nil {
		return nil, err
	}
	keys := res[0].([]string)
	return keys, nil
}