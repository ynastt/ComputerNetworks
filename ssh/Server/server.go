package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"os/exec"
	"strings"
	"golang.org/x/crypto/ssh"
)

func listenPort(config *ssh.ServerConfig) {
	log.Print("Ready to listen")
	listener, err := net.Listen("tcp", "127.0.0.1:2222")
	if err != nil {
		log.Fatal("Failed to listen on 2222 (%s)", err)
	}
	log.Print("Listening on 2222...")
	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Printf("Failed to accept incoming connection (%s)", err)
			continue
		}
		// Before use, a handshake must be performed on the incoming net.Conn.
		fmt.Printf("Handshaking for %s\n", conn.RemoteAddr())
		sshConn, chans, reqs, err := ssh.NewServerConn(conn, config)
		if err != nil {
			log.Print("Failed to handshake (%s)", err)
			continue
		}
		log.Printf("SSH connection from %s (%s)\n", sshConn.RemoteAddr(), sshConn.ClientVersion())
		// The incoming Request channel must be serviced.
		//DiscardRequest принимает и отклоняет все запросы переданного канала
		go ssh.DiscardRequests(reqs)
		handleChannels(chans)
	}
}

func handleChannels(chans <-chan ssh.NewChannel) {
	//работаем с входящим каналом в горутине
	for newChannel := range chans {
		go handleChannel(newChannel)
	}
}

func handleChannel(newChan ssh.NewChannel) {
		// тип канала должен быть session
		if typ := newChan.ChannelType(); typ != "session" {
			newChan.Reject(ssh.UnknownChannelType, fmt.Sprintf("unknown channel type: %s", typ))
			return
		}
		//возможность отказать запросу клиента
		channel, requests, err := newChan.Accept()
		if err != nil {
			log.Fatalf("Could not accept channel: %v", err)
		}
		//канал запросов - нам нужен  shell
		go func(in <-chan *ssh.Request) {
			for req := range in {
				req.Reply(req.Type == "shell" , nil)
			}
		}(requests)
		a := make([]byte, 1000)
		defer channel.Close()
		for {
			//channel.Write([]byte("Request "))
			//channel.Write([]byte(">"))
			channel.Read(a)
			cmd := string(a)
			i := strings.Index(cmd, "\n")
			cmd = cmd[:i]
			if cmd == "exit" {
				log.Println("Terminal is closed")
				channel.Close()
				break
			}
			splitted := strings.Split(cmd, " ")
			fmt.Println(splitted)
			res, err := exec.Command(splitted[0], splitted[1:]...).Output()
			if err != nil {
				log.Print(err)
				channel.Write([]byte(err.Error() + "\n"))
			}
			var s []byte
			b := []byte(">")
			s = append(s, b...)
			res = append(s, res...)
			log.Println("respond")
			channel.Write(res)
		}
}

func main() {
	//Функция, к выполняется, когда клиент пытается ввести пароль для входа
	config := ssh.ServerConfig{
		PasswordCallback: func(con ssh.ConnMetadata, password []byte) (*ssh.Permissions, error) {
			if con.User() == "abc" && string(password) == "1234" {
					return nil, nil
				
			}
			return nil, fmt.Errorf("password rejected for %q", con.User())
		},
		// можно явно разрешить аутентификацию анон. клиента
		// NoClientAuth: true,
	}
	//т.к. пара ключей сгенерирована заранее, то читаем из файла
	privateBytes, err := ioutil.ReadFile("/home/qso/.ssh/id_rsa")
	if err != nil {
		panic("Fail to load private key (./id_rsa)")
	}
	private, err := ssh.ParsePrivateKey(privateBytes)
	if err != nil {
		panic("Fail to parse private key")
	}
	config.AddHostKey(private)
	//начинаем слушать порт
	listenPort(&config)
}
