// Filters take place between the user request and the endpoints he wants to
// reach. There are used to check for authentification, permission for example,
// or any intermediate processing that would be necessary.
package filters


import (
    "launchpad.net/loggo"
)


var log = loggo.GetLogger("hivy.filters")
