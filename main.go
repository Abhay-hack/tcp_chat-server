package main

import (
	"bufio"
	"fmt"
	"net"
	"strings"
	"sync"
	"time"
)

// Client represents a connected user
type Client struct {
	conn     net.Conn
	name     string
	messages chan string
}

// Server holds the state of the chat application
type Server struct {
	clients    map[net.Conn]*Client
	broadcast  chan string
	register   chan *Client
	unregister chan *Client
	mu         sync.Mutex
}

func NewServer() *Server {
	return &Server{
		clients:    make(map[net.Conn]*Client),
		broadcast:  make(chan string),
		register:   make(chan *Client),
		unregister: make(chan *Client),
	}
}

// Run manages the lifecycle of all connections
func (s *Server) Run() {
	for {
		select {
		case client := <-s.register:
			s.mu.Lock()
			s.clients[client.conn] = client
			s.mu.Unlock()
			fmt.Printf("New client joined: %s\n", client.conn.RemoteAddr())

		case client := <-s.unregister:
			s.mu.Lock()
			delete(s.clients, client.conn)
			s.mu.Unlock()
			client.conn.Close()
			fmt.Printf("Client disconnected: %s\n", client.name)

		case msg := <-s.broadcast:
			s.mu.Lock()
			for _, client := range s.clients {
				// Send message to each client's specific channel
				client.messages <- msg
			}
			s.mu.Unlock()
		}
	}
}

func handleConnection(conn net.Conn, s *Server) {
	defer func() {
		conn.Close()
	}()

	// 1. Get Nickname
	conn.Write([]byte("Enter your nickname: "))
	nameReader := bufio.NewReader(conn)
	name, _ := nameReader.ReadString('\n')
	name = strings.TrimSpace(name)

	if name == "" {
		name = "Anonymous"
	}

	client := &Client{
		conn:     conn,
		name:     name,
		messages: make(chan string),
	}

	s.register <- client

	// 2. Start a writer goroutine for this specific client
	go func() {
		for msg := range client.messages {
			fmt.Fprintln(client.conn, msg)
		}
	}()

	s.broadcast <- fmt.Sprintf(">>> %s has joined the chat!", name)

	// 3. Read loop: Listen for messages from this client
	scanner := bufio.NewScanner(conn)
	for scanner.Scan() {
		text := scanner.Text()
		if text == "/quit" {
			break
		}
		t := time.Now().Format("15:04")
		// This format looks much cleaner: [14:20] [Abhay]: Hello!
		s.broadcast <- fmt.Sprintf("[%s] [%s]: %s", t, name, text)

	}

	s.broadcast <- fmt.Sprintf("<<< %s has left the chat.", name)
	s.unregister <- client
}

func main() {
	server := NewServer()
	go server.Run()

	listener, err := net.Listen("tcp", ":8080")
	if err != nil {
		fmt.Println("Error starting server:", err)
		return
	}
	defer listener.Close()

	fmt.Println("Chat Server started on :8080")

	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("Connection error:", err)
			continue
		}
		go handleConnection(conn, server)
	}
}
