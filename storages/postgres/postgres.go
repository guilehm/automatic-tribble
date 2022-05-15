package postgres

import (
	"context"
	"errors"
	"log"
	"os"
	"time"
	"tribble/models"

	"github.com/jackc/pgx/v4/pgxpool"
)

type Postgres struct {
	db *pgxpool.Pool
}

func connect() *pgxpool.Pool {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	db, err := pgxpool.Connect(ctx, os.Getenv("DATABASE_URL"))
	if err != nil {
		log.Fatalf("Unable to connect to database: %v", err.Error())
	}
	return db
}

func GetPostgres() Postgres {
	return Postgres{
		db: connect(),
	}
}

func (p Postgres) Close() {
	p.db.Close()
}

func (p Postgres) GetUser(ctx context.Context, ID int) (*models.User, error) {
	sql := `SELECT id, name, email, date_joined FROM users WHERE id=$1`

	var user models.User
	if err := p.db.QueryRow(ctx, sql, ID).Scan(
		&user.ID, &user.Name, &user.Email, &user.DateJoined,
	); err != nil {
		return &models.User{}, err
	}
	return &user, nil
}

func (p Postgres) GetUserByEmail(ctx context.Context, email string) (*models.User, error) {
	sql := `SELECT id, name, email, password, date_joined FROM users WHERE email=$1`

	var user models.User
	if err := p.db.QueryRow(ctx, sql, email).Scan(
		&user.ID, &user.Name, &user.Email, &user.Password, &user.DateJoined,
	); err != nil {
		return nil, err
	}
	return &user, nil
}

func (p Postgres) GetUserByRefresh(ctx context.Context, refresh string) (*models.User, error) {
	sql := `SELECT id, name, email, password, date_joined FROM users WHERE refresh_token=$1`

	var user models.User
	if err := p.db.QueryRow(ctx, sql, refresh).Scan(
		&user.ID, &user.Name, &user.Email, &user.Password, &user.DateJoined,
	); err != nil {
		return nil, err
	}
	return &user, nil
}

func (p Postgres) GetUserList(ctx context.Context) ([]*models.User, error) {
	users := make([]*models.User, 0)

	sql := `SELECT id, name, email, date_joined FROM users`
	rows, err := p.db.Query(ctx, sql)
	if err != nil {
		return users, err
	}

	for rows.Next() {
		var user models.User
		err = rows.Scan(&user.ID, &user.Name, &user.Email, &user.DateJoined)
		if err != nil {
			return users, err
		}
		users = append(users, &user)
	}
	return users, nil
}

func (p Postgres) CreateUser(ctx context.Context, user models.User) (*models.User, error) {

	sql := `INSERT INTO users (name, email, date_joined, password, token, refresh_token) 
			VALUES ($1, $2, $3, $4, $5, $6) 
			RETURNING id`

	var id int
	if err := p.db.QueryRow(
		ctx,
		sql,
		user.Name,
		user.Email,
		user.DateJoined,
		user.Password,
		user.Token,
		user.RefreshToken,
	).Scan(&id); err != nil {
		return &models.User{}, err
	}
	user.ID = id
	return &user, nil
}

func (p Postgres) UpdateUser(ctx context.Context, user models.User) (*models.User, error) {
	// TODO: updating only name for now
	sql := `UPDATE users SET name=$2 WHERE id=$1`
	res, err := p.db.Exec(ctx, sql, user.ID, user.Name)
	if err != nil {
		return &models.User{}, err
	}

	if rowsAffected := res.RowsAffected(); rowsAffected == 0 {
		return &models.User{}, errors.New("user not found")
	}

	return &user, nil
}

func (p Postgres) UpdateUserTokens(ctx context.Context, ID int, token, refresh string) error {
	sql := `UPDATE users SET token=$1, refresh_token=$2 WHERE id=$3`
	_, err := p.db.Exec(ctx, sql, token, refresh, ID)
	return err
}

func (p Postgres) DeleteUser(ctx context.Context, ID int) error {
	sql := `DELETE FROM users WHERE id=$1`
	res, err := p.db.Exec(ctx, sql, ID)
	if err != nil {
		return err
	}
	if rowsAffected := res.RowsAffected(); rowsAffected == 0 {
		return errors.New("user not found")
	}
	return nil
}
