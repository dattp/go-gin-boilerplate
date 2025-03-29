package service

import (
	"context"

	"go-gin-boilerplate/internal/app/worker/handler"
	"go-gin-boilerplate/internal/eventbus"
	"go-gin-boilerplate/internal/model"
	"go-gin-boilerplate/internal/repository"

	"github.com/hibiken/asynq"
	"github.com/sirupsen/logrus"
)

// EventNoteCreated is the event name for note creation
const EventNoteCreated = "note:created"

// NoteService handles business logic for notes
type NoteService struct {
	repo        *repository.NoteRepository
	logger      *logrus.Logger
	asynqClient *asynq.Client
	eventBus    *eventbus.EventBus
}

// NewNoteService creates a new note service
func NewNoteService(
	repo *repository.NoteRepository,
	logger *logrus.Logger,
	asynqClient *asynq.Client,
	eventBus *eventbus.EventBus,
) *NoteService {
	service := &NoteService{
		repo:        repo,
		logger:      logger,
		asynqClient: asynqClient,
		eventBus:    eventBus,
	}
	eventBus.Subscribe(EventNoteCreated, service.handleNoteCreatedEvent)
	return service
}

// CreateNote creates a new note and enqueues a background task
func (s *NoteService) CreateNote(ctx context.Context, note *model.Note) error {
	if err := s.repo.Create(ctx, note); err != nil {
		return err
	}

	// Publish note creation event to event bus
	s.eventBus.Publish(EventNoteCreated, note)

	return nil
}

// GetNote retrieves a note by ID
func (s *NoteService) GetNote(ctx context.Context, id string) (*model.Note, error) {
	return s.repo.GetByID(ctx, id)
}

// GetAllNotes retrieves all notes
func (s *NoteService) GetAllNotes(ctx context.Context) ([]model.Note, error) {
	return s.repo.GetAll(ctx)
}

// UpdateNote updates an existing note
func (s *NoteService) UpdateNote(ctx context.Context, note *model.Note) error {
	return s.repo.Update(ctx, note)
}

// DeleteNote removes a note
func (s *NoteService) DeleteNote(ctx context.Context, id string) error {
	return s.repo.Delete(ctx, id)
}

func (s *NoteService) handleNoteCreatedEvent(note *model.Note) {
	// Enqueue note creation task
	if err := handler.EnqueueNoteCreatedTask(s.asynqClient, note); err != nil {
		s.logger.WithError(err).Error("Failed to enqueue note creation task")
		// Don't return error as the note was created successfully
	}
}
