package main

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
	"github.com/myselfBZ/sat-jade/internal/auth"
	"github.com/myselfBZ/sat-jade/internal/store"
	"go.uber.org/zap"
)

var adminUser = store.User{
	ID: "3812cd57-1148-4876-830f-f36b53eb9805",
	Email: "admin@admin.com",
	Role: store.ROLE_ADMIN,
}



func newTestJWTToken(t *testing.T, a *api, userId string) string {
	claims := jwt.MapClaims{
		"sub": userId,
		"exp": time.Now().Add(a.config.auth.exp).Unix(),
		"iat": time.Now().Unix(),
		"nbf": time.Now().Unix(),
		"iss": a.config.auth.aud,
		"aud": a.config.auth.aud,
	}
	token, err := a.auth.GenerateToken(claims)

	if err != nil {
		t.Fatalf("couldn't generate a jwt token: %v", err)
	}

	return token
}

func newTestApi() *api {
	api := api{
		config: config{
			addr: ":8080",
			auth: authConfig{
				secret: "test-secret",
				exp:    time.Hour * time.Duration(1),
				aud:    "test-aud",
			},
		},
	}

	logger := zap.Must(zap.NewProduction(zap.AddCaller())).Sugar()
	defer logger.Sync()
	api.logger = logger

	api.storage = store.NewMockStorage()
	api.auth = auth.NewJWTAuthenticator(
		api.config.auth.secret,
		api.config.auth.aud,
		api.config.auth.aud,
	)

	api.storage.Users.Create(context.TODO(), &adminUser)

	return &api
}

func executeRequest(req *http.Request, e *echo.Echo) *httptest.ResponseRecorder {
	rr := httptest.NewRecorder()
	e.ServeHTTP(rr, req)

	return rr
}




func newTestRequest(t *testing.T ,method string, path string, body interface{}) *http.Request {
	var byteData []byte
	var err error
	if body != nil {
		byteData, err = json.Marshal(body)
	}

	if err != nil {
		t.Fatal("json.Marshal() in newTestRequest:", err)
	}
	buff := bytes.NewBuffer(byteData)
	req, err := http.NewRequest(method, path, buff)
	if err != nil {
		t.Fatal("newTestRequest couldn't create a request:", err)
	}
	return req
}


func checkResponseCode(t *testing.T, expected, actual int) {
	if expected != actual {
		t.Errorf("Expected response code %d. Got %d", expected, actual)
	}
}
