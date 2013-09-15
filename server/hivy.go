package main

import (
    "os"
    "os/signal"
    "fmt"
    "flag"
    //"os/exec"
    "./state_engine"

    "github.com/coreos/etcd/store"
    "github.com/coreos/go-etcd/etcd"
    "strings"

    "io/ioutil"
    "launchpad.net/goyaml"
    "path/filepath" 

    "os/exec"

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

type authority struct {
}



type persistence struct {
	serv_port string
	serv_ip string
	serv_name string
	serv_certificate string
	killer chan bool
}

func NewPersistence() *persistence {
	return &persistence{}
}

func (p *persistence) InitPersistence(name string, port string, ip string, cert string) {
	if port != "" {
		p.serv_port = port
	}
	if ip != "" {
		p.serv_ip = ip
	}
	if name != "" {
		p.serv_name = name
	}
	if cert != "" {
		p.serv_certificate = cert
	}
	//init
	//local := make(chan string)
	cmd := exec.Command("./utils","run",p.serv_name)
	//go func() {}
	fmt.Println("run :",p.serv_name)
	cmd.Start()

	//cmd = exec.Command("sleep","10")
	//cmd.Run()

	cmd = exec.Command("./utils","stop",p.serv_name)
	cmd.Run()
	cmd = exec.Command("./utils","init",p.serv_name,p.serv_port,p.serv_ip)
	cmd.Run()
	fmt.Println("init:",p.serv_name," done")
	killer = make(chan bool)
}

func (p *persistence) On() {
	cmd := exec.Command("./utils","run",p.serv_name)
	cmd.Start()
}
func (p *persistence) Off() {
	cmd := exec.Command("./utils","stop",p.serv_name)
	cmd.Run()
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

	//name string, port int,  ip string cert string,
	var hivy = NewPersistence()
	hivy.InitPersistence("serv_hivy","4001","127.0.0.1","")



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

	go func() {
		hivy.On()
		<-c // u need it to keep the etcd serv alive
	}()
	
	receiver := make(chan *store.Response)
	prompt := make(chan string)
	var machine = state_engine.NewState_engine()
	machine.InitState_engine()

	go func() {
		for {
			mess := <- prompt
			fmt.Println(mess)
			cont := strings.Fields(mess)
			if cont[0] == "SET"{
				req := strings.Fields(strings.Replace(cont[1], "/"," ",-1))
				machine.Dispatcher(req,cont[2])
			}



		}
	}()

	//prompt access for multiple process
	//from mess define the state status of machine
	

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
	hivy.Off()
}

