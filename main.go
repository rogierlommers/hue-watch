package main

import (
	"time"

	"github.com/amimof/huego"
	"github.com/sirupsen/logrus"
)

type sensorState struct {
	ID       int
	presence bool
}

func main() {
	logrus.SetLevel(logrus.DebugLevel)

	bridge := huego.New("192.168.1.7", "") // enter Bridge IP and API password here

	// savedState contains in-memory state of all sensors
	var savedState = make(map[int]sensorState)

	for {
		sensors, err := bridge.GetSensors()
		if err != nil {
			logrus.Errorf("bridge: %q", err)
		}

		checkChanges(savedState, sensors)
		time.Sleep(1 * time.Second)
	}

}

func alertLight(bridge *huego.Bridge) {
	lightState := huego.State{On: true, Alert: "select"}
	bridge.SetLightState(3, lightState)
}

func checkChanges(savedState map[int]sensorState, sensors []huego.Sensor) {

	for _, sensor := range sensors {
		if sensor.Type == "ZLLPresence" {

			if _, exists := savedState[sensor.ID]; !exists {
				logrus.Infof("found sensor %q, saving as ID %d", sensor.Name, sensor.ID)
				savedState[sensor.ID] = sensorState{ID: sensor.ID, presence: sensor.State["presence"].(bool)}
			} else {

				if savedState[sensor.ID].presence != sensor.State["presence"].(bool) {
					logrus.Debugf("change detected for sensor %q --> %v", sensor.Name, sensor.State["presence"].(bool))
					savedState[sensor.ID] = sensorState{ID: sensor.ID, presence: sensor.State["presence"].(bool)}
				}
			}

		}

	}
}
