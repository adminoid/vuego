package user

import (
	"context"
	"encoding/json"
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
	router.POST("/auth/register", h.UserRegister)
	router.POST("/auth/login", h.Login)
	router.POST("/jwt/refresh", h.RefreshJwt)
}

type RtStruct struct {
	Rt string `json:"rt"`
}

func (h *handler) RefreshJwt(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Content-Type", "application/json")

	// todo:
	//	 1. get rt from request
	var rt RtStruct
	err := json.NewDecoder(r.Body).Decode(&rt)
	if err != nil {
		fmt.Printf("json.Decode error is -- %v", err)
	}

	//	 2. check token for expiring
	valid, err := checkToken(rt.Rt)
	if err != nil {
		http.Error(w, fmt.Sprintf("%v", err), 400)
		return
	}
	//fmt.Println(valid)

	//	 3. get user from db by rt
	if valid {
		user, err := h.repository.GetUserByRt(context.TODO(), rt.Rt)
		if err != nil {
			http.Error(w, fmt.Sprintf("%v", err), 400)
			return
		}

		fmt.Println(user.RefreshToken)

		//	 4. run generateTokenPair (will save new rt into db)
		pairMap, err := h.generateTokenPair(user.ID)
		if err != nil {
			http.Error(w, fmt.Sprintf("%v", err), 400)
			return
		}

		jwtTokensJson, err := json.Marshal(pairMap)
		if err != nil {
			http.Error(w, fmt.Sprintf("%v", err), 500)
			return
		}

		w.WriteHeader(http.StatusOK)
		_, err = w.Write(jwtTokensJson)
		if err != nil {
			http.Error(w, fmt.Sprintf("%v", err), 500)
			return
		}
	}

	// 	 5. return new pair
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
		http.Error(w, fmt.Sprintf("%v", err), 400)
		return
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
	_, err = w.Write(jwtTokensJson)
	if err != nil {
		http.Error(w, fmt.Sprintf("%v", err), 500)
		return
	}
}
