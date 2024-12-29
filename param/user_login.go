package param

type LoginRequest struct {
	PhoneNumber string `json:"phone_number"`
	Password    string `json:"password"`
}

type LoginResponse struct {
	Tokens Tokens   `json:"tokens"`
	User   UserInfo `json:"user"`
}
