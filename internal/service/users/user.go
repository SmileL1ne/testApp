package users

import (
	"context"
	"net/http"
	"strconv"
	"testApp/internal/entity"
	"testApp/internal/repository/users"
)

type UserService interface {
	GetAll(context.Context, string, string) (entity.PaginatedUserData, error)
	Save(context.Context, entity.User) (int, error)
	Update(context.Context, string, entity.User) (int, error)
	Delete(context.Context, string) (int, error)
}

type userService struct {
	userRepo users.UserRepository
}

func NewUserService(r users.UserRepository) *userService {
	return &userService{
		userRepo: r,
	}
}

func (us *userService) GetAll(ctx context.Context, pageStr string, pageSizeStr string) (entity.PaginatedUserData, error) {
	page := 1
	pageSize := 10

	if pageStr != "" {
		pageInt, err := strconv.Atoi(pageStr)
		if err == nil && pageInt > 0 {
			page = pageInt
		}
	}
	if pageSizeStr != "" {
		pageSizeInt, err := strconv.Atoi(pageSizeStr)
		if err == nil && pageSizeInt > 0 {
			pageSize = pageSizeInt
		}
	}

	users, err := us.userRepo.GetAll(ctx, page, pageSize)
	if err != nil {
		return entity.PaginatedUserData{}, err
	}

	paginatedUserData := entity.PaginatedUserData{
		Users:    users,
		Page:     page,
		PageSize: pageSize,
	}

	return paginatedUserData, nil
}

func (us *userService) Save(ctx context.Context, user entity.User) (int, error) {
	userExtraInfo, err := us.userRepo.GetExtraInfo(user.Name)
	if err != nil {
		return http.StatusInternalServerError, err
	}

	userDto := entity.UserDto{
		Name:        user.Name,
		Surname:     user.Surname,
		Patronymic:  user.Patronymic,
		Age:         userExtraInfo.Age,
		Gender:      userExtraInfo.Gender,
		Nationality: userExtraInfo.Country,
	}

	err = us.userRepo.Save(ctx, userDto)
	if err != nil {
		return http.StatusInternalServerError, err
	}

	return http.StatusOK, nil
}

func (us *userService) Update(ctx context.Context, idStr string, user entity.User) (int, error) {
	id, err := strconv.Atoi(idStr)
	if err != nil || id < 1 {
		return http.StatusNotFound, err
	}

	userExtraInfo, err := us.userRepo.GetExtraInfo(user.Name)
	if err != nil {
		return http.StatusInternalServerError, err
	}

	userDto := entity.UserDto{
		Name:        user.Name,
		Surname:     user.Surname,
		Patronymic:  user.Patronymic,
		Age:         userExtraInfo.Age,
		Gender:      userExtraInfo.Gender,
		Nationality: userExtraInfo.Country,
	}

	err = us.userRepo.Update(ctx, id, userDto)
	if err != nil {
		return http.StatusInternalServerError, err
	}

	return http.StatusOK, nil
}

func (us *userService) Delete(ctx context.Context, idStr string) (int, error) {
	id, err := strconv.Atoi(idStr)
	if err != nil || id < 1 {
		return http.StatusNotFound, err
	}

	err = us.userRepo.Delete(ctx, id)
	if err != nil {
		return http.StatusInternalServerError, err
	}

	return http.StatusOK, nil
}
