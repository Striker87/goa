package jwt

import (
	//"context"
	"crypto/rsa"
	"fmt"
	"io/ioutil"
	//"net/http"
	"os"
	"path/filepath"
	"time"

	jwtgo "github.com/dgrijalva/jwt-go"
	//"github.com/goadesign/goa"
	//"github.com/goadesign/goa/middleware/security/jwt"
	uuid "github.com/satori/go.uuid"
)

// LoadJWTPrivateKeys loads PEM encoded RSA private key.
func LoadJWTPrivateKey(keyPath string) (*rsa.PrivateKey, error) {
	keyFile, err := filepath.Glob(keyPath)
	if err != nil {
		return nil, err
	}
	pem, err := ioutil.ReadFile(keyFile[0])
	if err != nil {
		return nil, err
	}
	key, err := jwtgo.ParseRSAPrivateKeyFromPEM([]byte(pem))
	if err != nil {
		return nil, fmt.Errorf("failed to load key %s: %s", keyFile, err)
	}
	return key, nil
}

func CreateJWTToken(email string) (string, error) {
	token := GenerateJWT(email)
	privKey, err := LoadJWTPrivateKey(os.Getenv("JWT_PRIVATE_KEY"))
	if err != nil {
		panic(err)
	}
	return token.SignedString(privKey)
}

func GenerateJWT(email string) *jwtgo.Token {
	token := jwtgo.New(jwtgo.SigningMethodRS512)
	oneMonth := time.Now().Add(time.Duration(24*30) * time.Hour).Unix()
	uuid, _ := uuid.NewV4()
	token.Claims = jwtgo.MapClaims{
		"iss":        "Issuer",
		"aud":        "Audience",
		"exp":        oneMonth,
		"jti":        uuid.String(),
		"iat":        time.Now().Unix(),
		"nbf":        2,
		"sub":        "subject",
		"user.email": email,
	}
	return token
}