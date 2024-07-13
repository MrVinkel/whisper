package whisper

import (
	"fmt"
	"os/exec"
	"runtime"
)

func Open(link string) error {
	var err error

	switch runtime.GOOS {
	case "linux":
		err = exec.Command("xdg-open", link).Start()
		if err != nil {
			err = exec.Command("www-browser", link).Start()
		}
	case "windows":
		err = exec.Command("rundll32", "url.dll,FileProtocolHandler", link).Start()
	case "darwin":
		err = exec.Command("open", link).Start()
	default:
		err = fmt.Errorf("unsupported platform")
	}

	if err != nil {
		return fmt.Errorf("failed to open browser: %w", err)
	}
	return nil

}
