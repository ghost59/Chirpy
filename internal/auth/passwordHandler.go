package auth

import ("github.com/alexedwards/argon2id"
		"log")

func HashPassword(password string) (string, error){
	hash, err := argon2id.CreateHash(password, argon2id.DefaultParams)
	if err != nil {
		log.Printf("%s", err)
		return "", err
	}

	return hash, err
}

func CheckPasswordHash(password, hash string) (bool, error){
	 match, err := argon2id.ComparePasswordAndHash(password, hash)
	if err != nil {
		log.Printf("%s", err)
		return false, err
	}


	return match, nil
	
}
