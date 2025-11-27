// platform/databases/repositories/messages/repository.go
package messages

import (
	"context"
	"database/sql"
)

type Message struct {
	ID        int    `db:"id"`
	Codigo    string `db:"codigo"`
	Tipo      string `db:"tipo"`
	Categoria string `db:"categoria"`
	Contenido string `db:"contenido"`
}

type MessageRepository interface {
	GetAllMessages(ctx context.Context) ([]Message, error)
	GetMessageByCode(ctx context.Context, codigo string) (*Message, error)
}

type repository struct {
	db *sql.DB
}

func NewMessageRepository(db *sql.DB) MessageRepository {
	return &repository{db: db}
}

func (r *repository) GetAllMessages(ctx context.Context) ([]Message, error) {
	query := "SELECT id, codigo, tipo, categoria, contenido FROM mensajes"
	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var messages []Message
	for rows.Next() {
		var msg Message
		if err := rows.Scan(&msg.ID, &msg.Codigo, &msg.Tipo, &msg.Categoria, &msg.Contenido); err != nil {
			return nil, err
		}
		messages = append(messages, msg)
	}
	return messages, rows.Err()
}

func (r *repository) GetMessageByCode(ctx context.Context, codigo string) (*Message, error) {
	query := "SELECT id, codigo, tipo, categoria, contenido FROM mensajes WHERE codigo = ?"
	var msg Message
	err := r.db.QueryRowContext(ctx, query, codigo).Scan(
		&msg.ID, &msg.Codigo, &msg.Tipo, &msg.Categoria, &msg.Contenido,
	)
	if err != nil {
		return nil, err
	}
	return &msg, nil
}
