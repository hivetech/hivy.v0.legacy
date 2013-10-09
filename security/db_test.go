package security

import (
	. "launchpad.net/gocheck"
	//"testing"
)

// Hook up gocheck into the "go test" runner.
//func Test(t *testing.T) { TestingT(t) }

//type SecuritySuite struct{}
//var _ = Suite(&SecuritySuite{})

type dbBridgeTest func(username string, hash string) (bool, error)

//FIXME Mysql is still waiting for a hash, although etcd use plain password
func testCheck(dbCallback dbBridgeTest, hash string, c *C) {
	username := "xav"
	//FIXME The database is hardcoded in db.go !!
	isOK, err := dbCallback(username, hash)
	c.Assert(err, IsNil)
	c.Check(isOK, Equals, true)

	incorrectHash := hash[2:]
	isOK, err = dbCallback(username, incorrectHash)
	c.Assert(err, IsNil)
	c.Check(isOK, Equals, false)

	incorrectUsername := "CheckNorris"
	isOK, err = dbCallback(incorrectUsername, hash)
	c.Assert(err, NotNil)
	c.Check(isOK, Equals, false)
}

/*
 *func (s *SecuritySuite) TestMysqlCheck(c *C) {
 *    hash := "eGF2OmJvc3M="
 *    testCheck(MysqlCheckCredentials, hash, c)
 *}
 *
 *
 *func (s *SecuritySuite) TestEtcdCheck(c *C) {
 *    hash := "boss"
 *    testCheck(EtcdCheckCredentials, hash, c)
 *}
 */
