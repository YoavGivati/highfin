package util

import (
	"code.highf.in/chalkhq/shared/log"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"os/exec"
	"strings"
)

func Clone(account string, project string, email string, server string) {
	// generate if needed and copy dev keys
	ApplyKeys(email)
	// not necessary when cloning
	// check if we have access
	//has_access := ValidateKey(account, project, email)

	//if has_access == true {

	_ = exec.Command("rm", "-R", "/vagrant/code").Run()
	cmd := exec.Command("git", "clone", account+"_"+project+"@"+server+":code.git", "/vagrant/code")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	if err != nil {
		log.Log(err.Error())
		return
	}
	log.Log("Successfully cloned " + account + "'s " + project + " project to /vagrant/code")
	//return
	//}
}

func Create(account string, project string, email string, server string) {
	log.Log(server)
	public_key, _ := ioutil.ReadFile(KEY_PATH + "/guppy.pub")
	res, err := http.PostForm("http://"+server+"/project/create", url.Values{"account": {account}, "project": {project}, "key": {string(public_key)}, "force": {"true"}})
	if err != nil {
		log.Log("request error: ", err.Error())
		return
	}
	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	log.Log(string(body))
}

func Push(message string) {
	// todo: run unit tests locally, only continue if passes, unless "commit_if_passes = false" in -.json
	// wrapper around common git workflow in a frybox
	log.Log("adding")
	cmd := exec.Command("git", "add", ".")
	cmd.Dir = "/vagrant/code"
	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout
	cmd.Run()
	log.Log("committing")
	cmd = exec.Command("git", "commit", "-a", "-m", message)
	cmd.Dir = "/vagrant/code"
	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout
	cmd.Run()

	log.Log("pushing")
	// todo: if commit fails, don't push
	cmd = exec.Command("git", "push", "origin", "dev-next")
	cmd.Dir = "/vagrant/code"
	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout
	cmd.Run()

}

func Deploy(account string, project string, branch string, server string) {
	res, err := http.PostForm("http://"+server+"/project/deploy", url.Values{"account": {account}, "project": {project}, "branch": {branch}})
	if err != nil {
		log.Log("request error: ", err.Error())
		return
	}
	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	log.Log(string(body))
}

func CloneSalmon() {
	ApplyKeys("dev@chalkhq.com")

	_ = exec.Command("rm", "-r", "-f", "/vagrant/salmon").Run()
	//_ = exec.Command("mkdir", "-p", "/vagrant/salmon").Run()

	// todo(yoav) if it exists ask the user if they want to overwrite it

	// todo(yoav) should also checkout dev branch, and should be executed by vagrant install

	// todo(yoav) update all exec.Commands to pipe stdout and std error to os.out/err

	//cmd := exec.Command("git", "clone", "https://github.com/YoavGivati/salmon", "/vagrant/salmon")
	cmd := exec.Command("git", "clone", "git@github.com:YoavGivati/salmon.git", "/vagrant/salmon")
	log.Log(strings.Join(cmd.Args, " "))
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	_ = cmd.Run()

}
