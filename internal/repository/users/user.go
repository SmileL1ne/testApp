package users

import (
	"context"
	"fmt"
	"sync"
	"testApp/internal/entity"

	"github.com/jackc/pgx/v5"
)

const (
	urlAge     = "https://api.agify.io/?name="
	urlGender  = "https://api.genderize.io/?name="
	urlCountry = "https://api.nationalize.io/?name="
)

type UserRepository interface {
	GetAll(context.Context, int, int, string) ([]entity.UserDto, error)
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

func (r *userRepository) GetAll(ctx context.Context, page int, pageSize int, sorting string) ([]entity.UserDto, error) {
	query := fmt.Sprintf(`SELECT * FROM users ORDER BY %s LIMIT $1 OFFSET $2`, sorting)

	rows, err := r.DB.Query(ctx, query, pageSize, (page-1)*pageSize)
	if err != nil {
		return nil, err
	}

	var users []entity.UserDto

	for rows.Next() {
		var user entity.UserDto
		err := rows.Scan(&user.ID, &user.Name, &user.Surname, &user.Patronymic, &user.Age, &user.Gender, &user.Nationality)
		if err != nil {
			return nil, err
		}
		users = append(users, user)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return users, nil
}

func (r *userRepository) GetExtraInfo(name string) (entity.UserExtraInfo, error) {
	userAge := &entity.UserAge{}
	userGender := &entity.UserGender{}
	userCountry := &entity.UserNationality{}

	relations := map[string]any{
		urlAge + name:     &userAge,
		urlGender + name:  &userGender,
		urlCountry + name: &userCountry,
	}

	var wg sync.WaitGroup
	var errCh = make(chan error, len(relations))

	wg.Add(len(relations))
	for url, target := range relations {
		go func(url string, target any) {
			defer wg.Done()
			if err := r.fetchURL(url, target); err != nil {
				errCh <- err
			}
		}(url, target)
	}

	wg.Wait()
	close(errCh)

	if err := <-errCh; err != nil {
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

	_, err := r.DB.Exec(ctx, query,
		user.Name, user.Surname, user.Patronymic,
		user.Age, user.Gender, user.Nationality)

	if err != nil {
		return err
	}

	return nil
}

func (r *userRepository) Update(ctx context.Context, id int, user entity.UserDto) error {
	query := `UPDATE users SET name = $1, surname = $2, patronymic = $3, age = $4, 
		gender = $5, nationality = $6 WHERE id = $7`

	_, err := r.DB.Exec(ctx, query, user.Name, user.Surname, user.Patronymic,
		user.Age, user.Gender, user.Nationality, id)
	if err != nil {
		return err
	}

	return nil
}

func (r *userRepository) Delete(ctx context.Context, id int) error {
	query := "DELETE FROM users WHERE id=$1"

	_, err := r.DB.Exec(ctx, query, id)
	if err != nil {
		return err
	}

	return nil
}
