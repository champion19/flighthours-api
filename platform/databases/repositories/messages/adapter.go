// platform/databases/repositories/messages/adapter.go
package messages

import (
	"context"

	"github.com/champion19/flighthours-api/core/interactor/services/domain"
)

// MessageRepositoryAdapter adapts the repository MessageRepository to domain MessageRepository
type MessageRepositoryAdapter struct {
	repo MessageRepository
}

func NewMessageRepositoryAdapter(repo MessageRepository) domain.MessageRepository {
	return &MessageRepositoryAdapter{repo: repo}
}

func (a *MessageRepositoryAdapter) GetAllMessages(ctx context.Context) ([]domain.Message, error) {
	repoMessages, err := a.repo.GetAllMessages(ctx)
	if err != nil {
		return nil, err
	}

	domainMessages := make([]domain.Message, len(repoMessages))
	for i, msg := range repoMessages {
		domainMessages[i] = domain.Message{
			Codigo:    msg.Codigo,
			Contenido: msg.Contenido,
		}
	}

	return domainMessages, nil
}
