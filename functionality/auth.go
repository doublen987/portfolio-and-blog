package functionality

import (
	"errors"
	"time"

	"github.com/dgrijalva/jwt-go"
)

var (
	ErrUsernamePasswordNotFound = errors.New("username or password don't match")
	ErrCreatingJWT              = errors.New("can't create jwt token")
	ErrTokenNotValid            = errors.New("jwt token not valid")
	ErrParsingClaims            = errors.New("can't parse claims")
)

var jwtKey = []byte("OMEGALUL")

var users = map[string]string{
	"Nikola": "Nesovic",
}

type User struct {
	Username string `json:"username"`
}

type LoginOutput struct {
	Token     string    `json:"token"`
	ExpiresAt time.Time `json:"expiresAt"`
	username  string    `json:"username"`
}

type Claims struct {
	Username string `json:"username"`
	jwt.StandardClaims
}

func Login(username string, password string) (LoginOutput, error) {
	var out LoginOutput

	expectedPassword, ok := users[username]

	if !ok || expectedPassword != password {
		return out, ErrUsernamePasswordNotFound
	}

	expirationTime := time.Now().Add(24 * 60 * time.Minute)

	claims := &Claims{
		Username: username,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(jwtKey)
	out.Token = tokenString
	if err != nil {
		return out, ErrCreatingJWT
	}

	return out, nil
}

func Authenticate(tokenString string) (User, error) {

	var user User
	// Initialize a new instance of `Claims`
	claims := &Claims{}

	// Parse the JWT string and store the result in `claims`.
	// Note that we are passing the key in this method as well. This method will return an error
	// if the token is invalid (if it has expired according to the expiry time we set on sign in),
	// or if the signature does not match
	tkn, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return jwtKey, nil
	})
	if err != nil {
		return user, ErrParsingClaims
	}
	if !tkn.Valid {
		return user, ErrTokenNotValid
	}
	user.Username = claims.Username
	return user, nil
}
