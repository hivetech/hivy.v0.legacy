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
    is_ok, err := dbCallback(username, hash)
    c.Assert(err, IsNil)
    c.Check(is_ok, Equals, true)

    incorrect_hash := hash[2:]
    is_ok, err = dbCallback(username, incorrect_hash)
    c.Assert(err, IsNil)
    c.Check(is_ok, Equals, false)

    incorrect_username := "CheckNorris"
    is_ok, err = dbCallback(incorrect_username, hash)
    c.Assert(err, NotNil)
    c.Check(is_ok, Equals, false)
}


func (s *SecuritySuite) TestMysqlCheck(c *C) {
    hash := "eGF2OmJvc3M="
    testCheck(MysqlCheckCredentials, hash, c)
}


func (s *SecuritySuite) TestEtcdCheck(c *C) {
    hash := "boss"
    testCheck(EtcdCheckCredentials, hash, c)
}
