package context

//import "fmt"

type Context struct {
	User string
	Project string
	Target []string //example : 1 row : series : base, N rows : package : [a,b,c] c-plugin : [d,e,f]
	Value string
	State string
}

//define allow contexts
/*
type StateContexts struct {
	users []string{}
	projects []string{} // generic : "tag" link 
	utop map[string][]string{}
}
*/

//le contexte peut etre def a part mais la gestion de ma machine a etat doit aller dans un paquet state_engine
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
