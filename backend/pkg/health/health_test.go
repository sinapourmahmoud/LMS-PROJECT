package health

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
)

type fakePinger struct {
	err error
}

func (f fakePinger) PingContext(ctx context.Context) error { return f.err }

func TestHandler(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name       string
		pingErr    error
		wantStatus int
		wantBody   string
	}{
		{name: "database is up", pingErr: nil, wantStatus: http.StatusOK, wantBody: `{"status":"ok"}`},
		{name: "database is down", pingErr: errors.New("connection refused"), wantStatus: http.StatusServiceUnavailable, wantBody: `{"status":"unavailable"}`},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			router := gin.New()
			router.GET("/health", Handler(fakePinger{err: tt.pingErr}))

			rec := httptest.NewRecorder()
			req := httptest.NewRequest(http.MethodGet, "/health", nil)
			router.ServeHTTP(rec, req)

			if rec.Code != tt.wantStatus {
				t.Errorf("status = %d, want %d", rec.Code, tt.wantStatus)
			}
			if rec.Body.String() != tt.wantBody {
				t.Errorf("body = %s, want %s", rec.Body.String(), tt.wantBody)
			}
		})
	}
}
