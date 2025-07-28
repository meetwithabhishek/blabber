package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/gemalto/requester"
	"github.com/meetwithabhishek/blabber"
	"golang.org/x/net/websocket"
)

const Username = "abhishek"

type model struct {
	messages []blabber.MessageResponse
	input    string
	ws       *websocket.Conn
}

var teaProgram *tea.Program

var styleUsername = lipgloss.NewStyle().
	Bold(true).
	Foreground(lipgloss.Color("#FAFAFA")).
	Background(lipgloss.Color("#7D56F4"))
	// Background(lipgloss.ANSIColor(rand.Intn(7)))

var styleText = lipgloss.NewStyle().
	Bold(true).
	Foreground(lipgloss.Color("#FAFAFA"))

func initModel() *model {
	// Fetch initial messages
	var r []blabber.MessageResponse
	_, _, err := requester.ReceiveContext(context.Background(), &r,
		requester.Get("http://localhost:8080/messages"),
	)
	if err != nil {
		panic(err)
	}

	// Establish WebSocket connection
	wsURL := "ws://localhost:8080/ws?username=" + Username
	ws, err := websocket.Dial(wsURL, "", "http://localhost/")
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
	var data blabber.WebSocketMessage
	var dataInBytes []byte

	for {
		if err := websocket.Message.Receive(m.ws, &dataInBytes); err != nil {
			log.Fatalf("failed to receive ws message: %v", err)
		}

		if err := json.Unmarshal(dataInBytes, &data); err != nil {
			log.Fatalf("failed to unmarshal websocket message: %v", err)
		}

		switch data.MessageType {
		case blabber.ServerDataMessage:
			var response blabber.MessageResponse
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
			data, err := json.Marshal(blabber.WebSocketMessage{MessageType: blabber.ClientDataMessage, Message: []byte(m.input)})
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
		s += styleUsername.Render(v.Username) + ": " + styleText.Render(v.Message) + "\n\n"
	}
	s += styleUsername.Render(Username) + ": " + styleText.Render(m.input)
	return s
}

func main() {
	ensureConfigExists()

	p := tea.NewProgram(initModel())
	teaProgram = p
	if _, err := p.Run(); err != nil {
		fmt.Printf("Alas, there's been an error: %v", err)
		os.Exit(1)
	}
}

// GetPlayPath gives the absolute path for the safe directory inside the tool's config directory.
func GetPlayPath(elem ...string) string {
	h := os.Getenv("HOME")
	pl := path.Join(h, "."+blabber.AppName)

	return path.Join(append([]string{pl}, elem...)...)
}
