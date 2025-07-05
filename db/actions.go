package db

import (
	"benchmarker/models"
	"context"
	"database/sql"
	"go.mongodb.org/mongo-driver/bson"
	"os"
)

var dbType string

func init() {
	dbType = os.Getenv("DB_TYPE")
}

func SaveMessage(message *models.Message) (int64, error) {
	if dbType == "mongodb" {
		collection := MongoDB.Collection("messages")
		ctx := context.Background()
		res, err := collection.InsertOne(ctx, message)
		if err != nil {
			return 0, err
		}
		if oid, ok := res.InsertedID.(int64); ok {
			return oid, nil
		}
		return 0, nil // MongoDB's ObjectID is not int64, so return 0
	}

	query := `
        INSERT INTO messages (content, created_at, host_ip)
        VALUES ($1, $2, $3)
        RETURNING id`

	var messageID int64
	err := DB.QueryRow(
		query,
		message.Content,
		message.CreatedAt,
		message.HostIP,
	).Scan(&messageID)

	if err != nil {
		return 0, err
	}

	return messageID, nil
}

func GetMessage(id int64) (*models.Message, error) {
	if os.Getenv("DB_TYPE") == "mongodb" {
		collection := MongoDB.Collection("messages")
		ctx := context.Background()
		var result models.Message
		err := collection.FindOne(ctx, bson.M{"id": id}).Decode(&result)
		if err != nil {
			return nil, err
		}
		return &result, nil
	}

	query := `
        SELECT id, content, created_at, host_ip
        FROM messages
        WHERE id = $1`

	message := &models.Message{}
	err := DB.QueryRow(query, id).Scan(
		&message.ID,
		&message.Content,
		&message.CreatedAt,
		&message.HostIP,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return message, nil
}

// Get all messages
func GetMessages() ([]models.Message, error) {
	if os.Getenv("DB_TYPE") == "mongodb" {
		collection := MongoDB.Collection("messages")
		ctx := context.Background()
		cur, err := collection.Find(ctx, bson.M{})
		if err != nil {
			return nil, err
		}
		defer cur.Close(ctx)

		var messages []models.Message
		for cur.Next(ctx) {
			var msg models.Message
			if err := cur.Decode(&msg); err != nil {
				return nil, err
			}
			messages = append(messages, msg)
		}
		if err := cur.Err(); err != nil {
			return nil, err
		}
		return messages, nil
	}

	query := `
        SELECT id, content, created_at, host_ip
        FROM messages
        ORDER BY id DESC`

	rows, err := DB.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var messages []models.Message
	for rows.Next() {
		var msg models.Message
		err := rows.Scan(
			&msg.ID,
			&msg.Content,
			&msg.CreatedAt,
			&msg.HostIP,
		)
		if err != nil {
			return nil, err
		}
		messages = append(messages, msg)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return messages, nil
}

// DeleteAllMessages deletes all messages from the database
func DeleteAllMessages() error {
	if os.Getenv("DB_TYPE") == "mongodb" {
		collection := MongoDB.Collection("messages")
		ctx := context.Background()
		_, err := collection.DeleteMany(ctx, bson.M{})
		return err
	}

	query := `DELETE FROM messages`

	_, err := DB.Exec(query)
	return err
}
