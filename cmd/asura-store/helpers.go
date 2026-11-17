package main

import (
	"fmt"
	"os/exec"
)

func have(name string) bool {
	_, err := exec.LookPath(name)
	return err == nil
}

func emitStage(stage string, percent int) {
	fmt.Printf("{\"stage\":%q,\"percent\":%d}\n", stage, percent)
}
