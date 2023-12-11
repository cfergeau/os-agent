package led

import (
	"os"
	"path/filepath"
	"strings"
)

// unit tests override this
var sysClassLedsPath = "/sys/class/leds/"

type LED struct {
	Name           string
	DefaultTrigger string
}

func (led LED) getSysfsPath(elems ...string) string {
	relPath := filepath.Join(elems...)

	return filepath.Join(sysClassLedsPath, led.Name, relPath)
}

func (led LED) GetTrigger() (string, error) {
	ledTriggerFilePath := led.getSysfsPath("trigger")
	ledTrigger, err := os.ReadFile(ledTriggerFilePath)
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(ledTrigger)), err
}

func (led LED) SetTrigger(newState bool) error {
	ledTriggerFilePath := led.getSysfsPath("trigger")
	var newTrigger []byte

	if newState {
		newTrigger = []byte(led.DefaultTrigger)
	} else {
		newTrigger = []byte("none")
	}

	err := os.WriteFile(ledTriggerFilePath, newTrigger, 0600)
	if err != nil {
		return err
	}

	return nil
}
