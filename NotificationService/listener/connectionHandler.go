package listener

import (
	"fmt"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
)

type ConnectionWithName struct {
	connection *websocket.Conn
	name       string
}

var connections []ConnectionWithName

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func Listen() {

	fmt.Println("Starting service")
	setupRouter()
	log.Fatal(http.ListenAndServe(":8080", nil))

}

func homePage(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Home page")
}
func wsEndpoint(w http.ResponseWriter, r *http.Request) {
	upgrader.CheckOrigin = func(r *http.Request) bool { return true }

	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
	}

	log.Println("Client connected")

	err = ws.WriteMessage(1, []byte("Hi client"))
	if err != nil {
		log.Println(err)
	}

	_, p, err := ws.ReadMessage()
	if err != nil {
		log.Println(err)
		return
	}
	connections = append(connections, ConnectionWithName{
		name:       string(p),
		connection: ws,
	})

	fmt.Println("ID OF USER THAT CONNECTED:", string(p))
	reader(ws)
}

func SendMessage(message string, recipient string) {
	for _, connection := range connections {
		if connection.name == recipient {
			err := connection.connection.WriteMessage(1, []byte(fmt.Sprintf("Message for: %s [%s]", recipient, message)))
			if err != nil {
				log.Println(err)
				err := connection.connection.Close()
				if err != nil {
					fmt.Println("Error occurred with connection while closing")
				}
			}
		}
	}
}

func reader(conn *websocket.Conn) {
	for {
		messageType, p, err := conn.ReadMessage()
		if err != nil {
			log.Println(err)
			return
		}

		fmt.Println(string(p))

		if err := conn.WriteMessage(messageType, p); err != nil {
			log.Println(err)
			return
		}
	}
}

func setupRouter() {

	http.HandleFunc("/", homePage)
	http.HandleFunc("/ws", wsEndpoint)
}
