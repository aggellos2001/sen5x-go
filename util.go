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

package main

import "github.com/aggellos2001/sen5x-go/sen5xlib"

// Pass by reference to avoid copying the whole slice of measurements every time
func AverageMeasurement(sliceOfMeasurements *[]sen5xlib.SensorMeasurement) sen5xlib.SensorMeasurement {
	var avgMeasurement sen5xlib.SensorMeasurement
	for _, measurement := range *sliceOfMeasurements {
		avgMeasurement.PM1_0 += measurement.PM1_0
		avgMeasurement.PM2_5 += measurement.PM2_5
		avgMeasurement.PM4_0 += measurement.PM4_0
		avgMeasurement.PM10_0 += measurement.PM10_0
		avgMeasurement.Hum += measurement.Hum
		avgMeasurement.Temp += measurement.Temp
		avgMeasurement.VOC += measurement.VOC
		avgMeasurement.NOx += measurement.NOx
	}
	avgMeasurement.PM1_0 /= uint(len(*sliceOfMeasurements))
	avgMeasurement.PM2_5 /= uint(len(*sliceOfMeasurements))
	avgMeasurement.PM4_0 /= uint(len(*sliceOfMeasurements))
	avgMeasurement.PM10_0 /= uint(len(*sliceOfMeasurements))
	avgMeasurement.Hum /= int(len(*sliceOfMeasurements))
	avgMeasurement.Temp /= int(len(*sliceOfMeasurements))
	avgMeasurement.VOC /= int(len(*sliceOfMeasurements))
	avgMeasurement.NOx /= int(len(*sliceOfMeasurements))

	return avgMeasurement
}
