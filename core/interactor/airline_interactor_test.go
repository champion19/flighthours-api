package interactor

import (
	"context"
	"errors"
	"testing"

	"github.com/champion19/flighthours-api/core/interactor/services/domain"
	"github.com/champion19/flighthours-api/core/ports/input"
	"github.com/champion19/flighthours-api/core/ports/output"
)

type fakeAirlineService struct {
	getByIDFn      func(ctx context.Context, id string) (*domain.Airline, error)
	updateStatusFn func(ctx context.Context, id string, status bool) error
	activateFn     func(ctx context.Context, id string) error
	deactivateFn   func(ctx context.Context, id string) error
	beginTxFn      func(ctx context.Context) (output.Tx, error)
	listAirlinesFn func(ctx context.Context, filters map[string]interface{}) ([]domain.Airline, error)
}

var _ input.AirlineService = (*fakeAirlineService)(nil)

func (f *fakeAirlineService) BeginTx(ctx context.Context) (output.Tx, error) {
	if f.beginTxFn != nil {
		return f.beginTxFn(ctx)
	}
	return &fakeTx{}, nil
}

func (f *fakeAirlineService) GetAirlineByID(ctx context.Context, id string) (*domain.Airline, error) {
	if f.getByIDFn != nil {
		return f.getByIDFn(ctx, id)
	}
	return nil, errors.New("not implemented")
}

func (f *fakeAirlineService) UpdateAirlineStatus(ctx context.Context, id string, status bool) error {
	if f.updateStatusFn != nil {
		return f.updateStatusFn(ctx, id, status)
	}
	return errors.New("not implemented")
}

func (f *fakeAirlineService) ActivateAirline(ctx context.Context, id string) error {
	if f.activateFn != nil {
		return f.activateFn(ctx, id)
	}
	return errors.New("not implemented")
}

func (f *fakeAirlineService) DeactivateAirline(ctx context.Context, id string) error {
	if f.deactivateFn != nil {
		return f.deactivateFn(ctx, id)
	}
	return errors.New("not implemented")
}

func (f *fakeAirlineService) ListAirlines(ctx context.Context, filters map[string]interface{}) ([]domain.Airline, error) {
	if f.listAirlinesFn != nil {
		return f.listAirlinesFn(ctx, filters)
	}
	return nil, errors.New("not implemented")
}

func TestAirlineInteractor_GetAirlineByID(t *testing.T) {
	ctx := context.Background()

	t.Run("success => returns airline", func(t *testing.T) {
		expectedAirline := &domain.Airline{
			ID:          "airline-123",
			AirlineName: "Test Airlines",
			AirlineCode: "TST",
			Status:      "active",
		}
		svc := &fakeAirlineService{
			getByIDFn: func(context.Context, string) (*domain.Airline, error) {
				return expectedAirline, nil
			},
		}
		interactor := NewAirlineInteractor(svc, noopLogger{})

		result, err := interactor.GetAirlineByID(ctx, "airline-123")
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if result == nil {
			t.Fatal("expected airline, got nil")
		}
		if result.ID != expectedAirline.ID {
			t.Errorf("expected ID %s, got %s", expectedAirline.ID, result.ID)
		}
		if result.AirlineName != expectedAirline.AirlineName {
			t.Errorf("expected name %s, got %s", expectedAirline.AirlineName, result.AirlineName)
		}
	})

	t.Run("not found => returns ErrAirlineNotFound", func(t *testing.T) {
		svc := &fakeAirlineService{
			getByIDFn: func(context.Context, string) (*domain.Airline, error) {
				return nil, domain.ErrAirlineNotFound
			},
		}
		interactor := NewAirlineInteractor(svc, noopLogger{})

		_, err := interactor.GetAirlineByID(ctx, "non-existent")
		if !errors.Is(err, domain.ErrAirlineNotFound) {
			t.Fatalf("expected %v, got %v", domain.ErrAirlineNotFound, err)
		}
	})

	t.Run("service error => propagate error", func(t *testing.T) {
		serviceErr := errors.New("service unavailable")
		svc := &fakeAirlineService{
			getByIDFn: func(context.Context, string) (*domain.Airline, error) {
				return nil, serviceErr
			},
		}
		interactor := NewAirlineInteractor(svc, noopLogger{})

		_, err := interactor.GetAirlineByID(ctx, "airline-123")
		if !errors.Is(err, serviceErr) {
			t.Fatalf("expected %v, got %v", serviceErr, err)
		}
	})
}

func TestAirlineInteractor_ActivateAirline(t *testing.T) {
	ctx := context.Background()

	t.Run("success", func(t *testing.T) {
		activateCalled := false
		svc := &fakeAirlineService{
			getByIDFn: func(context.Context, string) (*domain.Airline, error) {
				return &domain.Airline{ID: "airline-123", Status: "inactive"}, nil
			},
			activateFn: func(context.Context, string) error {
				activateCalled = true
				return nil
			},
		}
		interactor := NewAirlineInteractor(svc, noopLogger{})

		err := interactor.ActivateAirline(ctx, "airline-123")
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if !activateCalled {
			t.Fatal("expected ActivateAirline to be called")
		}
	})

	t.Run("airline not found => returns error", func(t *testing.T) {
		svc := &fakeAirlineService{
			getByIDFn: func(context.Context, string) (*domain.Airline, error) {
				return nil, domain.ErrAirlineNotFound
			},
		}
		interactor := NewAirlineInteractor(svc, noopLogger{})

		err := interactor.ActivateAirline(ctx, "non-existent")
		if !errors.Is(err, domain.ErrAirlineNotFound) {
			t.Fatalf("expected %v, got %v", domain.ErrAirlineNotFound, err)
		}
	})

	t.Run("activate fails => returns error", func(t *testing.T) {
		activateErr := errors.New("failed to activate")
		svc := &fakeAirlineService{
			getByIDFn: func(context.Context, string) (*domain.Airline, error) {
				return &domain.Airline{ID: "airline-123"}, nil
			},
			activateFn: func(context.Context, string) error {
				return activateErr
			},
		}
		interactor := NewAirlineInteractor(svc, noopLogger{})

		err := interactor.ActivateAirline(ctx, "airline-123")
		if !errors.Is(err, activateErr) {
			t.Fatalf("expected %v, got %v", activateErr, err)
		}
	})
}

func TestAirlineInteractor_DeactivateAirline(t *testing.T) {
	ctx := context.Background()

	t.Run("success", func(t *testing.T) {
		deactivateCalled := false
		svc := &fakeAirlineService{
			getByIDFn: func(context.Context, string) (*domain.Airline, error) {
				return &domain.Airline{ID: "airline-123", Status: "active"}, nil
			},
			deactivateFn: func(context.Context, string) error {
				deactivateCalled = true
				return nil
			},
		}
		interactor := NewAirlineInteractor(svc, noopLogger{})

		err := interactor.DeactivateAirline(ctx, "airline-123")
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if !deactivateCalled {
			t.Fatal("expected DeactivateAirline to be called")
		}
	})

	t.Run("airline not found => returns error", func(t *testing.T) {
		svc := &fakeAirlineService{
			getByIDFn: func(context.Context, string) (*domain.Airline, error) {
				return nil, domain.ErrAirlineNotFound
			},
		}
		interactor := NewAirlineInteractor(svc, noopLogger{})

		err := interactor.DeactivateAirline(ctx, "non-existent")
		if !errors.Is(err, domain.ErrAirlineNotFound) {
			t.Fatalf("expected %v, got %v", domain.ErrAirlineNotFound, err)
		}
	})

	t.Run("deactivate fails => returns error", func(t *testing.T) {
		deactivateErr := errors.New("failed to deactivate")
		svc := &fakeAirlineService{
			getByIDFn: func(context.Context, string) (*domain.Airline, error) {
				return &domain.Airline{ID: "airline-123"}, nil
			},
			deactivateFn: func(context.Context, string) error {
				return deactivateErr
			},
		}
		interactor := NewAirlineInteractor(svc, noopLogger{})

		err := interactor.DeactivateAirline(ctx, "airline-123")
		if !errors.Is(err, deactivateErr) {
			t.Fatalf("expected %v, got %v", deactivateErr, err)
		}
	})
}
