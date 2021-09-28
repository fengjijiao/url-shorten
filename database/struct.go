package database

import (
	"encoding/gob"
	"goto/commonio"
	"io"
	"log"
	"os"
	"sync"
)

type URLStore struct {
	urls map[string]string
	mu sync.RWMutex
	file *os.File
}

type record struct {
	Key, URL string
}

func NewURLStore(filename string) *URLStore {
	s := &URLStore{urls: make(map[string]string)}
	f, err := os.OpenFile(filename, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		panic(err)
	}
	s.file = f
	if err := s.load(); err != nil {
		log.Panicln("Error loading data in URLStore: ", err)
	}
	return s
}

func (s *URLStore) Exist(key string) (isPresent bool) {
	_, isPresent = s.urls[key]
	return
}

func (s *URLStore) Get(key string) (url string) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	url = s.urls[key]
	return
}

func (s *URLStore) Set(key string, url string) (res bool) {
	res = false
	if !s.Exist(key) {
		s.mu.Lock()
		defer s.mu.Unlock()
		s.urls[key] = url
		res = true
	}
	return
}

func (s *URLStore) Put(url string) string {
	for {
		key := commonio.RandomString(9)
		if s.Set(key, url) {
			if err := s.save(); err != nil {
				log.Println("Error saving to URLStore: ", err)
			}
			return key
		}
	}
	panic("should not get here")
}

func (s *URLStore) Count() int {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return len(s.urls)
}

func (s *URLStore) save() error {
	e := gob.NewEncoder(s.file)
	var res []record = make([]record, len(s.urls))
	i := 0
	for key, value := range s.urls {
		res[i] = record{Key: key, URL: value}
		i++
	}
	return e.Encode(res)
}

func (s *URLStore) load() error {
	if _, err := s.file.Seek(0, 0); err != nil {
		return err
	}
	d := gob.NewDecoder(s.file)
	var err error
	for err == nil {
		var r []record
		if err = d.Decode(&r); err == nil {
			for _, r2 := range r {
				s.Set(r2.Key, r2.URL)
			}
		}
	}
	if err == io.EOF {
		return nil
	}
	return err
}