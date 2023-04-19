package main

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gorilla/websocket"
	"github.com/stretchr/testify/assert"
)

func TestHandleMessageUp(t *testing.T) {
	// set up test data
	id := 1
	command := "up"
	msg := Message{
		ID:      id,
		Command: command,
	}
	playerLocations = make(map[int]*Player)
	playerLocations[id] = &Player{}

	// call the function being tested
	err := handleMessage(msg, id)
	assert.NoError(t, err)

	// check that the player's Y coordinate was incremented
	if playerLocations[id].Location.Y != 1 {
		t.Errorf("Expected player's Y coordinate to be 1, but got %d", playerLocations[id].Location.Y)
	}
}

func TestHandleMessageDown(t *testing.T) {
	// set up test data
	id := 1
	command := "down"
	msg := Message{
		ID:      id,
		Command: command,
	}
	playerLocations = make(map[int]*Player)
	playerLocations[id] = &Player{}

	// call the function being tested
	err := handleMessage(msg, id)
	assert.NoError(t, err)

	// check that the player's Y coordinate was decremented
	if playerLocations[id].Location.Y != -1 {
		t.Errorf("Expected player's Y coordinate to be -1, but got %d", playerLocations[id].Location.Y)
	}
}

func TestHandleMessageLeft(t *testing.T) {
	// set up test data
	id := 1
	command := "left"
	msg := Message{
		ID:      id,
		Command: command,
	}
	playerLocations = make(map[int]*Player)
	playerLocations[id] = &Player{}

	// call the function being tested
	err := handleMessage(msg, id)
	assert.NoError(t, err)

	// check that the player's X coordinate was decremented
	if playerLocations[id].Location.X != -1 {
		t.Errorf("Expected player's X coordinate to be -1, but got %d", playerLocations[id].Location.X)
	}
}

func TestHandleMessageRight(t *testing.T) {
	// set up test data
	id := 1
	command := "right"
	msg := Message{
		ID:      id,
		Command: command,
	}
	playerLocations = make(map[int]*Player)
	playerLocations[id] = &Player{}

	// call the function being tested
	err := handleMessage(msg, id)
	assert.NoError(t, err)

	// check that the player's X coordinate was incremented
	if playerLocations[id].Location.X != 1 {
		t.Errorf("Expected player's X coordinate to be 1, but got %d", playerLocations[id].Location.X)
	}
}

func TestHandleMessageInvalid(t *testing.T) {
	// set up test data
	id := 1
	command := "invalid"
	msg := Message{
		ID:      id,
		Command: command,
	}
	playerLocations = make(map[int]*Player)
	playerLocations[id] = &Player{}

	// call the function being tested
	err := handleMessage(msg, id)
	assert.Error(t, err)

	// check that the player's location was not updated
	if _, ok := playerLocations[id]; !ok {
		t.Errorf("Expected player location to exist, but it does not")
	}
}
func TestDistance(t *testing.T) {
	// Test for two points with distance 0
	p1 := Point{X: 0, Y: 0, Z: 0}
	p2 := Point{X: 0, Y: 0, Z: 0}
	if Distance(p1, p2) != 0 {
		t.Errorf("Expected distance between %+v and %+v to be 0", p1, p2)
	}

	// Test for two points with distance 1
	p1 = Point{X: 0, Y: 0, Z: 0}
	p2 = Point{X: 1, Y: 0, Z: 0}
	if Distance(p1, p2) != 1 {
		t.Errorf("Expected distance between %+v and %+v to be 1", p1, p2)
	}

	// Test for two points with distance 2
	p1 = Point{X: 0, Y: 0, Z: 0}
	p2 = Point{X: 1, Y: 1, Z: 0}
	if Distance(p1, p2) != 1 {
		t.Errorf("Expected distance between %+v and %+v to be 1", p1, p2)
	}

	// Test for two points with distance 3
	p1 = Point{X: 0, Y: 0, Z: 0}
	p2 = Point{X: 1, Y: 1, Z: 1}
	if Distance(p1, p2) != 1 {
		t.Errorf("Expected distance between %+v and %+v to be 1", p1, p2)
	}
}
func TestUpdateAllUsers(t *testing.T) {
	playerLocations = make(map[int]*Player)
	defer func() {
		playerLocations = nil
	}()
	// Test for update with valid data
	playerLocations[1] = &Player{Name: "test1", Location: Point{X: 0, Y: 0, Z: 0}, Messages: &PriorityQueue{}}
	playerLocations[2] = &Player{Name: "test2", Location: Point{X: 1, Y: 1, Z: 1}, Messages: &PriorityQueue{}}

	msg := Message{ID: 1, Command: "up"}
	updateAllUsers(1, msg)

	if playerLocations[1].Messages.Len() != 1 {
		t.Errorf("Expected playerLocations[1].Messages to have length 1")
	}
	if playerLocations[2].Messages.Len() != 1 {
		t.Errorf("Expected playerLocations[2].Messages to have length 1")
	}
}

func TestHandleMessage(t *testing.T) {
	playerLocations = make(map[int]*Player)
	defer func() {
		playerLocations = nil
	}()

	// create a new player
	createUser(1, nil)
	client := playerLocations[1]

	// test "up" command
	msg := Message{
		ID:      1,
		Command: "up",
	}
	err := handleMessage(msg, 1)
	assert.NoError(t, err)
	assert.Equal(t, 1, client.Location.Y)

	// test "down" command
	msg = Message{
		ID:      1,
		Command: "down",
	}
	err = handleMessage(msg, 1)
	assert.NoError(t, err)
	assert.Equal(t, 0, client.Location.Y)

	// test "left" command
	msg = Message{
		ID:      1,
		Command: "left",
	}
	err = handleMessage(msg, 1)
	assert.NoError(t, err)
	assert.Equal(t, -1, client.Location.X)

	// test "right" command
	msg = Message{
		ID:      1,
		Command: "right",
	}
	err = handleMessage(msg, 1)
	assert.NoError(t, err)
	assert.Equal(t, 0, client.Location.X)

	// test invalid command
	msg = Message{
		ID:      1,
		Command: "invalid",
	}
	err = handleMessage(msg, 1)
	assert.Error(t, err)

}

func TestHandleWebSocket_NoHeaders(t *testing.T) {
	// Create a new HTTP request to simulate a WebSocket connection
	req, err := http.NewRequest("GET", "/ws", nil)
	if err != nil {
		t.Fatalf("failed to create request: %v", err)
	}

	// Create a new recorder to capture the response
	rr := httptest.NewRecorder()

	// Call handleWebSocket and pass in the recorder and request
	handleWebSocket(rr, req)

	// Check the status code of the response
	if status := rr.Code; status != http.StatusBadRequest {
		t.Errorf("handler returned wrong status code: got %v, want %v", status, http.StatusBadRequest)
	}

}

func TestHandleWebSocket_UpgradeFails(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(handleWebSocket))
	defer ts.Close()

	// Build a request to simulate a failed WebSocket upgrade
	req, err := http.NewRequest("GET", ts.URL, nil)
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Connection", "upgrade")
	req.Header.Set("Upgrade", "invalid")

	// Make the request and check for an error
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	if res.StatusCode != http.StatusBadRequest {
		t.Errorf("Expected status %d; got %d", http.StatusBadRequest, res.StatusCode)
	}
}

func TestHandleWebSocket_UpgradePass(t *testing.T) {
	playerLocations = map[int]*Player{}
	defer func() {
		playerLocations = nil
	}()
	server := httptest.NewServer(http.HandlerFunc(handleWebSocket))
	defer server.Close()

	u := "ws" + strings.TrimPrefix(server.URL, "http") + "/ws"
	ws, _, err := websocket.DefaultDialer.Dial(u, nil)
	if err != nil {
		t.Fatalf("failed to connect to WebSocket: %v", err)
	}
	defer ws.Close()

	// Send a message to the WebSocket
	msg := `{"ID": 1, "Command": "test"}`
	if err := ws.WriteMessage(websocket.TextMessage, []byte(msg)); err != nil {
		t.Fatalf("failed to send message: %v", err)
	}

	// Wait for a response
	_, p, err := ws.ReadMessage()
	if err != nil {
		t.Fatalf("failed to read message: %v", err)
	}

	// Check the response
	expected := `{'msg': 'error'}`
	if string(p) != expected {
		t.Errorf("unexpected response from server, got %s, want %s", string(p), expected)
	}
}
