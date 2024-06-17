package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockRateLimiter struct {
	mock.Mock
}

func (m *MockRateLimiter) Allow(ip string, token string) (bool, error) {
	args := m.Called(ip, token)
	return args.Bool(0), args.Error(1)
}

func (m *MockRateLimiter) BlockDuration() time.Duration {
	args := m.Called()
	return args.Get(0).(time.Duration)
}

func TestRateLimitMiddleware_Allow(t *testing.T) {
	mockRateLimiter := new(MockRateLimiter)

	ip := "192.168.1.1"
	token := "abc123"

	mockRateLimiter.On("Allow", ip, token).Return(true, nil)

	middleware := RateLimitMiddleware(mockRateLimiter)

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.RemoteAddr = ip + ":1234"
	req.Header.Set("API_KEY", token)
	w := httptest.NewRecorder()
	middleware(handler).ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)

	mockRateLimiter.AssertExpectations(t)
}

func TestRateLimitMiddleware_Block(t *testing.T) {
	mockRateLimiter := new(MockRateLimiter)

	ip := "192.168.1.1"
	token := "abc123"

	mockRateLimiter.On("Allow", ip, token).Return(false, nil)
	mockRateLimiter.On("BlockDuration").Return(60 * time.Second)

	middleware := RateLimitMiddleware(mockRateLimiter)

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.RemoteAddr = ip + ":1234"
	req.Header.Set("API_KEY", token)
	w := httptest.NewRecorder()
	middleware(handler).ServeHTTP(w, req)
	assert.Equal(t, http.StatusTooManyRequests, w.Code)
	assert.Equal(t, "1m0s", w.Header().Get("Retry-After"))

	mockRateLimiter.AssertExpectations(t)
}

func TestRateLimitMiddleware_InternalServerError(t *testing.T) {
	mockRateLimiter := new(MockRateLimiter)

	mockRateLimiter.On("Allow", mock.Anything, mock.Anything).Return(false, assert.AnError)

	middleware := RateLimitMiddleware(mockRateLimiter)

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.RemoteAddr = "192.168.1.1:1234"
	w := httptest.NewRecorder()
	middleware(handler).ServeHTTP(w, req)
	assert.Equal(t, http.StatusInternalServerError, w.Code)

	mockRateLimiter.AssertExpectations(t)
}
