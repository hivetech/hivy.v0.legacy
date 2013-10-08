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


func (t *testSuite) TestRegisterGet() {
    r := NewRouter(testFilter, testFilter, false)
    err := r.Map("GET", "test/{id}", testEndpoint)
    t.Equal(err, nil)

    // Bad http method
    err = r.Map("BOOM", "test/{id}", testEndpoint)
    t.Check(err, gocheck.NotNil)

    // Not "/" in the path, not allowed
    err = r.Map("PUT", "test", testEndpoint)
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

    // Bad format, missing last "/"
    //path, param, err = parsePath("/user")
    //t.Check(err, gocheck.NotNil)
    //t.Check(path, gocheck.Equals, "")
    //t.Check(param, gocheck.Equals, "")
}
