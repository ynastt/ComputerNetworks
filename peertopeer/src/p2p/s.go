package main

import (
	"encoding/json"
	"fmt"
	"github.com/skorobogatov/input"
	"net"
	"github.com/mgutz/logxi/v1"
)

var(
	globalIP = "127.0.0.1"
	post string = ""
	MyAddr string = ""
	encoder *json.Encoder = nil
	decoder *json.Decoder = nil
)

type Request struct {
	Data *json.RawMessage `json:"data"`
}

type MyStr struct {
	IP string
	Sum string
}

type Peer struct {
	logger log.Logger    // Объект для печати логов
	conn   *net.TCPConn  // Объект TCP-соединения
	enc    *json.Encoder // Объект для кодирования и отправки сообщений
	dec    *json.Decoder // Объект для декодирования сообщений
}


func NewPeer(conn *net.TCPConn) *Peer { 
	return &Peer{
		logger: log.New(fmt.Sprintf("peer %s", conn.RemoteAddr().String())),
		conn:   conn,
		enc:    json.NewEncoder(conn),
		dec:    json.NewDecoder(conn),
	}
}


func send_request(encoder *json.Encoder, data interface{}) {
	var raw json.RawMessage
	raw, _ = json.Marshal(data)
	encoder.Encode(&Request{&raw})
}

func interact(conns[] lya) {
	for {
		var command string
		command = input.Gets()

		switch command {
			case "post":
				var str string
				fmt.Print("enter new post: ")
				str = input.Gets()
				post = str
			case "remove":
				post = ""
			case "print":
				for _, e := range conns {
					send_request(e.encoder, MyStr{MyAddr, ""})
					var resp Request
					if err := e.decoder.Decode(&resp); err != nil {
						fmt.Printf("error: %v\n", err)
						break
					}
					var count MyStr
					if err := json.Unmarshal(*resp.Data, &count); err != nil {
						fmt.Printf("error: malformed data field in response\n")
					} else {
						fmt.Printf("result: " + count.Sum + "\n")
					}
				}
			default:
				fmt.Printf("error: unknown command\n")
		}
		
	}
}


func (peer *Peer) serve() {
	defer peer.conn.Close()
	for {
		var req Request
		if err := peer.dec.Decode(&req); err != nil {
			peer.logger.Error("cannot decode message", "reason", err)
			break
		} else {
			peer.logger.Info("received message")
			peer.handleRequest(&req)
		}
	}
}

func (client *Peer) respond(data interface{}) {
	var raw json.RawMessage
	raw, _ = json.Marshal(data)
	client.enc.Encode(&Request{&raw})
}

func (peer *Peer) handleRequest(req *Request) {
		errorMsg := ""
		if req.Data == nil {
			errorMsg = "data field is absent"
		} else {
			g := new(MyStr)
			if err := json.Unmarshal(*req.Data, g); err != nil {
				errorMsg = "malformed data field"
			} else {
				if post != "" {
					peer.respond(MyStr{g.IP, fmt.Sprintf("%s%s:\n%s\n\n", g.Sum, MyAddr, post)})
				} else {
					peer.respond(MyStr{g.IP, fmt.Sprintf("%s%s:\nnothin to report\n", g.Sum, MyAddr)})
				}
			}
		}
		if errorMsg == "" {
			peer.logger.Info("information from peer added succesfully")
		} else {
			peer.logger.Error("add failed", "reason", errorMsg)
		}
}

func listen(addrStr string) {
	var listener *net.TCPListener
	defer listener.Close()

	if addr, err := net.ResolveTCPAddr("tcp", addrStr); err != nil {
		log.Error("address resolution failed", "address", addrStr)
	} else {
		log.Info("resolved TCP address", "address", addr.String())

		if listener, err = net.ListenTCP("tcp", addr); err != nil {
			log.Error("listening failed", "reason", err)
		} else {
            // Цикл приёма входящих соединений.
			for {
				if conn, err := listener.AcceptTCP(); err != nil {
					log.Error("cannot accept connection", "reason", err)
				} else {
					log.Info("accepted connection", "address", conn.RemoteAddr().String())

					go NewPeer(conn).serve()
				}
			}
		}
	}
}

type lya struct{
	conn *net.TCPConn
	encoder *json.Encoder
	decoder *json.Decoder
}

func main() {
	ps := [...]string {
		"127.0.0.1:6000",
		"127.0.0.1:6001",
		"127.0.0.1:6002",
	}

	var conns [] lya
	var port string
	

	fmt.Scan(&port)
	MyAddr = fmt.Sprintf("%s:%s", globalIP, port)
  	go listen(MyAddr)

  	var s string
  	for {
  		fmt.Println("continue? y/n")
  		fmt.Scan(&s)
  		if s == "y" { 
  			break
  		}
	}


	for {
		for _, e := range ps {
			if e != MyAddr{
				if addr, err := net.ResolveTCPAddr("tcp", e); err != nil {
					fmt.Printf("error: %v\n", err)
				} else if conn, err := net.DialTCP("tcp", nil, addr); err == nil {
					encoder = json.NewEncoder(conn)
					decoder = json.NewDecoder(conn)
					conns = append(conns, lya{conn, encoder, decoder})
				}
			}
		}
		interact(conns)
	}
}
