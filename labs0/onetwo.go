package main

import (
    "fmt"// пакет для форматированного ввода вывода
    "log"// пакет для логирования
    "net/http"// пакет для поддержки HTTP протокола
)

func handler(w http.ResponseWriter, r *http.Request) {
    fmt.Fprintf(w, "Hi there, I love %s!", get(r))
}

func get(r *http.Request) string{
	return r.URL.Path[1:]
}

func main() {
    http.HandleFunc("/", handler)
    err := http.ListenAndServe(":3000", nil) // задаем слушать порт
 	 if err != nil {
   	 	log.Fatal("ListenAndServe: ", err)
   	 }
}

