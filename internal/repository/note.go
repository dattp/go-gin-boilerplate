package repository

import (
	"context"

	"go-gin-boilerplate/internal/config"
	"go-gin-boilerplate/internal/database"
	"go-gin-boilerplate/internal/model"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

// NoteRepository handles database operations for notes
type NoteRepository struct {
	collection *mongo.Collection
}

// NewNoteRepository creates a new note repository
func NewNoteRepository(db *database.MongoDBClient, cfg *config.Config) *NoteRepository {
	return &NoteRepository{
		collection: db.GetMongoDBClient().Database(cfg.MongoDB.Database).Collection("notes"),
	}
}

// Create inserts a new note into the database
func (r *NoteRepository) Create(ctx context.Context, note *model.Note) error {
	note.BeforeCreate()
	result, err := r.collection.InsertOne(ctx, note)
	if err != nil {
		return err
	}
	note.ID = result.InsertedID.(primitive.ObjectID)
	return nil
}

// GetByID retrieves a note by its ID
func (r *NoteRepository) GetByID(ctx context.Context, id string) (*model.Note, error) {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	var note model.Note
	err = r.collection.FindOne(ctx, bson.M{"_id": objectID}).Decode(&note)
	if err != nil {
		return nil, err
	}
	return &note, nil
}

// GetAll retrieves all notes
func (r *NoteRepository) GetAll(ctx context.Context) ([]model.Note, error) {
	cursor, err := r.collection.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var notes []model.Note
	if err := cursor.All(ctx, &notes); err != nil {
		return nil, err
	}
	return notes, nil
}

// Update updates an existing note
func (r *NoteRepository) Update(ctx context.Context, note *model.Note) error {
	note.BeforeUpdate()
	_, err := r.collection.UpdateOne(
		ctx,
		bson.M{"_id": note.ID},
		bson.M{"$set": note},
	)
	return err
}

// Delete removes a note by its ID
func (r *NoteRepository) Delete(ctx context.Context, id string) error {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	_, err = r.collection.DeleteOne(ctx, bson.M{"_id": objectID})
	return err
}
