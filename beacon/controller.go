package beacon

import (
  "fmt"
  "path/filepath"
  "strconv"
  "strings"

  "github.com/coreos/go-etcd/etcd"
)

const (
  // MaxMachinesRule restritcs number of deployments allowed
  // Note that multi-machines deployment counts for one
  MaxMachinesRule int = 10
)

// Controller knows about rules, constraints and has methods to check them
type Controller struct {
  db          *etcd.Client
  maxMachines int
  user        string
}

// NewController setup debug mode and return an initialized controller
func NewController(user string, debug bool) *Controller {
  // Print out requests etcd is processing
  if debug {
    etcd.OpenDebug()
  }
  // Default config
  //TODO machines ip https://github.com/coreos/go-etcd/blob/master/etcd/client.go
  machines := []string{"http://127.0.0.1:4001"}
  return &Controller{
    db:          etcd.NewClient(machines),
    maxMachines: MaxMachinesRule,
    user:        user,
  }
}

// SetUser changes the user currently controlled
func (c *Controller) SetUser(user string) {
  c.user = user
}

// Set wraps go-etcd Set method
func (c *Controller) Set(key, value string, ttl uint64) (*etcd.Response, error) {
  return c.db.Set(key, value, ttl)
}

// Get wraps go-etcd Get method
func (c *Controller) Get(path string) ([]*etcd.Response, error) {
  return c.db.Get(path)
}

// Delete wraps go-etcd Delete method
func (c *Controller) Delete(path string) (*etcd.Response, error) {
  return c.db.Delete(path)
}

func (c *Controller) setMethodPermission(method, permission string) error {
  //TODO Temporary permission with ttl ?
  var ttl uint64
  // Setting /hivy/security/{user}/methods/{method} to "true" makes {method}
  // available for {user]
  feedback, err := c.db.Set(filepath.Join("hivy/security", c.user, "methods", method), permission, ttl)
  if err != nil {
    return fmt.Errorf("unable to set method permission (%s)\n", method)
  }
  log.Debugf("%v\n", feedback)
  return nil
}

// EnableMethod makes the given method available for the currently controlled
// user (i.e. c.user)
func (c *Controller) EnableMethod(method string) error {
  log.Infof("enabling method %s\n", method)
  return c.setMethodPermission(method, "true")
}

// DisableMethod makes the given method forbidden for the currently controlled
// user (i.e. c.user)
func (c *Controller) DisableMethod(method string) error {
  log.Infof("disabling method %s\n", method)
  return c.setMethodPermission(method, "false")
}

// CheckMethod makes sure the given one is like GET/path/to/endpoint (without parameters)
func (c *Controller) CheckMethod(method string) (bool, error) {
  // Permission needs /hivy/security/{user}/methods/{method} to exist
  result, err := c.db.Get(filepath.Join("hivy/security", c.user, "methods", method))

  if err != nil {
    log.Errorf("[controller.CheckMethod] %v\n", err)
    return false, err
  } else if result[0].Value == "false" {
    log.Infof("[controller.CheckMethod] Method forbidden\n")
    return false, nil
  }

  log.Infof("Method allowed, processing (%s:%s)", c.user, method)
  return true, nil
}

// Ressource the current state or value of the given ressource name
// Localtion: v1/keys/hivy/security/{user}/ressources/{ressource}
func (c *Controller) Ressource(name string) (string, error) {
  result, err := c.db.Get(filepath.Join("hivy/security", c.user, "ressources", name))
  if err != nil {
    log.Errorf("[controller.CheckMethod] %v\n", err)
    return "", err
  }
  return result[0].Value, nil
}

func (c *Controller) updateMachinesRessource(operation int) error {
  // Get current state of the machine ressources (i.e. currently deployed machines)
  machinesStr, err := c.Ressource("machines")
  if err != nil {
    return err
  }

  // Add the requested number of machines to deploy
  machines, err := strconv.Atoi(machinesStr)
  if err != nil {
    return err
  }
  newMachinesStr := strconv.Itoa(machines + operation)

  // Update machines state in database
  feedback, err := c.db.Set(filepath.Join("hivy/security", c.user, "ressources", "machines"), newMachinesStr, 0)
  if err != nil {
    log.Errorf("[controller.updateMachinesRessource] %v\n", err)
    return err
  }
  log.Debugf("%v\n", feedback)

  // Check if it reached a rule
  if (machines + operation) > c.maxMachines {
    c.DisableMethod("GET/juju/deploy")
  } else {
    c.EnableMethod("GET/juju/deploy")
  }

  return nil
}

// Update actualize states in the database regarding the given request
func (c *Controller) Update(method string) error {
  switch {
  case strings.Contains(method, "deploy"):
    log.Infof("Updating machines ressource (+1)\n")
    return c.updateMachinesRessource(1)
  }
  return nil
}
