package entities

type Login struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type RegisterData struct {
	Username   string `json:"username"`
	Password   string `json:"password"`
	Name       string `json:"name"`
	Surname    string `json:"surname"`
	Patronymic string `json:"patronymic"`
}
