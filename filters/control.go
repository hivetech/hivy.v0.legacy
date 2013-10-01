package filters


import (
    "fmt"
    "net/http"
    "path/filepath"
    "strings"

	"github.com/emicklei/go-restful"
    "github.com/coreos/go-etcd/etcd"

    "github.com/hivetech/hivy/security"
)


func EtcdControl(request *restful.Request, response *restful.Response, chain *restful.FilterChain) {
    username, _, _ := security.Credentials(request)
    method := fmt.Sprintf("%s:%s", request.Request.Method, request.Request.URL)
    param_less_method := strings.Split(method, "?")[0]

    etcd.OpenDebug()
    defer etcd.CloseDebug()
    storage := etcd.NewClient()

    result, err := storage.Get(filepath.Join("hivy/security", username, "methods", param_less_method))
    fmt.Println(result)
    if err != nil {
        log.Errorf("[controlGate] %v\n", err)
        response.WriteError(http.StatusUnauthorized, err)
        return
    }
    //TODO Something to do with value ? Just 1 for now
    //permission := result[0].Value

    log.Infof("Method allowed, processing (%s:%s)", username, method)
	chain.ProcessFilter(request, response)
}
