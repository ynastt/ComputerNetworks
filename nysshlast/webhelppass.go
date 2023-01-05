package main

import (
	"flag"
	"html/template"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)



var addrs = flag.String("addr", "localhost:8088", "http service address")

var upgraderr = websocket.Upgrader{} // use default options

func echos(w http.ResponseWriter, r *http.Request) {
	c, err := upgraderr.Upgrade(w, r, nil)
	if err != nil {
		log.Print("upgrade:", err)
		return
	}
	defer c.Close()
	for {
		mt, message, err := c.ReadMessage()
		if err != nil {
			log.Println("Error with reading message: ", err)
			break
		}
		
		usernameChannel <- string(message)
		log.Println("got username on webserver: %s", string(message))
		
		passwordChannel <- string(message)
		log.Println("got password on webserver: %s", string(message))
		
		go func() {
			log.Println("got response from server")
			message := []byte(<-answerChannel)
			if err != nil {
				log.Println("Error with reading message:", err)
			}
			log.Printf("recieved msg: %s ", message)
			err = c.WriteMessage(mt, message)
			if err != nil {
				log.Println("Error with writing message:", err)
				
			}
		}()
		
	}
}

var answerChannel = make(chan string)
var passwordChannel = make(chan string)

func homee(w http.ResponseWriter, r *http.Request) {
	homeTemplate.Execute(w, "ws://"+r.Host+"/echo")
}


func main() {
	go connectClientnServer()
	flag.Parse()
	http.HandleFunc("/echo", echos)
	http.HandleFunc("/", homee)
	log.Fatal(http.ListenAndServe(*addrs, nil))
}


var homeTemplate = template.Must(template.New("").Parse(`
<!DOCTYPE html>
<html lang="ru">

<head>
	<meta charset="utf-8">
	<title>log in page</title>
</head>

<body>
	<p><input type="password" placeholder="password" id="pass"></p>
	<button id="connect">send</button>
</body>

</html>`))
