package models

import "time"

type UserShowData struct {
	ID        int
	UserName  string
	NikeName  string
	Email     string
	Avatar    string
	City      string
	Age       int
	Introduce string
	CreatedAt time.Time
}
