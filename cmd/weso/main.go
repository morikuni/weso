package main

import (
	"fmt"
	"io"
	"log"
	"os"
	"regexp"
	"strings"

	"github.com/morikuni/weso"
	"github.com/peterh/liner"
	"golang.org/x/net/websocket"
)

func Parse(s string) []string {
	r, _ := regexp.Compile("(\".*\"|[^\\s]*)")
	ss := r.FindAllString(s+" ", -1)
	for i, s := range ss {
		if strings.HasPrefix(s, "\"") {
			ss[i] = s[1 : len(s)-1]
		}
	}
	return ss
}

func Receive(conn *websocket.Conn) {
	for {
		data := make([]byte, 1024)
		if _, err := conn.Read(data); err != nil {
			if err == io.EOF {
				fmt.Println("connection closed")
				os.Exit(0)
			}
			log.Fatalln("error on receive", err)
		}
		fmt.Println("\r<<", string(data))
	}
}

func Send(conn *websocket.Conn) {
	stdin := liner.NewLiner()
	defer stdin.Close()
	stdin.SetCtrlCAborts(true)

	template, err := weso.NewTemplateFile("./template.txt")
	if err != nil {
		log.Fatalln("error on loading template", err)
	}

	stdin.SetCompleter(weso.Completer(template))

	for {
		text, err := stdin.Prompt("> ")
		if err == io.EOF || err == liner.ErrPromptAborted {
			os.Exit(0)
		} else if err != nil {
			log.Fatalln("error on input", err)
		} else if text == "" {
			continue
		}

		if strings.HasPrefix(text, ".") {
			args := Parse(text)
			name := args[0][1:]
			if template.IsDefined(name) {
				msg, err := template.Apply(name, args[1:]...)
				if err != nil {
					fmt.Println("error on template", err)
					continue
				}
				if _, err := conn.Write(msg); err != nil {
					log.Fatalln("error on send", err)
				}
			} else {
				fmt.Println("no template for", name)
			}
		} else {
			if _, err := conn.Write([]byte(text)); err != nil {
				log.Fatalln("error on send", err)
			}
		}

		stdin.AppendHistory(text)
	}
}

func main() {
	conn, err := websocket.Dial("ws://localhost:9000/chat", "", "http://localhost:9001/chat")
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	defer func() { conn.Close() }()

	go Send(conn)
	Receive(conn)
}
