package models

import (
	"time"
)

type Message struct {
	ID        int64     `json:"id" bson:"id,omitempty"`
	Content   string    `json:"content" bson:"content"`
	CreatedAt time.Time `json:"created_at" bson:"created_at"`
	HostIP    string    `json:"host_ip" bson:"host_ip"`
}
