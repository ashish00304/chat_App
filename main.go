package main

import (
	"fmt"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

var (
	upgrader = websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool { return true },
	}
	clients   = make(map[*websocket.Conn]bool)
	broadcast = make(chan string)
	mutex     = sync.Mutex{}
)

// Handle incoming connections
func handleConnections(c *gin.Context) {
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		fmt.Println("WebSocket upgrade error:", err)
		return
	}
	defer conn.Close()

	mutex.Lock()
	clients[conn] = true
	mutex.Unlock()

	for {
		_, msg, err := conn.ReadMessage()
		if err != nil {
			fmt.Println("Client disconnected:", err)
			mutex.Lock()
			delete(clients, conn)
			mutex.Unlock()
			break
		}

		userMsg := string(msg)
		fmt.Println("Received:", userMsg)
		broadcast <- "🧑 You: " + userMsg

		// Generate AI reply
		go func(input string) {
			time.Sleep(1 * time.Second) // simulate thinking delay
			reply := getAIReply(input)
			broadcast <- "🤖 Bot: " + reply
		}(userMsg)
	}
}

// Handle broadcasting to all clients
func handleMessages() {
	for {
		msg := <-broadcast
		mutex.Lock()
		for client := range clients {
			err := client.WriteMessage(websocket.TextMessage, []byte(msg))
			if err != nil {
				fmt.Println("Write error:", err)
				client.Close()
				delete(clients, client)
			}
		}
		mutex.Unlock()
	}
}

// Very basic AI reply logic
func getAIReply(input string) string {
	input = strings.ToLower(input)

	switch {
	case strings.Contains(input, "hi"), strings.Contains(input, "hello"):
		return "Hey there! 👋 How are you today?"
	case strings.Contains(input, "how are you"):
		return "I'm great! Just chatting here 😄"
	case strings.Contains(input, "your name"):
		return "I'm GoBot, your friendly Go chat assistant 🤖"
	case strings.Contains(input, "time"):
		return "Current time is " + time.Now().Format("3:04 PM")
	case strings.Contains(input, "bye"):
		return "Goodbye! Have a great day 👋"
	default:
		return "Interesting... tell me more!"
	}
}

func main() {
	router := gin.Default()
	router.LoadHTMLGlob("tamplates/*")
	router.Static("./tamplates/index.js", "./tamplates/index.js")

	router.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.html", nil)
	})

	router.GET("/ws", handleConnections)

	go handleMessages()

	router.Run(":8080")
}
