package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"os/signal"

	"golang.org/x/net/websocket"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Fprintf(os.Stderr, "usage: %s [URL]", os.Args[0])
		os.Exit(1)
	}
	ws := connect(os.Args[1])
	trapCtrlC(ws)
	go write(ws)
	read(ws)
}

func connect(addr string) *websocket.Conn {
	log.Printf("connecting to %s...", addr)
	ws, err := websocket.Dial(addr, "", addr)
	if err != nil {
		log.Fatal(err)
	}
	log.Print("ready, exit with CTRL+C.")
	return ws
}

// Graceful shutdown
func trapCtrlC(c *websocket.Conn) {
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, os.Interrupt)
	go func() {
		for range ch {
			fmt.Println("\nexiting")
			c.Close()
			os.Exit(0)
		}
	}()
}

// Send STDIN lines to websocket server.
func write(ws *websocket.Conn) {
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		t := scanner.Text()
		ws.Write([]byte(t))
		fmt.Printf(">> %s\n", t)
	}
}

// Read from websocket and print messages to STDOUT
func read(ws *websocket.Conn) {
	msg := make([]byte, 512)
	for {
		n, err := ws.Read(msg)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("<< %s\n", msg[:n])
	}
}
