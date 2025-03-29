package model

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Note represents a note document in MongoDB
type Note struct {
	ID        primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Title     string            `bson:"title" json:"title" binding:"required"`
	Content   string            `bson:"content" json:"content" binding:"required"`
	CreatedAt time.Time         `bson:"created_at" json:"created_at"`
	UpdatedAt time.Time         `bson:"updated_at" json:"updated_at"`
}

// BeforeCreate sets the timestamps before creating a new note
func (n *Note) BeforeCreate() {
	now := time.Now()
	n.CreatedAt = now
	n.UpdatedAt = now
}

// BeforeUpdate sets the updated timestamp before updating a note
func (n *Note) BeforeUpdate() {
	n.UpdatedAt = time.Now()
} 