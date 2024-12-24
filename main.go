package main

import (
	"GoGameApp/repository/mysql"
	"GoGameApp/service/authservice"
	"GoGameApp/service/userservice"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
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

	mux := http.NewServeMux()
	mux.HandleFunc("/health", healthCheckHandler)
	mux.HandleFunc("/users/register", userRegisterHandler)
	mux.HandleFunc("/users/login", userLoginHandler)
	mux.HandleFunc("/users/profile", userProfileHandler)

	fmt.Println("http server running on port 8080...")
	server := http.Server{Addr: ":8080", Handler: mux}
	err := server.ListenAndServe()
	if err != nil {
		panic(err)
	}
}

func userRegisterHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		fmt.Fprint(w, `{ "error": "invalid method"}`)
	}

	data, err := io.ReadAll(r.Body)
	if err != nil {
		w.Write([]byte(fmt.Sprintf(`{ "error": "%s"}`, err.Error())))

		return
	}

	var req userservice.RegisterRequest
	err = json.Unmarshal(data, &req)
	if err != nil {
		w.Write([]byte(fmt.Sprintf(`{ "error": "%s"}`, err.Error())))

		return
	}

	authSvc := authservice.New(JWTSignKey, AccessTokenSubject, RefreshTokenSubject,
		AccessTokenExpireDuration, RefreshTokenExpireDuration)

	mysqlRepo := mysql.New()
	userSvc := userservice.New(mysqlRepo, authSvc)

	_, err = userSvc.Register(req)
	if err != nil {
		w.Write([]byte(fmt.Sprintf(`{ "error": "%s"}`, err.Error())))

		return
	}

	w.Write([]byte(`{ "message": "user created" }`))
}

func healthCheckHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, `{ "message": "All Good!"}`)
}

func userLoginHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		fmt.Fprint(w, `{ "error": "invalid method"}`)
	}

	data, err := io.ReadAll(r.Body)
	if err != nil {
		w.Write([]byte(fmt.Sprintf(`{ "error": "%s"}`, err.Error())))

		return
	}

	var req userservice.LoginRequest
	err = json.Unmarshal(data, &req)
	if err != nil {
		w.Write([]byte(fmt.Sprintf(`{ "error": "%s"}`, err.Error())))

		return
	}

	authSvc := authservice.New(JWTSignKey, AccessTokenSubject, RefreshTokenSubject,
		AccessTokenExpireDuration, RefreshTokenExpireDuration)

	mysqlRepo := mysql.New()
	userSvc := userservice.New(mysqlRepo, authSvc)

	resp, err := userSvc.Login(req)
	if err != nil {
		w.Write([]byte(fmt.Sprintf(`{ "error": "%s"}`, err.Error())))

		return
	}

	data, err = json.Marshal(resp)
	if err != nil {
		w.Write([]byte(fmt.Sprintf(`{ "error": "%s"}`, err.Error())))

		return
	}

	w.Write(data)
}

func userProfileHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		fmt.Fprint(w, `{ "error": "invalid method"}`)
	}

	authSvc := authservice.New(JWTSignKey, AccessTokenSubject, RefreshTokenSubject,
		AccessTokenExpireDuration, RefreshTokenExpireDuration)

	jwtToken := r.Header.Get("Authorization")
	claims, err := authSvc.ParseToken(jwtToken)
	if err != nil {
		fmt.Fprintf(w, `{"error": "invalid access token"}`)

		return
	}

	req := userservice.ProfileRequest{UserID: claims.UserID}
	mysqlRepo := mysql.New()
	userSvc := userservice.New(mysqlRepo, authSvc)

	resp, err := userSvc.Profile(req)
	if err != nil {
		w.Write([]byte(fmt.Sprintf(`{ "error": "%s"}`, err.Error())))

		return
	}

	data, err := json.Marshal(resp)
	if err != nil {
		w.Write([]byte(fmt.Sprintf(`{ "error": "%s"}`, err.Error())))

		return
	}

	w.Write(data)
}
