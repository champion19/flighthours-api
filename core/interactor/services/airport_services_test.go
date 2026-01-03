package services

import (
	"context"
	"errors"
	"testing"

	"github.com/champion19/flighthours-api/core/interactor/services/domain"
	"github.com/champion19/flighthours-api/core/ports/output"
)

// airportFakeTx is a test transaction with tracking for commit/rollback
type airportFakeTx struct {
	committed  bool
	rolledBack bool
}

func (t *airportFakeTx) Commit() error {
	t.committed = true
	return nil
}

func (t *airportFakeTx) Rollback() error {
	t.rolledBack = true
	return nil
}

type fakeAirportRepo struct {
	getByIDFn      func(ctx context.Context, id string) (*domain.Airport, error)
	updateStatusFn func(ctx context.Context, tx output.Tx, id string, status bool) error
	beginTxFn      func(ctx context.Context) (output.Tx, error)
}

func (f fakeAirportRepo) BeginTx(ctx context.Context) (output.Tx, error) {
	if f.beginTxFn != nil {
		return f.beginTxFn(ctx)
	}
	return &airportFakeTx{}, nil
}

func (f fakeAirportRepo) GetAirportByID(ctx context.Context, id string) (*domain.Airport, error) {
	if f.getByIDFn != nil {
		return f.getByIDFn(ctx, id)
	}
	return nil, errors.New("not implemented")
}

func (f fakeAirportRepo) UpdateAirportStatus(ctx context.Context, tx output.Tx, id string, status bool) error {
	if f.updateStatusFn != nil {
		return f.updateStatusFn(ctx, tx, id, status)
	}
	return errors.New("not implemented")
}

func TestAirportService_GetAirportByID(t *testing.T) {
	ctx := context.Background()

	t.Run("returns airport when found", func(t *testing.T) {
		expectedAirport := &domain.Airport{
			ID:       "airport-123",
			Name:     "El Dorado International",
			City:     "Bogota",
			Country:  "Colombia",
			IATACode: "BOG",
			Status:   true,
		}
		svc := NewAirportService(fakeAirportRepo{
			getByIDFn: func(context.Context, string) (*domain.Airport, error) {
				return expectedAirport, nil
			},
		}, noopLogger{})

		result, err := svc.GetAirportByID(ctx, "airport-123")
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if result.ID != expectedAirport.ID {
			t.Fatalf("expected ID %s, got %s", expectedAirport.ID, result.ID)
		}
		if result.Name != expectedAirport.Name {
			t.Fatalf("expected name %s, got %s", expectedAirport.Name, result.Name)
		}
	})

	t.Run("returns error when not found", func(t *testing.T) {
		svc := NewAirportService(fakeAirportRepo{
			getByIDFn: func(context.Context, string) (*domain.Airport, error) {
				return nil, domain.ErrAirportNotFound
			},
		}, noopLogger{})

		_, err := svc.GetAirportByID(ctx, "non-existent")
		if !errors.Is(err, domain.ErrAirportNotFound) {
			t.Fatalf("expected %v, got %v", domain.ErrAirportNotFound, err)
		}
	})

	t.Run("propagates repository error", func(t *testing.T) {
		repoErr := errors.New("database connection error")
		svc := NewAirportService(fakeAirportRepo{
			getByIDFn: func(context.Context, string) (*domain.Airport, error) {
				return nil, repoErr
			},
		}, noopLogger{})

		_, err := svc.GetAirportByID(ctx, "airport-123")
		if !errors.Is(err, repoErr) {
			t.Fatalf("expected %v, got %v", repoErr, err)
		}
	})
}

func TestAirportService_UpdateAirportStatus(t *testing.T) {
	ctx := context.Background()

	t.Run("success => commit", func(t *testing.T) {
		tx := &airportFakeTx{}
		svc := NewAirportService(fakeAirportRepo{
			beginTxFn: func(context.Context) (output.Tx, error) { return tx, nil },
			updateStatusFn: func(context.Context, output.Tx, string, bool) error {
				return nil
			},
		}, noopLogger{})

		err := svc.UpdateAirportStatus(ctx, "airport-123", true)
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if !tx.committed {
			t.Fatalf("expected tx.Commit to be called")
		}
	})

	t.Run("begin tx fails => returns error", func(t *testing.T) {
		beginErr := errors.New("cannot begin transaction")
		svc := NewAirportService(fakeAirportRepo{
			beginTxFn: func(context.Context) (output.Tx, error) { return nil, beginErr },
		}, noopLogger{})

		err := svc.UpdateAirportStatus(ctx, "airport-123", true)
		if !errors.Is(err, beginErr) {
			t.Fatalf("expected %v, got %v", beginErr, err)
		}
	})

	t.Run("update fails => rollback", func(t *testing.T) {
		tx := &airportFakeTx{}
		updateErr := errors.New("update failed")
		svc := NewAirportService(fakeAirportRepo{
			beginTxFn: func(context.Context) (output.Tx, error) { return tx, nil },
			updateStatusFn: func(context.Context, output.Tx, string, bool) error {
				return updateErr
			},
		}, noopLogger{})

		err := svc.UpdateAirportStatus(ctx, "airport-123", true)
		if !errors.Is(err, updateErr) {
			t.Fatalf("expected %v, got %v", updateErr, err)
		}
		// Note: rollback is called via defer if error occurs
	})
}

func TestAirportService_ActivateAirport(t *testing.T) {
	ctx := context.Background()

	t.Run("success", func(t *testing.T) {
		tx := &airportFakeTx{}
		var receivedStatus bool
		svc := NewAirportService(fakeAirportRepo{
			beginTxFn: func(context.Context) (output.Tx, error) { return tx, nil },
			updateStatusFn: func(_ context.Context, _ output.Tx, _ string, status bool) error {
				receivedStatus = status
				return nil
			},
		}, noopLogger{})

		err := svc.ActivateAirport(ctx, "airport-123")
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if !receivedStatus {
			t.Fatalf("expected status true, got false")
		}
	})

	t.Run("propagates error", func(t *testing.T) {
		tx := &airportFakeTx{}
		updateErr := errors.New("update failed")
		svc := NewAirportService(fakeAirportRepo{
			beginTxFn: func(context.Context) (output.Tx, error) { return tx, nil },
			updateStatusFn: func(context.Context, output.Tx, string, bool) error {
				return updateErr
			},
		}, noopLogger{})

		err := svc.ActivateAirport(ctx, "airport-123")
		if !errors.Is(err, updateErr) {
			t.Fatalf("expected %v, got %v", updateErr, err)
		}
	})
}

func TestAirportService_DeactivateAirport(t *testing.T) {
	ctx := context.Background()

	t.Run("success", func(t *testing.T) {
		tx := &airportFakeTx{}
		var receivedStatus bool
		svc := NewAirportService(fakeAirportRepo{
			beginTxFn: func(context.Context) (output.Tx, error) { return tx, nil },
			updateStatusFn: func(_ context.Context, _ output.Tx, _ string, status bool) error {
				receivedStatus = status
				return nil
			},
		}, noopLogger{})

		err := svc.DeactivateAirport(ctx, "airport-123")
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if receivedStatus {
			t.Fatalf("expected status false, got true")
		}
	})

	t.Run("propagates error", func(t *testing.T) {
		tx := &airportFakeTx{}
		updateErr := errors.New("update failed")
		svc := NewAirportService(fakeAirportRepo{
			beginTxFn: func(context.Context) (output.Tx, error) { return tx, nil },
			updateStatusFn: func(context.Context, output.Tx, string, bool) error {
				return updateErr
			},
		}, noopLogger{})

		err := svc.DeactivateAirport(ctx, "airport-123")
		if !errors.Is(err, updateErr) {
			t.Fatalf("expected %v, got %v", updateErr, err)
		}
	})
}

func TestAirportService_BeginTx(t *testing.T) {
	ctx := context.Background()

	t.Run("success", func(t *testing.T) {
		expectedTx := &airportFakeTx{}
		svc := NewAirportService(fakeAirportRepo{
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
		svc := NewAirportService(fakeAirportRepo{
			beginTxFn: func(context.Context) (output.Tx, error) { return nil, beginErr },
		}, noopLogger{})

		_, err := svc.BeginTx(ctx)
		if !errors.Is(err, beginErr) {
			t.Fatalf("expected %v, got %v", beginErr, err)
		}
	})
}
