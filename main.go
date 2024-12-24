package main

import (
	"GoGameApp/config"
	"GoGameApp/delivery/httpserver"
	"GoGameApp/repository/mysql"
	"GoGameApp/service/authservice"
	"GoGameApp/service/userservice"
	"time"
)

const (
	JWTSignKey                 = "very_secret_key"
	AccessTokenSubject         = "at"
	RefreshTokenSubject        = "rt"
	AccessTokenExpireDuration  = time.Hour * 24
	RefreshTokenExpireDuration = time.Hour * 24 * 7
)

func main() {

	cfg := config.Config{
		HTTPServer: config.HTTPServer{Port: 8080},
		Auth: authservice.Config{
			SignKey:               JWTSignKey,
			AccessSubject:         AccessTokenSubject,
			RefreshSubject:        RefreshTokenSubject,
			AccessExpirationTime:  AccessTokenExpireDuration,
			RefreshExpirationTime: RefreshTokenExpireDuration,
		},
		Mysql: mysql.Config{
			User:                 "root",
			Passwd:               "root",
			Net:                  "tcp",
			Addr:                 "localhost:3306",
			DBName:               "gameapp_local",
			AllowNativePasswords: true,
		},
	}

	authSvc, userSvc := setupServices(cfg)
	server := httpserver.New(cfg, authSvc, userSvc)

	server.Serve()
}

// func userLoginHandler(w http.ResponseWriter, r *http.Request) {
// 	if r.Method != http.MethodPost {
// 		fmt.Fprint(w, `{ "error": "invalid method"}`)
// 	}

// 	data, err := io.ReadAll(r.Body)
// 	if err != nil {
// 		w.Write([]byte(fmt.Sprintf(`{ "error": "%s"}`, err.Error())))

// 		return
// 	}

// 	var req userservice.LoginRequest
// 	err = json.Unmarshal(data, &req)
// 	if err != nil {
// 		w.Write([]byte(fmt.Sprintf(`{ "error": "%s"}`, err.Error())))

// 		return
// 	}

// 	authSvc := authservice.New(JWTSignKey, AccessTokenSubject, RefreshTokenSubject,
// 		AccessTokenExpireDuration, RefreshTokenExpireDuration)

// 	mysqlRepo := mysql.New()
// 	userSvc := userservice.New(mysqlRepo, authSvc)

// 	resp, err := userSvc.Login(req)
// 	if err != nil {
// 		w.Write([]byte(fmt.Sprintf(`{ "error": "%s"}`, err.Error())))

// 		return
// 	}

// 	data, err = json.Marshal(resp)
// 	if err != nil {
// 		w.Write([]byte(fmt.Sprintf(`{ "error": "%s"}`, err.Error())))

// 		return
// 	}

// 	w.Write(data)
// }

// func userProfileHandler(w http.ResponseWriter, r *http.Request) {
// 	if r.Method != http.MethodGet {
// 		fmt.Fprint(w, `{ "error": "invalid method"}`)
// 	}

// 	authSvc := authservice.New(JWTSignKey, AccessTokenSubject, RefreshTokenSubject,
// 		AccessTokenExpireDuration, RefreshTokenExpireDuration)

// 	jwtToken := r.Header.Get("Authorization")
// 	claims, err := authSvc.ParseToken(jwtToken)
// 	if err != nil {
// 		fmt.Fprintf(w, `{"error": "invalid access token"}`)

// 		return
// 	}

// 	req := userservice.ProfileRequest{UserID: claims.UserID}
// 	mysqlRepo := mysql.New()
// 	userSvc := userservice.New(mysqlRepo, authSvc)

// 	resp, err := userSvc.Profile(req)
// 	if err != nil {
// 		w.Write([]byte(fmt.Sprintf(`{ "error": "%s"}`, err.Error())))

// 		return
// 	}

// 	data, err := json.Marshal(resp)
// 	if err != nil {
// 		w.Write([]byte(fmt.Sprintf(`{ "error": "%s"}`, err.Error())))

// 		return
// 	}

// 	w.Write(data)
// }

func setupServices(cfg config.Config) (authservice.Service, userservice.Service) {
	authSvc := authservice.New(cfg.Auth)
	mysqlRepo := mysql.New(cfg.Mysql)
	userSvc := userservice.New(mysqlRepo, authSvc)

	return authSvc, userSvc
}
