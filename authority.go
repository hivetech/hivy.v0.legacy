// Authority package that takes care of user's https requests authentification
//
// This package builds a RESTful server checking user's credentials before
// processing the given callback at the given url path. The Register() function
// instructs the package this parameters. Login and password are provided
// through standard http mechanism and currently verified in mysql database
// after some base64 decoding.
package main

import (
	"github.com/emicklei/go-restful"
    "fmt"
    "strings"
    "encoding/base64"
)

// When registered to the authority server, user-defined callback function
// are processed when "path" is reached by authentified requests
// Example:
//      authority.Register("/hello", func() {fmt.Println("Hello world")})
func Register(path string, callback restful.RouteFunction) {
    ws := new(restful.WebService)
    ws.Path(path).
        // Json and xml answer format accepted for i/o
        Consumes(restful.MIME_XML, restful.MIME_JSON).
        Produces(restful.MIME_JSON, restful.MIME_XML)

    //FIXME This design prevents from custom parameters. However user and project could be standard options to set
    // Create pipeline pipeline get-request -> authentification -> callback
    ws.Route(ws.GET("/{parameter}").Filter(basicAuthenticate).To(callback))
    restful.Add(ws)
}

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

// Intermediate step that will check encoded credentials before processing the received request.
// This function is explicitely used in Register() as a filter in the request pipeline.
func basicAuthenticate(req *restful.Request, resp *restful.Response, chain *restful.FilterChain) {
    encoded := req.Request.Header.Get("Authorization")[6:]

    username, passwd, err := decodeCredentials(encoded)
    if err != nil {
        log.Errorf("Credentials decoding failed (%v)", encoded)
		resp.AddHeader("WWW-Authenticate", "Basic realm=Protected Area")
		resp.WriteErrorString(402, "402: Error decoding credentials")
		return
    }
    log.Infof("User %s trying to connect with %s\n", username, passwd)

    ok, err := CheckCredentials(username, encoded)
    if err != nil {
        log.Errorf("[basicAuthenticate] %v", err)
		resp.AddHeader("WWW-Authenticate", "Basic realm=Protected Area")
		resp.WriteErrorString(403, "403: Error checking credentials")
        return 
    }
    if ! ok {
        log.Warningf("Authentification failed (%v)", encoded)
		resp.AddHeader("WWW-Authenticate", "Basic realm=Protected Area")
		resp.WriteErrorString(401, "401: Not Authorized")
		return
	}
    log.Infof("Authentification granted, processing (%v)", encoded)
	chain.ProcessFilter(req, resp)
}
