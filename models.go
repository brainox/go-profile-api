
package main
import "time"

// CatFactResponse represents the response from the Cat Facts API
type CatFactResponse struct {
	Fact string `json:"fact"`
}

// ProfileResponse represents the response for the /me endpoint
type ProfileResponse struct {
	Status    string   `json:"status"`
	User      UserInfo `json:"user"`
	Timestamp string   `json:"timestamp"`
	Fact      string   `json:"fact"`
}

// UserInfo contains user details
type UserInfo struct {
	Email string `json:"email"`
	Name  string `json:"name"`
	Stack string `json:"stack"`
}

// ErrorResponse represents an error response
type ErrorResponse struct {
	Status  string `json:"status"`
	Message string `json:"message"`
}

// Config holds application configuration
type Config struct {
	Email             string
	Name              string
	Stack             string
	CatFactAPIURL     string
	CatFactAPITimeout time.Duration
	Port              string
}