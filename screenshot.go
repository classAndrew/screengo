package main

import (
	"fmt"
	"os/exec"
	"time"
)

// Screenshot Screenshots the current screen
func Screenshot(seconds int) {
	time.Sleep(time.Duration(seconds) * time.Second)
	exec.Command("import", "-window", "root", "screenshot.png").Run()
	fmt.Println("Screenshot taken")
}
