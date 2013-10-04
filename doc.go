// Router between user requests and Hive jobs. The client sends authentified
// http requests to reach endpoints defined in the endpoints directory.
// Through the process a centralized configuration server, backed by etcd,
// stores user-defined and hive settings, and a central authority is asked for
// methods permissions. Login and password are provided through standard
// http mechanism and currently verified in etcd database after some base64
// decoding.
//
// Usage example:
//      $ go run hivy --verbose --listen 0.0.0.0:8080
// Client usage example
//      $ GET -C name:pass http://localhost:8080/login/name
//      $ curl --user name:pass http://127.0.0.1:8080/login/name
//      $ python -m "import requests; requests.get('http://127.0.0.1:8080/login/name', auth=('user', 'pass')'"
package main
