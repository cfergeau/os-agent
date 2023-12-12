package led

import (
	"os"
	"slices"
	"testing"

	"github.com/stretchr/testify/require"
)

func testGetTriggerForLed(t *testing.T, ledName string, expectedTrigger string) {
	led := LED{Name: ledName}

	trigger, err := led.GetTrigger()
	require.NoError(t, err)
	require.Equal(t, expectedTrigger, trigger)
}

func TestGetTrigger(t *testing.T) {
	sysClassLedsPath = "testdata"
	testGetTriggerForLed(t, "activity", "activity")
	testGetTriggerForLed(t, "power", "default-on")
	testGetTriggerForLed(t, "user", "heartbeat")
}

func requireEqualFileContent(t *testing.T, path, content string) {
	data, err := os.ReadFile(path)
	require.NoError(t, err)
	require.Equal(t, content, string(data))
}

func TestSetTrigger(t *testing.T) {
	sysClassLedsPath = t.TempDir()
	ledPower := LED{Name: "power", DefaultTrigger: "default-on"}
	err := os.Mkdir(ledPower.getSysfsPath(), 0755)
	require.NoError(t, err)

	err = ledPower.SetTrigger(false)
	require.NoError(t, err)

	ledTriggerFilePath := ledPower.getSysfsPath("trigger")
	requireEqualFileContent(t, ledTriggerFilePath, "none")

	err = ledPower.SetTrigger(true)
	require.NoError(t, err)

	requireEqualFileContent(t, ledTriggerFilePath, "default-on")

}

func TestGetLeds(t *testing.T) {
	sysClassLedsPath = "testdata"
	leds, err := GetLeds()
	require.NoError(t, err)
	ledNames := []string{}
	for _, led := range leds {
		ledNames = append(ledNames, led.Name)
	}
	slices.Sort(ledNames)
	require.Equal(t, []string{"activity", "power", "user"}, ledNames)

}
