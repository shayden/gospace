package main

import "github.com/gorilla/websocket"

type PlayerUpdate struct {
	//The recipient of this update
	WhoAmI                  int    `json:"who_am_i"`
	WhoDidIt                int    `json:"who_did_it"`
	WhenDidTheyDoIt         int    `json:"when"`
	HowFarAwayAreTheyFromMe int    `json:"distance"`
	WhatDidTheyDo           string `json:"what"`
	WhenShouldTheyGetThis   int    `json:"-"`
	Location                Point  `json:"where"`
}

type Message struct {
	ID      int
	Command string
}

type Player struct {
	Name        string
	Location    Point
	Messages    *PriorityQueue
	Conn        *websocket.Conn
	UpdatesChan chan PlayerUpdate
}

type Point struct {
	X int
	Y int
	Z int
}
