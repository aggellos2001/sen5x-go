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

package sen5xlib

/*
#include "sen5x_i2c.h"
#include "sensirion_common.h"
#include "sensirion_i2c_hal.h"
*/
import "C"

import (
	"fmt"
	"math"
	"strings"
	"time"
	"unsafe"
)

// Different models of the SEN5x sensor family may have different features.
// The negative max value of a non-existing feature is used to indicate that
// the feature is not available. (e.g. NOx for SEN54 which is not unsigned as you can notice)
type SensorMeasurement struct {
	PM1_0  uint
	PM2_5  uint
	PM4_0  uint
	PM10_0 uint
	Hum    int
	Temp   int
	VOC    int
	NOx    int
}

// Start Measurement (0x0021)
func StartMeasurement() error {
	err := C.sen5x_start_measurement()
	if err == 1 {
		return fmt.Errorf("error starting measurement")
	}
	return nil

}

// Start Measurement in RHT/Gas-Only Measurement Mode (0x0037)
func StartMeasurementWithoutPM() error {
	err := C.sen5x_start_measurement_without_pm()
	if err == 1 {
		return fmt.Errorf("error starting measurement without PM")
	}
	return nil
}

// Stop Measurement (0x0104)
func StopMeasurement() error {
	err := C.sen5x_stop_measurement()
	if err == 1 {
		return fmt.Errorf("error stopping measurement")
	} else {
		return nil
	}
}

// Read Data-Ready Flag (0x0202)
func ReadDataReady() (bool, error) {
	var ready C.bool
	error := C.sen5x_read_data_ready((*C.bool)(unsafe.Pointer(&ready)))
	if error == 1 {
		return false, fmt.Errorf("error reading data ready")
	}
	return bool(ready), nil
}

// Read Measured Values (0x03C4)
func ReadMeasuredValues() (SensorMeasurement, error) {
	var mass_concentration_pm1p0 C.float
	var mass_concentration_pm2p5 C.float
	var mass_concentration_pm4p0 C.float
	var mass_concentration_pm10p0 C.float
	var ambient_humidity C.float
	var ambient_temperature C.float
	var voc_index C.float
	var nox_index C.float

	err := C.sen5x_read_measured_values((*C.float)(unsafe.Pointer(&mass_concentration_pm1p0)), (*C.float)(unsafe.Pointer(&mass_concentration_pm2p5)), (*C.float)(unsafe.Pointer(&mass_concentration_pm4p0)), (*C.float)(unsafe.Pointer(&mass_concentration_pm10p0)), (*C.float)(unsafe.Pointer(&ambient_humidity)), (*C.float)(unsafe.Pointer(&ambient_temperature)), (*C.float)(unsafe.Pointer(&voc_index)), (*C.float)(unsafe.Pointer(&nox_index)))
	if err == 1 {
		return SensorMeasurement{}, fmt.Errorf("error reading measured values")
	}

	product, _ := ReadProductName()
	if strings.Contains(product, "SEN55") {
		return SensorMeasurement{
			PM1_0:  uint(mass_concentration_pm1p0),
			PM2_5:  uint(mass_concentration_pm2p5),
			PM4_0:  uint(mass_concentration_pm4p0),
			PM10_0: uint(mass_concentration_pm10p0),
			Hum:    int(ambient_humidity),
			Temp:   int(ambient_temperature),
			VOC:    int(voc_index),
			NOx:    int(nox_index),
		}, nil
	} else if strings.Contains(product, "SEN54") {
		return SensorMeasurement{
			PM1_0:  uint(mass_concentration_pm1p0),
			PM2_5:  uint(mass_concentration_pm2p5),
			PM4_0:  uint(mass_concentration_pm4p0),
			PM10_0: uint(mass_concentration_pm10p0),
			Hum:    int(ambient_humidity),
			Temp:   int(ambient_temperature),
			VOC:    int(voc_index),
			NOx:    -math.MaxInt,
		}, nil
	} else {
		return SensorMeasurement{
			PM1_0:  uint(mass_concentration_pm1p0),
			PM2_5:  uint(mass_concentration_pm2p5),
			PM4_0:  uint(mass_concentration_pm4p0),
			PM10_0: uint(mass_concentration_pm10p0),
			Hum:    -math.MaxInt,
			Temp:   -math.MaxInt,
			VOC:    -math.MaxInt,
			NOx:    -math.MaxInt,
		}, nil
	}
	// return SensorMeasurement{
	// 	PM1_0:  uint(mass_concentration_pm1p0),
	// 	PM2_5:  uint(mass_concentration_pm2p5),
	// 	PM4_0:  uint(mass_concentration_pm4p0),
	// 	PM10_0: uint(mass_concentration_pm10p0),
	// 	Hum:    int(ambient_humidity),
	// 	Temp:   int(ambient_temperature),
	// 	VOC:    int(voc_index),
	// 	NOx:    int(nox_index),
	// }, nil
}

// Read / Write Temperature Compensation Parameters (0x60B2)
func ReadTemperatureCompensationParameters() (int, int, int, error) {
	var temp_offset C.short
	var slope C.short
	var time_constant C.ushort

	err := C.sen5x_get_temperature_offset_parameters((*C.short)(unsafe.Pointer(&temp_offset)), (*C.short)(unsafe.Pointer(&slope)), (*C.ushort)(unsafe.Pointer(&time_constant)))
	if err == 1 {
		return 0, 0, 0, fmt.Errorf("error reading temperature compensation parameters")
	}
	return int(temp_offset), int(slope), int(time_constant), nil
}
func SetTemperatureCompensationParameters(temp_offset int, slope int, time_constant int) error {
	err := C.sen5x_get_temperature_offset_parameters((*C.short)(unsafe.Pointer(&temp_offset)), (*C.short)(unsafe.Pointer(&slope)), (*C.ushort)(unsafe.Pointer(&time_constant)))
	if err == 1 {
		return fmt.Errorf("error reading temperature compensation parameters")
	}
	return nil
}

// Read/ Write Warm Start Parameter (0x60C6)
func ReadWarmStartParameter() (int, error) {
	var warm_start C.uint16_t
	err := C.sen5x_get_warm_start_parameter((*C.uint16_t)(unsafe.Pointer(&warm_start)))
	if err == 1 {
		return 0, fmt.Errorf("error reading warm start parameter")
	}
	return int(warm_start), nil
}
func SetWarmStartParameter(warm_start int) error {
	err := C.sen5x_set_warm_start_parameter((C.ushort)(warm_start))
	if err == 1 {
		return fmt.Errorf("error setting warm start parameter")
	}
	return nil
}

// Read/ Write VOC Algorithm Tuning Parameters (0x60D0)
func ReadVOCTuningParameters() (index_offset int, learning_time_offset_hours int, learning_time_gain_hours int, gating_max_duration_minutes int, std_initial int, gain_factor int, err error) {
	cerror := C.sen5x_get_voc_algorithm_tuning_parameters(
		(*C.short)(unsafe.Pointer(&index_offset)),
		(*C.short)(unsafe.Pointer(&learning_time_offset_hours)),
		(*C.short)(unsafe.Pointer(&learning_time_gain_hours)),
		(*C.short)(unsafe.Pointer(&gating_max_duration_minutes)),
		(*C.short)(unsafe.Pointer(&std_initial)),
		(*C.short)(unsafe.Pointer(&gain_factor)))
	if cerror == 1 {
		return 0, 0, 0, 0, 0, 0, fmt.Errorf("error reading VOC tuning parameters")
	}
	return int(index_offset), int(learning_time_offset_hours), int(learning_time_gain_hours), int(gating_max_duration_minutes), int(std_initial), int(gain_factor), nil
}
func SetVOCTuningParameters(index_offset int, learning_time_offset_hours int, learning_time_gain_hours int, gating_max_duration_minutes int, std_initial int, gain_factor int) error {
	cerror := C.sen5x_set_voc_algorithm_tuning_parameters(
		(C.short)(index_offset),
		(C.short)(learning_time_offset_hours),
		(C.short)(learning_time_gain_hours),
		(C.short)(gating_max_duration_minutes),
		(C.short)(std_initial),
		(C.short)(gain_factor))
	if cerror == 1 {
		return fmt.Errorf("error setting VOC tuning parameters")
	}
	return nil
}

// Read/ Write NOx Algorithm Tuning Parameters (0x60E1)
func ReadNOxTuningParameters() (index_offset, learning_time_offset_hours,
	learning_time_gain_hours, gating_max_duration_minutes, std_initial,
	gain_factor int, err error) {
	cerror := C.sen5x_get_nox_algorithm_tuning_parameters(
		(*C.short)(unsafe.Pointer(&index_offset)),
		(*C.short)(unsafe.Pointer(&learning_time_offset_hours)),
		(*C.short)(unsafe.Pointer(&learning_time_gain_hours)),
		(*C.short)(unsafe.Pointer(&gating_max_duration_minutes)),
		(*C.short)(unsafe.Pointer(&std_initial)),
		(*C.short)(unsafe.Pointer(&gain_factor)))
	if cerror == 1 {
		return 0, 0, 0, 0, 0, 0, fmt.Errorf("error reading NOx tuning parameters")
	}
	return int(index_offset), int(learning_time_offset_hours), int(learning_time_gain_hours), int(gating_max_duration_minutes), int(std_initial), int(gain_factor), nil
}
func SetNOxTuningParameters(index_offset, learning_time_offset_hours, learning_time_gain_hours, gating_max_duration_minutes, std_initial, gain_factor int) error {
	cerror := C.sen5x_set_nox_algorithm_tuning_parameters(
		(C.short)(index_offset),
		(C.short)(learning_time_offset_hours),
		(C.short)(learning_time_gain_hours),
		(C.short)(gating_max_duration_minutes),
		(C.short)(std_initial),
		(C.short)(gain_factor))
	if cerror == 1 {
		return fmt.Errorf("error setting NOx tuning parameters")
	}
	return nil
}

// Read/ Write RH/T Acceleration Mode (0x60F7)
func ReadAccelerationMode() (mode int, err error) {
	cerror := C.sen5x_get_rht_acceleration_mode(
		(*C.ushort)(unsafe.Pointer(&mode)))
	if cerror == 1 {
		return 0, fmt.Errorf("error reading acceleration mode")
	}
	return int(mode), nil
}
func SetAccelerationMode(mode int) error {
	cerror := C.sen5x_set_rht_acceleration_mode(
		(C.ushort)(mode))
	if cerror == 1 {
		return fmt.Errorf("error setting acceleration mode")
	}
	return nil
}

// Read/ Write VOC Algorithm State (0x6181)
func ReadVOCAlgorithmState() (state, statesize int, err error) {
	cerror := C.sen5x_get_voc_algorithm_state(
		(*C.uchar)(unsafe.Pointer(&state)),
		(C.uchar)(statesize))
	if cerror == 1 {
		return 0, 0, fmt.Errorf("error reading VOC algorithm state")
	}
	return int(state), int(statesize), nil
}
func SetVOCAlgorithmState(state, statesize int) error {
	cerror := C.sen5x_set_voc_algorithm_state(
		(*C.uchar)(unsafe.Pointer(&state)),
		(C.uchar)(statesize))
	if cerror == 1 {
		return fmt.Errorf("error setting VOC algorithm state")
	}
	return nil
}

// Start Fan Cleaning (0x5607)
func StartFanCleaning() error {
	err := C.sen5x_start_fan_cleaning()
	if err == 1 {
		return fmt.Errorf("error starting fan cleaning")
	}
	return nil
}

// Read/ Write Fan Cleaning Interval (0x5603)
func ReadFanCleaningInterval() (interval int, err error) {
	cerror := C.sen5x_get_fan_auto_cleaning_interval(
		(*C.uint)(unsafe.Pointer(&interval)))
	if cerror == 1 {
		return 0, fmt.Errorf("error reading fan cleaning interval")
	}
	return int(interval), nil
}
func SetFanCleaningInterval(interval int) error {
	cerror := C.sen5x_set_fan_auto_cleaning_interval(
		(C.uint)(interval))
	if cerror == 1 {
		return fmt.Errorf("error setting fan cleaning interval")
	}
	return nil
}

// Read Product Name (0xD014)
func ReadProductName() (string, error) {
	var productName [32]C.uchar
	cerror := C.sen5x_get_product_name((*C.uchar)(unsafe.Pointer(&productName)), 32)
	if cerror == 1 {
		return "", fmt.Errorf("error getting product name")
	}
	return string(productName[:]), nil

}

// Read Serial Number (0xD033)
func ReadSerialNumber() (string, error) {
	var serialNumber [32]C.uchar
	err := C.sen5x_get_serial_number((*C.uchar)(unsafe.Pointer(&serialNumber)), 32)
	if err == 1 {
		return "", fmt.Errorf("error getting serial number")
	}
	return string(serialNumber[:]), nil

}

// Read Firmware Version (0xD100)
func ReadFirmwareVersion() (string, error) {
	var firmware_major C.uchar
	var firmware_minor C.uchar
	var firmware_debug C.bool
	var hardware_major C.uchar
	var hardware_minor C.uchar
	var protocol_major C.uchar
	var protocol_minor C.uchar

	err := C.sen5x_get_version((*C.uchar)(unsafe.Pointer(&firmware_major)), (*C.uchar)(unsafe.Pointer(&firmware_minor)), (*C.bool)(unsafe.Pointer(&firmware_debug)), (*C.uchar)(unsafe.Pointer(&hardware_major)), (*C.uchar)(unsafe.Pointer(&hardware_minor)), (*C.uchar)(unsafe.Pointer(&protocol_major)), (*C.uchar)(unsafe.Pointer(&protocol_minor)))
	if err == 1 {
		return "", fmt.Errorf("error getting firmware version")
	}
	return fmt.Sprintf("Firmware: %d.%d\n Hardware %d.%d", int(firmware_major), int(firmware_minor), int(hardware_major), int(hardware_minor)), nil

}

// Read Device Status (0xD206)
func ReadDeviceStatus() (status int, err error) {
	cerror := C.sen5x_read_device_status((*C.uint)(unsafe.Pointer(&status)))
	if cerror == 1 {
		return 0, fmt.Errorf("error getting device status")
	}
	return int(status), nil
}

// Clear Device Status (0xD210)
func ClearDeviceStatus() error {
	var dummyStatus C.uint
	err := C.sen5x_read_and_clear_device_status((*C.uint)(unsafe.Pointer(&dummyStatus)))
	if err == 1 {
		return fmt.Errorf("error clearing device status")
	}
	return nil
}

// Device Reset (0xD304)
func ResetDevice() error {
	err := C.sen5x_device_reset()
	if err == 1 {
		return fmt.Errorf("error resetting device")
	}
	return nil

}

func InitializeHal() {
	C.sensirion_i2c_hal_init()
}

func FreeHal() {
	C.sensirion_i2c_hal_free()
}

func SleepHal(seconds uint32) {
	useconds := time.Duration(seconds) * time.Microsecond
	C.sensirion_i2c_hal_sleep_usec((C.uint32_t)(useconds))
	time.Sleep(time.Duration(seconds) * time.Second)
}
