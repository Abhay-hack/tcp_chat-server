package main

import (
    "bufio"
    "fmt"
    "net"
    "os"
)

func main() {
    conn, err := net.Dial("tcp", "localhost:8080")
    if err != nil {
        fmt.Println("Could not connect to server.")
        return
    }

    // Goroutine to receive messages from server
    go func() {
        scanner := bufio.NewScanner(conn)
        for scanner.Scan() {
            // This clears the current line and prints the server message
            fmt.Printf("\r%s\n> ", scanner.Text()) 
        }
    }()

    // Main loop to send messages
    fmt.Print("> ")
    scanner := bufio.NewScanner(os.Stdin)
    for scanner.Scan() {
        fmt.Fprintln(conn, scanner.Text())
        fmt.Print("> ") // Visual prompt for the user
    }
}