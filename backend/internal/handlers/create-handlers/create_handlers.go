package create_handlers

import (
	model "REST_project/internal/models"
	"errors"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"io"
	"log/slog"
	"net/http"
)

type Server interface {
	CreatePost(content string, event_id int) (int, error)
	CreateComment(postID int, participantID int, content string) (int, error)
	GetPosts() ([]model.Post, error)
	GetComments() ([]model.Comment, error)
}

// ... существующие структуры RequestPostCreate и RequestCommentCreate ...

func GetPosts(log *slog.Logger, s Server) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "internal.handlers.create-handlers.GetPosts"
		log = log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		log.Info("getting posts")

		posts, err := s.GetPosts()
		if err != nil {
			log.Error("failed to get posts", slog.String("error", err.Error()))
			render.JSON(w, r, model.Response{
				Status: "Error",
				Error:  "failed to get posts",
			})
			return
		}

		render.JSON(w, r, model.Response{
			Status: "OK",
			Data:   posts,
		})
	}
}

func GetComments(log *slog.Logger, s Server) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "internal.handlers.create-handlers.GetComments"
		log = log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		log.Info("getting comments")

		comments, err := s.GetComments()
		if err != nil {
			log.Error("failed to get comments", slog.String("error", err.Error()))
			render.JSON(w, r, model.Response{
				Status: "Error",
				Error:  "failed to get comments",
			})
			return
		}

		render.JSON(w, r, model.Response{
			Status: "OK",
			Data:   comments,
		})
	}
}


type RequestPostCreate struct {
	Content string `json:"content"`
	EventID int    `json:"event_id"`
}

func respOk(w http.ResponseWriter, r *http.Request) {
	render.JSON(w, r, model.Response{
		Status: "OK",
	})
}

func CreatePost(log *slog.Logger, s Server) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "internal.handlers.register-handlers.CreatePost"
		log = log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		var req RequestPostCreate
		err := render.DecodeJSON(r.Body, &req)
		if errors.Is(err, io.EOF) {
			log.Error("request body is empty")
			render.JSON(w, r, model.Response{
				Status: "Error",
				Error:  "empty request",
			})
			return
		}
		if err != nil {
			log.Error("failed to decode request body", slog.String("error", err.Error()))
			render.JSON(w, r, model.Response{
				Status: "Error",
				Error:  "invalid request format",
			})
			return
		}

		if req.Content == "" || req.EventID <= 0 {
			log.Error("invalid request data")
			render.JSON(w, r, model.Response{
				Status: "Error",
				Error:  "invalid data provided",
			})
			return
		}

		log.Info("creating post", slog.Any("request", req))

		_, err = s.CreatePost(req.Content, req.EventID)
		if err != nil {
			log.Error("failed to create post", slog.String("error", err.Error()))
			render.JSON(w, r, model.Response{
				Status: "Error",
				Error:  "failed to create post",
			})
			return
		}

		respOk(w, r)
	}
}

type RequestCommentCreate struct {
	PostID        int    `json:"post_id"`
	ParticipantID int    `json:"participant_id"`
	Content       string `json:"content"`
}

func CreateComment(log *slog.Logger, s Server) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "internal.handlers.register-handlers.CreateComment"
		log = log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		var req RequestCommentCreate
		err := render.DecodeJSON(r.Body, &req)
		if errors.Is(err, io.EOF) {
			log.Error("request body is empty")
			render.JSON(w, r, model.Response{
				Status: "Error",
				Error:  "empty request",
			})
			return
		}
		if err != nil {
			log.Error("failed to decode request body", slog.String("error", err.Error()))
			render.JSON(w, r, model.Response{
				Status: "Error",
				Error:  "invalid request format",
			})
			return
		}

		if req.Content == "" || req.PostID <= 0 || req.ParticipantID <= 0 {
			log.Error("invalid request data")
			render.JSON(w, r, model.Response{
				Status: "Error",
				Error:  "invalid data provided",
			})
			return
		}

		log.Info("creating comment", slog.Any("request", req))

		_, err = s.CreateComment(req.PostID, req.ParticipantID, req.Content)
		if err != nil {
			log.Error("failed to create comment", slog.String("error", err.Error()))
			render.JSON(w, r, model.Response{
				Status: "Error",
				Error:  "failed to create comment",
			})
			return
		}
		respOk(w, r)
	}
}
