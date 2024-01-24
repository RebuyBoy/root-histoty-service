package repository

import (
	"context"
	"fmt"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jmoiron/sqlx"
	"root-histoty-service/config"
	"root-histoty-service/internal/model"
)

type UserRepo struct {
	db  *sqlx.DB
	cfg *config.DBConfig
}

func NewUserRepo(db *sqlx.DB, cfg *config.DBConfig) *UserRepo {
	return &UserRepo{
		db:  db,
		cfg: cfg,
	}
}

func (u *UserRepo) SaveUser(ctx context.Context, player model.Player) error {
	query := `
				INSERT INTO player(id, name, avatar, registered_at, pin_code)
				VALUES (:id, :name, :avatar, :registered_at, :pin_code)
			 `

	_, err := u.db.NamedExecContext(ctx, query, &player)
	if err != nil {
		return fmt.Errorf("failed to insert user: %w", err)
	}

	return nil
}

func (u *UserRepo) IsExists(ctx context.Context, name string) (bool, error) {
	params := map[string]interface{}{
		"name": name,
	}

	rows, err := u.db.NamedQuery(`SELECT * FROM player WHERE player.name=:name`, params)

	if err != nil {
		return false, fmt.Errorf("failed to check if player with name %s already exists: %w", name, err)
	}
	return rows.Next(), nil
}

func (u *UserRepo) GetUserByName(ctx context.Context, name string) (*model.Player, error) {
	exists, err := u.IsExists(ctx, name)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, fmt.Errorf("player with name %s does not exists", name)
	}

	params := map[string]interface{}{
		"name": name,
	}

	rows, err := u.db.NamedQuery(`SELECT * FROM player WHERE player.name=:name`, params)

	if err != nil {
		return nil, fmt.Errorf("failed to get player with name %s : %w", name, err)
	}

	var player model.Player

	for rows.Next() {
		err := rows.StructScan(&player)
		if err != nil {
			return nil, fmt.Errorf("failed to parse player with name %s : %w", name, err)
		}
	}
	return &player, nil
}

func (u *UserRepo) GetUserById(ctx context.Context, id string) (*model.Player, error) {

	return nil, nil
}
