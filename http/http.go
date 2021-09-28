package http

import (
	"errors"
	"fmt"
	"goto/database"
	"io"
	"net/http"
)

var store *database.URLStore

func init() {
	store = database.NewURLStore("test.gob")
}

func Run() {
	http.HandleFunc("/", redirect)
	http.HandleFunc("/add", add)
	http.HandleFunc("/count", count)
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		panic(err)
	}
}

func redirect(w http.ResponseWriter, req *http.Request) {
	var key string
	err := func() error {
		key = req.URL.Path[1:]
		if key == "" {
			return errors.New("key is empty")
		}
		if !store.Exist(key) {
			return errors.New("the key does not exist")
		}
		return nil
	}()
	if err != nil {
		w.WriteHeader(400)
		io.WriteString(w, err.Error())
		//http.NotFound(w, r)
		return
	}
	http.Redirect(w, req, store.Get(key), http.StatusFound)
}

const addForm = `<html><head><title>add</title></head><body><form method="POST" action="/add">URL: <input type="text" name="url"> <input type="submit" value="add"></body></html>`
func add(w http.ResponseWriter, req *http.Request) {
	switch req.Method {
	case "GET":
		fmt.Fprint(w, addForm)
	case "POST":
		url := req.FormValue("url")
		err := func() error {
			if url == "" {
				return errors.New("missing parameters")
			}
			return nil
		}()
		if err != nil {
			w.WriteHeader(400)
			io.WriteString(w, err.Error())
			return
		}
		key := store.Put(url)
		fmt.Fprintf(w, "ok, short url is: http://192.168.188.147:8080/%s", key)
	default:
		io.WriteString(w, "missing parameters")
	}
}

func count(w http.ResponseWriter, req *http.Request) {
	fmt.Fprintf(w, "count: %d", store.Count())
}