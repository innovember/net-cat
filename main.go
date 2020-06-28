package main

import (
	"bufio"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net"
	"os"
	"sync"
	"time"
)

type User struct {
	Name string
	Conn net.Conn
}

var (
	users       map[string]*User
	allMessages string
	mutex       sync.Mutex
)

func init() {
	users = make(map[string]*User)
}

func main() {
	args := os.Args[1:]
	if len(args) != 1 {
		fmt.Println("[USAGE]: ./net-cat $port")
		return
	}
	port := args[0]
	server, err := net.Listen("tcp", ":"+port)
	if err != nil {
		log.Fatalln(err.Error())
	}
	defer server.Close()
	fmt.Printf("Listening on port %s\n", port)
	fmt.Printf("Command: \"nc localhost %s\"\n", port)

	for {
		conn, err := server.Accept()
		if len(users) == 10 {
			fmt.Println(conn, "Chat already filled")
			conn.Close()
			break
		}
		if err != nil {
			log.Fatalln(err.Error())
		}
		mutex.Lock()
		user := &User{
			Conn: conn,
		}
		mutex.Unlock()
		user.newUser(user.Conn)
		go user.handle()
		defer conn.Close()
	}
}

func (user *User) newUser(conn net.Conn) {
	txt, errTxt := ioutil.ReadFile("greetings.txt")
	if errTxt != nil {
		os.Exit(0)
	}
	io.WriteString(conn, string(txt)+"\n")
	io.WriteString(conn, "[ENTER YOUR NAME]: ")
	scanner := bufio.NewScanner(conn)
	scanner.Scan()
	user.Name = scanner.Text()
	users[user.Name] = user
	mutex.Lock()
	io.WriteString(conn, allMessages)
	mutex.Unlock()
	io.WriteString(conn, fmt.Sprintf("[%s][%s]: ", getFormatTime(), user.Name))
	greetings := fmt.Sprintf("%s has joined our chat...", user.Name)
	mutex.Lock()
	allMessages += greetings + "\n"
	for key := range users {
		if key != user.Name {
			io.WriteString(users[key].Conn, "\n"+greetings+"\n")
			io.WriteString(users[key].Conn, fmt.Sprintf("[%s][%s]: ", getFormatTime(), users[key].Name))
		}
	}
	mutex.Unlock()

}

func (user *User) handle() {
	scanner := bufio.NewScanner(user.Conn)
	for scanner.Scan() {
		text := scanner.Text()
		if !checkMsg(text) {
			msg := fmt.Sprintf("[%s][%s]:%s ", getFormatTime(), user.Name, text)
			mutex.Lock()
			allMessages += msg + "\n"
			for key := range users {
				if key != user.Name {
					io.WriteString(users[key].Conn, "\n"+msg+"\n")
					io.WriteString(users[key].Conn, fmt.Sprintf("[%s][%s]: ", getFormatTime(), users[key].Name))
				}
			}
			mutex.Unlock()
		}
		io.WriteString(user.Conn, fmt.Sprintf("[%s][%s]: ", getFormatTime(), user.Name))

	}
	user.disconnect()
}

func (user *User) disconnect() {
	msg := fmt.Sprintf("%s has left our chat...", user.Name)
	mutex.Lock()
	allMessages += msg + "\n"
	for key := range users {
		if key != user.Name {
			io.WriteString(users[key].Conn, "\n"+msg+"\n")
			io.WriteString(users[key].Conn, fmt.Sprintf("[%s][%s]: ", getFormatTime(), users[key].Name))
		} else {
			delete(users, key)
		}
	}
	mutex.Unlock()
}
func getFormatTime() string {
	now := time.Now()
	return fmt.Sprintf("%d-%02d-%02d %02d:%02d:%02d", now.Year(), now.Month(), now.Day(), now.Hour(), now.Minute(), now.Second())
}

func checkMsg(str string) bool {
	if str == "\n" || str == "" || str == "\t" || str == " " {
		return true
	}
	return false
}
