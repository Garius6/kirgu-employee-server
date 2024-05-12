package repository

import (
	"context"
	"embed"
	"errors"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5/stdlib"
	"github.com/pressly/goose/v3"
	pgxUUID "github.com/vgarvardt/pgx-google-uuid/v5"
	"kirgu.ru/employee/model"
)

var (
	ErrNotFound          = errors.New("not found")
	ErrIncorrectPassword = errors.New("incorrect password")
	ErrUsernameIsTaken   = errors.New("username already taken")
)

type PostgresRepository struct {
	pool *pgxpool.Pool
}

//go:embed migrations/*.sql
var embedMigrations embed.FS

func NewPostgresRepository(dbUrl string) (*PostgresRepository, error) {
	pgxConfig, err := pgxpool.ParseConfig(dbUrl)
	if err != nil {
		panic(err)
	}

	pgxConfig.AfterConnect = func(ctx context.Context, conn *pgx.Conn) error {
		pgxUUID.Register(conn.TypeMap())
		return nil
	}

	pool, err := pgxpool.NewWithConfig(context.Background(), pgxConfig)
	if err != nil {
		return nil, err
	}

	err = execMigrations(pool)
	if err != nil {
		return nil, err
	}

	return &PostgresRepository{pool: pool}, nil
}

func execMigrations(pool *pgxpool.Pool) error {
	goose.SetBaseFS(embedMigrations)

	if err := goose.SetDialect("postgres"); err != nil {
		return err
	}

	db := stdlib.OpenDBFromPool(pool)
	if err := goose.Up(db, "migrations"); err != nil {
		return err
	}

	return nil
}

func (p *PostgresRepository) SignIn(username string, password string) (*model.User, error) {
	user, err := p.getUserByUsername(username)
	if err != nil {
		return nil, err
	}

	if user.Password != password {
		return nil, ErrIncorrectPassword
	}

	return user, nil
}

func (p *PostgresRepository) SignUp(username string, password string, passwordConfirmation string) error {

	if password != passwordConfirmation {
		return ErrIncorrectPassword
	}

	// Check that user with this username doesn`t exists
	_, err := p.getUserByUsername(username)
	if err == nil {
		return ErrUsernameIsTaken
	} else if !errors.Is(err, ErrNotFound) {
		return err
	}

	uid := uuid.New()
	_, err = p.pool.Exec(context.Background(), "INSERT INTO users(id, username, password) VALUES ($1, $2, $3)", uid, username, password)
	if err != nil {
		return err
	}

	return nil
}

func (p *PostgresRepository) getUserByUsername(username string) (*model.User, error) {
	var user model.User

	err := p.pool.QueryRow(context.Background(), "SELECT * FROM users WHERE username = $1;", username).Scan(
		&user.Id, &user.Username, &user.Password,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, err
	}
	return &user, nil
}
