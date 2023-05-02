package db

import (
	"database/sql"

	"hamza.achi/auth/models"
)

func (db Database) CreateUser(User *models.User) error {
	var id int

	stm, err := db.Conn.Prepare("INSERT INTO users (username, password) VALUES ($1, $2) RETURNING id")

	if err != nil {
		return err
	}

	defer stm.Close()

	err = stm.QueryRow(User.Username, User.Password).Scan(&id)
	if err != nil {
		return err
	}

	User.ID = id

	return nil
}

func (db Database) GetUserById(Username string) (models.User, error) {

	user := models.User{}
	stm, err := db.Conn.Prepare("SELECT * FROM users WHERE username = $1;")

	if err != nil {
		return user, err
	}

	defer stm.Close()

	row := stm.QueryRow(Username)
	switch err := row.Scan(&user.ID, &user.Username, &user.Password); err {
	case sql.ErrNoRows:
		return user, ErrNoMatch
	default:
		return user, err
	}

}
