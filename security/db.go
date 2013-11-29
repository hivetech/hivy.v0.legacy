package security

import (
  "fmt"
	"path/filepath"

	"github.com/coreos/go-etcd/etcd"
)

// EtcdCheckCredentials queries the etcd database to compare given and stored hashes.
func EtcdCheckCredentials(username, hash string, debug bool) (bool, error) {
	//TODO This is no longer hash but clear passwd for now
  sorted := true
	if debug {
		etcd.OpenDebug()
		defer etcd.CloseDebug()
	}
  //FIXME Should I use controller here ?
  machines := []string{"http://127.0.0.1:4001"}
  storage := etcd.NewClient(machines)
	// Global settings
	response, err := storage.Get(filepath.Join("hivy/security", username, "password"), sorted)
  fmt.Println(response)
  fmt.Println(username)
	if err != nil {
		return false, fmt.Errorf("[db.EtcdCheckCredentials::storage.Get] %v\n", err)
	}
	return (hash == response.Value), nil
}
