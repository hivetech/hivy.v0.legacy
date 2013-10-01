package main


import (
    "fmt"

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
    a := NewAuthority(testFilter, testFilter)
    a.RegisterGET("/test/{id}", testEndpoint)
}
