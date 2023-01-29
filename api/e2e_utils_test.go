package api

import (
	"bytes"
	"context"
	"database/sql"
	"io"
	"mime/multipart"
	"os"
	"testing"
	"time"

	"github.com/google/uuid"
	_ "github.com/lib/pq"
)

type User struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	Email     string    `json:"email"`
	FirstName string    `json:"first_name"`
	LastName  string    `json:"last_name"`
	Activated bool      `json:"activated"`
	Version   int       `json:"version"`
}

type Profile struct {
	ID             uuid.UUID `json:"id"`
	CreatedAt      time.Time `json:"created_at"`
	ProfileUser    uuid.UUID `json:"profile_user"`
	ProfileName    string    `json:"profile_name"`
	ProfilePicture string    `json:"profile_picture"`
	Version        int       `json:"-"`
}

func (app *Application) testFormProfile(t *testing.T) (io.Reader, string) {
	bodyBuf := &bytes.Buffer{}
	bodyWriter := multipart.NewWriter(bodyBuf)

	// Add team name
	bodyWriter.WriteField("profile_name", "John Doe")

	// Add team picture
	filename := "./test/images/profile.jpg"
	fileWriter, err := bodyWriter.CreateFormFile("profile_picture", filename)
	if err != nil {
		t.Fatal(err)
	}

	// Open file
	fileHandler, err := os.Open(filename)
	if err != nil {
		t.Fatal(err)
	}

	// Copy file
	_, err = io.Copy(fileWriter, fileHandler)
	if err != nil {
		t.Fatal(err)
	}

	// Put on body
	contentType := bodyWriter.FormDataContentType()
	bodyWriter.Close()

	return bodyBuf, contentType
}

func OpenDB(dsn string) (*sql.DB, error) {
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, err
	}

	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(25)

	duration, err := time.ParseDuration("15m")
	if err != nil {
		return nil, err
	}

	db.SetConnMaxIdleTime(duration)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err = db.PingContext(ctx)
	if err != nil {
		return nil, err
	}

	return db, nil
}
