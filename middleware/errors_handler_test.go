package middleware

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/champion19/flighthours-api/core/interactor/services/domain"
	"github.com/champion19/flighthours-api/platform/cache/messaging"
	cachetypes "github.com/champion19/flighthours-api/platform/cache/types"
	"github.com/champion19/flighthours-api/platform/logger"
	"github.com/gin-gonic/gin"
)

type noopLogger struct{}

func (noopLogger) Info(string, ...any)    {}
func (noopLogger) Error(string, ...any)   {}
func (noopLogger) Debug(string, ...any)   {}
func (noopLogger) Warn(string, ...any)    {}
func (noopLogger) Success(string, ...any) {}
func (noopLogger) Fatal(string, ...any)   {}
func (noopLogger) Panic(string, ...any)   {}
func (noopLogger) WithTraceID(string) logger.Logger {
	return noopLogger{}
}

type fakeMessageCacheRepo struct {
	messages []cachetypes.CachedMessage
}

func (r fakeMessageCacheRepo) GetAllActiveForCache(context.Context) ([]cachetypes.CachedMessage, error) {
	return r.messages, nil
}
func (r fakeMessageCacheRepo) GetByCodeForCache(context.Context, string) (*cachetypes.CachedMessage, error) {
	return nil, nil
}
func (r fakeMessageCacheRepo) GetByCodeWithStatusForCache(context.Context, string) (*cachetypes.CachedMessage, error) {
	return nil, nil
}

func newTestCache(t *testing.T) *messaging.MessageCache {
	t.Helper()

	repo := fakeMessageCacheRepo{messages: []cachetypes.CachedMessage{
		{Code: domain.MsgUserDuplicate, Type: cachetypes.TypeError, Content: "duplicate"},
		{Code: domain.MsgValFieldRequired, Type: cachetypes.TypeError, Content: "missing ${0}"},
	}}
	c := messaging.NewMessageCache(repo, 0)
	if err := c.LoadMessages(context.Background()); err != nil {
		t.Fatalf("load cache: %v", err)
	}
	return c
}

func TestErrorHandler_MappedError(t *testing.T) {
	gin.SetMode(gin.TestMode)
	cache := newTestCache(t)
	errHandler := NewErrorHandler(cache)

	r := gin.New()
	r.Use(errHandler.Handle())
	r.GET("/", func(c *gin.Context) {
		c.Error(domain.ErrDuplicateUser)
	})

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	r.ServeHTTP(w, req)

	if w.Code != http.StatusConflict {
		t.Fatalf("expected %d got %d body=%s", http.StatusConflict, w.Code, w.Body.String())
	}

	var resp ErrorResponse
	if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
		t.Fatalf("invalid json: %v", err)
	}
	if resp.Success {
		t.Fatalf("expected success=false")
	}
	if resp.Code != domain.MsgUserDuplicate {
		t.Fatalf("expected code %q got %q", domain.MsgUserDuplicate, resp.Code)
	}
}

func TestErrorHandler_ValidationFieldsParam(t *testing.T) {
	gin.SetMode(gin.TestMode)
	cache := newTestCache(t)
	errHandler := NewErrorHandler(cache)

	r := gin.New()
	r.Use(errHandler.Handle())
	r.GET("/", func(c *gin.Context) {
		c.Set("validation_fields", []string{"a", "b"})
		c.Error(domain.ErrSchemaFieldRequired)
	})

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	r.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected %d got %d body=%s", http.StatusBadRequest, w.Code, w.Body.String())
	}

	var resp ErrorResponse
	if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
		t.Fatalf("invalid json: %v", err)
	}
	if resp.Message != "missing a, b" {
		t.Fatalf("expected substituted message, got %q", resp.Message)
	}
}
