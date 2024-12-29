package config

var defaultConfig = map[string]any{
	"auth.access_subject":  AccessTokenSubject,
	"auth.refresh_subject": RefreshTokenSubject,
}
