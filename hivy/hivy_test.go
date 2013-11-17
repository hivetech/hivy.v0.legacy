package main

import (
    "time"
    "net/http"
    "io/ioutil"
    "strings"

    "launchpad.net/gocheck"
    "github.com/mreiferson/go-httpclient"
)

const url string = "http://127.0.0.1:8080"

func setupHivy() {
    //TODO Run automaticall hivy and etcd
    go hivy(url, false)
}

func sendHTTPRequest(method, endpoint, user, pass string) (*http.Response, error){
    transport := &httpclient.Transport{
        ConnectTimeout:        1*time.Second,
        RequestTimeout:        10*time.Second,
        ResponseHeaderTimeout: 5*time.Second,
    }
    defer transport.Close()

    client := &http.Client{Transport: transport}
    prefix := "/v0/actions"
    req, _ := http.NewRequest(method, url + prefix + endpoint, nil)
    req.SetBasicAuth(user, pass)

    //resp, err := client.Get(req)
    return client.Do(req)
}

func (t *testSuite) TestHivyDummy() {
    resp, err := sendHTTPRequest("GET", "/dummy", "xav", "boss")
    defer resp.Body.Close()
    t.Check(err, gocheck.IsNil)

    contents, err := ioutil.ReadAll(resp.Body)
    t.Check(err, gocheck.IsNil)

    t.True(strings.Contains(string(contents), "dummy"))
    //t.Equal(string(contents), "dummy")
}

func (t *testSuite) TestBadAuthentification() {
    resp, err := sendHTTPRequest("GET", "/dummy", "wrong", "login")
    defer resp.Body.Close()
    t.Check(err, gocheck.IsNil)

    contents, err := ioutil.ReadAll(resp.Body)
    t.Check(err, gocheck.IsNil)

    t.True(strings.Contains(string(contents), "Key Not Found"))
    //t.Equal(string(contents), "Key Not Found")
}
