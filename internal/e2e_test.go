package internal

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func setupTestServer() (*httptest.Server, *partRepository) {
	repo := NewPartRepository()
	service := NewPartService(repo)
	h := NewHandler(service)

	mux := http.NewServeMux()
	h.RegisterRoutes(mux)

	return httptest.NewServer(mux), repo
}

func testClient() *http.Client {
	return &http.Client{
		Timeout: 5 * time.Second,
	}
}

func TestAPI(t *testing.T) {
	server, repo := setupTestServer()
	defer server.Close()

	client := testClient()

	t.Run("GET /parts возвращает пустой список", func(t *testing.T) {
		// Создаём отдельный сервер для теста пустого списка
		emptyServer, _ := setupTestServer()
		defer emptyServer.Close()

		resp, err := client.Get(emptyServer.URL + "/parts")
		require.NoError(t, err)
		defer func() { _ = resp.Body.Close() }()

		require.Equal(t, http.StatusOK, resp.StatusCode)

		var parts []Part
		err = json.NewDecoder(resp.Body).Decode(&parts)
		require.NoError(t, err)
		require.Empty(t, parts)
	})

	t.Run("GET /parts возвращает список деталей", func(t *testing.T) {
		repo.Create(Part{Name: "Ионный двигатель", Type: "engine", Quantity: 1, Weight: 10})

		resp, err := client.Get(server.URL + "/parts")
		require.NoError(t, err)
		defer func() { _ = resp.Body.Close() }()

		require.Equal(t, http.StatusOK, resp.StatusCode)

		var parts []Part
		err = json.NewDecoder(resp.Body).Decode(&parts)
		require.NoError(t, err)
		require.NotEmpty(t, parts)
	})

	t.Run("POST /parts создаёт деталь", func(t *testing.T) {
		body := bytes.NewBufferString(`{"name":"Новый двигатель","type":"engine","quantity":5,"weight":100.5}`)

		resp, err := client.Post(server.URL+"/parts", "application/json", body)
		require.NoError(t, err)
		defer func() { _ = resp.Body.Close() }()

		require.Equal(t, http.StatusCreated, resp.StatusCode)

		var part Part
		err = json.NewDecoder(resp.Body).Decode(&part)
		require.NoError(t, err)
		require.NotZero(t, part.ID)
		require.Equal(t, "Новый двигатель", part.Name)
	})

	t.Run("POST /parts возвращает 400 для невалидного JSON", func(t *testing.T) {
		body := bytes.NewBufferString(`{invalid json}`)

		resp, err := client.Post(server.URL+"/parts", "application/json", body)
		require.NoError(t, err)
		defer func() { _ = resp.Body.Close() }()

		require.Equal(t, http.StatusBadRequest, resp.StatusCode)
	})

	t.Run("DELETE /parts/{id} удаляет деталь", func(t *testing.T) {
		created := repo.Create(Part{Name: "Титановая обшивка", Type: "fuel", Quantity: 1, Weight: 50})
		url := fmt.Sprintf("%s/parts/%d", server.URL, created.ID)

		req, err := http.NewRequest(http.MethodDelete, url, nil)
		require.NoError(t, err)

		resp, err := client.Do(req)
		require.NoError(t, err)
		defer func() { _ = resp.Body.Close() }()

		require.Equal(t, http.StatusNoContent, resp.StatusCode)
	})

	t.Run("DELETE /parts/{id} возвращает 404 для несуществующей детали", func(t *testing.T) {
		url := fmt.Sprintf("%s/parts/%d", server.URL, 99999)

		req, err := http.NewRequest(http.MethodDelete, url, nil)
		require.NoError(t, err)

		resp, err := client.Do(req)
		require.NoError(t, err)
		defer func() { _ = resp.Body.Close() }()

		require.Equal(t, http.StatusNotFound, resp.StatusCode)
	})

	t.Run("DELETE /parts/{id} возвращает 400 для невалидного ID", func(t *testing.T) {
		url := fmt.Sprintf("%s/parts/%s", server.URL, "not_a_number")

		req, err := http.NewRequest(http.MethodDelete, url, nil)
		require.NoError(t, err)

		resp, err := client.Do(req)
		require.NoError(t, err)
		defer func() { _ = resp.Body.Close() }()

		require.Equal(t, http.StatusBadRequest, resp.StatusCode)
	})
}
