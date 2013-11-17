package beacon

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

	err := SetupLog("test", verbose, filename)
	t.Nil(err)
	err = SetupLog("test", notVerbose, filename)
	t.Nil(err)
	loggo.RemoveWriter("hivy.main")
}

func (t *testSuite) TestVersion() {
	version := StableVersion()
	t.Equal(version.major, 0)
	t.Equal(version.minor, 1)
	t.Equal(version.fix, 5)
}

func (t *testSuite) TestVersionString() {
    version := Version{1, 2, 3}
	t.Equal(version.String(), "1.2.3")
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
