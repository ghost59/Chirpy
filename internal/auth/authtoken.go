package auth

import ("github.com/golang-jwt/jwt/v5"
		"github.com/google/uuid"
		"time"
		"net/http"
		"strings"
		"errors"
		"crypto/rand"
		"encoding/hex"
		)


func MakeJWT(userID uuid.UUID, tokenSecret string, expiresIn time.Duration) (string, error){
	times := time.Now().UTC().Add(expiresIn)
	
	token := jwt.NewWithClaims(jwt.SigningMethodHS256,jwt.RegisteredClaims{
		Issuer: "chirpy-access",
		IssuedAt: jwt.NewNumericDate(time.Now().UTC()),
		ExpiresAt:jwt.NewNumericDate(times),
		Subject: userID.String(),
	})

	tokenString, err := token.SignedString([]byte(tokenSecret))
	return tokenString, err
}


func ValidateJWT( tokenString, tokenSecret string) (uuid.UUID, error){
	mytoke := &jwt.RegisteredClaims{}
	_, err := jwt.ParseWithClaims(tokenString,mytoke,func(token *jwt.Token) (any, error){
		return []byte(tokenSecret), nil
	}  )
	if err != nil {
		return uuid.Nil, err
	}
	userID, err := uuid.Parse(mytoke.Subject)

	if err != nil {
		return uuid.Nil, err
	}

	return userID, nil	 
}

func GetBearerToken(headers http.Header) (string, error){
  authHeader := headers.Get("Authorization")
  if authHeader == ""{
  	return "", errors.New("No authorization header included")
  }

  splitAuth := strings.Fields(authHeader)

  if len(splitAuth) < 2 {
  	return "", errors.New("malformed authorization header")
  }

  if splitAuth[0] != "Bearer"{
  	return "", errors.New("malformed authorization header")
  }
  
  return splitAuth[1], nil
}


func MakeRefreshToken() string {
	key := make([]byte, 32)
	rand.Read(key)
	encodedStr := hex.EncodeToString(key)

	return encodedStr
	
}


func GetAPIKey(headers http.Header) (string, error){
authHeader := headers.Get("Authorization")
  if authHeader == ""{
  	return "", errors.New("No authorization header included")
  }

  splitAuth := strings.Fields(authHeader)

  if len(splitAuth) < 2 {
  	return "", errors.New("malformed authorization header")
  }

  if splitAuth[0] != "ApiKey"{
  	return "", errors.New("malformed authorization header")
  }
  
  return splitAuth[1], nil
}
