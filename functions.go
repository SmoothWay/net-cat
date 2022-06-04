package main

import (
	"bufio"
	"errors"
	"io/ioutil"
	"log"
	"net"
	"strings"
	"sync"
	"time"
)

// Handling connection
func handleConnection(conn net.Conn, mut *sync.Mutex) {
	newconn := bufio.NewReader(conn)

	// Checking clients nickname
	name := ""
	for {
		name, _ = newconn.ReadString('\n')
		if len(name) > 0 {
			if isValid(name[:len(name)-1]) {
				conn.Write([]byte("\n[ENTER YOUR NAME]:"))
			} else {
				break
			}
		} else {
			conn.Close()
			return
		}
	}
	name = strings.TrimSpace(name)
	newTime := strings.Split(time.Now().String(), ".")
	finalTime := "[" + newTime[0] + "]"
	finalyName := "[" + name + "]:"
	greetings := " has joined our chat..."
	parting := " has left our chat..."

	// loging
	log.Printf("connection from %v as %v", conn.RemoteAddr(), finalyName)
	// sending chat history
	conn.Write(tempHistory)

	// showing client name and greet others
	// Using mutex to synchronize processes. It used to manage integrity of resources (file, data, ram)
	mut.Lock()
	users[conn] = name
	for user, value := range users {
		if user != conn {
			user.Write([]byte("\n" + name + greetings + "\n"))
			user.Write([]byte(finalTime + "[" + value + "]:"))
		}
	}
	mut.Unlock()
	defer conn.Close()

	for {
		// initializing time
		newTime = strings.Split(time.Now().String(), ".")
		finalTime = "[" + newTime[0] + "]"

		conn.Write([]byte(finalTime + finalyName))

		message := ""
		var err error

		for {
			message, err = newconn.ReadString('\n')
			if err != nil {
				err = errors.New("Error")
				break
			} else if isValid(message) {
				conn.Write([]byte(finalTime + finalyName))
			} else {
				break
			}
		}

		// Info members if someone disconnected
		if err != nil {
			for user, value := range users {
				if user != conn {
					user.Write([]byte("\n" + name + parting + "\n"))
					user.Write([]byte(finalTime + "[" + value + "]:"))
				} else {
					delete(users, user) // deleting user from map of users if he disconnects
					log.Printf("%v has left chat. Currently %d users in chat\n", name, len(users))
				}
			}
			break
		}
		mut.Lock()

		// send message to all members of chat

		for user, value := range users {
			if user != conn {
				if message == "\n" {
					break
				}
				user.Write([]byte("\n" + finalTime + finalyName + message))
				user.Write([]byte(finalTime + "[" + value + "]:"))
			}
		}
		mut.Unlock()
		log.Print(finalyName + message)
		// save histroy of chat
		tempHistory = append(tempHistory, finalTime+finalyName+message...)

		// Write history to file

		ioutil.WriteFile("assets/chatHistory.txt", tempHistory, 0666)
	}
}

func validPort(p string) bool {
	for i := 0; i < len(p); i++ {
		if p[i] >= '0' && p[i] <= '9' {
			continue
		}
		return false
	}
	return true
}

func isValid(s string) bool {
	for i := range s {
		if s[i] != ' ' && s[i] != '\n' {
			return false
		}
	}
	return true
}
