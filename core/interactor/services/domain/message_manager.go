// core/interactor/services/domain/message_manager.go
package domain

import (
	"context"
	"errors"
	"sync"
)

type MessageManager struct {
	messages map[string]string // codigo -> contenido
	mu       sync.RWMutex
	repo     MessageRepository
}

type MessageRepository interface {
	GetAllMessages(ctx context.Context) ([]Message, error)
}

type Message struct {
	Codigo    string
	Contenido string
}

func NewMessageManager(repo MessageRepository) *MessageManager {
	return &MessageManager{
		messages: make(map[string]string),
		repo:     repo,
	}
}

// LoadMessages carga todos los mensajes desde la BD
func (m *MessageManager) LoadMessages(ctx context.Context) error {
	messages, err := m.repo.GetAllMessages(ctx)
	if err != nil {
		return err
	}

	m.mu.Lock()
	defer m.mu.Unlock()

	// Limpiar mapa actual
	m.messages = make(map[string]string)

	// Cargar nuevos mensajes
	for _, msg := range messages {
		m.messages[msg.Codigo] = msg.Contenido
	}

	return nil
}

// GetMessage obtiene un mensaje por código
func (m *MessageManager) GetMessage(codigo string) string {
	m.mu.RLock()
	defer m.mu.RUnlock()

	if msg, ok := m.messages[codigo]; ok {
		return msg
	}
	return codigo
}

// GetError retorna un error con el mensaje correspondiente
func (m *MessageManager) GetError(codigo string) error {
	return errors.New(m.GetMessage(codigo))
}

// Reload recarga los mensajes (para usar sin reiniciar)
func (m *MessageManager) Reload(ctx context.Context) error {
	return m.LoadMessages(ctx)
}
