//Package hivy - REST framework
//============================
//
// Copyright 2013 Xavier Bruhiere
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//         http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
//
//
// Gate between user requests and Hive jobs. The client sends authentified
// http requests to reach endpoints defined in the endpoints directory.
// Through the process a centralized configuration server, backed by etcd,
// stores user-defined and hive settings, and a central authority is asked for
// methods permissions. Login and password are provided through standard
// http mechanism and currently verified in etcd database after some base64
// decoding.
//
// Usage example:
//      $ go run hivy.hivy.go --verbose --listen 0.0.0.0:8080
// Client usage example
//      $ GET -C name:pass http://localhost:8080/login/
//      $ curl --user name:pass http://127.0.0.1:8080/juju/deploy?project=trading
//      $ python -m "import requests; requests.delete('http://127.0.0.1:8080/user', data={'user': 'Chuck'}, auth=('user', 'pass')'"
package hivy
