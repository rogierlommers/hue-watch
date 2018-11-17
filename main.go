package main

import (
	"os"
	"time"

	"github.com/amimof/huego"
	"github.com/sirupsen/logrus"
)

type sensorState struct {
	ID       int
	presence bool
}

// enter Bridge IP and API password here
var bridge = huego.New("192.168.1.7", os.Getenv("HUE_API_KEY"))

// savedState contains in-memory state of all sensors
var savedState = make(map[int]sensorState)

func main() {
	logrus.SetLevel(logrus.DebugLevel)

	if err := showLights(); err != nil {
		logrus.Error(err)
	}

	for {
		sensors, err := bridge.GetSensors()
		if err != nil {
			logrus.Errorf("bridge: %q", err)
		}

		checkChanges(savedState, sensors)
		time.Sleep(1 * time.Second)
	}

}

func alertLight(sensor huego.Sensor) {

	if sensor.ID == 12 {
		light, err := bridge.GetLight(3)
		if err != nil {
			logrus.Error(err)
		}

		lightState := huego.State{Alert: "select"}
		if light.IsOn() {
			lightState.On = true
		} else {
			lightState.On = false
		}
		bridge.SetLightState(3, lightState)
	}

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
					alertLight(sensor)
					savedState[sensor.ID] = sensorState{ID: sensor.ID, presence: sensor.State["presence"].(bool)}
				}
			}

		}

	}
}

func showLights() error {
	lights, err := bridge.GetLights()
	if err != nil {
		return err
	}

	for _, light := range lights {
		logrus.Debugf("lamp detected: %q (id: %d, turned on: %v)", light.Name, light.ID, light.IsOn())
	}

	return nil
}
