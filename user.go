package main

import (
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"strings"
)

type Vector3 struct {
	X float64
	Y float64
	Z float64
}

type User struct {
	id     string
	Name   string
	Pos    Vector3
	Rot    Vector3
	Active bool
	conn   *websocket.Conn
}

func NewUser(name string, conn *websocket.Conn) *User {
	return &User{
		id:     uuid.NewString(),
		Name:   name,
		Pos:    Vector3{X: 0, Y: 0, Z: 0},
		Rot:    Vector3{X: 0, Y: 0, Z: 0},
		Active: true,
		conn:   conn,
	}
}

func (u *User) HandleMessage(m []byte) {
	message := string(m)
	if strings.Contains(message, "|") {
		cmd := strings.Split(message, "|")
		switch cmd[0] {
		case "pos":
			var pos Vector3
			json.Unmarshal([]byte(cmd[1]), &pos)
			u.Pos = pos
		case "rot":
			var rot Vector3
			json.Unmarshal([]byte(cmd[1]), &rot)
			u.Rot = rot
		}
	}
}

func (u *User) Broadcast(m interface{}) {
	message, err := json.Marshal(m)
	if err != nil {
		fmt.Println("Could not marshal object", m)
	}
	u.BroadcastBytes(message)
}

func (u *User) BroadcastBytes(m []byte) {
	err := u.conn.WriteMessage(websocket.BinaryMessage, m)
	if err != nil {
		fmt.Println("Error broadcasting message", err)
	}
}
