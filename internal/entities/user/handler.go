package user

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/adminoid/vuego/pkg/logging"
	"github.com/julienschmidt/httprouter"
	"net/http"
)

type handler struct {
	logger     *logging.Logger
	repository RepositoryUser
}

type Handler interface {
	Register(router *httprouter.Router)
}

func NewHandler(repository RepositoryUser, logger *logging.Logger) *handler {
	return &handler{
		logger:     logger,
		repository: repository,
	}
}

func (h *handler) Register(router *httprouter.Router) {
	router.Handler(http.MethodGet, "/users", checkAuth(h.GetList))
	router.POST("/register", h.UserRegister)
	router.POST("/login", h.Login)
}

type CredentialsRegister struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (h *handler) UserRegister(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Content-Type", "application/json")

	var (
		userRegister CredentialsRegister
		user         User
	)

	err := json.NewDecoder(r.Body).Decode(&userRegister)
	if err != nil {
		fmt.Printf("json.Decode error is -- %v", err)
	}

	user.Name = userRegister.Name
	user.Email = userRegister.Email
	user.PasswordHash = HashAndSalt(userRegister.Password)

	dbResponse := h.repository.Create(context.TODO(), &user)

	fmt.Println(dbResponse)

	w.WriteHeader(http.StatusOK)
}

func (h *handler) GetList(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Content-Type", "application/json")

	all, err := h.repository.FindAll(context.TODO())
	if err != nil {
		http.Error(w, "Forbidden", 400)
	}

	allBytes, err := json.Marshal(all)
	if err != nil {
		http.Error(w, "Parsing json error", 500)
	}

	w.WriteHeader(http.StatusOK)
	_, err = w.Write(allBytes)
	if err != nil {
		http.Error(w, "write header error", 500)
	}
}

func (h *handler) Login(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Content-Type", "application/json")

	var u User
	err := json.NewDecoder(r.Body).Decode(&u)
	if err != nil {
		fmt.Printf("json.Decode error is -- %v", err)
	}

	jwtTokensByte, err := h.LoginCheck(u)
	if err != nil {
		fmt.Printf("login check error is -- %v", err)
	}

	jwtTokensJson, err := json.Marshal(jwtTokensByte)
	if err != nil {
		fmt.Printf("marshal error is -- %v", err)
	}

	w.WriteHeader(http.StatusOK)
	w.Write(jwtTokensJson)
}

func (h *handler) LoginCheck(u User) (map[string]string, error) {

	// TODO Check user here

	if u.Name != "Petr" || string(u.PasswordHash) != "pwd" {
		return map[string]string{}, errors.New("incorrect credentials")
	}

	validTokens, err := generateTokenPair()
	if err != nil {
		return map[string]string{}, errors.New("token generation error")
	}

	return validTokens, nil
}
