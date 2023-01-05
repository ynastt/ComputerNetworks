package main

import (
	"flag"
	"html/template"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)



var addrs = flag.String("addr", "localhost:8000", "http service address")

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
		
		go func() {
			log.Println("got response from server")
			message := []byte(<-answerChannel1)
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

var answerChannel1 = make(chan string)
var usernameChannel = make(chan string)

func homee(w http.ResponseWriter, r *http.Request) {
	homeTemplate.Execute(w, "ws://"+r.Host+"/echos")
}


func main() {
	go compareUserame()
	flag.Parse()
	http.HandleFunc("/echos", echos)
	http.HandleFunc("/", homee)
	log.Fatal(http.ListenAndServe(*addrs, nil))
}


var homeTemplate = template.Must(template.New("").Parse(`
<!DOCTYPE html>
<html lang="ru">

<head>
    <meta charset="utf-8">
    <title>username page</title>
    <script>
        window.addEventListener("load", function(evt) {
            var output = document.getElementById("output");
            var input = document.getElementById("input");
            var ws;
            var print = function(message) {
                var d = document.createElement("div");
                d.textContent = message;
                output.appendChild(d);
                output.scroll(0, output.scrollHeight);
            };
            document.getElementById("open").onclick = function(evt) {
                if (ws) {
                    return false;
                }
                ws = new WebSocket("{{.}}");
                ws.onopen = function(evt) {
                    print("THE CONNECTION IS OPENED");
                }
                ws.onclose = function(evt) {
                    print("THE CONNECTION IS CLOSED");
                    ws = null;
                }
                ws.onmessage = function(evt) {
                    print("RESPONSE: " + evt.data);
                }
                ws.onerror = function(evt) {
                    print("ERROR: " + evt.data);
                }
                return false;
            };
            document.getElementById("send").onclick = function(evt) {
                if (!ws) {
                    return false;
                }
                print("REQUEST: " + input.value);
                ws.send(input.value);
                return false;
            };
            document.getElementById("close").onclick = function(evt) {
                if (!ws) {
                    return false;
                }
                ws.close();
                return false;
            };
        });

    </script>
</head>

<body>
	<div class="form-container2" id="form-container2">
        <form action="" method="post" id="form2">
            <p><button id="open">Open</button>
                <button id="close">Close</button>
            </p>
            <p><input type="text" placeholder="username" id="name"></p>
            <p><button id="send">ok</button></p>
        </form>
        <p>is it right?</p>
        <div id="output" style="max-height: 70vh;overflow-y: scroll;"></div>
    </div>
</body>

</html>`))
