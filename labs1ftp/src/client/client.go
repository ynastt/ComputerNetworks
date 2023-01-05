package main
 
import (
	"fmt"
    	"strings"
    	"github.com/skorobogatov/input"
    	"github.com/jlaffaye/ftp"
    	"os"
    	"io"
)
 
func connectToServer(conn *ftp.ServerConn) {
	for {
        	cur, _ := conn.CurrentDir() //путь к текущей директории
        	fmt.Print(cur, ": command: ")
        	command := input.Gets()
 
        	switch command {
        		case "store": //загрузка файла go ftp клиентом на ftp сервер
            			fmt.Printf("local file:  ")
            			localFile := input.Gets()
            			fmt.Printf("store as:  ")
            			remoteFile := input.Gets()
            			if len(remoteFile) == 0 {
					remoteFile, _ = conn.CurrentDir()
				}
            			file, err := os.Open(localFile);
            			if err != nil {
            			    fmt.Println("error: ", err)
            			}
            			err = conn.Stor(remoteFile, file) //загрузка на удаленный сервер
            			if err != nil {
                			fmt.Println("error: ", err)
            			}
            		case "fetch": //скачивание файла go ftp клиентом с ftp сервера;
        			fmt.Printf("file:  ")
				name := input.Gets()
				res, err := conn.Retr(name)
				if err != nil {
					fmt.Println("error: ", err)
				}
				defer res.Close()
 
				outFile, err := os.Create(name)
				if err != nil {
					fmt.Println("error: ", err)
				}
				defer outFile.Close()
 
				if _, err = io.Copy(outFile, res); err != nil {
					fmt.Println("error: ", err)
				}
				res.Close()
			case "mkdir": //создание директории go ftp клиентом на ftp сервере
            			fmt.Printf("directory:  ")
            			path := input.Gets()
            			pathArray := strings.Split(path, "/")
            			dir := pathArray[len(pathArray) - 1]
            			pathArray = pathArray[:len(pathArray) - 1]
            			for _, dir := range pathArray {
                			if err := conn.ChangeDir(dir); err != nil {
                    				conn.MakeDir(dir)
                    				conn.ChangeDir(dir)
                			}
            			}
            			if err := conn.MakeDir(dir); err != nil {
                			fmt.Println("error: ", err)
            			} 
            			conn.ChangeDir(cur)
            		case "delete": //удаление go ftp клиентом  файла на ftp сервере;
            			fmt.Printf("file:  ")
            			path := input.Gets()
            			if err := conn.Delete(path); err != nil {
                			fmt.Println("error: ", err)
            			}			
        		case "list": //получение go ftp клиентом содержимого директории на ftp сервере
            			fmt.Printf("directory:  ")	
            			path := input.Gets()
            			if len(path) == 0 {
                			path = cur //не вводим заново путь к текущему каталогу
                		}		   
              			if err := conn.ChangeDir(path); err != nil {//текующая директория -> 
              			                                            //другая директория по заданному пути
				}
            			if entries, err := conn.List(path); err != nil {
                			fmt.Println("error: ", err)
            			} else {
                			fmt.Println("\ncurrent directory: ", path, "\n")
                			for _, e := range entries {
                    				fmt.Printf("Name: %-20s Type: %-10s Size: %-10d %s %s\n", e.Name, e.Type, e.Size, e.Time.Format("2006/1/2"), e.Time.Format("15:04"))
                			}
                			fmt.Println()
            			}
        		case "cd"://смена каталога для удобства
            			fmt.Printf("directory:  ")
            			path := input.Gets()
            			if err := conn.ChangeDir(path); err != nil {
                			fmt.Println("error: ", err)
            			}
            		case "rmdir": //удаление go ftp клиентом директории на ftp сервере
				fmt.Printf("directory:  ")
				path := input.Gets()
				if err := conn.RemoveDir(path); err != nil {
					fmt.Println("error: ", err)
				}
        		case "quit": //прощание с сервером
        			fmt.Println("goodbye!")
            			return
        		default:
        			fmt.Println("UNKNOWN COMMAND")
        	}	
    	}
}
 
func main() {
    //установка соединения
    if conn, err := ftp.Dial("students.yss.su:21"); err != nil {
        fmt.Println("ERROR: ", err)
    } else {
        fmt.Print("login = ")
        login := input.Gets()
        fmt.Print("password = ")
        password := input.Gets()
        if err = conn.Login(login, password); err != nil {
            fmt.Println("ERROR: ", err)
        } else {
            fmt.Println("welcome!")
            connectToServer(conn)
        }
    }
}
