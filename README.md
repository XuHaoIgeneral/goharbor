# goharbor


## v 1.0.6
```golang
package main

import (
	"context"
	"fmt"
	"github.com/XuHaoIgeneral/goharbor"
)

const (
	username = "admin"
	password = "admin"
	host     = "http://myharbor.company.com"
)

func main() {
	// create harbor client
	c, err := harbor.NewClient(nil, host, username, password, true)
	if err != nil {
		panic(err)
	}

	// list project
	ps, err := c.ListProjects(context.Background(), nil)
	if err != nil {
		panic(err)
	}

	// dump projects
	for _, p := range ps {
		fmt.Printf("%+v\n", p)
	}
}

```
