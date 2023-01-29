package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestE2E(t *testing.T) {
	// Configuration
	var cfg Config
	cfg.Limiter.Enabled = true
	cfg.Limiter.Rps = 2
	cfg.Limiter.Burst = 6
	cfg.SolrURL = "http://localhost:8983"
	cfg.SolrProfile = "profiles"

	// Logger
	logger := New(os.Stdout, LevelInfo)

	// Application
	app := Application{
		Config: cfg,
		Logger: logger,
	}

	// API Routes
	ts := httptest.NewTLSServer(app.Routes())
	defer ts.Close()

	// Database
	db, err := OpenDB(
		"postgres://postgres:postgres@localhost:5432/postgres?sslmode=disable")
	if err != nil {
		t.Fatal(err.Error())
	}
	defer db.Close()

	// Read SQL file
	script, err := os.ReadFile("./test/sql/delete_all.sql")
	if err != nil {
		t.Fatal(err)
	}

	// Delete Records
	_, err = db.Exec(string(script))
	if err != nil {
		t.Fatal(err)
	}

	// Delete the Profile Indexing on the Solr
	res, err := http.Post(
		"http://localhost:8983/solr/"+cfg.SolrProfile+"/update?commit=true",
		"application/json",
		bytes.NewReader([]byte("{'delete': {'query': '*:*'}}")))
	if err != nil {
		t.Fatal(err)
	}
	if res.StatusCode != http.StatusOK {
		t.Fatal(res)
	}

	// Initial
	email := "jon@doe.com"
	password := "pa55word"

	// Initial
	var userResponse map[string]User
	client := &http.Client{}

	t.Run("Register User", func(t *testing.T) {
		data := fmt.Sprintf(
			`{"email": "%v", "password": "%v", "first_name": "Jon", "last_name": "Doe"}`,
			email,
			password)
		req, _ := http.NewRequest(
			"POST",
			"http://localhost:8000/service/users",
			bytes.NewReader([]byte(data)))
		req.Header.Add("Content-Type", "application/json")

		res, err := client.Do(req)
		assert.Nil(t, err)
		assert.Equal(t, http.StatusCreated, res.StatusCode)

		body, err := io.ReadAll(res.Body)
		defer res.Body.Close()
		assert.Nil(t, err)

		err = json.Unmarshal(body, &userResponse)
		assert.Nil(t, err)
		assert.Equal(t, email, userResponse["user"].Email)
	})

	// Initial
	type authType struct {
		Token string `json:"token"`
	}
	var authentication authType

	t.Run("Login User", func(t *testing.T) {
		data := fmt.Sprintf(
			`{"email": "%v", "password": "%v"}`,
			email,
			password)
		req, _ := http.NewRequest(
			"POST",
			"http://localhost:8000/service/users/authentication",
			bytes.NewReader([]byte(data)))
		req.Header.Add("Content-Type", "application/json")

		res, err := client.Do(req)
		assert.Nil(t, err)
		assert.Equal(t, http.StatusOK, res.StatusCode)

		body, err := io.ReadAll(res.Body)
		defer res.Body.Close()
		assert.Nil(t, err)

		err = json.Unmarshal(body, &authentication)
		assert.Nil(t, err)
		assert.NotNil(t, authentication.Token)
	})

	// Initial
	var profileResponse map[string]Profile

	t.Run("Create Profile", func(t *testing.T) {
		tBody, tContentType := app.testFormProfile(t)
		req, _ := http.NewRequest(
			"POST",
			"http://localhost:8000/service/profiles",
			tBody)
		req.Header.Add("Content-Type", tContentType)

		bearer := fmt.Sprintf("Bearer %v", authentication.Token)
		req.Header.Set("Authorization", bearer)

		res, err := client.Do(req)
		assert.Nil(t, err)
		assert.Equal(t, http.StatusCreated, res.StatusCode)

		body, err := io.ReadAll(res.Body)
		defer res.Body.Close()
		assert.Nil(t, err)

		err = json.Unmarshal(body, &profileResponse)
		assert.Nil(t, err)
		assert.Equal(t,
			profileResponse["profile"].ProfileUser,
			userResponse["user"].ID)
	})

	t.Run("Search Profile", func(t *testing.T) {
		// We use gRPC to update a Profile on the Solr,
		// so it needs to wait a couple seconds until the updating is done
		time.Sleep(2 * time.Second)

		req, _ := http.NewRequest(
			"GET",
			ts.URL+"/service/search/profiles/?q=*:*",
			nil)

		res, err := ts.Client().Do(req)
		assert.Nil(t, err)
		assert.Equal(t, http.StatusOK, res.StatusCode)

		body, err := io.ReadAll(res.Body)
		defer res.Body.Close()
		assert.Nil(t, err)

		var result map[string]interface{}
		err = json.Unmarshal([]byte(body), &result)
		assert.Nil(t, err)

		response := result["response"].(map[string]interface{})
		assert.NotNil(t, response)
		assert.Equal(t, response["numFound"], float64(1))

		docs := response["docs"].([]interface{})
		assert.NotNil(t, docs)

		doc := docs[0].(map[string]interface{})
		assert.NotNil(t, doc)
		assert.NotNil(t, doc["profile_name"])
	})
}
