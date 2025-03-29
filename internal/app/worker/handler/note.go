package handler

import (
	"context"
	"encoding/json"

	"go-gin-boilerplate/internal/config"
	"go-gin-boilerplate/internal/model"

	"github.com/hibiken/asynq"
	"github.com/sirupsen/logrus"
)

const (
	TypeNoteCreated = "note:created"
)

// NoteHandler handles note-related background tasks
type NoteHandler struct {
	logger *logrus.Logger
}

// NewNoteHandler creates a new note handler
func NewNoteHandler(logger *logrus.Logger) *NoteHandler {
	return &NoteHandler{
		logger: logger,
	}
}

// HandleNoteCreatedTask processes the note creation task
func (h *NoteHandler) HandleNoteCreatedTask(ctx context.Context, t *asynq.Task) error {
	var note model.Note
	if err := json.Unmarshal(t.Payload(), &note); err != nil {
		return err
	}

	h.logger.WithFields(logrus.Fields{
		"id":         note.ID,
		"title":      note.Title,
		"content":    note.Content,
		"created_at": note.CreatedAt,
	}).Info("Processing note creation task")

	return nil
}

// EnqueueNoteCreatedTask enqueues a new note creation task
func EnqueueNoteCreatedTask(client *asynq.Client, note *model.Note) error {
	payload, err := json.Marshal(note)
	if err != nil {
		return err
	}

	task := asynq.NewTask(TypeNoteCreated, payload, asynq.Queue(config.QueueHigh))
	_, err = client.Enqueue(task)
	return err
}
