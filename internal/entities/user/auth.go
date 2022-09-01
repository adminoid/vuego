package user

import (
	"crypto/rand"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"golang.org/x/crypto/scrypt"
	"io"
	"log"
	"net/http"
	"time"
)

var mySignKey = []byte("@#$ASDf9324$@%#sdafBSDFRR$$@3493n3SDF")

const (
	PwSaltBytes = 32
	PwHashBytes = 64
)

func HashAndSalt(password string) []byte {
	salt := make([]byte, PwSaltBytes)
	_, err := io.ReadFull(rand.Reader, salt)
	if err != nil {
		log.Fatal(err)
	}

	hash, err := scrypt.Key([]byte(password), salt, 1<<14, 8, 1, PwHashBytes)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("%x\n", hash)

	return hash
}

func generateTokenPair() (map[string]string, error) {
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

	return map[string]string{
		"access_token":  t,
		"refresh_token": rt,
	}, nil
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
