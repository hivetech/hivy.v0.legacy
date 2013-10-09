package main

import (
	"testing"

	"github.com/remogatto/prettytest"
	"launchpad.net/gocheck"
	"launchpad.net/loggo"
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

func (t *testSuite) TestLogger() {
	//NOTE Logfile is not implemented yet
	filename := ""
	verbose := true
	notVerbose := false

	err := SetupLog(verbose, filename)
	t.Nil(err)
	err = SetupLog(notVerbose, filename)
	t.Nil(err)
	loggo.RemoveWriter("hivy.main")
}

func (t *testSuite) TestVersion() {
	version := Version()
	t.Equal("0.1.0", version)
}

func (t *testSuite) TestAllTheSame() {
	testArray := []string{"same", "same", "same"}
	common, err := allTheSame(testArray)
	t.Check(err, gocheck.IsNil)
	t.Check(common, gocheck.Equals, "same")

	testWrongArray := []string{"same", "different", "same"}
	common, err = allTheSame(testWrongArray)
	t.Check(err, gocheck.NotNil)
	t.Check(common, gocheck.Equals, "")
}
