package main

/*
Octopus is a Shark & Sqid co-ordinator.
- spins up new instances
- aggregates and respondds to health checks
- maintains open connection to Squids and Sharks at all times
- allocates Sharks to deploy containers
- manages Coral and builds containers before sending tar to a given Shark to deploy
*/

import (
	"bytes"
	"code.highf.in/chalkhq/shared/types"
	"code.highf.in/chalkhq/shark/api/project"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/exec"
	"strings"
)

type Paths struct {
	Base      string
	Config    string
	Templates string
	Public    string
	Routes    string
}

var (
	paths Paths
)

func init() {
	var whoami bytes.Buffer
	cmd := exec.Command("whoami")
	cmd.Stdout = &whoami
	_ = cmd.Run()
	if strings.Index(whoami.String(), "root") == -1 {
		fmt.Println("shark must be run as root")
		os.Exit(1)
		return
	}
}

func main() {
	fmt.Println("running shark")

	config := GetConfig()

	fmt.Println("listening on " + config.Port)

	http.HandleFunc("/", func(w http.ResponseWriter, req *http.Request) {
		fmt.Println("/")
		r := types.Response{}

		r.Response.Meta.Status = 200

		r.Response.Meta.Errors = make([]string, 0)

		r.Response.Data = make(map[string]interface{})

		r.W = w

		r.Req = *req
		r.Req.ParseForm()
		//r.Req.ParseMultipartForm(64)

		r.Fields = r.Req.Form
		if len(r.Hashbang) > 0 {
			r.Segments = strings.Split(r.Hashbang, "/")
		}

		router(r)

		//fmt.Println("/", r.Req.PostForm, r.Req.MultipartForm, r.Req.PostFormValue("project"))

	})

	//http.Handle("/", http.FileServer(http.Dir(paths.Config)))

	// Simple static webserver:
	log.Fatal(http.ListenAndServe(":"+config.Port, nil))
	//http.FileServer(http.Dir("./public")

}

func router(r types.Response) {
	fmt.Println("router")
	switch r.Req.URL.Path {

	case "/project/deploy":
		// receives a tar file, the command to run, and resource limits for each container
		// needs to run a docker container
		project.Deploy(r)

	case "/project/remove":
		// removes a running instance
		project.Remove(r)

	default:
		r.Response.Meta.Status = 404
		r.AddError("Path: " + r.Req.URL.Path + " is not valid")
		r.Kill(200)
	}
}
