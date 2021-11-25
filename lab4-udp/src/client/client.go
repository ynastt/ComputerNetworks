package main

import (
	"flag"
	"fmt"
	"github.com/mgutz/logxi/v1"
	"net"
	"os"
	"strconv"
	"encoding/json"
	"github.com/skorobogatov/input"
)

import "proto"

// interact - функция, взаимодействующая с сервером.
// index - номер запроса соответственно его индекс
// command - quit/add/height
// data - данные (скорость, угол в случае add)
// Данные могут быть пустыми (data == nil) - случаи quit & height
func interact(conn *net.UDPConn, index uint, command string, data interface{}) {
	id := strconv.Itoa(int(index))
	var raw json.RawMessage
	raw, _ = json.Marshal(data) 
	req, _ := json.Marshal(&proto.Request{command, &raw, id})
	buf := make([]byte,3000)
	for {
        	if _, err := conn.Write(req); err != nil {	// Отправляем запрос по udp соединению
			log.Error("sending request to server", "error", err)
			log.Info("try to send the request again later")
			continue
		} 
		if bytesRead, err := conn.Read(buf); err != nil {	// Читаем ответ с сервера
			log.Error("receiving answer from server", "error", err)
			continue
		} else {
			// Разбираем ответ с сервера
			// Если не смогли, то сообщаем об ошибке в терминал
			var resp proto.Response
			if err := json.Unmarshal(buf[:bytesRead], &resp); err != nil {
				log.Error("cannot parse answer", "answer", buf, "error", err)
			} else {
				// Если смогли прочитать ответ с сервера
				// Обработка и вывод ответа сервера 
				switch resp.Status {
				case "bye":
					if resp.Id == id {
						log.Info("client is off")
						return
					}	
				case "ok":
					var elem proto.Elem
					if err := json.Unmarshal(*resp.Data, &elem); err != nil {
						log.Error("cannot parse answer", "answer", resp.Data, "error", err)
					} else {
						if resp.Id == id {
							log.Info("successful interaction with server", "new V", elem.Velocity, "new alpha", elem.Angle)
							return
						}
					}	
				case "failed":
					var errorMsg string
					if err := json.Unmarshal(*resp.Data, &errorMsg); err != nil {
							log.Error("cannot parse answer", "answer", resp.Data, "error", err)
					} else {
						if resp.Id == id {
							log.Info("failed", "errorMessage", errorMsg)
							return
						}
					}
				case "result":
					var res string // наиб. высота
					if err := json.Unmarshal(*resp.Data, &res); err != nil {
							log.Error("cannot parse answer", "answer", resp.Data, "error", err)
					} else {
						if resp.Id == id { 					
							log.Info("successful interaction with server", "max height", res)
							fmt.Printf("result: %s\n", res)
							return
						}	
					}
				default:
					log.Info("error: server reports unknown status %q\n", resp.Status)
				}
			}
		}
	}		
}

func main() {
	var (
		serverAddrStr string
		n             uint
		helpFlag      bool
	)
	flag.StringVar(&serverAddrStr, "server", "127.0.0.1:6000", "set server IP address and port")
	flag.UintVar(&n, "n", 10, "set the number of requests")
	flag.BoolVar(&helpFlag, "help", false, "print options list")

	// Разбор адреса, установка соединения с сервером и
	// запуск взаимодействия с сервером.
	if flag.Parse(); helpFlag {
		fmt.Fprint(os.Stderr, "client [options]\n\nAvailable options:\n")
		flag.PrintDefaults()
	} else if serverAddr, err := net.ResolveUDPAddr("udp", serverAddrStr); err != nil {
		log.Error("resolving server address", "error", err)
	} else if conn, err := net.DialUDP("udp", nil, serverAddr); err != nil {
		log.Error("creating connection to server", "error", err)
	} else {
		defer conn.Close()
		for i := uint(0); i < n; i++ {
			// Чтение команды из стандартного потока ввода
			fmt.Printf("command = ")
			command := input.Gets()
			// Определяем какой именно запрос
			switch command {
			case "quit":
				interact(conn, i, "quit", nil)
				return
			case "add":
				var pair proto.Elem
				fmt.Printf("V = ")
				pair.Velocity = input.Gets()
				fmt.Printf("alpha = ")
				pair.Angle = input.Gets()
				interact(conn, i, "add", &pair)
			case "height":
				interact(conn, i, "height", nil)
        		default:
            			log.Info("error: unknown command\n")
            			i -= 1 // Оставим число запросов прежним (число n)
            			continue
			}
		}
	}
}