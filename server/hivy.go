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
    //"./server"
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
	go func() {
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
	}()

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

