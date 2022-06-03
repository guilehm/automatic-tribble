package postgres

import (
	"context"
	"errors"
	"log"
	"os"
	"strings"
	"time"
	"tribble/models"

	"github.com/jackc/pgx/v4/pgxpool"
)

type Postgres struct {
	DB *pgxpool.Pool
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
		DB: connect(),
	}
}

func (p Postgres) Close() {
	p.DB.Close()
}

func (p Postgres) GetUser(ctx context.Context, ID int) (*models.User, error) {
	sql := `SELECT id, username, email, date_joined FROM users WHERE id=$1`

	var user models.User
	if err := p.DB.QueryRow(ctx, sql, ID).Scan(
		&user.ID, &user.Username, &user.Email, &user.DateJoined,
	); err != nil {
		return nil, err
	}
	return &user, nil
}

func (p Postgres) GetUserByEmail(ctx context.Context, email string) (*models.User, error) {
	sql := `SELECT id, username, email, password, date_joined FROM users WHERE email=$1`

	var user models.User
	lowerEmail := strings.ToLower(*user.Email)
	if err := p.DB.QueryRow(ctx, sql, strings.ToLower(email)).Scan(
		&user.ID, &user.Username, &lowerEmail, &user.Password, &user.DateJoined,
	); err != nil {
		return nil, err
	}
	return &user, nil
}

func (p Postgres) GetUserByRefresh(ctx context.Context, refresh string) (*models.User, error) {
	sql := `SELECT id, username, email, password, date_joined FROM users WHERE refresh_token=$1`

	var user models.User
	if err := p.DB.QueryRow(ctx, sql, refresh).Scan(
		&user.ID, &user.Username, &user.Email, &user.Password, &user.DateJoined,
	); err != nil {
		return nil, err
	}
	return &user, nil
}

func (p Postgres) GetUserList(ctx context.Context) ([]*models.User, error) {
	users := make([]*models.User, 0)

	sql := `SELECT id, username, email, date_joined FROM users`
	rows, err := p.DB.Query(ctx, sql)
	if err != nil {
		return users, err
	}

	for rows.Next() {
		var user models.User
		err = rows.Scan(&user.ID, &user.Username, &user.Email, &user.DateJoined)
		if err != nil {
			return users, err
		}
		users = append(users, &user)
	}
	return users, nil
}

func (p Postgres) CreateUser(ctx context.Context, user models.User) (*models.User, error) {

	sql := `INSERT INTO users (username, email, date_joined, password, token, refresh_token) 
			VALUES ($1, $2, $3, $4, $5, $6) 
			RETURNING id`

	var id int
	if err := p.DB.QueryRow(
		ctx,
		sql,
		user.Username,
		strings.ToLower(*user.Email),
		user.DateJoined,
		user.Password,
		user.Token,
		user.RefreshToken,
	).Scan(&id); err != nil {
		return nil, err
	}
	user.ID = id
	if user.Email != nil {
		loweredEmail := strings.ToLower(*user.Email)
		user.Email = &loweredEmail
	}
	return &user, nil
}

func (p Postgres) UpdateUser(ctx context.Context, user models.User) (*models.User, error) {
	// TODO: updating only username for now
	sql := `UPDATE users SET username=$2 WHERE id=$1`
	res, err := p.DB.Exec(ctx, sql, user.ID, user.Username)
	if err != nil {
		return nil, err
	}

	if rowsAffected := res.RowsAffected(); rowsAffected == 0 {
		return nil, errors.New("user not found")
	}

	return &user, nil
}

func (p Postgres) UpdateUserTokens(ctx context.Context, ID int, token, refresh string) error {
	sql := `UPDATE users SET token=$1, refresh_token=$2 WHERE id=$3`
	_, err := p.DB.Exec(ctx, sql, token, refresh, ID)
	return err
}

func (p Postgres) DeleteUser(ctx context.Context, ID int) error {
	sql := `DELETE FROM users WHERE id=$1`
	res, err := p.DB.Exec(ctx, sql, ID)
	if err != nil {
		return err
	}
	if rowsAffected := res.RowsAffected(); rowsAffected == 0 {
		return errors.New("user not found")
	}
	return nil
}

func (p Postgres) GetPlayerList(ctx context.Context, ID int) ([]*models.Player, error) {
	sql := `SELECT name, xp, sprite, position_x, position_y FROM players WHERE user_id=$1`
	rows, err := p.DB.Query(ctx, sql, ID)
	if err != nil {
		return []*models.Player{}, err
	}

	players := make([]*models.Player, 0)
	for rows.Next() {
		var player models.Player
		err = rows.Scan(
			&player.Name,
			&player.XP,
			&player.Sprite,
			&player.PositionX,
			&player.PositionY,
		)
		if err != nil {
			return []*models.Player{}, err
		}
		players = append(players, &player)
	}

	return players, nil
}

func (p Postgres) CreatePlayer(ctx context.Context, player models.Player) (*models.Player, error) {
	sql := `INSERT INTO players (user_id, name, xp, sprite, position_x, position_y)
			VALUES ($1, $2, $3, $4, $5, $6)
			RETURNING id`

	var playerID int
	err := p.DB.QueryRow(
		ctx,
		sql,
		player.UserID,
		player.Name,
		player.XP,
		player.Sprite,
		player.PositionX,
		player.PositionY,
	).Scan(&playerID)
	player.ID = playerID

	if err != nil {
		return &player, err
	}

	return &player, nil
}

func (p Postgres) ValidateToken(ctx context.Context, refresh string) (bool, error) {
	sql := `SELECT id FROM users WHERE refresh_token=$1`
	var userId int
	err := p.DB.QueryRow(ctx, sql, refresh).Scan(&userId)
	if err != nil {
		return false, err
	}
	return true, nil
}
