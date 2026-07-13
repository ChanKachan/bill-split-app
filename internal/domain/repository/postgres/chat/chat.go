package chat

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

type ChatRepository interface {
	CreateMessage(ctx context.Context, messageData CreateMessangeRequest) (int, error)
	GetLastMessage(ctx context.Context, chatId int) ([]GetMessagesResponse, error)
}

type chatRepository struct {
	db *pgxpool.Pool
}

func NewChatRepository(db *pgxpool.Pool) ChatRepository {
	return &chatRepository{
		db: db,
	}
}

func (c *chatRepository) CreateMessage(ctx context.Context, messageData CreateMessangeRequest) (int, error) {
	ctx, cancel := context.WithTimeout(ctx, 15*time.Second)
	defer cancel()

	var id int

	query := `INSERT INTO chat_message (
              	user_id,
                message,
                chat_id,
                date_create,
                date_update
              ) 
			  VALUES ($1, $2, $3, $4, $5) RETURNING id;`

	err := c.db.QueryRow(
		ctx,
		query,
		messageData.UserID,
		messageData.Message,
		messageData.ChatID,
		messageData.DateCreate,
		messageData.DateUpdate,
	).Scan(&id)

	if err != nil {
		return 0, fmt.Errorf("error creating chat message to repository: %w", err)
	}

	return id, nil
}

// Возращает последние 50 сообщений
func (c *chatRepository) GetLastMessage(ctx context.Context, chatId int) ([]GetMessagesResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, 15*time.Second)
	defer cancel()

	query := `
	SELECT 
    	id,
    	user_id,
    	message,
    	data_create,
    	date_update,
    	chat_id
	FROM chat_message 
	WHERE chat_id = $1
	ORDER BY date_create DESC
	LIMIT 50 OFFSET 0`

	rows, err := c.db.Query(ctx, query, chatId)
	if err != nil {
		return nil, fmt.Errorf("error getting last message from repository: %w", err)
	}

	defer rows.Close()

	var datas []GetMessagesResponse
	for rows.Next() {
		var data GetMessagesResponse

		err := rows.Scan(
			&data.MessageID,
			&data.UserID,
			&data.Message,
			&data.DateUpdate,
			&data.ChatID,
		)
		if err != nil {
			return nil, fmt.Errorf("error getting last message: %w", err)
		}

		datas = append(datas, data)
	}
	return datas, nil
}
