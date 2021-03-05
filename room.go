package main

import "github.com/google/uuid"

type Room struct {
	id    string
	Users []*User
}

func NewRoom() *Room {
	return &Room{
		id:    uuid.NewString(),
		Users: make([]*User, 0),
	}
}

func (r *Room) filtered() *Room {
	return &Room{
		id:    r.id,
		Users: r.activeUsers(),
	}
}

func (r *Room) AddUser(u *User) {
	r.Users = append(r.Users, u)
	r.handleUserMessages(u)
}

func (r *Room) BroadcastState() {
	for _, u := range r.activeUsers() {
		u.Broadcast(r.filtered())
	}
}

func (r *Room) removeUser(u *User) {
	u.conn.Close()
	u.Active = false
}

func (r *Room) handleUserMessages(u *User) {
	defer r.removeUser(u)

	for {
		_, m, err := u.conn.ReadMessage()
		if err != nil {
			break
		}

		u.HandleMessage(m)
	}
}

func (r *Room) activeUsers() []*User {
	var users = make([]*User, 0)
	for _, u := range r.Users {
		if u.Active {
			users = append(users, u)
		}
	}
	return users
}
