package main

import (
	"net/http"
	"testing"
)


func TesCreateUserCreateToken(t *testing.T) {
	app := newTestApi()


	e := app.registerRoutes()

	// Creating the user
	req := newTestRequest(t, http.MethodPost, "/v1/auth/users", userPayload{
		FullName: "Jade the penguin",
		Email: "jade@satjade.live",
		Password: "hack me please",
	})

	rr := executeRequest(req, e)

	if rr.Code != http.StatusCreated {
		t.Fatalf("expected 201 got: %d", rr.Code)
	}

	// Logging in
	req = newTestRequest(t, http.MethodPost, "/v1/auth/token", loginPayload{
		Email: "jade@satjade.live",
		Password: "hack me please",
	})

	rr = executeRequest(req, e)

	if rr.Code != http.StatusOK {
		t.Fatalf("expected 200 got: %d", rr.Code)
	}
}
