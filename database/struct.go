package database

import "sync"

type URLStore struct {
	urls map[string]string
	mu sync.RWMutex
}

var urlStore *URLStore

func init() {
	urlStore = NewURLStore()
}

func NewURLStore() *URLStore {
	return &URLStore{urls: make(map[string]string)}
}

func Exist(key string) (isPresent bool) {
	_, isPresent = urlStore.urls[key]
	return
}

func Get(key string) (url string) {
	urlStore.mu.RLock()
	defer urlStore.mu.RUnlock()
	url = urlStore.urls[key]
	return
}

func Add(key string, url string) (res bool) {
	res = false
	if !Exist(key) {
		urlStore.mu.Lock()
		defer urlStore.mu.Unlock()
		urlStore.urls[key] = url
		res = true
	}
	return
}