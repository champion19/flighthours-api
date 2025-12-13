package services

import (
	"context"
	"errors"
	"testing"

	"github.com/champion19/flighthours-api/core/interactor/services/domain"
	"github.com/champion19/flighthours-api/core/ports/output"
)

type fakeMsgRepo struct {
	getByCodeFn   func(ctx context.Context, code string) (*domain.Message, error)
	getByIDFn     func(ctx context.Context, id string) (*domain.Message, error)
	getAllActive  func(ctx context.Context) ([]domain.Message, error)
	saveMsgFn     func(ctx context.Context, tx output.Tx, message domain.Message) error
	updateMsgFn   func(ctx context.Context, tx output.Tx, message domain.Message) error
	deleteMsgFn   func(ctx context.Context, tx output.Tx, id string) error
	beginTxFn     func(ctx context.Context) (output.Tx, error)
	getByModuleFn func(ctx context.Context, module string) ([]domain.Message, error)
	getByTypeFn   func(ctx context.Context, msgType string) ([]domain.Message, error)
}

func (f fakeMsgRepo) BeginTx(ctx context.Context) (output.Tx, error) {
	if f.beginTxFn != nil {
		return f.beginTxFn(ctx)
	}
	return nil, errors.New("not implemented")
}
func (f fakeMsgRepo) SaveMessage(ctx context.Context, tx output.Tx, message domain.Message) error {
	if f.saveMsgFn != nil {
		return f.saveMsgFn(ctx, tx, message)
	}
	return errors.New("not implemented")
}
func (f fakeMsgRepo) UpdateMessage(ctx context.Context, tx output.Tx, message domain.Message) error {
	if f.updateMsgFn != nil {
		return f.updateMsgFn(ctx, tx, message)
	}
	return errors.New("not implemented")
}
func (f fakeMsgRepo) DeleteMessage(ctx context.Context, tx output.Tx, id string) error {
	if f.deleteMsgFn != nil {
		return f.deleteMsgFn(ctx, tx, id)
	}
	return errors.New("not implemented")
}
func (f fakeMsgRepo) GetAllActive(ctx context.Context) ([]domain.Message, error) {
	if f.getAllActive != nil {
		return f.getAllActive(ctx)
	}
	return nil, errors.New("not implemented")
}
func (f fakeMsgRepo) GetByID(ctx context.Context, id string) (*domain.Message, error) {
	if f.getByIDFn != nil {
		return f.getByIDFn(ctx, id)
	}
	return nil, errors.New("not implemented")
}
func (f fakeMsgRepo) GetByCode(ctx context.Context, code string) (*domain.Message, error) {
	if f.getByCodeFn != nil {
		return f.getByCodeFn(ctx, code)
	}
	return nil, errors.New("not implemented")
}
func (f fakeMsgRepo) GetByType(ctx context.Context, msgType string) ([]domain.Message, error) {
	if f.getByTypeFn != nil {
		return f.getByTypeFn(ctx, msgType)
	}
	return nil, errors.New("not implemented")
}
func (f fakeMsgRepo) GetByModule(ctx context.Context, module string) ([]domain.Message, error) {
	if f.getByModuleFn != nil {
		return f.getByModuleFn(ctx, module)
	}
	return nil, errors.New("not implemented")
}

func TestMessageService_ValidateMessage(t *testing.T) {
	ctx := context.Background()

	t.Run("empty code => ErrMessageCodeRequired", func(t *testing.T) {
		svc := NewMessageService(fakeMsgRepo{}, noopLogger{})
		err := svc.ValidateMessage(ctx, domain.Message{Code: ""})
		if !errors.Is(err, domain.ErrMessageCodeRequired) {
			t.Fatalf("expected %v, got %v", domain.ErrMessageCodeRequired, err)
		}
	})

	t.Run("create and code exists => ErrMessageCodeDuplicate", func(t *testing.T) {
		svc := NewMessageService(fakeMsgRepo{getByCodeFn: func(context.Context, string) (*domain.Message, error) {
			return &domain.Message{ID: "1", Code: "C1"}, nil
		}}, noopLogger{})

		err := svc.ValidateMessage(ctx, domain.Message{ID: "", Code: "C1"})
		if !errors.Is(err, domain.ErrMessageCodeDuplicate) {
			t.Fatalf("expected %v, got %v", domain.ErrMessageCodeDuplicate, err)
		}
	})

	t.Run("update (ID not empty) => does not check duplicate", func(t *testing.T) {
		svc := NewMessageService(fakeMsgRepo{getByCodeFn: func(context.Context, string) (*domain.Message, error) {
			return &domain.Message{ID: "1", Code: "C1"}, nil
		}}, noopLogger{})

		err := svc.ValidateMessage(ctx, domain.Message{ID: "existing", Code: "C1"})
		if err != nil {
			t.Fatalf("expected nil, got %v", err)
		}
	})
}

func TestMessageService_GetMessageByID(t *testing.T) {
	ctx := context.Background()

	t.Run("repo returns nil,nil => ErrMessageNotFound", func(t *testing.T) {
		svc := NewMessageService(fakeMsgRepo{getByIDFn: func(context.Context, string) (*domain.Message, error) {
			return nil, nil
		}}, noopLogger{})

		_, err := svc.GetMessageByID(ctx, "x")
		if !errors.Is(err, domain.ErrMessageNotFound) {
			t.Fatalf("expected %v, got %v", domain.ErrMessageNotFound, err)
		}
	})

	t.Run("repo returns error => propagate", func(t *testing.T) {
		repoErr := errors.New("db error")
		svc := NewMessageService(fakeMsgRepo{getByIDFn: func(context.Context, string) (*domain.Message, error) {
			return nil, repoErr
		}}, noopLogger{})

		_, err := svc.GetMessageByID(ctx, "x")
		if !errors.Is(err, repoErr) {
			t.Fatalf("expected %v, got %v", repoErr, err)
		}
	})
}
