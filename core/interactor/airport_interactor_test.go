package interactor

import (
	"context"
	"errors"
	"testing"

	"github.com/champion19/flighthours-api/core/interactor/services/domain"
	"github.com/champion19/flighthours-api/core/ports/input"
	"github.com/champion19/flighthours-api/core/ports/output"
)

type fakeAirportService struct {
	getByIDFn      func(ctx context.Context, id string) (*domain.Airport, error)
	updateStatusFn func(ctx context.Context, id string, status bool) error
	activateFn     func(ctx context.Context, id string) error
	deactivateFn   func(ctx context.Context, id string) error
	beginTxFn      func(ctx context.Context) (output.Tx, error)
}

var _ input.AirportService = (*fakeAirportService)(nil)

func (f *fakeAirportService) BeginTx(ctx context.Context) (output.Tx, error) {
	if f.beginTxFn != nil {
		return f.beginTxFn(ctx)
	}
	return &fakeTx{}, nil
}

func (f *fakeAirportService) GetAirportByID(ctx context.Context, id string) (*domain.Airport, error) {
	if f.getByIDFn != nil {
		return f.getByIDFn(ctx, id)
	}
	return nil, errors.New("not implemented")
}

func (f *fakeAirportService) UpdateAirportStatus(ctx context.Context, id string, status bool) error {
	if f.updateStatusFn != nil {
		return f.updateStatusFn(ctx, id, status)
	}
	return errors.New("not implemented")
}

func (f *fakeAirportService) ActivateAirport(ctx context.Context, id string) error {
	if f.activateFn != nil {
		return f.activateFn(ctx, id)
	}
	return errors.New("not implemented")
}

func (f *fakeAirportService) DeactivateAirport(ctx context.Context, id string) error {
	if f.deactivateFn != nil {
		return f.deactivateFn(ctx, id)
	}
	return errors.New("not implemented")
}

func TestAirportInteractor_GetAirportByID(t *testing.T) {
	ctx := context.Background()

	t.Run("success => returns airport", func(t *testing.T) {
		expectedAirport := &domain.Airport{
			ID:       "airport-123",
			Name:     "El Dorado International",
			City:     "Bogota",
			Country:  "Colombia",
			IATACode: "BOG",
			Status:   true,
		}
		svc := &fakeAirportService{
			getByIDFn: func(context.Context, string) (*domain.Airport, error) {
				return expectedAirport, nil
			},
		}
		interactor := NewAirportInteractor(svc, noopLogger{})

		result, err := interactor.GetAirportByID(ctx, "airport-123")
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if result == nil {
			t.Fatal("expected airport, got nil")
		}
		if result.ID != expectedAirport.ID {
			t.Errorf("expected ID %s, got %s", expectedAirport.ID, result.ID)
		}
		if result.Name != expectedAirport.Name {
			t.Errorf("expected name %s, got %s", expectedAirport.Name, result.Name)
		}
	})

	t.Run("not found => returns ErrAirportNotFound", func(t *testing.T) {
		svc := &fakeAirportService{
			getByIDFn: func(context.Context, string) (*domain.Airport, error) {
				return nil, domain.ErrAirportNotFound
			},
		}
		interactor := NewAirportInteractor(svc, noopLogger{})

		_, err := interactor.GetAirportByID(ctx, "non-existent")
		if !errors.Is(err, domain.ErrAirportNotFound) {
			t.Fatalf("expected %v, got %v", domain.ErrAirportNotFound, err)
		}
	})

	t.Run("service error => propagate error", func(t *testing.T) {
		serviceErr := errors.New("service unavailable")
		svc := &fakeAirportService{
			getByIDFn: func(context.Context, string) (*domain.Airport, error) {
				return nil, serviceErr
			},
		}
		interactor := NewAirportInteractor(svc, noopLogger{})

		_, err := interactor.GetAirportByID(ctx, "airport-123")
		if !errors.Is(err, serviceErr) {
			t.Fatalf("expected %v, got %v", serviceErr, err)
		}
	})
}

func TestAirportInteractor_ActivateAirport(t *testing.T) {
	ctx := context.Background()

	t.Run("success", func(t *testing.T) {
		activateCalled := false
		svc := &fakeAirportService{
			getByIDFn: func(context.Context, string) (*domain.Airport, error) {
				return &domain.Airport{ID: "airport-123", Status: false}, nil
			},
			activateFn: func(context.Context, string) error {
				activateCalled = true
				return nil
			},
		}
		interactor := NewAirportInteractor(svc, noopLogger{})

		err := interactor.ActivateAirport(ctx, "airport-123")
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if !activateCalled {
			t.Fatal("expected ActivateAirport to be called")
		}
	})

	t.Run("airport not found => returns error", func(t *testing.T) {
		svc := &fakeAirportService{
			getByIDFn: func(context.Context, string) (*domain.Airport, error) {
				return nil, domain.ErrAirportNotFound
			},
		}
		interactor := NewAirportInteractor(svc, noopLogger{})

		err := interactor.ActivateAirport(ctx, "non-existent")
		if !errors.Is(err, domain.ErrAirportNotFound) {
			t.Fatalf("expected %v, got %v", domain.ErrAirportNotFound, err)
		}
	})

	t.Run("activate fails => returns error", func(t *testing.T) {
		activateErr := errors.New("failed to activate")
		svc := &fakeAirportService{
			getByIDFn: func(context.Context, string) (*domain.Airport, error) {
				return &domain.Airport{ID: "airport-123"}, nil
			},
			activateFn: func(context.Context, string) error {
				return activateErr
			},
		}
		interactor := NewAirportInteractor(svc, noopLogger{})

		err := interactor.ActivateAirport(ctx, "airport-123")
		if !errors.Is(err, activateErr) {
			t.Fatalf("expected %v, got %v", activateErr, err)
		}
	})
}

func TestAirportInteractor_DeactivateAirport(t *testing.T) {
	ctx := context.Background()

	t.Run("success", func(t *testing.T) {
		deactivateCalled := false
		svc := &fakeAirportService{
			getByIDFn: func(context.Context, string) (*domain.Airport, error) {
				return &domain.Airport{ID: "airport-123", Status: true}, nil
			},
			deactivateFn: func(context.Context, string) error {
				deactivateCalled = true
				return nil
			},
		}
		interactor := NewAirportInteractor(svc, noopLogger{})

		err := interactor.DeactivateAirport(ctx, "airport-123")
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if !deactivateCalled {
			t.Fatal("expected DeactivateAirport to be called")
		}
	})

	t.Run("airport not found => returns error", func(t *testing.T) {
		svc := &fakeAirportService{
			getByIDFn: func(context.Context, string) (*domain.Airport, error) {
				return nil, domain.ErrAirportNotFound
			},
		}
		interactor := NewAirportInteractor(svc, noopLogger{})

		err := interactor.DeactivateAirport(ctx, "non-existent")
		if !errors.Is(err, domain.ErrAirportNotFound) {
			t.Fatalf("expected %v, got %v", domain.ErrAirportNotFound, err)
		}
	})

	t.Run("deactivate fails => returns error", func(t *testing.T) {
		deactivateErr := errors.New("failed to deactivate")
		svc := &fakeAirportService{
			getByIDFn: func(context.Context, string) (*domain.Airport, error) {
				return &domain.Airport{ID: "airport-123"}, nil
			},
			deactivateFn: func(context.Context, string) error {
				return deactivateErr
			},
		}
		interactor := NewAirportInteractor(svc, noopLogger{})

		err := interactor.DeactivateAirport(ctx, "airport-123")
		if !errors.Is(err, deactivateErr) {
			t.Fatalf("expected %v, got %v", deactivateErr, err)
		}
	})
}
