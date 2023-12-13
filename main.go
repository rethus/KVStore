package main

import (
	"fmt"
	"github.com/labstack/echo/v4"
	"net/http"
	"sync"
)

type Storer[K comparable, V any] interface {
	Put(K, V) error
	Get(K) (V, error)
	Update(K, V) error
	Delete(K) (V, error)
}

type KVStore[K comparable, V any] struct {
	mu   sync.RWMutex
	data map[K]V
}

func NewKVStore[K comparable, V any]() *KVStore[K, V] {
	return &KVStore[K, V]{
		data: make(map[K]V),
	}
}

// check if the given key is in the store
func (s *KVStore[K, V]) Has(key K) bool {
	_, ok := s.data[key]
	return ok
}

func (s *KVStore[K, V]) Put(key K, value V) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	//if s.Has(key) {
	//	return fmt.Errorf("[Put] the key (%v) exists", key)
	//}

	s.data[key] = value

	return nil
}

func (s *KVStore[K, V]) Get(key K) (V, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	value, ok := s.data[key]
	if !ok {
		return value, fmt.Errorf("[Get] the key (%v) does not exists", key)
	}

	return value, nil
}

func (s *KVStore[K, V]) Update(key K, value V) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if !s.Has(key) {
		return fmt.Errorf("[Update] the key (%v) dose not exists", key)
	}

	s.data[key] = value

	return nil
}

func (s *KVStore[K, V]) Delete(key K) (V, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	value, ok := s.data[key]
	if !ok {
		return value, fmt.Errorf("[Delete] the key (%v) does not exists", key)
	}

	delete(s.data, key)

	return value, nil
}

type User struct {
	ID        string
	FirstName string
	Age       int
	Gender    string
}

type Server struct {
	//Storage    Storer[int, *User]
	Storage    Storer[string, string]
	ListenAddr string
}

func NewServer(listenAddr string) *Server {
	return &Server{
		//Storage:    NewKVStore[int, *User](),
		Storage:    NewKVStore[string, string](),
		ListenAddr: listenAddr,
	}
}

func (s *Server) handlePut(c echo.Context) error {
	key := c.Param("key")
	value := c.Param("value")

	s.Storage.Put(key, value)

	return c.JSON(http.StatusOK, map[string]string{"msg": "ok"})
}

func (s *Server) handleGet(c echo.Context) error {
	key := c.Param("key")
	value, err := s.Storage.Get(key)
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, map[string]string{
		"msg":   "ok",
		"value": value,
	})
}

func (s *Server) Start() {
	fmt.Printf("HTTP server is running on port %s", s.ListenAddr)

	e := echo.New()

	e.GET("/put/:key/:value", s.handlePut)
	e.GET("/get/:key", s.handleGet)

	e.Start(s.ListenAddr)
}

func main() {
	s := NewServer(":3000")
	s.Start()
}
