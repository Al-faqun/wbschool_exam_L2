package main

import (
	"os/exec"
	"testing"
)

// There is nothing really to test in task.go, this test is here just for fun
func TestMain(t *testing.T) {
	// todo: test realMain func instead
	cmd := exec.Command("./dev01.exe")
	out, err := cmd.CombinedOutput()

	sout := string(out)

	if err != nil {
		t.Errorf("Command dev01.exe exited with error %s", err)
	} else {
		t.Logf("Command exited successfully, out=%s", sout)
	}
}
