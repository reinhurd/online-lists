package main

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-resty/resty/v2"
	"github.com/stretchr/testify/assert"
	"online-lists/internal/clients/yandex"
)

func TestPingRoute(t *testing.T) {
	//TODO DO SOME MOCK
	restyCl := resty.New()
	yaClient := yandex.NewClient(restyCl, "SOMETESTID")
	router := setupRouter(yaClient)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/ping", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "pong", w.Body.String())
}
