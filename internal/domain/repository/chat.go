package repository

import (
	"github.com/jackc/pgx/v5/pgxpool"
)

type ChatRepository interface {
}

type chatRepository struct {
	db *pgxpool.Pool
}

func NewChatRepository(db *pgxpool.Pool) ChatRepository {
	return &chatRepository{
		db: db,
	}
}

//func (c *chatRepository) CreateMessage(ctx context.Context, messageData chat.Message) (int, error) {
//	query := `INSERT INTO messages (id, message) VALUES ($1, $2)`
//	c.db.Exec(ctx)
//
//	return id, nil
//}
