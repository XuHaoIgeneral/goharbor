package main

import (
	"context"
	"fmt"
	harbor "github.com/XuHaoIgeneral/goharbor"
)

const (
	username = "xxxxx"
	password = "xxxxx"
	host     = "http://xxxx.xxxx.xxx:xxxx"
)

func main() {
	// create harbor client
	c, err := harbor.NewClient(nil, host, username, password, true)
	if err != nil {
		panic(err)
	}
	//list project
	ps, err := c.ListProjects(context.Background(), nil)
	if err != nil {
		panic(err)
	}

	// dump projects
	for _, p := range ps {
		fmt.Printf("%+v\n", p)
	}

}
