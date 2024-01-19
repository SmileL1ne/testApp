package entity

type User struct {
	Name       string `json:"name"`
	Surname    string `json:"surname"`
	Patronymic string `json:"patronymic"`
}

type UserAge struct {
	Age int `json:"age"`
}

type UserGender struct {
	Gender string `json:"gender"`
}

type UserNationality struct {
	Country []struct {
		CountryID string `json:"country_id"`
	} `json:"country"`
}

type UserExtraInfo struct {
	Age     int
	Gender  string
	Country string
}

type UserDto struct {
	ID          int
	Name        string
	Surname     string
	Patronymic  string
	Age         int
	Gender      string
	Nationality string
}

type PaginatedUserData struct {
	Users    []UserDto
	Page     int
	PageSize int
}
