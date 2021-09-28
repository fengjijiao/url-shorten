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
	saveChan chan record
}

type record struct {
	Key, URL string
}

func NewURLStore(filename string) *URLStore {
	s := &URLStore{urls: make(map[string]string), saveChan: make(chan record, 1000)}
	if err := s.load(filename); err != nil {
		log.Panicln("Error loading data in URLStore: ", err)
	}
	go s.saveLoop(filename)
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
			s.saveChan <- record{key, url}
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

//func (s *URLStore) save() error {
//	e := gob.NewEncoder(s.file)
//	s.mu.Lock()
//	defer s.mu.Unlock()
//	var res []record = make([]record, len(s.urls))
//	i := 0
//	for key, value := range s.urls {
//		res[i] = record{Key: key, URL: value}
//		i++
//	}
//	return e.Encode(res)
//}

func (s *URLStore) load(filename string) error {
	f, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer f.Close()
	if _, err := f.Seek(0, 0); err != nil {
		return err
	}
	d := gob.NewDecoder(f)
	for err == nil {
		var r record
		if err = d.Decode(&r); err == nil {
			s.Set(r.Key, r.URL)
		}
	}
	if err == io.EOF {
		return nil
	}
	return err
}

func (s *URLStore) saveLoop(filename string) {
	f, err := os.OpenFile(filename, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		panic(err)
	}
	defer f.Close()
	e := gob.NewEncoder(f)
	for {
		r := <-s.saveChan
		if err := e.Encode(r); err != nil {
			panic(err)
		}
	}
}