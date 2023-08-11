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
This package handles how the config is structrured
*/
package conf

import (
	"log"
	"os"

	"github.com/BurntSushi/toml"
)

var DefaultConfig = Config{
	Sensor: Sensor{
		OperationMode: OperationMode{
			Main:      "all",
			Secondary: "gas",
		},
		FanCleaningInterval: 604800,
		ForceCleanFan:       false,
		RthAccelerationMode: 0,
		TemperatureOffset:   0,
	},
	Measurement: Measurement{
		WaitBetweenMeasurements:           0,
		TakeMeasurementsFor:               300,
		SleepUntilNextBatchOfMeasurements: 300,
		IgnoreFirstXMeasurements:          30,
	},
	Console: Console{
		Enabled: true,
	},
	SensorCommunity: SensorCommunity{
		Enabled:      false,
		SensorNodeID: "raspi-123456789",
	},
	DataLogging: DataLogging{
		Enabled:  true,
		FileName: "",
	},
}

type Config struct {
	Sensor          Sensor
	Measurement     Measurement
	Console         Console
	DataLogging     DataLogging
	SensorCommunity SensorCommunity
}

type Sensor struct {
	OperationMode       OperationMode
	FanCleaningInterval int
	ForceCleanFan       bool
	RthAccelerationMode int
	TemperatureOffset   float32
}

type OperationMode struct {
	Main      string
	Secondary string
}

type Measurement struct {
	WaitBetweenMeasurements           int
	TakeMeasurementsFor               int
	SleepUntilNextBatchOfMeasurements int
	IgnoreFirstXMeasurements          int
}

type Console struct {
	Enabled bool
}

type SensorCommunity struct {
	Enabled      bool
	SensorNodeID string
}

type DataLogging struct {
	Enabled  bool
	FileName string
}

func LoadConfig() (Config, error) {
	var config Config = DefaultConfig
	_, err := toml.DecodeFile("config.toml", &config)
	if err != nil {
		if _, ok := err.(*os.PathError); ok {
			// the file does not exist, create it with default values
			log.Println("config file does not exist, creating it with default values...")
			file, err := os.Create("config.toml")
			if err != nil {
				return Config{}, err
			}
			defer file.Close()
			toml.NewEncoder(file).Encode(DefaultConfig)
			log.Println("created config file successfully!")
			log.Println("please edit the config file and run the program again!")
			os.Exit(0)
		} else {
			return Config{}, err
		}
	}

	// open the config file and write the already parsed config to it
	// this help to write the default values to the config file if
	//it does not exist for example after an update to the program
	file, err := os.OpenFile("config.toml", os.O_WRONLY, 0644)
	if err != nil {
		return Config{}, err
	}
	err = toml.NewEncoder(file).Encode(config)
	if err != nil {
		return Config{}, err
	}
	err = file.Close()
	if err != nil {
		return Config{}, err
	}

	return config, nil
}
