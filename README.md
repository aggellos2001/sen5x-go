# sen5x-go

> [!NOTE]  
> This project is not affiliated with Sensirion in any way shape or form!

This project read and writes the measurements to a CSV file and sends them to the cloud to [sensor.community (optionally)](https://sensor.community/en/)
from a connected SEN5x sensors family through the I2C protocol.

The project is powered by a CGO wrapper that calls the C functions provided by Sensirion, and you can find here:
[https://github.com/aggellos2001/sen5xlib](https://github.com/aggellos2001/sen5xlib).


#### Why Go?

Go is very fast and easily to write language that can be compiled to a single binary, unlike Python which needs a lot of
dependencies and also needs to be installed on the target system. With Go, you can just copy the binary to your Raspberry
Pi and run it. No need to install anything else.


### Compatibility

Tested with a Raspberry Pi Zero 2W, but should work with the PI 4 as well.

> [!WARNING]
> **Using the 64-bit version of Raspberry PI OS is required to run t\he precompiled binary.**


## Installation

The installation is very easy. Just copy the binary to your Raspberry Pi and run it. You can find the latest binary in
the GitHub releases page.


### systemd service
To run the binary automatically, even after a power loss, you can easily create a systemd service. Just create a file in
`/etc/systemd/system/sen5x.service` with the following content:

```txt
[Unit]
Description=Runs the data collection script for the SEN55, even after a reboot of the system.
After=network.target
Wants=network-online.target

[Service]
WorkingDirectory=PATH-TO-FOLDER-WITH-EXECUTABLE
Restart=always
Type=simple
ExecStart=PATH-TO-BINARY
 
[Install]
WantedBy=multi-user.target
```

Make sure to change replace the path placeholders to the correct location on your system.

After that you need to reload the systemd daemon with the command
```bash
sudo systemctl daemon-reload
```
and enable the service with
```bash
sudo systemctl enable sen5x.service
```

Make sure to start the service also for the first time with
```bash 
sudo systemctl start sen5x.service
```

If you run the script as a service, and you want to see the logs, you can use the command `sudo journalctl -u sen5x.service -f`
to see the logs in real time.

## Configuration

This script is higly configurable. It uses the TOML format, which you can learn more [here](https://toml.io/en/v1.0.0).

When you run the script for the first time, it will create a file called `config.toml` in the same directory as the
binary. You can edit this file to change the configuration of the script.

[The default **config.toml** configuration with some comments to help you understand what each option does is included in the repository.](config.toml)


## Sensor Community

Here's a short guide on how to get your SensorNodeID from sensor.community.

Firstly, you will need to make an account on [https://devices.sensor.community/](https://devices.sensor.community/).

Then, go to **My sensors** tab and click **Register new sensor**.

![Register new sensor](/image/1.png)

On the *Sensor ID* field, you must put something unique. It's recommended to put the serial number of your Raspberry Pi.

![Register new sensor](/image/2.png)

You can find easily your serial number by running:

```bash
cat /proc/cpuinfo | grep Serial | cut -d ' ' -f 2
```

![Getting the device serial](/image/3.png)

The rest of the options are self-explanatory, except the hardware configuration. For this script, you need to pick a
**SPS30** and a **SHT30** component.

The **SPS30** will be used to upload the particulate matter data and the **SHT30** will be used to upload the temperature
and humidity data.

![Register new sensor](/image/4.png)

Currently, the NOx and VOC data are not supported by the script to upload to sensor.community.

After you click **Save settings** your new sensor node will appear on your list of sensors in the **My sensors** tab.

All you need is to copy the **SensorUID** and put in the **config.toml** file in the **SensorCommunity** section (and of
course set the **Enabled** option to true).

![Register new sensor](/image/5.png)

After that, you can run the script, and it will upload the data to sensor.community.

## Updates

Check the releases page for updates of the script. When a new version contains changes to the config.toml file, it will
automatically update and remove or add default values to the config.toml file as need, so make sure to read the release
notes and check the config yourself after an update.

## Cross compiling for Raspberry Pi

If you want to cross compile the script from the source for the Raspberry Pi you can use the following command:

```bash
CGO_ENABLED=1 GOOS=linux GOARCH=arm64 CC=aarch64-linux-gnu-gcc CC_FOR_TARGET=gcc-aarch64-linux-gnu go build -o bin/sen5x-bin .
```

For this to work you need to have the appropriate compiler installed and also your system should have a compatible libc
version. Currently, Raspberry Pi OS is based on Debian 11, so you must have a GLIBC version 2.31 or older.

## Versioning

This project uses [Go's module version numbering rules.](https://go.dev/doc/modules/version-numbers)

## Contributing

If you want to contribute to this project, feel free to open a pull request. I will be happy to review it. If you have
any questions, feel free to open an issue and make sure to follow the appropriate template.