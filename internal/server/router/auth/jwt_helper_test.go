package auth

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/anyshake/observer/internal/dao/action"
	"github.com/anyshake/observer/internal/server/middleware/auth_jwt"
	"github.com/anyshake/observer/pkg/timesource"
	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/gin-gonic/gin"
)

func newTestJWT(handler *action.Handler) (*jwt.GinJWTMiddleware, error) {
	return auth_jwt.New(timesource.New(time.Now), handler, time.Hour)
}

func loginForToken(t *testing.T, middleware *jwt.GinJWTMiddleware, userID string) string {
	t.Helper()
	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.POST("/login", func(ctx *gin.Context) {
		ctx.Set(auth_jwt.UserIdKey, userID)
		middleware.LoginHandler(ctx)
	})
	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, httptest.NewRequest(http.MethodPost, "/login", nil))
	if recorder.Code != http.StatusOK {
		t.Fatalf("token login status = %d body = %s", recorder.Code, recorder.Body.String())
	}
	var body struct {
		Data struct {
			Token string `json:"token"`
		} `json:"data"`
	}
	if err := json.Unmarshal(recorder.Body.Bytes(), &body); err != nil || body.Data.Token == "" {
		t.Fatalf("token body = %s, %v", recorder.Body.String(), err)
	}
	return body.Data.Token
}
