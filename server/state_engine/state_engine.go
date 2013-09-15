package state_engine

import "fmt"
import "errors"

//for the go API
import "github.com/coreos/go-etcd/etcd"
import "strconv"
import "strings"

//
import "os/exec"

type buffer struct {
	data map[string]map[string][]string
}

func NewBuffer() *buffer {
	return &buffer{}
}
func (b *buffer) InitBuffer() {
	b.data = make(map[string]map[string][]string)
}
func (b *buffer) setkey(scope string, key string) {
	if scope != "" && key != "" {
		if b.data[scope] == nil {
			b.data[scope] = make(map[string][]string)
		}
		b.data[scope][key] = []string{}
	}
}
func (b *buffer) setdata(scope string, key string, values []string) {
	if values == nil { return }
	if scope != "" && key != "" {
		if b.data[scope] == nil {
			b.data[scope] = make(map[string][]string)
		}
		for _, val := range values {
			b.data[scope][key] = append(b.data[scope][key], val)
		}
	}
}
func (b *buffer) getdata(scope string, key string) []string {
	if scope != "" && key != "" {
		return b.data[scope][key]
	}
	return []string{}
}
func (b *buffer) flush(scope string, key string) {
	if scope != "" && key != "" {
		delete(b.data[scope],key)
	} else if scope != "" {
		delete(b.data,scope)
	}
}

type input struct {
	User string
	Project string
	Target string
	Value string
}

func NewInput() *input {
	return &input{}
}

type hook func(i *input) bool

func (s *Situation) projectconfigsetted(i *input) bool {
	//ToDO : check that no error raised during the 
	// config setting
	fmt.Println("TODO : projectconfigsetted %s\n",i)
	return true
}

func (s *Situation) jujuisconfigured(i *input) bool {
	//ToDO : check if juju is bootstraped and running
	fmt.Println("TODO : jujuisconfigured %s\n",i)
	return true
}

func (s *Situation) doesprevstateisconfigured(i *input) bool {
	//ToDO : check if the prev Value from etcd status 
	//
	fmt.Println("TODO : doesprevstateisconfigured %s\n",i)
	return true
}

func (s *Situation) run(i *input) bool {
	fmt.Println("Hello I'm run\n")
	/***** CONFIGURATION GETTER FROM ETCD ****/
	scope := "config"
	begin := i.User+"/"+i.Project+"/"+scope+"/"
	end := "/"
	err := errors.New("")
	str := ""

	for key,_ := range s.b.data[scope] {
		str = begin+key+end
		s.b.data["config"][key], err = gets(str)
		if err != nil { fmt.Println(err) } //always return an error as index out of range
	}
	/************** GETTER DONE *****************/

	/********** SERVICE VALIDATION ***************/
	//deploy// 0 : cell, 1..end : services
	services := []string{}
	for _, cell := range s.b.data[scope]["cell"] {
		services = append(services,cell)
	}
	for _, serv := range s.b.data[scope]["services"] {
		services = append(services,serv)
	}
	//test if this situation allows such services
	for _,serv := range services {
		if s.What("services",serv) {
		} else { return false }
	}
	/************ SERVICE VALIDATION DONE ************/ 

	/***************** CELL CMD **********************/
	cmd := exec.Command("juju","deploy",
			    //"--repository="+SETTINGS.Cellpath,
			    "--repository=/home/supercoin/cells",
			    "local:"+s.b.data["config"]["series"][0]+"/"+services[0],
			    i.User+"-"+i.Project+"-"+services[0])	
	err = cmd.Run()
	if err != nil { fmt.Println(err) }
	/************* CELL CMD DONE *********************/
	
	/************* SERVICES DEPLOY ********************/
	for _,serv := range services[1:] { //0 is cell
		//cmd = exec.Command(SETTINGS.Jujupath+"juju","deploy",
		cmd = exec.Command("juju","deploy",				
				   serv,
				   i.User+"-"+i.Project+"-"+serv)
	
		err = cmd.Run()
		if err != nil { fmt.Println(err) }
	}	
	/************* SERVICE DEPLOY DONE ***************/
	return true

}

type condition struct {
	utop map[string]map[string]bool //event listener	user/project map
	ttov map[string]map[string]bool //condition	 	target/value map
	abilities map[string]map[string]bool //ability segment	indicate active functionnalities like services
}

func NewCondition() *condition{
	return &condition{}
}

func (c *condition) InitCondition() {
	c.utop = make(map[string]map[string]bool)
	c.ttov = make(map[string]map[string]bool)
	c.abilities = make(map[string]map[string]bool)
}

// yes yes, this is 3 times the same stuff, but I need to keep it explicit
// mapmapmapmapmapmamapappppaammmmad
func (c *condition) Addutop (user string,project string) {
	if user != "" && project != "" {
		if c.utop[user] == nil {
			c.utop[user] = make(map[string]bool)
		}
		c.utop[user][project] = true
	}
}

func (c *condition) Addttov (target string,value string) {
	if target != "" && value != "" {
		if c.ttov[target] == nil {
			c.ttov[target] = make(map[string]bool)
		}
		c.ttov[target][value] = true
	}
}

func (c *condition) Addability (domaine string,name string) {
	if domaine != "" && name != "" {
		if c.abilities[domaine] == nil {
			c.abilities[domaine] = make(map[string]bool)
		}
		c.abilities[domaine][name] = true
	}
}

type Situation struct {
	c *condition
	a []hook
	b *buffer	
}

func NewSituation() *Situation {
	return &Situation{}
}

func (s *Situation) InitSituation() {
	s.c = NewCondition()
	s.c.InitCondition()
	s.b = NewBuffer()
	s.b.InitBuffer()
}

func (s *Situation) Addutop (u string,p string) {
	s.c.Addutop(u,p)
}

func (s *Situation) Addttov (t string,v string) {
	s.c.Addttov(t,v)
}

func (s *Situation) Addability (d string,n string) {
	s.c.Addability(d,n)
}

func (s *Situation) Addhook(fp hook) {
	s.a = append(s.a,fp)
}

func (s *Situation) Configbuffer(scope string, keys []string) {
	for _, scopekey := range keys {
		s.b.setkey(scope,scopekey)
	}
}

type State_engine struct {
	State []*Situation
	CurIn *input //current input build from the request
}

func NewState_engine() *State_engine {
	return &State_engine{}
}

func (s *State_engine) InitState_engine() {
	//get situation desc and stuff from db
	//TODO how to get associated etcd??????

	var setconfig = NewSituation()
	setconfig.InitSituation()
	setconfig.Addutop("user","project")
	setconfig.Addutop("patate","coin")
	setconfig.Addttov("state","loaded")
	setconfig.Addability("persistence","user") //if username != servername then have to add a link table
	setconfig.Addability("persistence","patate")
	setconfig.Addhook(setconfig.projectconfigsetted)
	//setconfig.Addhook(setconfig.jujuisconfigured)

	var jujurunner = NewSituation()
	jujurunner.InitSituation()
	jujurunner.Addutop("user","project")
	jujurunner.Addttov("state","run")
	jujurunner.Addability("services","hivelab")
	jujurunner.Addability("services","mysql")
	jujurunner.Addability("services","wordpress")
	buff := []string{"series","cell","version","services","packages"}
	jujurunner.Configbuffer("config",buff)
	jujurunner.Addhook(jujurunner.jujuisconfigured) 
	jujurunner.Addhook(jujurunner.doesprevstateisconfigured)
	jujurunner.Addhook(jujurunner.run)

	s.Addsituation(setconfig)
	s.Addsituation(jujurunner)
}

func (s *State_engine) Addsituation(sit *Situation) {
	if sit != nil {
		s.State = append(s.State,sit)
	}
}

//this is the go routine which should be run on the main
func (s *State_engine) Dispatcher(req []string,value string) {
	//I don't mind the multiple validation of situation yet (which should result a dynamic merge of state or
	// the dynamic creation of a new one)

	var in = NewInput()

	/************ build input from request *******/	
	in.User = req[0]
	in.Project = req[1]
	in.Target = req[2]
	if in.Target == "state" {
		if value != "" {
			in.Value = value
		}
	} 
	/************ input rdy *******************/
	//TODO : find a proper way for Target[0]
	var situ *Situation
	for _,sit := range s.State {
		if sit.Who(in.User,in.Project) && sit.Where(in.Target,in.Value) {
			situ = sit
		}
	}
	if situ == nil {return} //no situation match, exit

	//when the sit is defined, exec his hooks with the current set of input
	for _, act := range situ.a {
		if act(in) == false { return } //if an hook raise a false, then stop the process
	}
	//here, everythings went ok , yeay!
}

func (s *Situation) Who(origin string, key string) bool {
	if origin != "" && key != "" {
		return s.c.utop[origin][key]
	}
	return false
}
func (s *Situation) Where(origin string, key string) bool {
	if origin != "" && key != "" {
		return s.c.ttov[origin][key]
	}
	return false
}
func (s *Situation) What(origin string, key string) bool {
	if origin != "" && key != "" {
		return s.c.abilities[origin][key]
	}
	return false
}

/************** ETCD API ***********************/

func get(key string) (string,error) {
	client := etcd.NewClient() //select the etcd relative to user
	resps, err := client.Get(key)
	if err != nil {
		return "",err
	}
	for _, resp := range resps {
		if resp.Value != ""{
			//return resp.Value, nil
			//fmt.Println(resp.Value)
			ans := strings.Replace(resp.Value," ","",1)
			return ans,nil
		}
	}
	return "",nil
}

// TODO : this function will get all index from a config value
func gets(cut_key string) ([]string,error) {
	resp := []string{}
	echo := ""
	err := errors.New("")
	i:=0
	key := cut_key+strconv.Itoa(i)
	for {
		echo, err = get(key)
		if i == 0 && err != nil {
			fmt.Println("no value indexed", err)
			return nil, err
		} else if err != nil {
			fmt.Println ("tail or panic", err)
			return resp, err
		}
		resp = append(resp, echo)
		i=i+1
		key = cut_key+strconv.Itoa(i)
	}
}

/*****************************************************/
