package auth_jwt

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/anyshake/observer/internal/testsupport"
	"github.com/anyshake/observer/pkg/timesource"
	"github.com/gin-gonic/gin"
)

func TestJWTAndWebsocketAuth(t *testing.T) {
	gin.SetMode(gin.TestMode)
	_, handler := testsupport.OpenDAO(t)
	adminID, err := handler.SysUserCreate("admin_user", "Anyshake@12#$", true)
	if err != nil {
		t.Fatal(err)
	}
	userID, err := handler.SysUserCreate("plain_user", "Anyshake@12#$", false)
	if err != nil {
		t.Fatal(err)
	}
	middleware, err := New(timesource.New(time.Now), handler, time.Hour)
	if err != nil {
		t.Fatal(err)
	}

	router := gin.New()
	router.Use(NewWebsocketAuthAdapter())
	router.POST("/login", func(ctx *gin.Context) {
		ctx.Set(UserIdKey, ctx.Query("user"))
		middleware.LoginHandler(ctx)
	})
	router.POST("/login-bad-type", func(ctx *gin.Context) {
		ctx.Set(UserIdKey, 1)
		middleware.LoginHandler(ctx)
	})
	router.GET("/me", middleware.MiddlewareFunc(), func(ctx *gin.Context) {
		ctx.Status(http.StatusNoContent)
	})
	router.GET("/refresh", middleware.RefreshHandler)

	adminToken := login(t, router, adminID)
	userToken := login(t, router, userID)
	for i := 0; i < 2; i++ {
		if status := authorizedStatus(router, "/me", adminToken); status != http.StatusNoContent {
			t.Fatalf("admin request %d status = %d", i, status)
		}
	}
	if status := authorizedStatus(router, "/me", userToken); status != http.StatusNoContent {
		t.Fatalf("user status = %d", status)
	}
	if status := authorizedStatus(router, "/me", ""); status != http.StatusUnauthorized {
		t.Fatalf("anonymous status = %d", status)
	}
	if status := authorizedStatus(router, "/me", "not-a-token"); status != http.StatusUnauthorized {
		t.Fatalf("bad token status = %d", status)
	}
	missing := login(t, router, "missing-user")
	if status := authorizedStatus(router, "/me", missing); status != http.StatusUnauthorized {
		t.Fatalf("missing user status = %d", status)
	}
	if status := authorizedStatus(router, "/refresh", adminToken); status != http.StatusOK {
		t.Fatalf("refresh status = %d", status)
	}

	recorder := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodGet, "/me", nil)
	request.Header.Set("Upgrade", "websocket")
	request.Header.Set("Connection", "keep-alive, Upgrade")
	request.Header.Set("Sec-WebSocket-Protocol", "json, Bearer#"+adminToken)
	router.ServeHTTP(recorder, request)
	if recorder.Code != http.StatusNoContent {
		t.Fatalf("websocket protocol status = %d", recorder.Code)
	}

	plain := httptest.NewRecorder()
	router.ServeHTTP(plain, httptest.NewRequest(http.MethodGet, "/me", nil))
	if plain.Code != http.StatusUnauthorized {
		t.Fatalf("plain websocket upgrade status = %d", plain.Code)
	}

	badType := httptest.NewRecorder()
	router.ServeHTTP(badType, httptest.NewRequest(http.MethodPost, "/login-bad-type", nil))
	if badType.Code != http.StatusUnauthorized {
		t.Fatalf("bad identity type status = %d body = %s", badType.Code, badType.Body.String())
	}
}

func login(t *testing.T, router *gin.Engine, userID string) string {
	t.Helper()
	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, httptest.NewRequest(http.MethodPost, "/login?user="+userID, nil))
	if recorder.Code != http.StatusOK {
		t.Fatalf("login %s status = %d body = %s", userID, recorder.Code, recorder.Body.String())
	}
	var body struct {
		Data struct {
			Token string `json:"token"`
		} `json:"data"`
	}
	if err := json.Unmarshal(recorder.Body.Bytes(), &body); err != nil || body.Data.Token == "" {
		t.Fatalf("login token: %v body %s", err, recorder.Body.String())
	}
	return body.Data.Token
}

func authorizedStatus(router *gin.Engine, path, token string) int {
	recorder := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodGet, path, nil)
	if token != "" {
		request.Header.Set("Authorization", "Bearer "+token)
	}
	router.ServeHTTP(recorder, request)
	return recorder.Code
}
