package hivy

import (
	"launchpad.net/gocheck"
)

func (t *testSuite) TestNewController() {
	debug := false
	c := NewController("john", debug)
	t.Equal(c.user, "john")
	t.Equal(c.maxMachines, 5)
}

func (t *testSuite) TestSetUser() {
	c := NewController("john", false)
	t.Equal(c.user, "john")
	c.SetUser("doe")
	t.Equal(c.user, "doe")
	t.Check(c.user, gocheck.Equals, "doe")
}
