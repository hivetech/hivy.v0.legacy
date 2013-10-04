package main


import (
    "testing"
    //"time"

    "launchpad.net/loggo"
    "github.com/remogatto/prettytest"
)


// Start of setup
type testSuite struct {
    prettytest.Suite
}

func TestRunner(t *testing.T) {
    prettytest.RunWithFormatter(
        t,
        new(prettytest.TDDFormatter),
        new(testSuite),
    )
}
// End of setup


/*
 *func (t *testSuite) TestRunEtcd() {
 *    stop := make(chan bool)
 *    verbose := true
 *    force := false
 *    go RunEtcd(stop, "master", "node", force, verbose, "profile")
 *    time.Sleep(3 * time.Second)
 *    //TODO Check if etcd process is running
 *    //TODO Check if profile file and node directory exists
 *    stop <- true
 *}
 */


func (t *testSuite) TestLogger() {
    //NOTE Logfile is not implemented yet
    filename := ""
    verbose := true
    not_verbose := false

    err := SetupLog(verbose, filename)
    t.Nil(err)
    err = SetupLog(not_verbose, filename)
    t.Nil(err)
    loggo.RemoveWriter("hivy.main")
}


func (t *testSuite) TestVersion() {
    version := Version()
    t.Equal("0.1.0", version)
}
