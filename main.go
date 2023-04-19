package main

import (
	"container/heap"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{} // use default options

var ClientMutex sync.Mutex
var idIndex int
var playerLocations map[int]*Player

func handleWebSocket(w http.ResponseWriter, r *http.Request) {
	log.Printf("got a request")

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("error upgrading: %v", err)
		return
	}

	ClientMutex.Lock()
	idIndex = idIndex + 1
	myID := idIndex
	log.Printf("new customer: %d", myID)
	createUser(myID, conn)
	go sendUpdates(myID)
	ClientMutex.Unlock()

	updateAllUsers(0, Message{
		ID:      idIndex,
		Command: "newUser",
	})

	defer func() {
		ClientMutex.Lock()
		delete(playerLocations, myID)
		ClientMutex.Unlock()
		conn.Close()
		log.Printf("closing connection: %d", myID)
	}()
	for {
		messageType, message, err := conn.ReadMessage()
		if err != nil {
			log.Printf("error reading message: %v", err)
			return
		}

		log.Printf("msg received:  %s [%d]", message, messageType)
		msg := Message{}

		err = json.Unmarshal(message, &msg)
		if err != nil {
			log.Printf("got something fucky: %v", err)
		}
		if msg.ID == 0 {
			msg.ID = myID
		}
		err = handleMessage(msg, myID)
		if err != nil {
			log.Printf("error handling message: %v", err)
		}
	}
}
func handleMessage(msg Message, id int) error {
	_, ok := playerLocations[id]
	if !ok {
		return fmt.Errorf("no user: %d", id)

	}
	switch msg.Command {
	case "up":
		playerLocations[id].Location.Y++
		updateAllUsers(id, msg)
	case "down":
		playerLocations[id].Location.Y--
		updateAllUsers(id, msg)
	case "left":
		playerLocations[id].Location.X--
		updateAllUsers(id, msg)
	case "right":
		playerLocations[id].Location.X++
		updateAllUsers(id, msg)
	default:
		err := updateOneUser(id, Message{})
		log.Printf("invalid message received: %v", msg)
		return err
	}
	return nil
}

func updateOneUser(id int, msg Message) error {
	client, ok := playerLocations[id]
	if !ok {
		return fmt.Errorf("no such user %d", id)
	}
	if client.Conn == nil {
		return fmt.Errorf("client connection not initialized: %d", id)
	}
	err := client.Conn.WriteMessage(1, []byte("{'msg': 'error'}"))
	if err != nil {
		log.Printf("error writing: %v", err)
		return err
	}
	return nil
}
func TimeDistortion() {
	log.Print("lets do the timewarp again")
	for {
		time.Sleep(5 * time.Millisecond)
		for _, c := range playerLocations {
			if c.Messages == nil {
				continue
			}
			if c.Messages.Len() == 0 {
				continue
			}
			item := heap.Pop(c.Messages)
			if item != nil {
				update := item.(*TimeItem)
				if update != nil {
					c.UpdatesChan <- update.value
				}
			}
		}
	}

}
func sendUpdates(id int) {
	player := playerLocations[id]
	for update := range player.UpdatesChan {
		log.Printf("sending time lagged update to %d", id)
		player.Conn.WriteJSON(update)
	}
}
func updateAllUsers(id int, msg Message) {
	log.Print("==========================update all =========================")
	player, ok := playerLocations[id]
	if !ok {
		log.Printf("no player found: %d", id)
		return
	}
	update := PlayerUpdate{
		WhoAmI:                  -1,
		WhoDidIt:                id,
		WhenDidTheyDoIt:         int(time.Now().Unix()),
		WhenShouldTheyGetThis:   0,
		HowFarAwayAreTheyFromMe: 0,
		WhatDidTheyDo:           msg.Command,
		Location:                player.Location,
	}
	for k, c := range playerLocations {
		update.WhoAmI = k
		d := Distance(c.Location, player.Location)
		when := time.Now().Add(time.Duration(d) * time.Second / 2)
		update.WhenShouldTheyGetThis = int(time.Now().Add(time.Duration(d) * time.Second).Unix())
		update.HowFarAwayAreTheyFromMe = d
		ti := &TimeItem{
			value:    update,
			priority: when,
		}
		log.Printf("sending update: %d %v", id, update)
		if playerLocations[k].Messages == nil {
			log.Printf("player: %d isn't fully init", k)
			continue
		}
		heap.Push(playerLocations[k].Messages, ti)
	}
	log.Print("==========================update all done =========================")
}
func createUser(id int, conn *websocket.Conn) {
	player := &Player{
		Name: "",
		Conn: conn,
	}
	playerLocations[id] = player
	msgs := make(PriorityQueue, 0)
	playerLocations[id].Messages = &msgs
	location := Point{}
	playerLocations[id].Location = location
	heap.Init(&msgs)
	playerLocations[id].UpdatesChan = make(chan PlayerUpdate, 100)
}

func main() {
	log.Print("starting space server")
	playerLocations = make(map[int]*Player)
	http.HandleFunc("/ws", handleWebSocket)
	go TimeDistortion()
	log.Fatal(http.ListenAndServe(":8080", nil))
}
