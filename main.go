package main

import (
	"GoGameApp/repository/mysql"
	"GoGameApp/service/userservice"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

func main() {

	mux := http.NewServeMux()
	mux.HandleFunc("/health", healthCheckHandler)
	mux.HandleFunc("/users/register", userRegisterHandler)

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

	mysqlRepo := mysql.New()
	userSvc := userservice.New(mysqlRepo)

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
