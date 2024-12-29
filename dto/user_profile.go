package dto

type ProfileRequest struct {
	UserID uint
}

type ProfileResponse struct {
	User UserInfo `json:"user"`
}
