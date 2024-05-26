package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

/*
	Notes

	- write a test by creating a file with a name ending in _test.go that contains functions named
	TestXXX with signature func (t *testing.T)

	- run the tests with : go test
		- must be in the directory containing the tests


*/

func TestAddItem(t *testing.T) {
	method := "POST"
	url := "http://localhost:8090/addItem"
	jsonStr := []byte(`{"title":"Learn go","description":"complete intern project"}`)
	req, err := http.NewRequest(method, url, bytes.NewBuffer(jsonStr))
	assert.Nil(t, err)

	req.Header.Set("Content-Type", "application/json")
	responseRecorder := httptest.NewRecorder()
	handler := http.HandlerFunc(addItem)
	handler.ServeHTTP(responseRecorder, req)
	assert.Equal(t, http.StatusCreated, responseRecorder.Code)

	expected := "Item successfully added"
	assert.Equal(t, expected, responseRecorder.Body.String())
}

func TestGetAllItems(t *testing.T) {
	// reset global variable (not sure if this is the correct way to do it)
	items = []Item{
		{Id: 1, Title: "Learn Go", Description: "Complete the project"},
		{Id: 2, Title: "Write Tests", Description: "Ensure all code is tested"},
	}
	method := "GET"
	url := "http://localhost:8090/getAllItems"
	req, err := http.NewRequest(method, url, nil)
	assert.Nil(t, err)

	responseRecorder := httptest.NewRecorder()
	handler := http.HandlerFunc(getAllItems)
	handler.ServeHTTP(responseRecorder, req)
	assert.Equal(t, http.StatusOK, responseRecorder.Code)

	var responseItems []Item
	err = json.Unmarshal(responseRecorder.Body.Bytes(), &responseItems)
	assert.Nil(t, err)

	// Compare the decoded response with the expected items
	assert.Equal(t, items, responseItems)
}

func TestCompleteItem(t *testing.T) {
	items = []Item{
		{Id: 1, Title: "Learn Go", Description: "Complete the project"},
		{Id: 2, Title: "Write Tests", Description: "Ensure all code is tested"},
	}
	method := "PUT"
	url := "http://localhost:8090/completeItem?id=1"
	req, err := http.NewRequest(method, url, nil)
	assert.Nil(t, err)

	responseRecorder := httptest.NewRecorder()
	handler := http.HandlerFunc(completeItem)
	handler.ServeHTTP(responseRecorder, req)

	assert.Equal(t, responseRecorder.Code, http.StatusOK)
	assert.Equal(t, items[0].Completed, true)
}

// Would add more tests for all logic trees (e.g didnt add unit tests for failures)
