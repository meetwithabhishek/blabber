package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"path"
	"strings"
	"sync"

	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
	"github.com/meetwithabhishek/blabber"
)

func init() {
	// create a sample config file, if it doesn't exists
	_, err := os.Stat(GetPlayPath())
	if err != nil {
		err := os.MkdirAll(GetPlayPath(), 0755)
		if err != nil {
			panic(err)
		}
	}
}

var upgrader = websocket.Upgrader{CheckOrigin: func(r *http.Request) bool { return true }}

var conns sync.Map

func websocketConnHandler(w http.ResponseWriter, r *http.Request) {
	// Get the URL query parameters
	queryParams := r.URL.Query()

	// Retrieve a specific parameter using Get()
	username := queryParams.Get("username")

	c, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Print("upgrade:", err)
		return
	}

	conns.Store(username, c)

	defer func() {
		conns.Delete(username)
		if err := c.Close(); err != nil {
			log.Fatalf("failed to close websocket connection: %v", err)
		}
	}()

	var request blabber.WebSocketMessage

	for {
		_, message, err := c.ReadMessage()
		if err != nil {
			log.Println("read:", err)
			break
		}

		if err := json.Unmarshal(message, &request); err != nil {
			log.Fatalf("failed to unmarshal websocket message: %v", err)
		}

		switch request.MessageType {
		case blabber.ClientDataMessage:
			if len(strings.Trim(string(request.Message), " \n")) == 0 {
				log.Printf("received empty message, skipping")
				continue
			}

			m, err := NewMessage(Message{Username: username, Message: string(request.Message)})
			if err != nil {
				log.Fatalln("failed to create message in db")
			}

			log.Printf("recv data message: %s", request.Message)

			err = broadcastDataMessage(*m)
			if err != nil {
				log.Fatalf("failed to broadcast message: %v", err)
			}

		default:
			log.Fatalln("unknown message type")
		}

		// no need to acknowledge the received message back to client
		// err = c.WriteMessage(mt, message)
		// if err != nil {
		// 	log.Println("write:", err)
		// 	break
		// }
	}
}

func broadcastDataMessage(m Message) error {
	conns.Range(func(key, value interface{}) bool {
		conn, ok := value.(*websocket.Conn)
		if !ok {
			return false
		}

		data, err := json.Marshal(blabber.MessageResponse{Username: m.Username, Message: m.Message})
		if err != nil {
			log.Fatalf("failed to marshal message: %v", err)
		}

		data, err = json.Marshal(blabber.WebSocketMessage{MessageType: blabber.ServerDataMessage, Message: data})
		if err != nil {
			log.Fatalf("failed to marshal message: %v", err)
		}

		err = conn.WriteMessage(websocket.BinaryMessage, data)
		if err != nil {
			log.Fatalf("failed to write message: %v", err)
		}
		return true
	})
	return nil
}

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/messages", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		messages, err := ListMessages()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		var response []blabber.MessageResponse

		for _, v := range messages {
			response = append(response, blabber.MessageResponse{Username: v.Username, Message: v.Message})
		}

		w.Header().Set("Content-Type", "application/json")
		jsonData, err := json.Marshal(response)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		if _, err := w.Write(jsonData); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}))

	r.HandleFunc("/ws", websocketConnHandler)

	http.Handle("/", r)

	err := Initialize()
	if err != nil {
		log.Fatalf("Error initializing: %v", err)
	}

	fmt.Println("Starting WebSocket server on :8080")
	err = http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatalf("ListenAndServe: %v", err)
	}
}

// GetPlayPath gives the absolute path for the safe directory inside the tool's config directory.
func GetPlayPath(elem ...string) string {
	h := os.Getenv("HOME")
	pl := path.Join(h, "."+blabber.AppName)

	return path.Join(append([]string{pl}, elem...)...)
}
