package commands

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"code.cloudfoundry.org/cli/plugin"
	"github.com/go-openapi/runtime"
	"github.com/go-openapi/runtime/client"
	"github.com/golang-jwt/jwt/v5"
)

// DefaultTokenFactory factory for retrieving tokens
type DefaultTokenFactory struct {
	cliConnection   plugin.CliConnection
	cachedToken     string
	cachedTokenTime int64
}

// NewDefaultTokenFactory creates default token factory
func NewDefaultTokenFactory(cliConnection plugin.CliConnection) *DefaultTokenFactory {
	return &DefaultTokenFactory{
		cliConnection: cliConnection,
	}
}

// NewToken retrives outh token
func (t *DefaultTokenFactory) NewToken() (runtime.ClientAuthInfoWriter, error) {
	rawToken, err := t.NewRawToken()
	return client.BearerToken(rawToken), err
}

func (t *DefaultTokenFactory) NewRawToken() (string, error) {
	var expirationTime int64
	if t.cachedToken != "" {
		var err error
		expirationTime, err = getTokenExpirationTime(t.cachedToken)
		if err != nil {
			return "", err
		}
	}
	currentTimeInSeconds := time.Now().Unix()
	expirationTime = (expirationTime - currentTimeInSeconds) / 2
	if currentTimeInSeconds-t.cachedTokenTime >= expirationTime {
		tokenString, err := t.cliConnection.AccessToken()
		if err != nil {
			return "", fmt.Errorf("Could not get access token: %s", err)
		}
		t.cachedTokenTime = currentTimeInSeconds
		t.cachedToken = getTokenValue(tokenString)
	}
	return t.cachedToken, nil
}

func getTokenExpirationTime(tokenString string) (int64, error) {
	// Parse jwt token string
	parser := jwt.NewParser(jwt.WithoutClaimsValidation())
	token, err := parser.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return token, nil
	})

	if err != nil && !strings.Contains(err.Error(), "key is of invalid type") {
		return 0, err
	}

	// Try to get token expiration time
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return 0, errors.New("Could not read token claims")
	}
	if value, ok := claims["exp"]; ok {
		return int64(value.(float64)), nil
	}
	return 0, errors.New("Could not get token exipiration time")
}
