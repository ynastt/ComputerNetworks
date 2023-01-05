package main

import (
	"fmt"
	"log"
	"golang.org/x/crypto/ssh"
)

var (
	hostname := "127.0.0.1"
	port := "2222"
	username := "abc"
	password := "1234"
)

func connectClientnServer() {
	
	pswd := <-passwordChannel   
	if pswd == password && usr == username {
		an := "allowed!"
		answerChannel <- an
		log.Println("right password and username")
		log.Println("redirect to webserver page")
		//перенаправить на localhost:8080 webserver.go
	}
	
	} else {
		an := "Wrong username or password!"
		answerChannel <- an
	}
		
}
compareUserame(){
	usr := <-usernameChannel
	if usr == username {
		an := "right username!"
		answerChannel <- an
		log.Println("right username")
		log.Println("redirect to password page")
		//перенаправить на localhost:8080 webhelppass.go
	}
	
	} else {
		an := "Wrong username!"
		answerChannel <- an
	}
}

func connectClientnWeb() {
	
	config := &ssh.ClientConfig{
		User: username,
		Auth: []ssh.AuthMethod{
			ssh.Password(password),
		},
		//InsecureIgnoreHostKey возвращает функцию для принятия любого ключа хоста
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}

	//Подключаемся к хосту (т.е. серверу)
	client, err := ssh.Dial("tcp", hostname+":"+port, config)
	if err != nil {
		log.Fatal(err)
	}
	defer client.Close()

	// новая сессия
	s, err := client.NewSession()
	if err != nil {
		log.Fatal("Creating session is failed: ", err)
	}
	defer s.Close()

	// создаем пайп для ввода команд
	//StdinPipe returns a pipe that will be connected to the command's standard input when the command starts.
	//The pipe will be closed automatically after Wait sees the command exit.
	stdin, err := s.StdinPipe()
	if err != nil {
		log.Fatal(err)
	}
	
	//создаем пайп для вывода ответа сервера на запрос выполнения команды
	//StdoutPipe returns a pipe that will be connected to the command's standard output when the command starts.
	//Wait will close the pipe after seeing the command exit
	stdout, err := s.StdoutPipe()
	if err != nil {
		log.Fatal(err)
	}

	err = s.Shell()
	if err != nil {
		log.Fatal(err)
	}

	for {
		//  каналы созданы в webserver.go
		cmd := <-requestChannel   //получаем команду из канала с webserver.go
		_, err = fmt.Fprintf(stdin, "%s\n", cmd) //тут используется пайп для ввода команд
		if err != nil {
			log.Fatal(err)
		}
		if cmd == "exit" {
			s.Close()
			str := "Terminal is closed. Time to close the ws connection"
			responseChannel <- str
			break
		}
		//считываем результат выполнения операции в байтах, чтобы передать в канал обратно на webserver.go
		a := make([]byte, 1000)  //число условное, можно и больше/меньше
		num, _ := stdout.Read(a) //тут используется пайп для вывода команд
		log.Println("Response to the requested command-", cmd)
		responseChannel <- string(a[:num])
	}
}
