package repository

import (
	"auth/common/logger"
	"auth/models"
	"database/sql"
)

type UserRepositoryInterface interface {
	Register(user models.User) error
	SignIn(signInInfo models.SignInData) error
	FindByEmail(email string) (*models.User, error)
	UpdateProfile(user *models.User) error
	UpdateWebsites(urls []string, userID int) error
	FindByID(id int) (*models.User, error)
}

type UserRepository struct {
	Db     *sql.DB
	logger logger.LoggerInterface
}

func NewUserRepository(Db *sql.DB, logger logger.LoggerInterface) UserRepositoryInterface {
	return &UserRepository{Db: Db, logger: logger}
}

func (r *UserRepository) Register(user models.User) error {
	_, err := r.Db.Exec("INSERT INTO users (name, email, password) VALUES ($1, $2, $3)", user.Name, user.Email, user.Password)
	if err != nil {
		return err
	}
	return nil
}

func (r *UserRepository) SignIn(signInInfo models.SignInData) error {
	return nil
}

func (r *UserRepository) FindByEmail(email string) (*models.User, error) {

	logger.LogInfo(email)
	var user models.User
	err := r.Db.QueryRow(`
		SELECT id, email, password, name, COALESCE(user_name,''), COALESCE(phone,''), COALESCE(bio,''), COALESCE(gender,''), COALESCE(profile_pic,'')
		FROM users
		WHERE email = $1`, email).Scan(&user.ID, &user.Email, &user.Password, &user.Name, &user.UserName, &user.Phone, &user.Bio, &user.Gender, &user.ProfilePic)
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
	_, err := r.Db.Exec("UPDATE users SET name = $1, email = $2, password = $3, user_name = $4, phone = $5, bio = $6, gender = $7, profile_pic = $8 WHERE id = $9", user.Name, user.Email, user.Password, user.UserName, user.Phone, user.Bio, user.Gender, user.ProfilePic, user.ID)
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
		SELECT id, email, password, name, COALESCE(user_name,''), COALESCE(phone,''), COALESCE(bio,''), COALESCE(gender,''), COALESCE(profile_pic,'')
		FROM users
		WHERE id = $1`, userID).Scan(&user.ID, &user.Email, &user.Password, &user.Name, &user.UserName, &user.Phone, &user.Bio, &user.Gender, &user.ProfilePic)
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
