package jujudeployconfig

import "fmt"
import "github.com/coreos/go-etcd/etcd"
import "../context"

import "errors"
import "strconv"
//import "builtin"

//TODO : generic issue : build jdc from yaml file => heavy factorisable functions
type JujuDeployConfig struct {
	Series string
	Cell string
	Version string
	Services []string
	Packages []string
}

func get(key string) (string,error) {
	client := etcd.NewClient()
	resps, err := client.Get(key)
	if err != nil {
		return "",err
	}
	for _, resp := range resps {
		if resp.Value != ""{
			//return resp.Value, nil
			//fmt.Println(resp.Value)
			return resp.Value,nil
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

func (jdc *JujuDeployConfig) GetEtcdConfig(ct *context.Context) error {

	//what about filepath.JOIN
	begin := ct.User+"/"+ct.Project+"/"+"config"+"/"
	end := "/"+"0"

	//binding err here because of
	//https://groups.google.com/forum/#!topic/golang-nuts/ItYSNKqt_kA
	//damn it...
	err := errors.New("")

	str := begin+"series"+end
	jdc.Series, err = get(str)
	if err != nil { fmt.Println(err) }

	str = begin+"cell"+end
	jdc.Cell, err = get(str)
	if err != nil { fmt.Println(err) }

	str = begin+"version"+end
	jdc.Version, err = get(str)
	if err != nil { fmt.Println(err) }


	end = "/" //ll add the index in gets
	str = begin+"services"+end
	jdc.Services, err = gets(str)
	fmt.Println(jdc.Services)
	//gets always return an error :D	
	//if err != nil { fmt.Println(err) }

	str = begin+"packages"+end
	jdc.Packages, err = gets(str)
	fmt.Println(jdc.Packages)
	//if err != nil { fmt.Println(err) }


	return nil
}