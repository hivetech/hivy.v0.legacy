package security

import (
	"encoding/base64"
	"fmt"
	"strings"

	"github.com/emicklei/go-restful"
)

//TODO Check this out: https://github.com/dotcloud/docker/blob/master/auth/auth.go

// Credentials (formatted as user:password) sent throug http are base64 encoded.
// This function takes it and returns originals username and password.
func decodeCredentials(encoded string) (string, string, error) {
	// Decode the original hash
	data, err := base64.StdEncoding.DecodeString(encoded)
	if err != nil {
		log.Errorf("[decodeCredentials] %v", err)
		return "", "", fmt.Errorf("[decodeCredentials] %v", err)
	}
	log.Debugf("%s => %s\n", encoded, string(data))

	// Separate user and password informations
	user := strings.Split(string(data), ":")[0]
	passwd := strings.Split(string(data), ":")[1]
	return user, passwd, nil
}

// Credentials extracts from the header the authentification hash, and decodes
// it to return username and password informations.
func Credentials(request *restful.Request) (string, string, error) {
	encoded := request.Request.Header.Get("Authorization")
	if len(encoded) > 6 {
		// [6:] extracts the hash
		return decodeCredentials(encoded[6:])
	}
	return "", "", fmt.Errorf("[credentials] No credentials found (%v)\n", encoded)
}
