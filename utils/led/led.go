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
	triggers, err := led.GetTriggers()
	if err != nil {
		return "", err
	}
	for _, trigger := range triggers {
		if strings.HasPrefix(trigger, "[") && strings.HasSuffix(trigger, "]") {
			return trigger[1 : len(trigger)-1], nil
		}
	}

	// return an error?
	return "", nil
}

func (led LED) GetTriggers() ([]string, error) {
	ledTriggerFilePath := led.getSysfsPath("trigger")
	ledTrigger, err := os.ReadFile(ledTriggerFilePath)
	if err != nil {
		return nil, err
	}
	return strings.Split(string(ledTrigger), " "), err
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

/* Returns a list of all LEDs available on the system in /sys/class/leds */
func GetLeds() ([]LED, error) {
	leds := []LED{}

	dirs, err := os.ReadDir(sysClassLedsPath)
	if err != nil {
		return nil, err
	}
	for _, dirEntry := range dirs {
		if !dirEntry.IsDir() {
			continue
		}
		leds = append(leds, LED{Name: dirEntry.Name()})
	}

	return leds, nil
}
