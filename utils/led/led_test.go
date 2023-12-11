package led

import (
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
