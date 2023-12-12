package generic

import (
	"unicode"

	"github.com/godbus/dbus/v5"
	"github.com/godbus/dbus/v5/introspect"
	"github.com/godbus/dbus/v5/prop"

	"github.com/home-assistant/os-agent/utils/led"
	logging "github.com/home-assistant/os-agent/utils/log"
)

const (
	objectPath = "/io/hass/os/Boards/Generic"
	ifaceName  = "io.hass.os.Boards.Generic"
)

type genericBoard struct {
	conn  *dbus.Conn
	props *prop.Properties
}

func getTriggerLED(led led.LED) bool {
	value, err := led.GetTrigger()
	if err != nil {
		logging.Error.Print(err)
	}
	return value != "none"
}

func setTriggerLED(led led.LED, c *prop.Change) *dbus.Error {
	logging.Info.Printf("Set Generic Board %s LED to %t", led.Name, c.Value)
	err := led.SetTrigger(c.Value.(bool))
	if err != nil {
		return dbus.MakeFailedError(err)
	}
	return nil
}

func newSetLEDCallback(led led.LED) func(c *prop.Change) *dbus.Error {
	return func(c *prop.Change) *dbus.Error {
		return setTriggerLED(led, c)
	}
}

func dbusNameForLed(led led.LED) string {
	// utf8string.NewString(ledName).IsASCII()
	dbusName := ""
	upperCase := true // indicates the next rune we add to dbusName should be upper-case
	for _, rune := range led.Name {
		if !(unicode.IsLetter(rune) || unicode.IsNumber(rune)) {
			// Ignore ':', '_', '-', ... which can appear in led names
			upperCase = true
			continue
		}

		if !upperCase {
			dbusName = dbusName + string(unicode.ToLower(rune))
		} else {
			dbusName = dbusName + string(unicode.ToUpper(rune))
			upperCase = false
		}

	}

	return dbusName + "LED"
}

func dbusExportLeds(conn *dbus.Conn, objectPath dbus.ObjectPath, ifaceName string, leds ...led.LED) error {
	d := genericBoard{
		conn: conn,
	}

	ledProps := map[string]*prop.Prop{}
	for _, led := range leds {
		ledProps[dbusNameForLed(led)] = &prop.Prop{
			Value:    getTriggerLED(led),
			Writable: true,
			Emit:     prop.EmitTrue,
			Callback: func(c *prop.Change) *dbus.Error { return setTriggerLED(led, c) },
		}
	}

	propsSpec := map[string]map[string]*prop.Prop{
		ifaceName: ledProps,
	}

	props, err := prop.Export(conn, objectPath, propsSpec)
	if err != nil {
		return err
	}
	d.props = props

	err = conn.Export(d, objectPath, ifaceName)
	if err != nil {
		return err
	}

	node := &introspect.Node{
		Name: string(objectPath),
		Interfaces: []introspect.Interface{
			introspect.IntrospectData,
			prop.IntrospectData,
			{
				Name:       ifaceName,
				Methods:    introspect.Methods(d),
				Properties: props.Introspection(ifaceName),
			},
		},
	}

	err = conn.Export(introspect.NewIntrospectable(node), objectPath, "org.freedesktop.DBus.Introspectable")
	if err != nil {
		return err
	}
	logging.Info.Printf("Exposing object %s with interface %s ...", objectPath, ifaceName)

	return nil
}

func InitializeDBus(conn *dbus.Conn) {
	var (
		ledPower    = led.LED{Name: "power", DefaultTrigger: "default-on"}
		ledActivity = led.LED{Name: "activity", DefaultTrigger: "activity"}
		ledUser     = led.LED{Name: "user", DefaultTrigger: "heartbeat"}
	)

	err := dbusExportLeds(conn, objectPath, ifaceName, ledPower, ledActivity, ledUser)
	if err != nil {
		logging.Critical.Panic(err)
	}
}
