/*
*	Copyright (C) 2013 Mathieu Dous
*/


package main

import (
    "os"
    "os/signal"
    "fmt"
    "flag"
    //"os/exec"
    "./state_engine"

    "github.com/coreos/etcd/store"
    "github.com/zabka/go-etcd/etcd"
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

type handler func(*etcd.Client, string, chan *store.Response, chan bool, chan string)

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


type Hive struct {
	edao *pers_cluster
	state_machine *state_engine.State_engine
}
func NewHive() *Hive {
	return &Hive{}
}
//func (h *Hive) InitHive() {
func (h *Hive) InitHive(t chan string, k chan bool, r chan *store.Response) {

	h.edao = NewPers_cluster()
	h.edao.InitPers_cluster()
	h.state_machine = state_engine.NewState_engine()
	h.state_machine.InitState_engine()
	
	/**************************************************/
	var public = NewPersistence()
	public.InitPersistence("serv_hivy","4000","127.0.0.1","")
	public.AddEvent("status")
	public.Connect(t,k,r)

	var user1 = NewPersistence()
	user1.InitPersistence("user","4002","127.0.0.1","")
	user1.AddEvent("user/project/state")
	user1.Connect(t,k,r)

	var user2 = NewPersistence()
	user2.InitPersistence("patate","4003","127.0.0.1","")
	user2.AddEvent("patate/coin/state")
	user2.Connect(t,k,r)


	var setconfig = state_engine.NewSituation()
	setconfig.InitSituation()
	setconfig.Addutop("user","project")
	setconfig.Addutop("patate","coin")
	setconfig.Addttov("state","loaded")
	setconfig.Addhook(setconfig.Projectconfigsetted)

	setconfig.AddClient("user",user1.client)
	setconfig.AddClient("patate",user2.client)

	var jujurunner = state_engine.NewSituation()
	jujurunner.InitSituation()
	jujurunner.Addutop("user","project")
	jujurunner.Addttov("state","run")
	jujurunner.Addability("services","hivelab")
	jujurunner.Addability("services","mysql")
	jujurunner.Addability("services","wordpress")
	buff := []string{"series","cell","version","services","packages"}
	jujurunner.Configbuffer("config",buff)
	jujurunner.Addhook(jujurunner.Jujuisconfigured) 
	jujurunner.Addhook(jujurunner.Doesprevstateisconfigured)
	jujurunner.Addhook(jujurunner.Run)

	jujurunner.AddClient("user",user1.client)
	jujurunner.AddClient("patate",user2.client)
	/**************************************************/

	h.edao.Addserv(public)
	h.edao.Addserv(user1)
	h.edao.Addserv(user2)

	h.state_machine.Addsituation(setconfig)
	h.state_machine.Addsituation(jujurunner)

}

type pers_cluster struct {
	cluster map[string]*persistence
}
func NewPers_cluster() *pers_cluster {
	return &pers_cluster{}
}
func (c *pers_cluster) InitPers_cluster() {
	c.cluster = make(map[string]*persistence) // should break here
}
func (c *pers_cluster) Addserv(p *persistence) {
	if p != nil && p.serv_port != "" {
		if c.cluster[p.serv_port] != nil {
			//there is already a living instance on this port
			return
		}
		c.cluster[p.serv_port] = p
	}
}
func (c *pers_cluster) Invocator(minion chan string) {
	alakazam := make(chan *persistence)
	for {
		port := <- minion
		fmt.Println(port)
		if c.cluster[port] != nil {
			go func() {
				if c.cluster[port].state == "ready" || c.cluster[port].state == "stopped"{
					c.cluster[port].On()
	
					for i,_ := range c.cluster[port].event_key {
						go c.cluster[port].Handler(i)
					}

					dispel := <- alakazam				
					for dispel != c.cluster[port] { 
						dispel = <- alakazam
					}
				} else if c.cluster[port].state == "started" {
					c.cluster[port].Off()
					alakazam <- c.cluster[port]
				}
			}()
		}
	}
}

type persistence struct {
	serv_port string
	serv_ip string
	serv_name string
	serv_certificate string
	state string
	client *etcd.Client
	pipe chan string
	kill chan bool
	receiv chan *store.Response
	event_key []string
	event_handler []handler
}
func NewPersistence() *persistence {
	return &persistence{state:"born"}
}
func (p *persistence) InitPersistence(name string, port string, ip string, cert string) {
	if port != "" { p.serv_port = port }
	if ip != "" { p.serv_ip = ip }
	if name != "" { p.serv_name = name }
	if cert != "" { p.serv_certificate = cert }
	//init
	//local := make(chan string)
	cmd := exec.Command("./utils","run",p.serv_name)
	fmt.Println("run :",p.serv_name)
	cmd.Start()
	cmd = exec.Command("./utils","stop",p.serv_name)
	cmd.Run()
	cmd = exec.Command("./utils","init",p.serv_name,p.serv_port,p.serv_ip)
	cmd.Run()
	fmt.Println("init:",p.serv_name," done")

	p.client = etcd.NewClient()
	p.client.InitClient(p.serv_ip,p.serv_port)

	p.state = "ready"
}
func (p *persistence) eventHandler(h handler) {
	if h != nil {
		p.event_handler = append(p.event_handler,h)
	}
}
func (p *persistence) AddEvent(e string) {
	if e != "" {
		p.event_key = append(p.event_key,e)
		p.eventHandler(watch) //default handler
	}
}
func (p *persistence) Handler(i int) {
	//watch(c.cluster[port].client, "/user/project/state", c.cluster[port].receiv , c.cluster[port].kill, c.cluster[port].pipe)
	p.event_handler[i](p.client, p.event_key[i], p.receiv, p.kill, p.pipe)
	<- p.kill
}
func (p *persistence) Connect(t chan string, k chan bool, r chan *store.Response) {
	p.pipe = t
	p.kill = k
	p.receiv = r
}
func (p *persistence) On() {
	cmd := exec.Command("./utils","run",p.serv_name)
	cmd.Start()
	p.state = "started"
}
func (p *persistence) Off() {
	cmd := exec.Command("./utils","stop",p.serv_name)
	cmd.Run()
	p.state = "stopped"
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
	receiver := make(chan *store.Response)
	fmt.Println("receiver",receiver)

	//name string, port int,  ip string cert string,

	//var STATE_ENGINE_CONFIG
	//var STATE_USER_CONFIG

	//exit method
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	stop := make(chan bool)
	go func() {
		<-c //stuck till c
		stop <- true //kill all whatchers
		os.Exit(0) //nothing to do here!
	}()

	//goto invocator
	//go func() {
	//	hivy.On()
	//	<-c // u need it to keep the etcd serv alive
	//}()
	
	var hivy = NewHive()
	
	
	prompt := make(chan string)
	
	hivy.InitHive(prompt,stop,receiver)
	//var machine = state_engine.NewState_engine()
	//machine.InitState_engine()

	//connect machine statement with the persistence handler
	//this looks archaic...actually, need to be more stronger
	go func() {
		for {
			mess := <- prompt
			fmt.Println("hello this is ",mess)
			cont := strings.Fields(mess)
			if cont[0] == "SET"{
				req := strings.Fields(strings.Replace(cont[1], "/"," ",-1))
				hivy.state_machine.Dispatcher(req,cont[2])
			}
		}
	}()

	//prompt access for multiple process
	//from mess define the state status of machine
	

	//here set the watched keys
	//TODO should be set from a user base & dynamically added
	//routines := []string{"patate/coin/state","user/project/state","an/other/state"}
	//init the bench of watcher
	//TODO set "personal" stop for each watcher to kill them dynamically

	


	//public.Connect(prompt,stop,receiver)

	//var hivy = NewPers_cluster()
	//hivy.InitPers_cluster()
	//hivy.Addserv(public)


	//client := public.client
	//go watch(client, "/user/project/state", receiver, stop, prompt)

	m := make(chan string)
	go hivy.edao.Invocator(m)
	m <- "4000"

	m <- "4002"

	m <- "4003"	


	//for i:=0;i<len(routines);i++ {
	//	go watch(client, routines[i], receiver, stop, prompt)
	//}
	//wait till the exit flag
	<-stop
	//public.Off()
}

