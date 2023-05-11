package main

import (
	"context"
	"fmt"
	"html/template"
	"log"
	"net/http"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/client"
)

func getContainers(cli *client.Client, network string) []types.Container {

	containers, err := cli.ContainerList(context.Background(), types.ContainerListOptions{
		Filters: filters.NewArgs(filters.KeyValuePair{Key: "network", Value: network}),
	})
	if err != nil {
		panic(err)
	}
	return containers
}

func getNetworks(cli *client.Client) []types.NetworkResource {

	networks, err := cli.NetworkList(context.Background(), types.NetworkListOptions{})
	if err != nil {
		panic(err)
	}

	return networks
}

type Page struct {
	Network    string
	Containers []types.Container
	Networks   []types.NetworkResource
}

func getRoot(cli *client.Client) func(w http.ResponseWriter, r *http.Request) {

	return func(w http.ResponseWriter, r *http.Request) {
		fmt.Printf("got / request\n")
		t, err := template.New("index.html").ParseFiles("index.html")

		if err != nil {
			log.Panicln(err)
		}
		network := r.URL.Query().Get("network")
		page := Page{
			Network:    network,
			Containers: getContainers(cli, network),
			Networks:   getNetworks(cli),
		}
		t.Execute(w, page)
	}
}

func main() {
	cli, err := client.NewClientWithOpts(client.FromEnv)
	if err != nil {
		log.Panicln(err)
	}
	defer cli.Close()
	log.Println("starting")
	http.HandleFunc("/", getRoot(cli))
	e := http.ListenAndServe("0.0.0.0:3333", nil)
	if e != nil {
		log.Fatalln(e)
	}
}
