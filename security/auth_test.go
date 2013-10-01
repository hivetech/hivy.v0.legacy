package security


import (
    . "launchpad.net/gocheck"
    "testing"
)


// Hook up gocheck into the "go test" runner.
func Test(t *testing.T) { TestingT(t) }

type SecuritySuite struct{}
var _ = Suite(&SecuritySuite{})


func (s *SecuritySuite) TestDecodeCredentials(c *C) {
    hash := "eGF2OmJvc3M="
    user, password, err := decodeCredentials(hash)
    c.Assert(err, IsNil)
    c.Check(user, Equals, "xav")
    c.Check(password, Equals, "boss")

    incorrect_hash := hash[6:]
    user, password, err = decodeCredentials(incorrect_hash)
    c.Assert(err, NotNil)
    c.Check(user, Equals, "")
    c.Check(password, Equals, "")
}
