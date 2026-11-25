package auth

import (
	"testing"
	"github.com/google/uuid"
	"time"
	)

func TestMake(t *testing.T) {
	user := uuid.New()
	secret := "kfdsjfkdsjflksdjflkjsdkflz"

	token, err := MakeJWT(user, secret, time.Minute * 5)
	if err != nil {
		t.Fatalf("unpexted error: %v", err)
	}
	valid, err := ValidateJWT(token, secret)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if valid != user {
		t.Errorf("got %v, want %v", valid, user)
	}

}


func Testfail(t *testing.T) {
	user := uuid.New()
	secret := "kfdsjfkdsjflksdjflkjsdkflz"

	token, err := MakeJWT(user, secret, time.Minute * 01)
	if err != nil {
		t.Fatalf("unpexted error: %v", err)
	}
	valid, err := ValidateJWT(token, secret)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if valid != user {
		t.Errorf("got %v, want %v", valid, user)
	}
	
	
}



 func Testsecret(t *testing.T) {
	user := uuid.New()
	secret := "kfdsjfkdsjflksdjflkjsdkflz"
	boo := "ksfl;kf;ldzkf"
	token, err := MakeJWT(user, secret, time.Minute * 01)
	if err != nil {
		t.Fatalf("unpexted error: %v", err)
	}
	valid, err := ValidateJWT(token, boo)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if valid != user {
		t.Errorf("got %v, want %v", valid, user)
	}
	
	
}



 
