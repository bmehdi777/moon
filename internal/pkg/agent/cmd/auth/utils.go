package auth

import (
	"errors"
	"os"
	"os/exec"
	"runtime"
)

// if needed somewhere else, it should be moved to a different package
func openInBrowser(url string) error {
	var cmd *exec.Cmd

	switch runtime.GOOS {
	case "linux":
		cmd = exec.Command("xdg-open", url)
	case "darwin":
		cmd = exec.Command("open", url)
	}

	if cmd != nil {
		// only show error, no stdout
		cmd.Stderr = os.Stderr
		err := cmd.Start()
		if err != nil {
			return err
		}
		err = cmd.Wait()
		if err != nil {
			return err
		}
		return nil
	} else {
		return errors.New("Unsupported platform")
	}
}
