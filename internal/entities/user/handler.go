package user

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/adminoid/vuego/pkg/logging"
	"github.com/julienschmidt/httprouter"
	"golang.org/x/crypto/bcrypt"
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
	user.PasswordHash = GetPasswordHash(userRegister.Password)

	dbResponse := h.repository.Create(context.TODO(), &user)

	fmt.Println(dbResponse)

	w.WriteHeader(http.StatusOK)
}

func (h *handler) GetList(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Content-Type", "application/json")

	all, err := h.repository.FindAll(context.TODO())
	if err != nil {
		http.Error(w, fmt.Sprintf("%v", err), 400)
		return
	}

	allBytes, err := json.Marshal(all)
	if err != nil {
		http.Error(w, fmt.Sprintf("%v", err), 500)
		return
	}

	w.WriteHeader(http.StatusOK)
	_, err = w.Write(allBytes)
	if err != nil {
		http.Error(w, fmt.Sprintf("%v", err), 500)
		return
	}
}

type CredentialsLogin struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (h *handler) Login(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Content-Type", "application/json")

	var u CredentialsLogin
	err := json.NewDecoder(r.Body).Decode(&u)
	if err != nil {
		http.Error(w, fmt.Sprintf("%v", err), 500)
		return
	}

	jwtTokensByte, err := h.LoginCheck(u)
	if err != nil {
		http.Error(w, fmt.Sprintf("%v", err), 500)
		return
	}

	jwtTokensJson, err := json.Marshal(jwtTokensByte)
	if err != nil {
		http.Error(w, fmt.Sprintf("%v", err), 500)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(jwtTokensJson)
}

func (h *handler) LoginCheck(u CredentialsLogin) (map[string]string, error) {

	// TODO Check user here
	user, err := h.repository.Get(context.TODO(), u.Email)
	if err != nil {
		return nil, err
	}

	res := bcrypt.CompareHashAndPassword(user.PasswordHash, []byte(u.Password))
	if res != nil {
		return nil, errors.New("wrong password")
	}

	validTokens, err := generateTokenPair()
	if err != nil {
		return map[string]string{}, errors.New("token generation error")
	}

	return validTokens, nil
}
