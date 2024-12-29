package param

type ProfileRequest struct {
	UserID uint
}

type ProfileResponse struct {
	User UserInfo `json:"user"`
}
