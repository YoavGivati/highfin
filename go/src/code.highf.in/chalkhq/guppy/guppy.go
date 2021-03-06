package main

import "os"
import "os/exec"
import "code.highf.in/chalkhq/guppy/util"
import "bytes"
import "strings"
import "code.highf.in/chalkhq/shared/nodejs"
import "code.highf.in/chalkhq/shared/log"

/*
This is the command-line interface to Guppy. Because Shark also uses Guppy, all of the utility functions are in the guppy/util package
The commandline version stores project info on disk so you don't have to keep configuring it
*/
var config util.Config

func init() {
	// if configuration is not set, it'll ask the user
	//config = util.GetConfig()
}

func main() {

	var command string
	// get the command from the commandline
	if len(os.Args) > 1 {
		command = os.Args[1]
	}

	switch command {

	case "bootstrap":
		// run as root
		var whoami bytes.Buffer
		cmd := exec.Command("whoami")
		cmd.Stdout = &whoami
		_ = cmd.Run()
		if strings.Index(whoami.String(), "root") == -1 {
			log.Log("guppy bootstrap must be run as root")
			return
		}
		// end run as root

		/*
			todo: install latest golang and node
		*/

		log.Log("bootstrapping...")

	case "config":
		// sets new config based on user input
		var new_config util.Config
		util.SetConfig(new_config)
	case "create":
		config = util.GetConfig()
		util.Create(config.Account, config.Project, config.Email, config.Server)
	case "get":
		// you can reclone and refresh your environment whenever you want.
		config = util.GetConfig()
		util.Clone(config.Account, config.Project, config.Email, config.Server)
	case "push":
		// should tests be run here?
		if len(os.Args) >= 3 {
			message := os.Args[2] // commit message
			util.Push(message)
		} else {
			log.Log("guppy push commit_message")
		}

	case "deploy":
		branch := os.Args[2]
		config = util.GetConfig()
		util.Deploy(config.Account, config.Project, branch, config.Server)

	case "validate-key":
		config = util.GetConfig()
		util.ValidateKey(config.Account, config.Project, config.Email)
		// todo(yoav) really only for testing, ValidateKey would only be called when cloning/pushing code

	case "test":
		util.Test(os.Args[1:])
		// todo(yoav) runs tests

	case "run":
		util.Run()

	case "npm-install":
		util.NpmInstall()
	case "npm":
		nodejs.Npm(os.Args[1:])
	case "grunt":
		nodejs.Grunt(os.Args[1:])
	case "test-salmon":
		// test salmon

	case "get-salmon":
		// clone salmon
		util.CloneSalmon()

	case "run-salmon":
		//util.Watch()

	case "set-server":
		config = util.GetConfig()
		if len(os.Args) < 3 {
			log.Log("guppy set-server [http://domain.com | http://ipaddress]")
			return
		}
		config.Server = os.Args[2]
		util.SetConfig(config)

	case "finish":
		// prepares the basebox to be compiled into a frybox
		// ie: when you're finished working on guppy, call guppy finish, then follow the steps to convert to a basebox
		cmd := exec.Command("cp", "-rfp", "/vagrant/go/bin/guppy", "/usr/bin/guppy")
		err := cmd.Run()
		if err != nil {
			log.Log(err.Error())
		}
		cmd = exec.Command("cp", "-rfp", "/vagrant/go/bin/godep", "/usr/bin/godep")
		err = cmd.Run()
		if err != nil {
			log.Log(err.Error())
		}

	case "update":
		// updates guppy to the latest version
		// todo: should output version before and after update
		// todo: should have a more concrete verioning/update system isntead of just downloading the latest build from github
		cmd := exec.Command("wget", "-o", "/guppy/guppy", "https://github.com/YoavGivati/highfin/raw/master/go/bin/guppy")
		err := cmd.Run()
		if err != nil {
			log.Log(err.Error())
		}

	default:
		log.Log("Try: guppy [bootstrap, config, get, test, run, push, deploy dev-next, npm-install]")
	}

}

// todo(yoav): make it accept command line arguments...
// watch function is going to be continuously running.. so can we still trigger commands somehow? how does that work?
