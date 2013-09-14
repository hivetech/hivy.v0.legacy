package main

import (
    "os"
    "os/signal"
    "fmt"
    "flag"
    "os/exec"

    jdc "./jujudeployconfig"
    "./context"

    //"github.com/codegangsta/cli"
    //"launchpad.net/loggo"

    "github.com/coreos/etcd/store"
    "github.com/coreos/go-etcd/etcd"
    "strings"


    "io/ioutil"
    "launchpad.net/goyaml"
    "path/filepath" 

)

/*
* Hivy manage the etcd base events by watcher
* combinaison.
*/

type settings struct {
	where string
	as string
	Cellpath string
	Jujupath string
}

func Newsettings() *settings {
	return &settings{where:"./",as:"settings.yaml",}
}

func (s *settings) GetYaml() (error) {
	src := filepath.Join(s.where,s.as)
	data, err := ioutil.ReadFile(src)
	if err != nil {
		if os.IsNotExist(err) {
			return fmt.Errorf("%s file can't be found in %s\n",s.as,s.where)
		}
		return err
	}
	err = goyaml.Unmarshal(data, &s)
	if err != nil {
		return fmt.Errorf("cannot parse %q: %v", src, err)
	}
	return nil
}

//from etcdctl advices
var watchFlag = flag.NewFlagSet("watch", flag.ExitOnError)
var index     = watchFlag.Int64("i", 0, "watch from the given index")

func watch(client *etcd.Client, key string, receiver chan *store.Response, stop chan bool, prompt chan string) {
	go client.Watch(key, uint64(*index), receiver, stop)
	for {
		resp := <-receiver
		prompt <- resp.Action+" "+resp.Key+" "+resp.Value
	}
}

//TODO add/remove dynamically keys, by watching input stream
func main() {
	//INIT
	var SETTINGS = Newsettings()
	err := SETTINGS.GetYaml()
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(SETTINGS)

	//var STATE_ENGINE_CONFIG
	//var STATE_USER_CONFIG

	client := etcd.NewClient()
	//exit method
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	stop := make(chan bool)
	go func() {
		<-c //stuck till c
		stop <- true //kill all whatchers
		os.Exit(0) //nothing to do here!
	}()
	
	receiver := make(chan *store.Response)
	prompt := make(chan string)

	//prompt access for multiple process
	//from mess define the state status of machine
	

/*







*/


	//here set the watched keys
	//TODO should be set from a user base & dynamically added
	routines := []string{"patate/coin/state","user/project/state","an/other/state"}
	//init the bench of watcher
	//TODO set "personal" stop for each watcher to kill them dynamically
	for i:=0;i<len(routines);i++ {
		go watch(client, routines[i], receiver, stop, prompt)
	}
	//wait till the exit flag
	<-stop
}

