package main

import (
	"fmt" // пакет для форматированного ввода вывода
	"net/http" // пакет для поддержки HTTP протокола
	"log" // пакет для логирования
	"github.com/mmcdole/gofeed" //RSS-парсер
)

//функция для извлечения новостных статей из RSS-канала
func fetchNews(url string) []*gofeed.Item {
  	fp := gofeed.NewParser()
  	feed, _ := fp.ParseURL(url)
  	return feed.Items
}

func handler(w http.ResponseWriter, r *http.Request) {
	const url string = "https://lenta.ru/rss"
	fp := gofeed.NewParser()
	feed, _ := fp.ParseURL(url)
	fmt.Println("in process") //в терминале отмечаем начало работы web-сервера
	fmt.Fprintf(w, "<h1>Новостной портал: %s</h1>", feed.Title)
	temp := fetchNews(url)
	for n := 1; n < 999; n++ {
		fmt.Println()
                currentNews := "<div><font size=\"5\"><a href=\"" + temp[n].Link + "\">" + temp[n].Title + "</a></font><br>" + 
			"<font size=\"3\">" + temp[n].Description + "</font><br><br>" + "</div>"
	fmt.Fprintf(w, currentNews)
	}
}

func main() {
	http.HandleFunc("/", handler) // установим роутер
   	err := http.ListenAndServe(":2000", nil) // задаем слушать порт
  	if err != nil {
   	    	log.Fatal("ListenAndServe: ", err)
   	}
}
