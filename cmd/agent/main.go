package main

import (
	"fmt"
	"os"

	"github.com/bmehdi777/moon/internal/pkg/agent/cmd"
)

func main() {
	if err := cmd.ExecuteClient(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
