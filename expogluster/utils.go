package expogluster

import (
	"bytes"
	"log"
	"os/exec"
	"strings"

	xml2json "github.com/samuelhug/goxml2json"
)

//Gluster executes oscommands
func gluster(vars ...string) (*bytes.Buffer, error) {
	if len(vars) < 1 {
		log.Println("incorrect url")
		return nil, nil
	}
	args := append(vars, "--xml")
	gCmd := exec.Command("gluster", args...)
	log.Println(vars, args, gCmd.Path, gCmd.Args)
	output, err := gCmd.CombinedOutput()
	if err != nil {
		log.Println(string(output))
	}

	xml := strings.NewReader(string(output))
	json, err := xml2json.Convert(xml)
	if err != nil {
		log.Println(err.Error())
		return json, err
	}
	return json, err
}
