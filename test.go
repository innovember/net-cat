package main

// FIXME: Delay between channels when user will disconnect, there is delay for 2 messages in queue for receiving output message about disconnection
// import (
// 	"bufio"
// 	"fmt"
// 	"io"
// 	"io/ioutil"
// 	"log"
// 	"net"
// 	"os"
// 	"time"
// )

// type User struct {
// 	Name   string
// 	Output chan Message
// }

// type Message struct {
// 	Time     string
// 	Username string
// 	Text     string
// }

// type ChatServer struct {
// 	Users map[string]User
// 	Join  chan User
// 	Leave chan User
// 	Input chan Message
// }

// var allMessages []string

// func (cs *ChatServer) Run() {
// 	for {
// 		select {
// 		case user := <-cs.Join:
// 			cs.Users[user.Name] = user
// 			go func() {
// 				res := ""
// 				for i := range allMessages {
// 					res += allMessages[i] + "\n"
// 				}
// 				cs.Input <- Message{
// 					Time:     fmt.Sprintf("[%s]: ", getFormatTime()),
// 					Username: user.Name,
// 					Text:     fmt.Sprintf("%s joined", user.Name) + "\n" + res,
// 				}
// 			}()
// 		case user := <-cs.Leave:
// 			delete(cs.Users, user.Name)
// 			go func() {
// 				cs.Input <- Message{
// 					Time:     fmt.Sprintf("[%s]: ", getFormatTime()),
// 					Username: user.Name,
// 					Text:     fmt.Sprintf("%s left", user.Name),
// 				}
// 			}()
// 		case msg := <-cs.Input:
// 			for _, user := range cs.Users {
// 				select {
// 				case user.Output <- msg:
// 				default:
// 				}
// 			}
// 		}
// 	}
// }

// func handleConn(chatServer *ChatServer, conn net.Conn) {
// 	defer conn.Close()
// 	if len(chatServer.Users) >= 2 {
// 		conn.Write([]byte("Sorry, there are already 10 users in here :(\nYou can check in later, bye!"))
// 		return
// 	}
// 	io.WriteString(conn, "Welcome to TCP-Chat!\n")

// 	txt, errTxt := ioutil.ReadFile("greetings.txt")
// 	if errTxt != nil {
// 		os.Exit(0)
// 	}

// 	io.WriteString(conn, string(txt)+"\n")
// 	io.WriteString(conn, "[ENTER YOUR NAME]: ")

// 	scanner := bufio.NewScanner(conn)
// 	scanner.Scan()
// 	user := User{
// 		Name:   scanner.Text(),
// 		Output: make(chan Message),
// 	}
// 	chatServer.Join <- user
// 	defer func() {
// 		chatServer.Leave <- user
// 	}()

// 	// Read from conn
// 	go func() {
// 		for scanner.Scan() {

// 			text := scanner.Text()
// 			chatServer.Input <- Message{
// 				Time:     fmt.Sprintf("[%s]: ", getFormatTime()),
// 				Username: user.Name,
// 				Text:     text,
// 			}
// 		}
// 	}()

// 	// write to conn
// 	for msg := range user.Output {
// 		if !checkMsg(msg.Text) {
// 			_, err := io.WriteString(conn, fmt.Sprintf("[%s]: ", getFormatTime())+"["+msg.Username+"]"+": "+msg.Text+"\n")
// 			allMessages = append(allMessages, fmt.Sprintf("[%s]: ", getFormatTime())+"["+msg.Username+"]"+": "+msg.Text+"\n")
// 			if err != nil {
// 				break
// 			}
// 		}
// 	}
// }

// func main() {
// 	args := os.Args[1:]
// 	if len(args) != 1 {
// 		fmt.Println("[USAGE]: ./net-cat $port")
// 		return
// 	}
// 	port := args[0]
// 	server, err := net.Listen("tcp", ":"+port)
// 	if err != nil {
// 		log.Fatalln(err.Error())
// 	}
// 	defer server.Close()
// 	fmt.Printf("Listening on port %s\n", port)
// 	fmt.Printf("Command: \"nc localhost %s\"\n", port)
// 	chatServer := &ChatServer{
// 		Users: make(map[string]User),
// 		Join:  make(chan User),
// 		Leave: make(chan User),
// 		Input: make(chan Message),
// 	}
// 	go chatServer.Run()

// 	for {
// 		conn, err := server.Accept()
// 		if err != nil {
// 			log.Fatalln(err.Error())
// 		}
// 		go handleConn(chatServer, conn)
// 	}
// }

// func getFormatTime() string {
// 	now := time.Now()
// 	return fmt.Sprintf("%d-%02d-%02d %02d:%02d:%02d", now.Year(), now.Month(), now.Day(), now.Hour(), now.Minute(), now.Second())
// }

// func checkMsg(str string) bool {
// 	if str == "\n" || str == "" || str == "\t" || str == " " {
// 		return true
// 	}
// 	return false
// }
