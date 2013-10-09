package security

import (
	"database/sql"
	"fmt"
	"path/filepath"

	// Requested by the lib
	"github.com/coreos/go-etcd/etcd"
	_ "github.com/go-sql-driver/mysql"
)

// MysqlCheckCredentials fetchs 'hive.login' table 'username' hash and, if found, compare
// it with the given one.
func MysqlCheckCredentials(username string, hash string) (bool, error) {
	//FIXME Hard coded
	db, err := sql.Open("mysql", "xavier:boss@/hive")
	if err != nil {
		return false, fmt.Errorf("[db.CheckCredentials::sql.Open] %v\n", err)
	}
	defer db.Close()

	// Open doesn't open a connection. Validate DSN data:
	if err = db.Ping(); err != nil {
		return false, fmt.Errorf("[db.CheckCredentials::db.Ping] %v\n", err)
	}

	// Prepare statement for reading data
	stmtOut, err := db.Prepare("SELECT hash FROM logins WHERE username = ?")
	if err != nil {
		return false, fmt.Errorf("[db.MysqlCheckCredentials::db.Prepare] %v\n", err)
	}
	defer stmtOut.Close()

	// Actually fetch back 'username' associated hash
	var hashFound string
	err = stmtOut.QueryRow(username).Scan(&hashFound)
	if err != nil {
		return false, fmt.Errorf("[db.MysqlCheckCredentials::QueryRow.Scan] %v\n", err)
	}
	return (hashFound == hash), nil
}

// EtcdCheckCredentials queries the etcd database to compare given and stored hashes.
func EtcdCheckCredentials(username, hash string, debug bool) (bool, error) {
	//TODO This is no longer hash but clear passwd for now
	if debug {
		etcd.OpenDebug()
		defer etcd.CloseDebug()
	}
	storage := etcd.NewClient()
	// Global settings
	response, err := storage.Get(filepath.Join("hivy/security", username, "password"))
	if err != nil || len(response) != 1 {
		return false, fmt.Errorf("[db.EtcdCheckCredentials::storage.Get] %v\n", err)
	}
	return (hash == response[0].Value), nil
}
