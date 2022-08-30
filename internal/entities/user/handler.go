package user

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/adminoid/vuego/pkg/logging"
	"github.com/dgrijalva/jwt-go"
	"github.com/julienschmidt/httprouter"
	"net/http"
	"time"
)

var mySignKey = []byte("@#$ASDf9324$@%#sdafBSDFRR$$@3493n3SDF")

type handler struct {
	logger     *logging.Logger
	repository RepositoryUser
}

type Handler interface {
	Register(router *httprouter.Router)
}

func NewHandler(repository RepositoryUser, logger *logging.Logger) Handler {
	return &handler{
		repository: repository,
		logger:     logger,
	}
}

func (h *handler) Register(router *httprouter.Router) {
	//router.Handler("/users", checkAuth(h.GetList))
	router.Handler(http.MethodGet, "/users", checkAuth(h.GetList))
	router.POST("/login", h.Login)
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
	w.Write(allBytes)
}

func (h *handler) Login(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Content-Type", "application/json")

	var u User
	json.NewDecoder(r.Body).Decode(&u)

	result, err := h.LoginCheck(u)
	if err != nil {
		fmt.Printf("error is -- %v", err)
	}

	//fmt.Println(result)

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(result))
}

func (h *handler) LoginCheck(u User) (string, error) {

	// TODO: Check user here

	if u.Name != "Petr" || u.Password != "pwd" {
		return "", errors.New("incorrect credentials")
	}

	validToken, err := h.GenerateJWT()
	if err != nil {
		return "", errors.New("token generation error")
	}

	return validToken, nil
}

func (h *handler) GenerateJWT() (string, error) {
	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)
	claims["exp"] = time.Now().Add(time.Minute * 1).Unix()
	tokenString, err := token.SignedString(mySignKey)
	if err != nil {
		h.logger.Fatal(err)
	}
	return tokenString, nil
}

func checkAuth(endpoint func(http.ResponseWriter, *http.Request)) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		if r.Header["Token"] != nil {

			token, err := jwt.Parse(r.Header["Token"][0], func(token *jwt.Token) (interface{}, error) {
				if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
					return nil, fmt.Errorf("there was an error")
				}
				return mySignKey, nil
			})

			if err != nil {
				fmt.Fprintf(w, err.Error())
			}

			if token.Valid {
				endpoint(w, r)
			}
		} else {

			fmt.Fprintf(w, "Not Authorized")
		}
	})
}
