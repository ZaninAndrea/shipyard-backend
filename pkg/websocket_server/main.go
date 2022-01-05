package main

import (
	"flag"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{} // use default options

type PubSub struct{
	topicToConnections map[string][]*websocket.Conn
	connectionToTopics map[*websocket.Conn][]string
}

func (pb *PubSub) Subscribe(topic string, c *websocket.Conn){
	connections, ok := pb.topicToConnections[topic]
	if !ok{
		connections = make([]*websocket.Conn, 0)
	}
	pb.topicToConnections[topic] = append(connections, c)

	topics, ok := pb.connectionToTopics[c]
	if !ok{
		topics = make([]string, 0)
	}
	pb.connectionToTopics[c] = append(topics, topic)
}

func (pb *PubSub) Unsubscribe(topic string, c *websocket.Conn){
	newTopics := make([]string, len(pb.connectionToTopics[c]) - 1)
	for _, t := range pb.connectionToTopics[c]{
		if t!=topic{
			newTopics = append(newTopics, t)
		}
	}
	pb.connectionToTopics[c] = newTopics
}

func subscriptionsHandler(w http.ResponseWriter, r *http.Request) {
	c, err := upgrader.Upgrade(w, r, nil)

	if err != nil {
		log.Print("upgrade:", err)
		return
	}
	defer c.Close()
	for {
		mt, message, err := c.ReadMessage()
		if err != nil {
			log.Println("read:", err)
			break
		}
		log.Printf("recv: %s", message)
		err = c.WriteMessage(mt, message)
		if err != nil {
			log.Println("write:", err)
			break
		}
	}
}

func main() {
	flag.Parse()
	log.SetFlags(0)
	http.HandleFunc("/subscriptions", subscriptionsHandler)
	log.Fatal(http.ListenAndServe("localhost:8080", nil))
}