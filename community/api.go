/*
 Copyright (c) 2023 Michail Angelos Tsiantakis

 This program is free software: you can redistribute it and/or modify
 it under the terms of the GNU General Public License as published by
 the Free Software Foundation, either version 3 of the License, or
 (at your option) any later version.

 This program is distributed in the hope that it will be useful,
 but WITHOUT ANY WARRANTY; without even the implied warranty of
 MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
 GNU General Public License for more details.

 You should have received a copy of the GNU General Public License
 along with this program. If not, see <https://www.gnu.org/licenses/>.
*/

/*
This package contains the functions to send data to the sensor.community site.
*/
package community

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/aggellos2001/sen5x-go/conf"
	"github.com/aggellos2001/sen5x-go/sen5xlib"
)

const (
	softwareVersion = "1.0"
	endpoint        = "https://api.sensor.community/v1/push-sensor-data/"
)

type SensorCommunityRequest struct {
	SoftwareVersion  string                 `json:"software_version"`
	SensorDataValues []SensorCommunityValue `json:"sensordatavalues"`
}

type SensorCommunityValue struct {
	ValueType string `json:"value_type"`
	Value     string `json:"value"`
}

func SendDataToAPI(measurement *sen5xlib.SensorMeasurement, conf *conf.Config) {
	httpClient := &http.Client{
		Timeout: 30 * time.Second,
	}

	var sensorPmRequest = SensorCommunityRequest{
		SoftwareVersion: softwareVersion,
		SensorDataValues: []SensorCommunityValue{
			{
				ValueType: "P1",
				Value:     fmt.Sprintln(measurement.PM10_0),
			},
			{
				ValueType: "P2",
				Value:     fmt.Sprintln(measurement.PM2_5),
			},
		},
	}

	jsonBody, err := json.Marshal(sensorPmRequest)
	if err != nil {
		log.Println(err)
		return
	}

	req, err := http.NewRequest("POST", endpoint, bytes.NewBuffer(jsonBody))
	if err != nil {
		log.Println(err)
		return
	}
	req.Header.Add("X-Sensor", conf.SensorCommunity.SensorNodeID)
	req.Header.Add("X-Pin", "1") // For sensirion i've seen both 1 and 7 on the data api ?
	req.Header.Add("Content-Type", "application/json")

	resp, err := httpClient.Do(req)
	if err != nil {
		log.Println(err)
		return
	}
	resp.Body.Close()

	var sensorWeatherRequest = SensorCommunityRequest{
		SoftwareVersion: softwareVersion,
		SensorDataValues: []SensorCommunityValue{
			{
				ValueType: "temperature",
				Value:     fmt.Sprintln(measurement.Temp),
			},
			{
				ValueType: "humidity",
				Value:     fmt.Sprintln(measurement.Hum),
			},
		},
	}

	jsonBody, err = json.Marshal(sensorWeatherRequest)
	if err != nil {
		log.Println(err)
		return
	}

	req, err = http.NewRequest("POST", endpoint, bytes.NewBuffer(jsonBody))
	if err != nil {
		log.Println(err)
		return
	}
	req.Header.Add("X-Sensor", conf.SensorCommunity.SensorNodeID)
	req.Header.Add("X-Pin", "7")
	req.Header.Add("Content-Type", "application/json")

	resp, err = httpClient.Do(req)
	if err != nil {
		log.Println(err)
		return
	}
	resp.Body.Close()
}
