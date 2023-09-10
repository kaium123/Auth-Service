package repository

import (
	"auth/common/logger"
	"auth/models"
	"database/sql"
	"errors"

	"github.com/lib/pq"
)

type UserRepositoryInterface interface {
	Register(user models.User) (int, error)
	SignIn(signInInfo models.SignInData) error
	FindByEmail(email string) (*models.User, error)
	UpdateProfile(user *models.User) error
	UpdateWebsites(urls []string, userID int) error
	FindByID(id int) (*models.User, error)
	RequestSent(userID int, requestedID int) error
	RequestAccept(userID int, requestedID int) error
	ManageConnection(userID int, friendID int) error
	ViewFriends(userID int) ([]*models.User, error)
	IsAlreadyRequestSent(userID int, requestedID int) error
	IsAlreadyRequestAccepter(userID int, requestedID int) error
}

type UserRepository struct {
	Db     *sql.DB
	logger logger.LoggerInterface
}

func NewUserRepository(Db *sql.DB, logger logger.LoggerInterface) UserRepositoryInterface {
	return &UserRepository{Db: Db, logger: logger}
}

func (r *UserRepository) Register(user models.User) (int, error) {
	var lastInsertedID int
	query := "INSERT INTO users (name, email, password) VALUES ($1, $2, $3) RETURNING id"

	err := r.Db.QueryRow(query, user.Name, user.Email, user.Password).Scan(&lastInsertedID)
	if err != nil {
		return 0, err
	}

	return lastInsertedID, nil
}

func (r *UserRepository) SignIn(signInInfo models.SignInData) error {
	return nil
}

func (r *UserRepository) FindByEmail(email string) (*models.User, error) {

	logger.LogInfo(email)
	var user models.User
	err := r.Db.QueryRow(`
		SELECT id, email, password, name, COALESCE(user_name,''), COALESCE(phone,''), COALESCE(bio,''), COALESCE(gender,'')
		FROM users
		WHERE email = $1`, email).Scan(&user.ID, &user.Email, &user.Password, &user.Name, &user.UserName, &user.Phone, &user.Bio, &user.Gender)
	if err != nil {
		logger.LogError(err.Error())
		if err.Error() == "sql: no rows in result set" {
			return nil, nil
		}
		return nil, err
	}
	logger.LogError(user)

	return &user, nil
}

func (r *UserRepository) UpdateProfile(user *models.User) error {
	_, err := r.Db.Exec("UPDATE users SET name = $1, email = $2, password = $3, user_name = $4, phone = $5, bio = $6, gender = $7 WHERE id = $8", user.Name, user.Email, user.Password, user.UserName, user.Phone, user.Bio, user.Gender, user.ID)
	if err != nil {
		return err
	}
	return nil
}

func (r *UserRepository) UpdateWebsites(urls []string, userID int) error {
	for _, url := range urls {
		_, err := r.Db.Exec("INSERT INTO websites (url, user_id) VALUES ($1, $2)", url, userID)
		if err != nil {
			return err
		}
	}
	return nil
}

func (r *UserRepository) FindByID(userID int) (*models.User, error) {
	var user models.User
	err := r.Db.QueryRow(`
		SELECT id, email, password, name, COALESCE(user_name,''), COALESCE(phone,''), COALESCE(bio,''), COALESCE(gender,'')
		FROM users
		WHERE id = $1`, userID).Scan(&user.ID, &user.Email, &user.Password, &user.Name, &user.UserName, &user.Phone, &user.Bio, &user.Gender)
	if err != nil && err.Error() == "sql: no rows in result set" {
		return nil, errors.New("user not found")
	}
	if err != nil {
		logger.LogError(err.Error())
		return nil, err
	}

	var urls []string

	rows, err := r.Db.Query(`
    	SELECT COALESCE(url,'')
    	FROM websites
    	WHERE user_id = $1`, userID)
	if err != nil {
		logger.LogError(err.Error())
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var url string
		if err := rows.Scan(&url); err != nil {
			logger.LogError(err.Error())
			return nil, err
		}

		user.Websites = append(user.Websites, url)
	}

	if err := rows.Err(); err != nil {
		logger.LogError(err.Error())
		return nil, err
	}

	logger.LogError(urls)

	return &user, nil
}

func (r *UserRepository) RequestSent(userID int, requestedID int) error {
	logger.LogError(requestedID, " ", userID)
	_, err := r.Db.Exec("INSERT INTO sent_requests (\"from\", \"to\") VALUES ($1, $2)", userID, requestedID)
	if err != nil {
		return err
	}
	return nil
}

func (r *UserRepository) RequestAccept(userID int, requestedID int) error {
	_, err := r.Db.Exec("INSERT INTO friends (\"from\", \"to\") VALUES ($1, $2)", requestedID, userID)
	if err != nil {
		return err
	}

	_, err = r.Db.Exec("DELETE FROM sent_requests WHERE \"from\" = $1 AND \"to\" = $2", requestedID, userID)
	if err != nil {
		return err
	}

	return nil
}

func (r *UserRepository) ManageConnection(userID int, friendID int) error {
	_, err := r.Db.Exec("DELETE FROM friends WHERE (\"from\" = $1 AND \"to\" = $2) OR (\"from\" = $3 AND \"to\" = $4)", friendID, userID, userID, friendID)
	if err != nil {
		return err
	}
	return nil
}
func (r *UserRepository) ViewFriends(userID int) ([]*models.User, error) {
	friendIDs := []int{}
	rows, err := r.Db.Query(`
        SELECT "from"
        FROM friends
        WHERE "to" = $1`, userID)
	if err != nil {
		logger.LogError(err.Error())
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var friendID int
		if err := rows.Scan(&friendID); err != nil {
			logger.LogError(err.Error())
			return nil, err
		}
		friendIDs = append(friendIDs, friendID)
	}

	if err := rows.Err(); err != nil {
		logger.LogError(err.Error())
		return nil, err
	}

	if len(friendIDs) == 0 {
		return []*models.User{}, nil
	}

	query := `
        SELECT id, COALESCE(name,''), COALESCE(user_name,''), COALESCE(phone,''), COALESCE(bio,''), COALESCE(gender,'')
        FROM users
        WHERE id = ANY($1);`

	rows, err = r.Db.Query(query, pq.Array(friendIDs))
	if err != nil {
		logger.LogError(err.Error())
		return nil, err
	}
	defer rows.Close()

	users := []*models.User{}

	for rows.Next() {
		var user models.User
		if err := rows.Scan(&user.ID, &user.Name, &user.UserName, &user.Phone, &user.Bio, &user.Gender); err != nil {
			logger.LogError(err.Error())
			return nil, err
		}
		users = append(users, &user)
	}

	if err := rows.Err(); err != nil {
		logger.LogError(err.Error())
		return nil, err
	}

	logger.LogError(users)

	return users, nil
}

func (r *UserRepository) IsAlreadyRequestSent(userID int, requestedID int) error {

	id := 0
	err := r.Db.QueryRow(`
		SELECT id
		FROM sent_requests
		WHERE "to" = $1 and "from" = $2`, requestedID, userID).Scan(&id)

	logger.LogInfo(id)

	if err != nil {
		logger.LogError(err.Error())
		if err == sql.ErrNoRows {
			return nil
		}
		return err
	}

	return errors.New("exists")
}

func (r *UserRepository) IsAlreadyRequestAccepter(userID int, requestedID int) error {

	id := 0
	err := r.Db.QueryRow(`
		SELECT id
		FROM friends
		WHERE ("from" = $1 and "to" = $2) OR ("from" = $3 and "to" = $4)`, requestedID, userID, userID, requestedID).Scan(&id)
	if err != nil {
		logger.LogError(err.Error())
		if err.Error() == "sql: no rows in result set" {
			return nil
		}
		return err
	}

	return errors.New("exists")
}
