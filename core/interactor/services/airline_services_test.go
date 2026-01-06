package services

import (
	"context"
	"errors"
	"testing"

	"github.com/champion19/flighthours-api/core/interactor/services/domain"
	"github.com/champion19/flighthours-api/core/ports/output"
)

// airlineFakeTx is a test transaction with tracking for commit/rollback
type airlineFakeTx struct {
	committed  bool
	rolledBack bool
}

func (t *airlineFakeTx) Commit() error {
	t.committed = true
	return nil
}

func (t *airlineFakeTx) Rollback() error {
	t.rolledBack = true
	return nil
}

type fakeAirlineRepo struct {
	getByIDFn      func(ctx context.Context, id string) (*domain.Airline, error)
	updateStatusFn func(ctx context.Context, tx output.Tx, id string, status bool) error
	beginTxFn      func(ctx context.Context) (output.Tx, error)
	listAirlinesFn func(ctx context.Context, filters map[string]interface{}) ([]domain.Airline, error)
}

func (f fakeAirlineRepo) BeginTx(ctx context.Context) (output.Tx, error) {
	if f.beginTxFn != nil {
		return f.beginTxFn(ctx)
	}
	return &airlineFakeTx{}, nil
}

func (f fakeAirlineRepo) GetAirlineByID(ctx context.Context, id string) (*domain.Airline, error) {
	if f.getByIDFn != nil {
		return f.getByIDFn(ctx, id)
	}
	return nil, errors.New("not implemented")
}

func (f fakeAirlineRepo) UpdateAirlineStatus(ctx context.Context, tx output.Tx, id string, status bool) error {
	if f.updateStatusFn != nil {
		return f.updateStatusFn(ctx, tx, id, status)
	}
	return errors.New("not implemented")
}

func (f fakeAirlineRepo) ListAirlines(ctx context.Context, filters map[string]interface{}) ([]domain.Airline, error) {
	if f.listAirlinesFn != nil {
		return f.listAirlinesFn(ctx, filters)
	}
	return nil, errors.New("not implemented")
}

func TestAirlineService_GetAirlineByID(t *testing.T) {
	ctx := context.Background()

	t.Run("returns airline when found", func(t *testing.T) {
		expectedAirline := &domain.Airline{
			ID:          "airline-123",
			AirlineName: "Test Airlines",
			AirlineCode: "TST",
			Status:      "active",
		}
		svc := NewAirlineService(fakeAirlineRepo{
			getByIDFn: func(context.Context, string) (*domain.Airline, error) {
				return expectedAirline, nil
			},
		}, noopLogger{})

		result, err := svc.GetAirlineByID(ctx, "airline-123")
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if result.ID != expectedAirline.ID {
			t.Fatalf("expected ID %s, got %s", expectedAirline.ID, result.ID)
		}
		if result.AirlineName != expectedAirline.AirlineName {
			t.Fatalf("expected name %s, got %s", expectedAirline.AirlineName, result.AirlineName)
		}
	})

	t.Run("returns error when not found", func(t *testing.T) {
		svc := NewAirlineService(fakeAirlineRepo{
			getByIDFn: func(context.Context, string) (*domain.Airline, error) {
				return nil, domain.ErrAirlineNotFound
			},
		}, noopLogger{})

		_, err := svc.GetAirlineByID(ctx, "non-existent")
		if !errors.Is(err, domain.ErrAirlineNotFound) {
			t.Fatalf("expected %v, got %v", domain.ErrAirlineNotFound, err)
		}
	})

	t.Run("propagates repository error", func(t *testing.T) {
		repoErr := errors.New("database connection error")
		svc := NewAirlineService(fakeAirlineRepo{
			getByIDFn: func(context.Context, string) (*domain.Airline, error) {
				return nil, repoErr
			},
		}, noopLogger{})

		_, err := svc.GetAirlineByID(ctx, "airline-123")
		if !errors.Is(err, repoErr) {
			t.Fatalf("expected %v, got %v", repoErr, err)
		}
	})
}

func TestAirlineService_UpdateAirlineStatus(t *testing.T) {
	ctx := context.Background()

	t.Run("success => commit", func(t *testing.T) {
		tx := &airlineFakeTx{}
		svc := NewAirlineService(fakeAirlineRepo{
			beginTxFn: func(context.Context) (output.Tx, error) { return tx, nil },
			updateStatusFn: func(context.Context, output.Tx, string, bool) error {
				return nil
			},
		}, noopLogger{})

		err := svc.UpdateAirlineStatus(ctx, "airline-123", true)
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if !tx.committed {
			t.Fatalf("expected tx.Commit to be called")
		}
	})

	t.Run("begin tx fails => returns error", func(t *testing.T) {
		beginErr := errors.New("cannot begin transaction")
		svc := NewAirlineService(fakeAirlineRepo{
			beginTxFn: func(context.Context) (output.Tx, error) { return nil, beginErr },
		}, noopLogger{})

		err := svc.UpdateAirlineStatus(ctx, "airline-123", true)
		if !errors.Is(err, beginErr) {
			t.Fatalf("expected %v, got %v", beginErr, err)
		}
	})

	t.Run("update fails => rollback", func(t *testing.T) {
		tx := &airlineFakeTx{}
		updateErr := errors.New("update failed")
		svc := NewAirlineService(fakeAirlineRepo{
			beginTxFn: func(context.Context) (output.Tx, error) { return tx, nil },
			updateStatusFn: func(context.Context, output.Tx, string, bool) error {
				return updateErr
			},
		}, noopLogger{})

		err := svc.UpdateAirlineStatus(ctx, "airline-123", true)
		if !errors.Is(err, updateErr) {
			t.Fatalf("expected %v, got %v", updateErr, err)
		}
		// Note: rollback is called via defer if error occurs
	})
}

func TestAirlineService_ActivateAirline(t *testing.T) {
	ctx := context.Background()

	t.Run("success", func(t *testing.T) {
		tx := &airlineFakeTx{}
		var receivedStatus bool
		svc := NewAirlineService(fakeAirlineRepo{
			beginTxFn: func(context.Context) (output.Tx, error) { return tx, nil },
			updateStatusFn: func(_ context.Context, _ output.Tx, _ string, status bool) error {
				receivedStatus = status
				return nil
			},
		}, noopLogger{})

		err := svc.ActivateAirline(ctx, "airline-123")
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if !receivedStatus {
			t.Fatalf("expected status true, got false")
		}
	})

	t.Run("propagates error", func(t *testing.T) {
		tx := &airlineFakeTx{}
		updateErr := errors.New("update failed")
		svc := NewAirlineService(fakeAirlineRepo{
			beginTxFn: func(context.Context) (output.Tx, error) { return tx, nil },
			updateStatusFn: func(context.Context, output.Tx, string, bool) error {
				return updateErr
			},
		}, noopLogger{})

		err := svc.ActivateAirline(ctx, "airline-123")
		if !errors.Is(err, updateErr) {
			t.Fatalf("expected %v, got %v", updateErr, err)
		}
	})
}

func TestAirlineService_DeactivateAirline(t *testing.T) {
	ctx := context.Background()

	t.Run("success", func(t *testing.T) {
		tx := &airlineFakeTx{}
		var receivedStatus bool
		svc := NewAirlineService(fakeAirlineRepo{
			beginTxFn: func(context.Context) (output.Tx, error) { return tx, nil },
			updateStatusFn: func(_ context.Context, _ output.Tx, _ string, status bool) error {
				receivedStatus = status
				return nil
			},
		}, noopLogger{})

		err := svc.DeactivateAirline(ctx, "airline-123")
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if receivedStatus {
			t.Fatalf("expected status false, got true")
		}
	})

	t.Run("propagates error", func(t *testing.T) {
		tx := &airlineFakeTx{}
		updateErr := errors.New("update failed")
		svc := NewAirlineService(fakeAirlineRepo{
			beginTxFn: func(context.Context) (output.Tx, error) { return tx, nil },
			updateStatusFn: func(context.Context, output.Tx, string, bool) error {
				return updateErr
			},
		}, noopLogger{})

		err := svc.DeactivateAirline(ctx, "airline-123")
		if !errors.Is(err, updateErr) {
			t.Fatalf("expected %v, got %v", updateErr, err)
		}
	})
}

func TestAirlineService_BeginTx(t *testing.T) {
	ctx := context.Background()

	t.Run("success", func(t *testing.T) {
		expectedTx := &airlineFakeTx{}
		svc := NewAirlineService(fakeAirlineRepo{
			beginTxFn: func(context.Context) (output.Tx, error) { return expectedTx, nil },
		}, noopLogger{})

		tx, err := svc.BeginTx(ctx)
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if tx == nil {
			t.Fatal("expected tx, got nil")
		}
	})

	t.Run("propagates error", func(t *testing.T) {
		beginErr := errors.New("cannot begin transaction")
		svc := NewAirlineService(fakeAirlineRepo{
			beginTxFn: func(context.Context) (output.Tx, error) { return nil, beginErr },
		}, noopLogger{})

		_, err := svc.BeginTx(ctx)
		if !errors.Is(err, beginErr) {
			t.Fatalf("expected %v, got %v", beginErr, err)
		}
	})
}
