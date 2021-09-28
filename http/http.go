package http

import (
	"errors"
	"fmt"
	"goto/commonio"
	"goto/database"
	"io"
	"net/http"
)

func Run() {
	http.HandleFunc("/", home)
	http.HandleFunc("/add", add)
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		panic(err)
	}
}

func home(w http.ResponseWriter, req *http.Request) {
	var key string
	err := func() error {
		key = req.URL.Path[1:]
		if key == "" {
			return errors.New("key is empty")
		}
		if !database.Exist(key) {
			return errors.New("the key does not exist")
		}
		return nil
	}()
	if err != nil {
		w.WriteHeader(400)
		io.WriteString(w, err.Error())
		return
	}
	http.Redirect(w, req, database.Get(key), http.StatusFound)
}

func add(w http.ResponseWriter, req *http.Request) {
	switch req.Method {
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
		key := commonio.RandomString(9)
		database.Add(key, url)
		fmt.Fprintf(w, "ok, short url is: http://192.168.188.147:8080/%s", key)
	default:
		io.WriteString(w, "missing parameters")
	}
}