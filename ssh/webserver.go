package main

import (
	"flag"
	"html/template"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)



var addr = flag.String("addr", "127.0.0.1:8000", "http service address")

var upgrader = websocket.Upgrader{} // use default options

func echo(w http.ResponseWriter, r *http.Request) {
	c, err := upgrader.Upgrade(w, r, nil)
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
		
		requestChannel <- string(message)
		log.Println("got command on webserver: ", string(message))
		
		go func() {
			log.Println("got response from server")
			message := []byte(<-responseChannel)
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

var requestChannel = make(chan string)
var responseChannel = make(chan string)

func home(w http.ResponseWriter, r *http.Request) {
	homeTemplate.Execute(w, "ws://"+r.Host+"/echo")
}


func main() {
	go connectClientnWeb()
	flag.Parse()
	http.HandleFunc("/echo", echo)
	http.HandleFunc("/", home)
	log.Fatal(http.ListenAndServe(*addr, nil))
}


var homeTemplate = template.Must(template.New("").Parse(`
<!DOCTYPE html>
<html>

<head>
    <meta charset="utf-8">
    <style>
        #form {
            width: 200px;
            padding: 20px;
        }

        input[type=text] {
            padding: 10px;
            margin: 10px 0;
            border: 0;
            box-shadow: 0 0 15px 4px rgba(0, 0, 0, 0.06);
        }

        .button {
            border: 
            background-color: #E7E7E7;
            color: #000;
            color: black;
            padding: 10px 24px;
            text-align: center;
            text-decoration: none;
            display: inline-block;
            font-size: 14px;
            cursor:pointer;
        }

    </style>
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
    <div class="form-container" id="form-container">
        <form action="" method="post" id="form">
            <p>
                <button id="open" class="button">Open the ws connection</button>
            </p>
            <p>
                <button id="close" class="button">Close the ws connection</button>
            </p>
            <p>Input the command <input type="text" size="40" placeholder="ls" id="input" autocomplete="off"></p>
            <p><button id="send" class="button">ok</button></p>
        </form>
        <p></p>
        <div id="output" style="max-height: 70vh;overflow-y: scroll;"></div>
    </div>
</body>

</html>`))
