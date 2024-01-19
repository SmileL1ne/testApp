package users

import (
	"context"
	"testApp/internal/entity"

	"github.com/jackc/pgx/v5"
)

const (
	urlAge     = "https://api.agify.io/?name="
	urlGender  = "https://api.genderize.io/?name="
	urlCountry = "https://api.nationalize.io/?name="
)

type UserRepository interface {
	GetAll(context.Context, int, int) ([]entity.UserDto, error)
	Save(context.Context, entity.UserDto) error
	Update(context.Context, int, entity.UserDto) error
	Delete(context.Context, int) error
	GetExtraInfo(name string) (entity.UserExtraInfo, error)
}

type userRepository struct {
	DB *pgx.Conn
}

func NewUserRepository(db *pgx.Conn) *userRepository {
	return &userRepository{
		DB: db,
	}
}

func (r *userRepository) GetAll(ctx context.Context, page int, pageSize int) ([]entity.UserDto, error) {
	query := `SELECT * FROM users ORDER BY id LIMIT $1 OFFSET $2`

	stmt, err := r.DB.Prepare(ctx, "get_all_users", query)
	if err != nil {
		// r.logger.Error("preparing statement", "error", err.Error())
		return nil, err
	}

	rows, err := r.DB.Query(ctx, stmt.Name, pageSize, (page-1)*pageSize)
	if err != nil {
		// r.logger.Error("executing statement", "error", err.Error())
		return nil, err
	}

	var users []entity.UserDto

	for rows.Next() {
		var user entity.UserDto
		err := rows.Scan(&user.ID, &user.Name, &user.Surname, &user.Patronymic, &user.Age, &user.Gender, &user.Nationality)
		if err != nil {
			// r.logger.Error("scanning rows", "error", err.Error())
			// http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return nil, err
		}
		users = append(users, user)
	}

	if err := rows.Err(); err != nil {
		// r.logger.Error("reading from rows", "error", err.Error())
		// http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return nil, err
	}

	return users, nil
}

func (r *userRepository) GetExtraInfo(name string) (entity.UserExtraInfo, error) {
	userAge := &entity.UserAge{}
	userGender := &entity.UserGender{}
	userCountry := &entity.UserNationality{}

	if err := r.fetchURL(urlAge+name, &userAge); err != nil {
		return entity.UserExtraInfo{}, err
	}
	if err := r.fetchURL(urlGender+name, &userGender); err != nil {
		return entity.UserExtraInfo{}, err
	}
	if err := r.fetchURL(urlCountry+name, &userCountry); err != nil {
		return entity.UserExtraInfo{}, err
	}

	countryID := ""
	if len(userCountry.Country) != 0 {
		countryID = userCountry.Country[0].CountryID
	}

	userExtraInfo := entity.UserExtraInfo{
		Age:     userAge.Age,
		Gender:  userGender.Gender,
		Country: countryID,
	}

	return userExtraInfo, nil
}

func (r *userRepository) Save(ctx context.Context, user entity.UserDto) error {
	query := `INSERT INTO users (name, surname, patronymic, age, gender, nationality) 
	VALUES($1, $2, $3, $4, $5, $6)`

	stmt, err := r.DB.Prepare(ctx, "insert_user", query)
	if err != nil {
		// r.logger.Error("preparing statement", "error", err.Error())
		return err
	}

	_, err = r.DB.Exec(ctx, stmt.Name,
		user.Name, user.Surname, user.Patronymic,
		user.Age, user.Gender, user.Nationality)

	if err != nil {
		// r.logger.Error("executing statement", "error", err.Error())
		return err
	}

	return nil
}

func (r *userRepository) Update(ctx context.Context, id int, user entity.UserDto) error {
	query := `UPDATE users SET name = $1, surname = $2, patronymic = $3, age = $4, 
		gender = $5, nationality = $6 WHERE id = $7`

	stmt, err := r.DB.Prepare(ctx, "update_user", query)
	if err != nil {
		// r.logger.Error("preparing statement", "error", err.Error())
		return err
	}

	_, err = r.DB.Exec(ctx, stmt.Name, user.Name, user.Surname, user.Patronymic,
		user.Age, user.Gender, user.Nationality, id)
	if err != nil {
		// r.logger.Error("updating user", "error", err.Error())
		return err
	}

	return nil
}

func (r *userRepository) Delete(ctx context.Context, id int) error {
	query := "DELETE FROM users WHERE id=$1"

	stmt, err := r.DB.Prepare(ctx, "delete_user", query)
	if err != nil {
		// r.logger.Error("preparing statement", "error", err.Error())
		return err
	}

	_, err = r.DB.Exec(ctx, stmt.Name, id)
	if err != nil {
		// r.logger.Error("executing statement", "error", err.Error())
		return err
	}

	return nil
}
