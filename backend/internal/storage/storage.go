package storage

import (
	cfg "REST_project/config"
	"database/sql"
	"REST_project/internal/models"
	"fmt"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/lib/pq"
)

const (
	migrationPath = "file://migrations" 
)

type Storage struct {
	DB *sql.DB
}

func New(c cfg.DatabaseCfg) (*Storage, error) {
	const op = "storage.connection"
	connstr := fmt.Sprintf("user=%s password=%s dbname=%s host=%s port=%s sslmode=disable", c.User, c.Password, c.DBName, c.Host, c.Port)
	db, err := sql.Open("postgres", connstr)
	if err != nil {
		return nil, fmt.Errorf("%s : %w", op, err)
	}
	if err = db.Ping(); err != nil {
		return nil, fmt.Errorf("%s : %w", op, err)
	}
	s := &Storage{
		db,
	}
	return s, err
}

func RunMigrations(db *sql.DB) error {
	driver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		return fmt.Errorf("Error of migrate: %w", err)
	}
	m, err := migrate.NewWithDatabaseInstance(
		migrationPath,
		"postgres",
		driver,
	)
	if err != nil {
		return fmt.Errorf("Error of create migrate: %w", err)
	}
	if err = m.Up(); err != nil {
		if err != migrate.ErrNoChange {
			return fmt.Errorf("Error of migrate UP: %w", err)
		}
	}
	return nil
}

func (s *Storage) EnterpriseRegister(name string) (int, error) {
	const op = "storage.postgres.EnterpriseRegister"

	stmt, err := s.DB.Prepare("INSERT INTO enterprises (name) VALUES ($1) RETURNING id;")
	if err != nil {
		return 0, fmt.Errorf("%s: %w", op, err)
	}
	defer stmt.Close()

	var id int
	err = stmt.QueryRow(name).Scan(&id)
	if err != nil {
		return 0, fmt.Errorf("%s: %w", op, err)
	}

	return id, nil
}

func (s *Storage) EventRegister(name, description string, enterprise_id int) (int, error) {
	const op = "storage.postgres.EventRegister"
	stmt, err := s.DB.Prepare("INSERT INTO events (name, enterprise_id, description) VALUES ($1, $2, $3) RETURNING id;")
	if err != nil {
		return 0, fmt.Errorf("%s: %w", op, err)
	}
	defer stmt.Close()

	var id int
	err = stmt.QueryRow(name, enterprise_id, description).Scan(&id)
	if err != nil {
		return 0, fmt.Errorf("%s: %w", op, err)
	}
	return id, nil
}

func (s *Storage) ParticipantRegister(event_id int, name string) (int, error) {
	const op = "storage.postgres.EventRegister"
	stmt, err := s.DB.Prepare("INSERT INTO participants (name, event_id) VALUES ($1, $2) RETURNING id;")
	if err != nil {
		return 0, fmt.Errorf("%s: %w", op, err)
	}
	defer stmt.Close()

	var id int
	err = stmt.QueryRow(name, event_id).Scan(&id)
	if err != nil {
		return 0, fmt.Errorf("%s: %w", op, err)
	}
	return id, nil
}

func (s *Storage) CreatePost(content string, event_id int) (int, error) {
	const op = "storage.postgres.EventRegister"
	stmt, err := s.DB.Prepare("INSERT INTO posts (content, event_id) VALUES ($1, $2) RETURNING id;")
	if err != nil {
		return 0, fmt.Errorf("%s: %w", op, err)
	}
	defer stmt.Close()

	var id int
	err = stmt.QueryRow(content, event_id).Scan(&id)
	if err != nil {
		return 0, fmt.Errorf("%s: %w", op, err)
	}
	return id, nil
}

func (s *Storage) CreateComment(postID, participantID int, content string) (int, error) {
	const op = "storage.postgres.CreateComment"

	stmt, err := s.DB.Prepare("INSERT INTO comments (post_id, participant_id, content) VALUES ($1, $2, $3) RETURNING id;")
	if err != nil {
		return 0, fmt.Errorf("%s: %w", op, err)
	}
	defer stmt.Close()

	var id int
	err = stmt.QueryRow(postID, participantID, content).Scan(&id)
	if err != nil {
		return 0, fmt.Errorf("%s: %w", op, err)
	}
	return id, nil
}

func (s *Storage) GetComments() ([]models.Comment, error) {
	const op = "storage.GetComments"

	rows, err := s.DB.Query("SELECT id, post_id, participant_id, content, created_at FROM comments")
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	defer rows.Close()

	var comments []models.Comment
	for rows.Next() {
		var c models.Comment
		if err := rows.Scan(&c.ID, &c.PostID, &c.ParticipantID, &c.Content, &c.CreatedAt); err != nil {
			return nil, fmt.Errorf("%s: %w", op, err)
		}
		comments = append(comments, c)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return comments, nil
}

// GetEnterprises возвращает все предприятия из базы данных
func (s *Storage) GetEnterprises() ([]models.Enterprise, error) {
	const op = "storage.GetEnterprises"

	rows, err := s.DB.Query("SELECT id, name FROM enterprises")
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	defer rows.Close()

	var enterprises []models.Enterprise
	for rows.Next() {
		var e models.Enterprise
		if err := rows.Scan(&e.ID, &e.Name); err != nil {
			return nil, fmt.Errorf("%s: %w", op, err)
		}
		enterprises = append(enterprises, e)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return enterprises, nil
}

// GetEvents возвращает все события из базы данных
func (s *Storage) GetEvents() ([]models.Event, error) {
	const op = "storage.GetEvents"

	rows, err := s.DB.Query("SELECT id, enterprise_id, name, description, created_at FROM events")
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	defer rows.Close()

	var events []models.Event
	for rows.Next() {
		var e models.Event
		if err := rows.Scan(&e.ID, &e.EnterpriseID, &e.Name, &e.Description, &e.CreatedAt); err != nil {
			return nil, fmt.Errorf("%s: %w", op, err)
		}
		events = append(events, e)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return events, nil
}

// GetParticipants возвращает всех участников из базы данных
func (s *Storage) GetParticipants() ([]models.Participant, error) {
	const op = "storage.GetParticipants"

	rows, err := s.DB.Query("SELECT id, event_id, name FROM participants")
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	defer rows.Close()

	var participants []models.Participant
	for rows.Next() {
		var p models.Participant
		if err := rows.Scan(&p.ID, &p.EventID, &p.Name); err != nil {
			return nil, fmt.Errorf("%s: %w", op, err)
		}
		participants = append(participants, p)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return participants, nil
}

// GetPosts возвращает все посты из базы данных
func (s *Storage) GetPosts() ([]models.Post, error) {
	const op = "storage.GetPosts"

	rows, err := s.DB.Query("SELECT id, event_id, content, created_at FROM posts")
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	defer rows.Close()

	var posts []models.Post
	for rows.Next() {
		var p models.Post
		if err := rows.Scan(&p.ID, &p.EventID, &p.Content, &p.CreatedAt); err != nil {
			return nil, fmt.Errorf("%s: %w", op, err)
		}
		posts = append(posts, p)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return posts, nil
}