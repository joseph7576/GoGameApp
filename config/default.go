package config

var defaultConfig = map[string]any{
	"auth.access_subject":      AccessTokenSubject,
	"auth.refresh_subject":     RefreshTokenSubject,
	"auth.access_expire_time":  AccessTokenExpireDuration,
	"auth.refresh_expire_time": RefreshTokenExpireDuration,
}
