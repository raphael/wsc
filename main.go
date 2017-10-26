package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"os/signal"
	"strings"

	"golang.org/x/net/websocket"
)

type headers []string

func (h *headers) String() string {
	return strings.Join(*h, ", ")
}

func (h *headers) Set(value string) error {
	*h = append(*h, value)
	return nil
}

func main() {
	var (
		target  = flag.String("u", "", "The URL to connect to")
		origin  = flag.String("o", "", "The origin to use in the WS request")
		quiet   = flag.Bool("q", false, "Only read from the socket")
		quietMode bool
		h       headers
		origURL *url.URL
	)
	flag.Var(&h, "H", `Headers to use in the WS request, can be used to multiple times to specify multiple headers.`+
		` Example: -H "Sample-Header-1: foo" -H "Sample-Header-2: bar"`)
	flag.Parse()

	if *target == "" {
		fmt.Fprintf(os.Stderr, "missing url\n")
		os.Exit(1)
	}

	if *origin != "" {
		var err error
		origURL, err = url.Parse(*origin)
		if err != nil {
			fmt.Fprintf(os.Stderr, "failed to parse origin URL: %s", err.Error())
			os.Exit(1)
		}
	}
	if *quiet == true {
		quietMode = true
	} else {
		quietMode = false
	}
	ws := connect(*target, makeHeader(h), origURL, quietMode)
	trapCtrlC(ws, quietMode)
	go write(ws, quietMode)
	read(ws, quietMode)
}

func makeHeader(h headers) http.Header {
	httpH := make(http.Header)
	for _, hv := range h {
		splits := strings.SplitN(hv, ":", 2)
		httpH.Add(strings.TrimSpace(splits[0]), strings.TrimSpace(splits[1]))
	}
	return httpH
}

func connect(addr string, h http.Header, origin *url.URL, quietMode bool) *websocket.Conn {
	if !quietMode {
		log.Printf("connecting to %s...", addr)
	}
	conf, err := websocket.NewConfig(addr, addr)
	if err != nil {
		log.Fatal(err)
	}
	conf.Header = h
	conf.Origin = origin
	ws, err := websocket.DialConfig(conf)
	if err != nil {
		log.Fatal(err)
	}
	if !quietMode {
		log.Print("ready, exit with CTRL+C.")
	}
	return ws
}

// Graceful shutdown
func trapCtrlC(c *websocket.Conn, quietMode bool) {
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, os.Interrupt)
	go func() {
		for range ch {
			if !quietMode {
				fmt.Println("\nexiting")
			}
			c.Close()
			os.Exit(0)
		}
	}()
}

// Send STDIN lines to websocket server.
func write(ws *websocket.Conn, quietMode bool) {
	if !quietMode {
		return
	}

	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		t := scanner.Text()
		ws.Write([]byte(t))
		fmt.Printf(">> %s\n", t)
	}
}

// Read from websocket and print messages to STDOUT
func read(ws *websocket.Conn, quietMode bool) {
	msg := make([]byte, 16384)
	for {
		n, err := ws.Read(msg)
		if err != nil {
			log.Fatal(err)
		}
		if !quietMode {
			fmt.Printf("<< %s\n", msg[:n])
		} else {
			fmt.Printf("%s\n", msg[:n])
		}
	}
}
