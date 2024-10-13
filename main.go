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
This program uses the SEN55 sensor from Sensirion to take various measurements and
print them to a file and or the community site sensor.community if configured.
A configuration file is available at the GitHub repository that you may use to configure
all the available
The folder sen5xlib contains the C bindings and the functions to communicate with the sensor
that Sensirion provides at their GitHub repository. This program is not endorsed by Sensirion!
*/
package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"

	"os"
	"os/signal"
	"strings"
	"time"

	"github.com/aggellos2001/sen5x-go/community"
	"github.com/aggellos2001/sen5x-go/conf"
	"github.com/aggellos2001/sen5xlib"
)

func handlePanic() {
	if err := recover(); err != nil {
		sen5xlib.FreeHal()
		log.Fatalln(err)
	}
}

func main() {

	// this will handle any panics that may occur
	defer handlePanic()

	//handle ctrl-c from user in a different goroutine
	sigintch := make(chan os.Signal, 1)
	signal.Notify(sigintch, os.Interrupt)
	go func() {
		<-sigintch
		err := sen5xlib.StopMeasurement()
		if err != nil {
			log.Fatalln("exiting the program with errors, check connection to sensor...", err)
		}
		log.Println("exiting gracefully...")
		os.Exit(0)
	}()

	config, err := conf.LoadConfig()
	if err != nil {
		log.Fatalln("error while loading config:", err)
	}

	if !config.Console.Enabled {
		// if console is disabled, discard all logs
		log.SetOutput(io.Discard)
	}

	prettyPrintedConfig, _ := json.MarshalIndent(config, "", "  ")
	log.Println("read config successfully and found configuration: ", string(prettyPrintedConfig))

	sen5xlib.InitializeHal()
	defer sen5xlib.FreeHal()

	log.Println("Initialized HAL")

	err = sen5xlib.ResetDevice()
	if err != nil {
		log.Fatalln(err)
	}

	serial, err := sen5xlib.ReadSerialNumber()
	if err != nil {
		log.Fatalln(err)
	}

	log.Println("serial number:", serial)

	prodName, err := sen5xlib.ReadProductName()
	if err != nil {
		log.Fatalln(err)
	}
	log.Println("product name:", prodName)

	firmware, err := sen5xlib.ReadFirmwareVersion()
	if err != nil {
		log.Fatalln(err)
	}
	log.Println("firmware version:", firmware)

	var csvMeasurementFile *os.File
	if config.DataLogging.Enabled {
		var fileName string
		// check if the user entered a file name as an argument
		if len(os.Args) > 1 {
			if err != nil {
				log.Fatalln(err)
			}
			fileName = os.Args[1]
			// check if filename is defined in config
		} else if config.DataLogging.FileName != "" {
			fileName = config.DataLogging.FileName
		} else {
			//otherwise ask input from user
			log.Println("please enter csv file to write the measurements :")
			_, err := fmt.Scanln(&fileName)
			if err != nil {
				log.Fatalln(err)
			}
		}
		// if the user did not enter the .csv extension, add it
		if !strings.HasSuffix(fileName, ".csv") {
			fileName += ".csv"
		}
		csvMeasurementFile, err = os.OpenFile(fileName, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			log.Fatalln(err)
		}
		defer func(csvMeasurementFile *os.File) {
			err := csvMeasurementFile.Close()
			if err != nil {
				log.Fatal("could not close measurement file. aborting program!")
			}
		}(csvMeasurementFile)
	} else {
		log.Println("data logging to a file is disabled")
	}

	for {

		if config.Sensor.OperationMode.Main == "all" {
			err := sen5xlib.StartMeasurement()
			if err != nil {
				log.Fatalln(err)
			}
			log.Println("Started measurement in all mode")
		} else {
			err := sen5xlib.StartMeasurementWithoutPM()
			if err != nil {
				log.Fatalln(err)
			}
			log.Println("started measurement without PM")
		}

		if config.Sensor.OperationMode.Main == "all" &&
			config.Sensor.ForceCleanFan {
			log.Println("Starting fan cleaning. PROGRAM WILL EXIT AFTER FAN CLEANING IS DONE! Change config afterwards!!!")
			err := sen5xlib.StartFanCleaning()
			if err != nil {
				log.Fatalln(err)
			}
			os.Exit(0)
		}

		var sensorMeasurements []sen5xlib.SensorMeasurement

		for i := 0; i < config.Measurement.TakeMeasurementsFor; i++ {

			// read data from sensor. wait until data is ready or time that user specified
			if config.Measurement.WaitBetweenMeasurements > 0 {
				sen5xlib.SleepHal(uint32(config.Measurement.WaitBetweenMeasurements))
			} else {
				for {
					dataReady, err := sen5xlib.ReadDataReady()
					if err != nil {
						log.Println(err)
					}
					if dataReady {
						break
					}
				}
			}

			// if the user specified to ignore the first X measurements, do so
			if i < config.Measurement.IgnoreFirstXMeasurements {
				_, err := sen5xlib.ReadMeasuredValues()
				if err != nil {
					log.Println("error while reading values from sensor", err)
					continue
				}
				continue
			}

			// read data from sensor
			sensorMeasurement, err := sen5xlib.ReadMeasuredValues()
			if err != nil {
				log.Println(err)
			} else {
				// append measurement to slice if no errors occurred
				sensorMeasurements = append(sensorMeasurements, sensorMeasurement)
				log.Println("read measurement:", sensorMeasurement)
			}

		} //end of inner for loop

		if config.Sensor.OperationMode.Secondary == "gas" {
			err := sen5xlib.StartMeasurementWithoutPM()
			if err != nil {
				log.Println(err)
			}
		} else {
			err := sen5xlib.StopMeasurement()
			if err != nil {
				log.Println(err)
			}
		}

		// find the average of all measurements
		avgMeasurement := AverageMeasurement(&sensorMeasurements)

		// free the memory allocated for the slice of measurements
		sensorMeasurements = nil

		str := fmt.Sprintf("%d,%d,%d,%d,%d,%d,%d,%d,%d\n", time.Now().Unix(), avgMeasurement.PM1_0, avgMeasurement.PM2_5, avgMeasurement.PM4_0, avgMeasurement.PM10_0, avgMeasurement.Hum, avgMeasurement.Temp, avgMeasurement.VOC, avgMeasurement.NOx)
		log.Println("average measurement:", str)

		// write the average measurement to the csv file with a timestamp
		// the timestamp is the number of seconds since 1/1/1970
		if config.DataLogging.Enabled {
			_, err := csvMeasurementFile.WriteString(str)
			if err != nil {
				log.Fatalln(err)
			}
			log.Println("wrote average measurement to csv file")
		}

		if config.SensorCommunity.Enabled {
			community.SendDataToAPI(&avgMeasurement, &config)
			log.Println("send average measurement to sensor.community API")
		}

		// sleep until next batch of measurements
		log.Printf("sleeping until next batch of measurements %d seconds...\n", config.Measurement.SleepUntilNextBatchOfMeasurements)
		sen5xlib.SleepHal(uint32(config.Measurement.SleepUntilNextBatchOfMeasurements))

	}

}
