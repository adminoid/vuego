package user

import (
	"context"
	"errors"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"golang.org/x/crypto/bcrypt"
	"net/http"
	"time"
)

var mySignKey = []byte("@#$ASDf9324$@%#sdafBSDFRR$$@3493n3SDF")

func (h *handler) LoginCheck(u CredentialsLogin) (map[string]string, error) {

	user, err := h.repository.Get(context.TODO(), u.Email)
	if err != nil {
		return nil, err
	}

	res := bcrypt.CompareHashAndPassword(user.PasswordHash, []byte(u.Password))
	if res != nil {
		return nil, errors.New("wrong password")
	}

	validTokens, err := h.generateTokenPair(user.ID)
	if err != nil {
		return map[string]string{}, errors.New("token generation error")
	}

	return validTokens, nil
}

func (h *handler) generateTokenPair(userId string) (map[string]string, error) {
	// Create token
	token := jwt.New(jwt.SigningMethodHS256)

	// Set claims
	// This is the information which frontend can use
	// The backend can also decode the token and get admin etc.
	claims := token.Claims.(jwt.MapClaims)
	claims["sub"] = 1
	claims["name"] = "Jon Doe"
	claims["admin"] = true
	claims["exp"] = time.Now().Add(time.Minute * 15).Unix()

	// Generate encoded token and send it as response.
	// The signing string should be secret (a generated UUID works too)
	t, err := token.SignedString(mySignKey)
	if err != nil {
		return nil, err
	}

	refreshToken := jwt.New(jwt.SigningMethodHS256)
	rtClaims := refreshToken.Claims.(jwt.MapClaims)
	rtClaims["sub"] = 1
	rtClaims["exp"] = time.Now().Add(time.Hour * 24).Unix()

	rt, err := refreshToken.SignedString(mySignKey)
	if err != nil {
		return nil, err
	}

	err = h.repository.UpdateRefreshToken(context.TODO(), userId, rt)
	if err != nil {
		return nil, err
	}

	return map[string]string{
		"access_token":  t,
		"refresh_token": rt,
	}, nil
}

func checkAuth(endpoint func(http.ResponseWriter, *http.Request)) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		if r.Header["Token"] != nil {

			checkingResult, err := checkToken(r.Header["Token"][0])
			if err != nil {
				http.Error(w, fmt.Sprintf("%v", err), 500)
				return
			}
			if checkingResult {
				endpoint(w, r)
			}
		} else {
			http.Error(w, fmt.Sprintf("Not Authorized"), 403)
			return
		}
	})
}

func checkToken(tokenValue string) (bool, error) {
	token, err := jwt.Parse(tokenValue, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("there was an error")
		}
		return mySignKey, nil
	})

	if err != nil {
		return false, err
	}

	return token.Valid, err
}

func GetPasswordHash(password string) []byte {
	pwd := []byte(password)

	hashedPassword, err := bcrypt.GenerateFromPassword(pwd, bcrypt.DefaultCost)
	if err != nil {
		panic(err)
	}
	return hashedPassword
}
