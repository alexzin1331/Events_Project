package register_handlers

import (
	model "REST_project/internal/models"
	"errors"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"io"
	"log/slog"
	"net/http"
)

type RequestEntRegister struct {
    Name string `json:"name"` 
}

type RequestEventRegister struct {
    Name         string `json:"name"`
    Description  string `json:"description"`
    EnterpriseID int    `json:"enterprise_id"`
}

type RequestUserRegister struct {
    EventID int    `json:"event_id"`
    Name    string `json:"name"`
}

type Server interface {
    EnterpriseRegister(name string) (int, error)
    EventRegister(name string, description string, enterpriseID int) (int, error)
    ParticipantRegister(eventID int, name string) (int, error) 
    GetEnterprises() ([]model.Enterprise, error)
    GetEvents() ([]model.Event, error)
    GetParticipants() ([]model.Participant, error) 
}

func respOk(w http.ResponseWriter, r *http.Request) {
	render.JSON(w, r, model.Response{
		Status: "OK",
	})
}

func RegisterEnterprise(log *slog.Logger, s Server) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "internal.handlers.register-handlers"
		log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)
		var req RequestEntRegister
		err := render.DecodeJSON(r.Body, &req)
		if errors.Is(err, io.EOF) {
			log.Error("enterprise_name is empty")
			render.JSON(w, r, model.Response{
				Status: "Error",
				Error:  "empty request",
			})
			return
		}
		if err != nil {
			log.Error("failed to decode request body", slog.Attr{
				Key:   "Error",
				Value: slog.StringValue(err.Error()),
			})
			render.JSON(w, r, model.Response{
				Status: "Error",
				Error:  "empty request",
			})
			return
		}
		log.Info("request body decoded", slog.Any("request", req))
		_, err = s.EnterpriseRegister(req.Name)
		if err != nil {
			log.Error("failed to write enterprise name to database")
			render.JSON(w, r, model.Response{
				Status: "Error",
				Error:  "failed to write data",
			})
			return
		}
		respOk(w, r)
	}
}


func RegisterEvent(log *slog.Logger, s Server) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "internal.handlers.register-handlers.RegisterEvent"
		log = log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		var req RequestEventRegister
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

		if req.Name == "" || req.Description == "" || req.EnterpriseID <= 0 {
			log.Error("invalid request data")
			render.JSON(w, r, model.Response{
				Status: "Error",
				Error:  "invalid data provided",
			})
			return
		}

		log.Info("registering event", slog.Any("request", req))

		_, err = s.EventRegister(req.Name, req.Description, req.EnterpriseID)
		if err != nil {
			log.Error("failed to register event", slog.String("error", err.Error()))
			render.JSON(w, r, model.Response{
				Status: "Error",
				Error:  "failed to register event",
			})
			return
		}

		respOk(w, r)
	}
}


func RegisterUser(log *slog.Logger, s Server) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "internal.handlers.register-handlers.RegisterUser"
		log = log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		var req RequestUserRegister
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

		if req.Name == "" || req.EventID <= 0 {
			log.Error("invalid request data")
			render.JSON(w, r, model.Response{
				Status: "Error",
				Error:  "invalid data provided",
			})
			return
		}

		log.Info("registering user", slog.Any("request", req))

		_, err = s.ParticipantRegister(req.EventID, req.Name)
		if err != nil {
			log.Error("failed to register user", slog.String("error", err.Error()))
			render.JSON(w, r, model.Response{
				Status: "Error",
				Error:  "failed to register user",
			})
			return
		}

		respOk(w, r)
	}
}

func GetEnterprises(log *slog.Logger, s Server) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "internal.handlers.create-handlers.GetEnterprises"
		log = log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		log.Info("getting enterprises")

		enterprises, err := s.GetEnterprises()
		if err != nil {
			log.Error("failed to get enterprises", slog.String("error", err.Error()))
			render.JSON(w, r, model.Response{
				Status: "Error",
				Error:  "failed to get enterprises",
			})
			return
		}

		render.JSON(w, r, model.Response{
			Status: "OK",
			Data:   enterprises,
		})
	}
}

func GetEvents(log *slog.Logger, s Server) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "internal.handlers.create-handlers.GetEvents"
		log = log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		log.Info("getting events")

		events, err := s.GetEvents()
		if err != nil {
			log.Error("failed to get events", slog.String("error", err.Error()))
			render.JSON(w, r, model.Response{
				Status: "Error",
				Error:  "failed to get events",
			})
			return
		}

		render.JSON(w, r, model.Response{
			Status: "OK",
			Data:   events,
		})
	}
}

func GetUsers(log *slog.Logger, s Server) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "internal.handlers.create-handlers.GetUsers"
		log = log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		log.Info("getting users")

		users, err := s.GetParticipants()
		if err != nil {
			log.Error("failed to get users", slog.String("error", err.Error()))
			render.JSON(w, r, model.Response{
				Status: "Error",
				Error:  "failed to get users",
			})
			return
		}

		render.JSON(w, r, model.Response{
			Status: "OK",
			Data:   users,
		})
	}
}
