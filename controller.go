package main


import (
    "fmt"
    "path/filepath"
    "strings"
    "strconv"

    "github.com/coreos/go-etcd/etcd"
)


const (
    MAX_MACHINES int = 5
)


type Controller struct {
    db *etcd.Client
    max_machines int
    user string
}


func NewController(user string, debug bool) *Controller {
    if debug {etcd.OpenDebug()}
    return &Controller{
        db           : etcd.NewClient(),
        max_machines : MAX_MACHINES,
        user         : user,
    }
}


func (c *Controller) SetUser(user string) {
    c.user = user
}


func (c *Controller) setMethodPermission(method, permission string) error {
    //TODO Temporary permission with ttl ?
    var ttl uint64 = 0
    feedback, err := c.db.Set(filepath.Join("hivy/security", c.user, "methods", method), permission, ttl)
    if err != nil {
        return fmt.Errorf("Unable to set method permission (%s)\n", method)
    }
    log.Debugf("%v\n", feedback)
    return nil
}


func (c *Controller) EnableMethod(method string) error {
    log.Infof("Enabling method %s\n", method)
    return c.setMethodPermission(method, "1")
}


func (c *Controller) DisableMethod(method string) error {
    log.Infof("Disabling method %s\n", method)
    return c.setMethodPermission(method, "0")
}


// method must be of kind GET/path/to/endpoint (without parameters)
func (c *Controller) CheckMethod(method string) (bool, error) {
    // Permission needs /hivy/security/{user}/methods/{method} to exist
    result, err := c.db.Get(filepath.Join("hivy/security", c.user, "methods", method))

    if err != nil {
        log.Errorf("[controller.CheckMethod] %v\n", err)
        return false, err
    } else if result[0].Value == "0" {
        log.Infof("[controller.CheckMethod] Method forbidden\n")
        return false, nil
    }

    log.Infof("Method allowed, processing (%s:%s)", c.user, method)
    return true, nil
}


func (c *Controller) Ressource(name string) (string, error) {
    result, err := c.db.Get(filepath.Join("hivy/security", c.user, "ressources", name))
    if err != nil {
        log.Errorf("[controller.CheckMethod] %v\n", err)
        return "", err
    }
    return result[0].Value, nil
}


func (c *Controller) updateMachinesRessource(operation int) error {
    machines_str, err := c.Ressource("machines")
    if err != nil {
        return err
    }
    machines, err := strconv.Atoi(machines_str)
    if err != nil {
        return err
    }
    new_machines_str := strconv.Itoa(machines + operation)

    feedback, err := c.db.Set(filepath.Join("hivy/security", c.user, "ressources", "machines"), new_machines_str, 0)
    if err != nil {
        log.Errorf("[controller.updateMachinesRessource] %v\n", err)
        return err
    }
    log.Debugf("%v\n", feedback)

    if (machines + operation) > c.max_machines {
        c.DisableMethod("GET/juju/deploy")
    } else {
        c.EnableMethod("GET/juju/deploy")
    }

    return nil
}


func (c *Controller) Update(method string) error {
    switch {
    case strings.Contains(method, "deploy"):
        log.Infof("Updating machines ressource (+1)\n")
        return c.updateMachinesRessource(1)
    }
    return nil
}
