package main

import (
	"fmt"
	"log"
	"os/exec"
)

func main2() {
	cmd := exec.Command("OsuStatusAddr")
	out, err := cmd.CombinedOutput()
	if err != nil {
		log.Fatalf("cmd.Run() failed with %s\n", err)
	}
	fmt.Printf("combined out:\n%s\n", string(out))
}
