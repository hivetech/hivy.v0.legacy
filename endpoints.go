package hivy

import (
	"launchpad.net/loggo"
	"github.com/bitly/go-simplejson"
	"github.com/emicklei/go-restful"
)

var log = loggo.GetLogger("hivy.endpoint")

// JSON converts string to simple.Json object
func JSON(data string) *simplejson.Json {
	json, err := simplejson.NewJson([]byte(data))
	if err != nil {
		panic(err)
	}
	return json
}

// EmptyJSON is a commodity shortcut
func EmptyJSON() *simplejson.Json {
	return JSON("{}")
}

/*
 * Note:
 * Route{method:   restful.WebService.GET,
 *       path:     "GET /path/{with}/{parameters},
 *       endpoint: Dummy
 * }
 */

// Dummy is the hello-world endpoint, for demo and test purpose
func Dummy(request *restful.Request, response *restful.Response) {
	response.WriteEntity(JSON(`{"you": "dummy"}`))
}
