package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net"
	"os"
	"path"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/gemalto/requester"
	"github.com/meetwithabhishek/blabber/common"
	"golang.org/x/net/websocket"
)

type model struct {
	messages []common.MessageResponse
	input    string
	ws       *websocket.Conn
}

var teaProgram *tea.Program

var styleText = lipgloss.NewStyle().
	Bold(true).
	Foreground(lipgloss.Color("#FAFAFA"))

func initModel() *model {
	address := net.JoinHostPort(conf.ServerAddress, "8080")

	replaceAddress := func(s string) string {
		return strings.Replace(s, "<server-address>", address, 1)
	}

	// Fetch initial messages
	var r []common.MessageResponse
	_, _, err := requester.ReceiveContext(context.Background(), &r,
		requester.Get(replaceAddress("http://<server-address>/messages")),
	)
	if err != nil {
		panic(err)
	}

	// Establish WebSocket connection
	wsURL := replaceAddress("ws://<server-address>/ws?username=") + conf.Username
	ws, err := websocket.Dial(wsURL, "", replaceAddress("http://<server-address>/"))
	if err != nil {
		log.Fatalf("Failed to connect to WebSocket server: %v", err)
	}

	m := model{
		messages: r,
		ws:       ws,
	}

	go m.websocketLoop()

	return &m
}

func (m *model) websocketLoop() {
	var data common.WebSocketMessage
	var dataInBytes []byte

	for {
		if err := websocket.Message.Receive(m.ws, &dataInBytes); err != nil {
			log.Fatalf("failed to receive ws message: %v", err)
		}

		if err := json.Unmarshal(dataInBytes, &data); err != nil {
			log.Fatalf("failed to unmarshal websocket message: %v", err)
		}

		switch data.MessageType {
		case common.ServerDataMessage:
			var response common.MessageResponse
			err := json.Unmarshal(data.Message, &response)
			if err != nil {
				log.Fatalf(`failed to unmarshal message response: %v`, err)
			}
			m.messages = append(m.messages, response)
			teaProgram.Send(struct{}{})
		default:
			log.Fatalln("unknown message type")
		}

	}
}

func (m *model) Init() tea.Cmd {
	return nil
}

func (m *model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit
		case "enter":
			// Send the message to the WebSocket server
			data, err := json.Marshal(common.WebSocketMessage{MessageType: common.ClientDataMessage, Message: []byte(m.input)})
			if err != nil {
				log.Fatalf("Failed to marshal message: %v", err)
			}
			if err := websocket.Message.Send(m.ws, data); err != nil {
				log.Fatalf("Failed to send message: %v", err)
			}
			m.input = ""
		case "backspace":
			if len(m.input) > 0 {
				m.input = m.input[:len(m.input)-1]
			}
		default:
			m.input += msg.String()
		}
	}
	return m, nil
}

func (m *model) View() string {
	var s string
	for _, v := range m.messages {
		s += getUsernameStyle(v.Username).Render(v.Username) + ": " + styleText.Render(v.Message) + "\n\n"
	}
	s += getUsernameStyle(conf.Username).Render(conf.Username) + ": " + styleText.Render(m.input)
	return s
}

func main() {
	err := ensureConfigExists()
	if err != nil {
		log.Fatalf("failed to create config: %v", err)
	}

	p := tea.NewProgram(initModel())
	teaProgram = p
	if _, err := p.Run(); err != nil {
		fmt.Printf("Alas, there's been an error: %v", err)
		os.Exit(1)
	}
}

// PlayPath is a directory where tool can store its data files, config files, cache files, etc.
// Its a directory just for this tool. 
// GetPlayPath takes options elements to append to the base play path and returns the full path.
func GetPlayPath(elem ...string) string {
	h := os.Getenv("HOME")
	pl := path.Join(h, "."+common.AppName)

	return path.Join(append([]string{pl}, elem...)...)
}
