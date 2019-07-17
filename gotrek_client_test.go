package gotrekclient

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func setupTrekHttpClient() *TrekHttpClient {

	return &TrekHttpClient{
		Timeout:               1 * time.Second,
		BackoffInterval:       5 * time.Microsecond,
		MaximumJitterInterval: 5 * time.Microsecond,
		RetryCount:            4,
	}
}

func TestPublish(t *testing.T) {

	type someStruct struct {
		SomeField    string
		AnotherField int
	}
	var someMapInterface map[string]interface{}

	wantID := "some-id"
	wantTimestamp := time.Now().Unix()
	wantTag := "some-tag"
	wantTrail := &someStruct{
		SomeField:    "some-field",
		AnotherField: 1,
	}

	jsonByte, _ := json.Marshal(wantTrail)
	json.Unmarshal(jsonByte, &someMapInterface)

	testServer := httptest.NewServer(http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {

		var trail map[string]interface{}
		err := json.NewDecoder(req.Body).Decode(&trail)
		fmt.Printf("%v", trail)
		assert.Nil(t, err, "Error should be nil")
		assert.Equal(t, wantID, trail["id"], "Audit ID not equal")
		assert.EqualValues(t, wantTimestamp, trail["timestamp"], "Timestamp not equal")
		assert.Equal(t, wantTag, trail["tag"], "Tag not equal")
		nested := trail["trail"].(map[string]interface{})
		assert.Equal(t, wantTrail.SomeField, nested["SomeField"], "Trail not equal")
		assert.EqualValues(t, wantTrail.AnotherField, nested["AnotherField"], "Trail not equal")

		res.WriteHeader(201)
	}))
	defer func() { testServer.Close() }()

	trek := setupTrekHttpClient()
	url := testServer.URL
	trekClient := NewTrekClient(url, "some-secret", trek)

	err := trekClient.Publish(wantID, someMapInterface, wantTimestamp, wantTag)

	assert.Nil(t, err, "Error should be nil")
}
