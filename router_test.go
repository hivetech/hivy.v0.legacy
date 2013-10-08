package main


import (
    "fmt"

    "launchpad.net/gocheck"
	"github.com/emicklei/go-restful"
)


func testFilter(req *restful.Request, resp *restful.Response, chain *restful.FilterChain) {
    fmt.Println("Filtering...")
	chain.ProcessFilter(req, resp)
}


func testEndpoint(req *restful.Request, resp *restful.Response) {
     resp.WriteEntity("{you: dummy}")
}


func (t *testSuite) TestMap() {
    r := NewRouter(testFilter, testFilter, false)
    err := r.Map("GET test/{id}", testEndpoint)
    t.Equal(err, nil)

    // Bad http method
    err = r.Map("BOOM test/{id}", testEndpoint)
    t.Check(err, gocheck.NotNil)

    // Not "/" in the path, not allowed
    err = r.Map("PUT test", testEndpoint)
    t.Check(err, gocheck.NotNil)
}


func (t *testSuite) TestParsePath() {
    path, param, err := parsePath("user/{id}/{name}")
    t.Check(err, gocheck.IsNil)
    t.Check(path, gocheck.Equals, "/user")
    t.Check(param, gocheck.Equals, "/{id}/{name}")

    // Bad format, missing last "/"
    path, param, err = parsePath("user")
    t.Check(err, gocheck.NotNil)
    t.Check(path, gocheck.Equals, "")
    t.Check(param, gocheck.Equals, "")

    // However "/path" is allowed
    path, param, err = parsePath("/security")
    t.Check(err, gocheck.IsNil)
    t.Check(path, gocheck.Equals, "/security")
    t.Check(param, gocheck.Equals, "")

    path, param, err = parsePath("/security/{credentials}/{token}")
    t.Check(err, gocheck.IsNil)
    t.Check(path, gocheck.Equals, "/security")
    t.Check(param, gocheck.Equals, "/{credentials}/{token}")
}


func (t *testSuite) TestParseRequest() {
    httpMethod, path, err := parseRequest("GET help/")
    t.Check(err, gocheck.IsNil)
    t.Check(httpMethod, gocheck.Equals, "GET")
    t.Check(path, gocheck.Equals, "help/")

    httpMethod, path, err = parseRequest("GET help/ useless")
    t.Check(err, gocheck.NotNil)
    t.Check(httpMethod, gocheck.Equals, "")
    t.Check(path, gocheck.Equals, "")
}


func (t *testSuite) TestSetupProfiling() {
    // At instanciation, profiling is setup if true
    var profiling = true
    r := NewRouter(testFilter, testFilter, profiling)
    t.True(r.profiling)
}


func (t *testSuite) TestCheckMapping() {
    //TODO Factorize userMap as testMap in testSuite ?
    var userMap = map[string]restful.RouteFunction{
        "PUT user/": testEndpoint,
        "DELETE user/": testEndpoint,
    }
    rootPath, err := checkMapping(userMap)
    t.Check(err, gocheck.IsNil)
    t.Check(rootPath, gocheck.Equals, "/user")

    var wrongUserMap = map[string]restful.RouteFunction{
        "PUT user/": testEndpoint,
        "DELETE other/": testEndpoint,
    }
    rootPath, err = checkMapping(wrongUserMap)
    t.Check(err, gocheck.NotNil)
    t.Check(rootPath, gocheck.Equals, "")
}


func (t *testSuite) TestMultiMap() {
    r := NewRouter(testFilter, testFilter, false)
    var userMap = map[string]restful.RouteFunction{
        "PUT user/": testEndpoint,
        "DELETE user/": testEndpoint,
    }
    err := r.MultiMap(userMap)
    t.Equal(err, nil)
}
