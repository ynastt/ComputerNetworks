package main

import (
	"encoding/json"
	"math"
	"flag"
	"fmt"
	"github.com/mgutz/logxi/v1"
	"net"
	"os"
	"strconv"
)

import "proto"

// Client - состояние клиента
type Client struct {
	res	map[int]proto.Response // Хэш-таблица(id -> response) для сохранения обработанных запросов
	v    	float64      		// Скорость камня
	alp   	int      		// Угол, под которым брошен камень
	count   float64         	// Посчитанная высота
}

// NewClient - конструктор клиента
func NewClient() *Client {
	return &Client{
		res:	make(map[int]proto.Response),
		v:    	0,
		alp: 	0,
		count:  0,
	}
}

// respond - Вспомогательная функция для передачи 
// ответа с указанным статусом и данными
// возвращает true, если ответ отправлен успешно, false если ошибка
func respond(status string, data interface{}, conn *net.UDPConn, addr *net.UDPAddr, id string) bool {
	var raw json.RawMessage
	switch typ := data.(type) {
	case float64:
		dataconv := fmt.Sprintf("%.2f", typ) // Преобразуем в строку с точностью 2 знака после запятой
		data = dataconv
	default:
		
	}
	raw, _ = json.Marshal(data) 
	respond, _ := json.Marshal(&proto.Response{status, &raw, id})
	// Отправляем ответ по udp соединению на адрес клиента
	if _, err := conn.WriteToUDP(respond, addr); err != nil {
		log.Error("sending message to client", "error", err)
		return false
	}
	return true
}

// serveClients - Функция, в которой реализован цикл взаимодействия с клиентом
func serveClients(conn *net.UDPConn) {
	buf := make([]byte, 3000)
	// Хэш-таблица для обработки клиентов
	// Адрес клиента в строке -> клиент
	mapc := make(map[string]*Client)
	for {
		// Читаем запрос с клиента по udp соединению 
		if bytesRead, addr, err := conn.ReadFromUDP(buf); err != nil {
			log.Error("receiving message from client", "error", err)
		} else {
			// Успешно получили запрос
			clientAddrStr := addr.String() // Переводим в строку (т.к. изначально формат *net.UDPAddr)
			// Проверяем, есть ли уже этот клиент
			// Если нет - создаем новый клиент
			_, exists := mapc[clientAddrStr]
			if !exists {
				log.Info("client is working", "client", clientAddrStr)
				mapc[clientAddrStr] = NewClient()
			}
			// Если да, смотрим запрос этого клиента по id 
			var req proto.Request
			// Разбираем запрос
			// Если не смогли, то сообщаем об ошибке в терминал и отправляем ответ клиенту
			if err := json.Unmarshal(buf[:bytesRead], &req); err != nil {
				log.Error("cannot parse request", "request", buf[:bytesRead], "error", err)
				//отправляем ответ(статус и данные) клиенту
				respond("failed", err, conn, addr, "-777")		
			} else {
				// Смогли прочитать запрос
				i, _ := strconv.Atoi(req.Id) //id запроса
				// Находим существующий запрос и отправляем ответ клиенту, иначе обрабатываем новый запрос
				resp, exists := mapc[clientAddrStr].res[i]
				if exists {
					log.Info("oops this request already exists. Sending the response again")
					respond(resp.Status, resp.Data, conn, addr, resp.Id)
				} else {
					// Обработка нового запроса клиента
					if handleRequest(&req, conn, addr, mapc) {
						log.Info("shutting down connection")
						break
					}		
				}
			}
		}
	}
}

// handleRequest - Функция обработки запроса от клиента. Он возвращает true,
// если клиент передал команду "quit" и хочет завершить общение.
func handleRequest(req *proto.Request,conn *net.UDPConn, addr *net.UDPAddr, mapc map[string]*Client) bool {
	g := 9.81 
	clientAddrStr := addr.String()
	i, _ := strconv.Atoi(req.Id)
	switch req.Command {
	case "quit":
		// В таблице клиентов находим клиент по адресу,
		// в струтуре клиента в таблицу для сохранения обработанных 							
		// запросов по Id записываем соответственный ответ сервера
		mapc[clientAddrStr].res[i] = proto.Response{"bye", nil, req.Id}
		// Отправляем ответ клиенту
		if respond("bye", nil, conn, addr, req.Id) {
			log.Info("client is off now")
		}
		return true	
	case "height":
		h := mapc[clientAddrStr].count
		var raw json.RawMessage
		raw, _ = json.Marshal(h)
		mapc[clientAddrStr].res[i] = proto.Response{"result", &raw, req.Id}
		if respond("result", h, conn, addr, req.Id) {
			log.Info("successful interaction with client", "height", h,
			         "client", clientAddrStr)
		}
	case "add": 
		errorMsg := ""
		var elem proto.Elem
		if req.Data == nil {
			errorMsg = "data field is absent"
		} else {
			if err := json.Unmarshal(*req.Data, &elem); err != nil {
				errorMsg = "cannot parse request"
			} else {
				if vel, err := strconv.ParseFloat(elem.Velocity, 64); err != nil {
					errorMsg = "malformed data field in velocity"
				} else {
					if alpha, err := strconv.Atoi(elem.Angle); err != nil {
						errorMsg = "malformed data field in angle"
					} else {
						log.Info("performing addition new parameters", "Velocity", elem.Velocity, "Angle", elem.Angle)
						log.Info("start counting max height")
						mapc[clientAddrStr].v = vel
						mapc[clientAddrStr].alp = alpha
						mapc[clientAddrStr].count = 0
						var help float64
						// Считаем по формуле H = (v^2*sin(fi)^2)/(2 * g)
						fi := float64(alpha) * (math.Pi / 180.0)
						num := vel * vel * math.Sin(fi) * math.Sin(fi)
						den:= 2 * g
						help = num / den
						mapc[clientAddrStr].count += help
					}
				}
			}
		}
		if errorMsg == "" {
			var raw json.RawMessage
			raw, _ = json.Marshal(elem)
			mapc[clientAddrStr].res[i] = proto.Response{"ok", &raw, req.Id}
			if respond("ok", elem, conn, addr, req.Id) {
				log.Info("successful interaction with server", "new V", elem.Velocity, "new alpha", elem.Angle, "client", clientAddrStr)
			}
		} else {
			log.Error("addition failed", "error message", errorMsg)
			respond("failed", errorMsg, conn, addr, req.Id)
		}
	}
	return false
}

func main() {
	var (
		serverAddrStr string
		helpFlag      bool
	)
	// Работа с командной строкой, в которой может указываться необязательный ключ -addr.
	flag.StringVar(&serverAddrStr, "addr", "127.0.0.1:6000", "set server IP address and port")
	flag.BoolVar(&helpFlag, "help", false, "print options list")
	// 1 - Разбор адреса, строковое представление которого находится в переменной serverAddrStr.
	// 2- Инициация слушания сети на заданном адресе.
	// 3 - Запуск функции для обслуживания клиентов.
	if flag.Parse(); helpFlag {
		fmt.Fprint(os.Stderr, "server [options]\n\nAvailable options:\n")
		flag.PrintDefaults()
	} else if serverAddr, err := net.ResolveUDPAddr("udp", serverAddrStr); err != nil { //1
		log.Error("resolving server address", "error", err)
	} else if conn, err := net.ListenUDP("udp", serverAddr); err != nil { //2
		log.Error("creating listening connection", "error", err)
	} else {
		log.Info("server listens incoming messages from clients")
		serveClients(conn) //3
	}
}