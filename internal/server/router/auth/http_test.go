package auth

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/anyshake/observer/internal/testsupport"
	"github.com/anyshake/observer/pkg/logger"
	"github.com/gin-gonic/gin"
)

func TestAuthRoutes(t *testing.T) {
	gin.SetMode(gin.TestMode)
	logger.Init()
	_, handler := testsupport.OpenDAO(t)
	userID, err := handler.SysUserCreate("admin_user", "Anyshake@12#$", true)
	if err != nil {
		t.Fatal(err)
	}
	jwtMiddleware, err := newTestJWT(handler)
	if err != nil {
		t.Fatal(err)
	}

	router := gin.New()
	group := router.Group("/")
	Setup(group, handler, jwtMiddleware.MiddlewareFunc(), func(ctx *gin.Context) {
		ctx.Status(http.StatusCreated)
	}, func(ctx *gin.Context) {
		ctx.Status(http.StatusAccepted)
	})

	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, httptest.NewRequest(http.MethodGet, "/auth", nil))
	if recorder.Code != http.StatusUnauthorized {
		t.Fatalf("anonymous auth status = %d", recorder.Code)
	}

	token := loginForToken(t, jwtMiddleware, userID)
	recorder = httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodGet, "/auth", nil)
	request.Header.Set("Authorization", "Bearer "+token)
	router.ServeHTTP(recorder, request)
	if recorder.Code != http.StatusOK {
		t.Fatalf("valid auth status = %d body = %s", recorder.Code, recorder.Body.String())
	}

	recorder = postJSON(router, "/auth", map[string]any{})
	if recorder.Code != http.StatusBadRequest {
		t.Fatalf("empty auth status = %d", recorder.Code)
	}
	recorder = postJSON(router, "/auth", map[string]any{"action": "nope"})
	if recorder.Code != http.StatusBadRequest {
		t.Fatalf("invalid action status = %d", recorder.Code)
	}
	recorder = postJSON(router, "/auth", map[string]any{"action": "login"})
	if recorder.Code != http.StatusUnauthorized && recorder.Code != http.StatusBadRequest {
		t.Fatalf("incomplete login status = %d body = %s", recorder.Code, recorder.Body.String())
	}
	loginBody := map[string]any{
		"action": "login", "session": "session", "secret": "secret", "nonce": "!!!",
		"challenge_id": "missing", "challenge_solution": "1:aa", "captcha_id": "captcha",
		"captcha_val": "1234", "payload": "e30=",
	}
	if postJSON(router, "/auth", loginBody).Code != http.StatusBadRequest {
		t.Fatal("invalid nonce encoding was accepted")
	}
	loginBody["nonce"] = "YQ=="
	if postJSON(router, "/auth", loginBody).Code != http.StatusUnauthorized {
		t.Fatal("unknown challenge was accepted")
	}
	if postJSON(router, "/auth", loginBody).Code != http.StatusForbidden {
		t.Fatal("replayed nonce was accepted")
	}
	recorder = postJSON(router, "/auth", map[string]any{"action": "refresh"})
	if recorder.Code != http.StatusAccepted {
		t.Fatalf("refresh status = %d", recorder.Code)
	}
	recorder = postJSON(router, "/auth", map[string]any{"action": "preauth"})
	if recorder.Code != http.StatusOK {
		t.Fatalf("preauth status = %d body = %s", recorder.Code, recorder.Body.String())
	}
	var preauth struct {
		Data map[string]any `json:"data"`
	}
	if err := json.Unmarshal(recorder.Body.Bytes(), &preauth); err != nil || preauth.Data["public_key"] == "" {
		t.Fatalf("preauth body = %s, %v", recorder.Body.String(), err)
	}
	again := postJSON(router, "/auth", map[string]any{"action": "preauth"})
	if again.Code != http.StatusOK {
		t.Fatalf("second preauth status = %d", again.Code)
	}
	_ = time.Second
}

func postJSON(router *gin.Engine, path string, body any) *httptest.ResponseRecorder {
	payload, _ := json.Marshal(body)
	recorder := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodPost, path, bytes.NewReader(payload))
	request.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(recorder, request)
	return recorder
}
