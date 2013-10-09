package endpoints

const (
	// GlobalHelp shows a sum up of the app
	GlobalHelp string = `
RESTful framework for busy people. Effectively expose secured, fully configurable background jobs.
`

	// JujuHelp explains juju exposed endpoints
	JujuHelp string = `
Format: GET /juju/{command}?project={project}
With command one of bootstrap, status or deploy method.
The later needs the project parameter as it will read {project} sepcific
configuration and setup accordingly your private cells.  Port exposure and
charms relationships are automatically processed.
`

	// LoginHelp explains what happens when logging in
	LoginHelp string = `
Format: GET /login
If well authentified, hivy returns a certificate for further secured interactions.
`

	// UserHelp details endpoints relative to user management
	UserHelp string = `
Format: GET /createuser?user={user}&pass={pass}
Store a new user and its credentials, allowing him to access the rest of the API, restricted to his method permissions.
`

	// HelpHelp explains the help endpoint
	HelpHelp string = `
Format: GET /help?method={method}
Will return an help message on the method if provided, global otherwise.
`

	// ConfigHelp explains etcd interaction
	ConfigHelp string = `
Format: GET v1/keys/{path}
Interact with the etcd database, with classic set and get methods.
Note: Currently accessible through an other port
`
)
