package main

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/myselfBZ/sat-jade/internal/store"
	"github.com/stretchr/testify/assert"
)



func TestGetSelf(t *testing.T) {
	api := newTestApi()
	e := api.registerRoutes()
	targetUser := &store.User{
		ID: "user-123", 
		Email: "test@satjade.live", 
		FullName: "Test User",
	}

	t.Run("Success 200 fetch", func(t *testing.T) {
		req := newTestRequest(t, http.MethodGet, "/v1/users/self", nil)
		rr := httptest.NewRecorder()

		c := e.NewContext(req, rr)

		c.Set(userCtxKey, targetUser)

		if err := api.getUserSelfHandler(c); err != nil {
			t.Fatal("error getUserSelfHandler(): ", err)
		}

		user := &store.User{}

		if err := json.NewDecoder(rr.Body).Decode(&user); err != nil {
			t.Fatal("error decoding response: ", err)
		}

		assert.Equal(t, targetUser, user, "getUserSelfHandler did not return the same user")
	})

	t.Run("401 Unauthorized fetch", func(t *testing.T) {
		req := newTestRequest(t, http.MethodGet, "/v1/users/self", nil)
		rr := httptest.NewRecorder()

		c := e.NewContext(req, rr)

		err := api.getUserSelfHandler(c)
		assert.Error(t, err)
		he, ok := err.(*echo.HTTPError)
		assert.True(t, ok)
		assert.Equal(t, http.StatusUnauthorized, he.Code)
	})
}


func TestDeleteUser(t *testing.T) {
	api := newTestApi()
	e := api.registerRoutes()


	t.Run("Delete target user 204 no content", func(t *testing.T) {
		targetUser := &store.User{
			Email: "test@satjade.live", 
			FullName: "Test User",
		}

		api.storage.Users.Create(context.TODO(), targetUser)
		req := newTestRequest(t, http.MethodPost, "/v1/users/"+targetUser.ID, nil)		
		rr := httptest.NewRecorder()

		c := e.NewContext(req, rr)

		c.SetPath("/v1/users/:id")
		c.SetParamNames("id")
		c.SetParamValues(targetUser.ID)

		err := api.deleteUserHandler(c)

		assert.NoError(t, err, "unsuccessfull delete")
		assert.Equal(t, http.StatusNoContent, rr.Code)
	})



	t.Run("Status not found 404", func(t *testing.T) {
		req := newTestRequest(t, http.MethodPost, "/v1/users/non-existent-user", nil)		
		rr := httptest.NewRecorder()

		c := e.NewContext(req, rr)

		c.SetPath("/v1/users/:id")
		c.SetParamNames("id")
		c.SetParamValues("non-existent-user")

		err := api.deleteUserHandler(c)

		assert.Error(t, err)

		he, ok := err.(*echo.HTTPError)
		assert.True(t, ok)
		assert.Equal(t, http.StatusNotFound, he.Code)

	})



	// Runs with middleware that checks against admin role
	t.Run("Unauthorized delete 401", func(t *testing.T) {
		req := newTestRequest(t, http.MethodPost, "/v1/users/test-user-id", nil)		
		rr := executeRequest(req, e)
		assert.Equal(t, http.StatusUnauthorized, rr.Code)
	})


	t.Run("Authorized delete 204", func(t *testing.T) {
		targetUser := &store.User{
			Email: "test2@satjade.live", 
			FullName: "Test User",
		}

		api.storage.Users.Create(context.TODO(), targetUser)

		token := newTestJWTToken(t, api, adminUser.ID)
		req := newTestRequest(t, http.MethodDelete, "/v1/users/"+targetUser.ID, nil)
		req.Header.Add("Authorization", "Bearer "+token)

		rr := executeRequest(req, e)

		assert.Equal(t, http.StatusNoContent, rr.Code)
	})

}


func TestGetUsers(t *testing.T) {
	api := newTestApi()
	e := api.registerRoutes()

	existingUsers := []store.User{
		{
			ID:       "user-200",
			Email:    "jane.doe@example.com",
			FullName: "Jane Doe",
		},
		{
			ID:       "user-300",
			Email:    "bob.smith@example.com",
			FullName: "Bob Smith",
		},
		{
			ID:       "user-400",
			Email:    "alice.wonder@example.com",
			FullName: "Alice Wonder",
		},
	}


	for _, u := range existingUsers {
		api.storage.Users.Create(context.TODO(), &u)
	}

	t.Run("200 fetch all users", func(t *testing.T) {
		req := newTestRequest(t, http.MethodGet, "/v1/users/", nil)
		rr := httptest.NewRecorder()
		c := e.NewContext(req, rr)
		err := api.getUsersHandler(c)
		assert.NoError(t, err)
		var users []store.User
		err = json.NewDecoder(rr.Body).Decode(&users)
		assert.NoError(t, err)

		// taking the admin into consideration
		assert.Equal(t, len(existingUsers) + 1, len(users))
	})

	t.Run("200 fetch all users with middleware", func(t *testing.T) {
		token := newTestJWTToken(t, api, adminUser.ID)
		req := newTestRequest(t, http.MethodGet, "/v1/users/", nil)
		req.Header.Add("Authorization", "Bearer "+token)

		rr := executeRequest(req, e)

		var users []store.User
		err := json.NewDecoder(rr.Body).Decode(&users)
		assert.NoError(t, err)

		assert.Equal(t, http.StatusOK, rr.Code)

		// taking the admin into consideration
		assert.Equal(t, len(existingUsers) + 1, len(users))
	})
}





