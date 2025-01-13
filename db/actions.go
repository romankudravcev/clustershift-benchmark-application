package db

import (
	"benchmarker/models"
	"database/sql"
)

func SaveMessage(message *models.Message) (int64, error) {
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
	query := `DELETE FROM messages`

	_, err := DB.Exec(query)
	return err
}
