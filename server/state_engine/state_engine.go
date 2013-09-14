package state_engine

import "fmt"

type input struct {
	User string
	Project string
	Target []string //example : 1 row : series : base, N rows : package : [a,b,c] c-plugin : [d,e,f]
	Value string
	State string
}

type action func(i *input) bool

//****//

func projectconfigsetted(i *input) bool {
	//TODO : check that no error raised during the 
	// config setting
	fmt.Println("TODO : projectconfigsetted %s\n",i)
	return true
}

func jujuisconfigured(i *input) bool {
	//TODO : check if juju is bootstraped and running
	fmt.Println("TODO : jujuisconfigured %s\n",i)
	return true
}

func doesprevstateisconfigured(i *input) bool {
	//TODO : check if the prev Value from etcd status 
	//
	fmt.Println("TODO : doesprevstateisconfigured %s\n",i)
	return true
}

func run(i *input) bool {
	
}

type Condition struct {
	utop map[string]map[string]bool //user/project map
	ttov map[string]map[string]bool //target/value map
	abilities map[string]map[string]bool //indicate active functionnalities like services
}

func NewCondition() *Condition{
	return &Condition{}
}

func (c *Condition) InitCondition() {
	c.utop = make(map[string]map[string]bool)
	c.ttov = make(map[string]map[string]bool)
	c.abilities = make(map[string]map[string]bool)
}

// yes yes, this is 3 times the same stuff, but I need to keep it explicit
// mapmapmapmapmapmamapappppaammmmad
func (c *Condition) Addutop (user string,project string) {
	if user != "" && project != "" {
		if c.utop[user] == nil {
			c.utop[user] = make(map[string]bool)
		}
		c.utop[user][project] = true
	}
}

func (c *Condition) Addttov (target string,value string) {
	if target != "" && value != "" {
		if c.ttov[target] == nil {
			c.ttov[target] = make(map[string]bool)
		}
		c.ttov[target][value] = true
	}
}

func (c *Condition) Addability (domaine string,name string) {
	if domaine != "" && name != "" {
		if c.abilities[domaine] == nil {
			c.abilities[domaine] = make(map[string]bool)
		}
		c.abilities[domaine][name] = true
	}
}

type State struct {
	c *Condition
	a []action	
}

func NewState() *State {
	return &State{}
}

func (s *State) InitState() {
	s.c = NewCondition()
	s.c.InitCondition()
}

func (s *State) Addutop (u string,p string) {
	s.c.Addutop(u,p)
}

func (s *State) Addttov (t string,v string) {
	s.c.Addttov(t,v)
}

func (s *State) Addability (d string,n string) {
	s.c.Addability(d,n)
}

func (s *State) Addaction(fp action) {
	s.a = append(s.a,fp)
}

func init() {

	var setconfig = NewState()
	setconfig.InitState()
	setconfig.Addutop("user","project")
	setconfig.Addutop("patate","coin")
	setconfig.Addttov("state","loaded")
	setconfig.Addaction(projectconfigsetted)

	var jujurunner = NewState()
	jujurunner.InitState()
	jujurunner.Addutop("user","project")
	jujurunner.Addability("service","mysql")
	jujurunner.Addability("service","wordpress")
	jujurunner.Addaction(jujuisconfigured)


}

//to get the current state and actions liked to it
func (c *Context) State_location() string {
	//what about a map[string]map[string]bool/int{} ? this could make the matching more efficient
	users := []string{"patate","user"}
	projects := []string{"project","pomme","coin"}	
	//NB : very easy to extract and build from yaml
	utop := map[string][]string{}
	utop[users[0]]=[]string{projects[2]}
	utop[users[1]]=[]string{projects[0]}
	
	if utop[c.User] != nil {
		for i:= range utop[c.User] {
			if utop[c.User][i] == c.Project {
				//matching
				//some state stuff here...damn, fix the state_engine
				if c.Target[0] == "state" && c.Value == "loaded" { return "configloaded" }
				if c.Target[0] == "state" && c.Value == "run" { return "running" }
			}
		}
	}
	return ""
}


func () {
//TODO add buffer effect
	for {
		mess := <- prompt
		fmt.Println(mess)
		cont := strings.Fields(mess)
		if cont[0] == "SET" {
			//well, cos I cant' manage fucking split dumbdumb
			ext := strings.Fields(strings.Replace(cont[1], "/"," ",-1))
			ct := context.Context{}
			ct.User = ext[0]
			ct.Project = ext[1]
			ct.Target = ext[2:]
			if ext[len(ext)-1] == "state" {
				if cont[2] != "" { ct.Value = cont[2] }
			}
			//THIS IS THE STATE_ENGINE WATCHER
			switch ct.State_location() { // NB we can have previous info from the etcd server with the prevalue state
				case "configloaded": //NB fallthrough could be use for sequences
					//set watcher for run/value = True
					fmt.Println("configloaded")

				case "running":
					// $cell, $series from config
					// juju deploy --repo=HERE local:$series/$cell $user-$project-$cell
					//TODO test if charm is already deployed
					fmt.Println("running")

					//TODO test here from etcd answer that the previous state was loaded!!
					//REST logik : need to reload conf for each request
					conf := jdc.JujuDeployConfig{}
					err := conf.GetEtcdConfig(&ct)// TODO force read only for ct!!
					if err != nil {
						fmt.Println(err)
					}

					//conf.Packages are used to custom the hivelab
					//find a way to say to ansible "put this shit 
					//in the next container

					if SETTINGS.Jujupath == " " { SETTINGS.Jujupath = "" } // "\ " forbidden
					cmd := exec.Command(SETTINGS.Jujupath+"juju","deploy",
					"--repository="+SETTINGS.Cellpath,
					"local:"+conf.Series+"/"+conf.Cell,
					ct.User+"-"+ct.Project+"-"+conf.Cell)
				
					err = cmd.Run()
					if err != nil { fmt.Println(err) }

					//myTiny dbhash
					availableservices := map[string]bool{}//could meme unix mod/own
					availableservices["mysql"]=true
					availableservices["wordpress"]=true

					for i:= range conf.Services {
						if availableservices[conf.Services[i]] {
							cmd = exec.Command(SETTINGS.Jujupath+"juju","deploy",
							/* no repo for services at MVP */
							conf.Services[i],
							ct.User+"-"+ct.Project+"-"+conf.Services[i])
						
							err = cmd.Run()
							if err != nil { fmt.Println(err) }
						}
					}

					//TODO
					//when services rdy, connect them to hivelab
					//and between them if needed
					//well, hivelab need some hooks for connection


				default:
					fmt.Println("can't reach the state or the state is not registered on action purpose")
			}
		}
    	}
}
