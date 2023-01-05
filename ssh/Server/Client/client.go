package main

import (
	"fmt"
	"log"
	"os"
	
	"github.com/skorobogatov/input"
	"golang.org/x/crypto/ssh"
)

func main() {

	hostname := "127.0.0.1"
	port := "2222"
	username := "abc"
	password := "1234"

	config := &ssh.ClientConfig{
		User: username,
		Auth: []ssh.AuthMethod{
			ssh.Password(password),
		},
		//InsecureIgnoreHostKey возвращает функцию для принятия любого ключа хоста
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}

	//Подключаемся к хосту
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

	// Enable system stdout
	s.Stdout = os.Stdout
	s.Stderr = os.Stderr

	err = s.Shell()
	if err != nil {
		log.Fatal(err)
	}

	for {
		cmd := input.Gets() //вводим команду
		_, err = fmt.Fprintf(stdin, "%s\n", cmd) //тут используется пайп для ввода команд
		if err != nil {
			log.Fatal(err)
		}
		if cmd == "exit" {
			s.Close()
			break
		}
	}
}
