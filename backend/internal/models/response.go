package models

import "time"

type Enterprise struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

type Event struct {
	ID           int       `json:"id"`
	EnterpriseID int       `json:"enterprise_id"`
	Name         string    `json:"name"`
	Description  string    `json:"description,omitempty"`
	CreatedAt    time.Time `json:"created_at"`
}

type Participant struct {
	ID      int    `json:"id"`
	EventID int    `json:"event_id"`
	Name    string `json:"name"`
}

type Post struct {
	ID        int       `json:"id"`
	EventID   int       `json:"event_id"`
	Content   string    `json:"content"`
	CreatedAt time.Time `json:"created_at"`
}

type Comment struct {
	ID           int       `json:"id"`
	PostID       int       `json:"post_id"`
	ParticipantID int       `json:"participant_id"`
	Content      string    `json:"content"`
	CreatedAt    time.Time `json:"created_at"`
}

type Response struct {
	Status string      `json:"status"`
	Data   interface{} `json:"data,omitempty"`
	Error  string      `json:"error,omitempty"`
}

// Request структуры для входящих запросов

type CreateEnterpriseRequest struct {
	Name string `json:"name"`
}

type CreateEventRequest struct {
	EnterpriseID int    `json:"enterprise_id"`
	Name         string `json:"name"`
	Description  string `json:"description,omitempty"`
}

type CreateParticipantRequest struct {
	EventID int    `json:"event_id"`
	Name    string `json:"name"`
}

type CreatePostRequest struct {
	EventID int    `json:"event_id"`
	Content string `json:"content"`
}

type CreateCommentRequest struct {
	PostID       int    `json:"post_id"`
	ParticipantID int    `json:"participant_id"`
	Content      string `json:"content"`
}