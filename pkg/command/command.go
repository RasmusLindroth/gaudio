package command

import (
	"io/ioutil"
	"os/exec"
)

func RunCmd(name string, parts ...string) {
	exec.Command(name, parts...).Start()
}

func RunWithOutput(name string, parts ...string) (string, error) {
	cmd := exec.Command(name, parts...)
	outpipe, err := cmd.StdoutPipe()
	if err != nil {
		return "", err
	}

	err = cmd.Start()
	if err != nil {
		return "", err
	}

	obuf, err := ioutil.ReadAll(outpipe)
	if err != nil {
		return "", err
	}

	err = cmd.Wait()
	if err != nil {
		return "", err
	}

	return string(obuf), nil
}
