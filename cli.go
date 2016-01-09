package weso

import (
	"errors"
	"fmt"
	"io"
	"os"
	"regexp"
	"strings"

	"github.com/peterh/liner"
	"golang.org/x/net/websocket"
)

// Parse parse string to template name and args.
func Parse(s string) []string {
	r, _ := regexp.Compile("(\".+\"|[^\\s]+)")
	ss := r.FindAllString(s+" ", -1)
	for i, s := range ss {
		if strings.HasPrefix(s, "\"") && strings.HasSuffix(s, "\"") && len(s) >= 2 {
			ss[i] = s[1 : len(s)-1]
		}
	}
	return ss
}

// CLI is interface for websocket.
type CLI struct {
	conn     *websocket.Conn
	liner    *liner.State
	template *Template
	debug    bool
}

// NewCLI create CLI from Config.
func NewCLI(conf *Config) (*CLI, error) {
	origin := "http://127.0.0.1"
	if conf.Origin != "" {
		origin = conf.Origin
	}
	conn, err := websocket.Dial(conf.URL, "", origin)
	if err != nil {
		return nil, err
	}
	template := EmptyTemplate
	if _, err := os.Stat(conf.Template); err == nil {
		template, err = NewTemplateFile(conf.Template)
		if err != nil {
			return nil, err
		}
	}

	cli := &CLI{
		conn,
		liner.NewLiner(),
		template,
		conf.Debug,
	}

	cli.liner.SetCtrlCAborts(true)
	cli.liner.SetCompleter(Completer(template))

	return cli, nil
}

// QuitRequest is used when websocket should be closed.
var QuitRequest = errors.New("quit request")

// Close close connection.
func (c *CLI) Close() {
	c.conn.Close()
	c.liner.Close()
}

// Run start receive and send loop.
func (c *CLI) Run() error {
	ec := make(chan error)
	go func() {
		for {
			err := c.Send()
			if err != nil {
				ec <- err
				return
			}
		}
	}()

	go func() {
		for {
			err := c.Receive()
			if err != nil {
				ec <- err
				return
			}
		}
	}()

	err := <-ec
	if err == QuitRequest {
		return nil
	}
	return err
}

// Send read command from Stdin and send message.
func (c *CLI) Send() error {
	text, err := c.liner.Prompt("> ")
	if err == io.EOF || err == liner.ErrPromptAborted {
		return QuitRequest
	} else if err != nil {
		return err
	} else if text == "" {
		return nil
	}

	c.liner.AppendHistory(text)

	if strings.HasPrefix(text, ".") {
		args := Parse(text)
		name := args[0][1:]
		if c.template.IsDefined(name) {
			msg, err := c.template.Apply(name, args[1:]...)
			if err != nil {
				fmt.Println("! error on template", name, err)
				return nil
			}
			if c.debug {
				fmt.Println("! expanded template:", string(msg))
			}
			if _, err := c.conn.Write(msg); err != nil {
				return err
			}
		} else {
			fmt.Println("! no template for", name)
		}
	} else {
		if _, err := c.conn.Write([]byte(text)); err != nil {
			return err
		}
	}

	return nil
}

// Receive receive messages and print them.
func (c *CLI) Receive() error {
	data := make([]byte, 1024)
	if _, err := c.conn.Read(data); err != nil {
		if err == io.EOF {
			fmt.Println("\r! connection closed")
			return QuitRequest
		}
		return err
	}
	fmt.Println("\r<<", string(data))
	return nil
}
