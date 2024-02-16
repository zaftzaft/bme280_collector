package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/alecthomas/kingpin/v2"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/quhar/bme280"
	"golang.org/x/exp/io/i2c"
)

var (
	i2cbus   = kingpin.Flag("i2cbus", "i2c bus device").Short('y').Default("/dev/i2c-1").String()
	i2caddr  = kingpin.Flag("addr", "i2c chip address. default 0x76(118)").Short('a').Default("118").Int()
	filename = kingpin.Flag("output", "output filename").Short('o').Default("/tmp/node_exporter/bme280.prom").String()
	mkdirp   = kingpin.Flag("mkdirp", "mkdirp").Short('p').Bool()
)

func Run() int {
	registry := prometheus.NewRegistry()

	bme280TemperatureGauge := prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "bme280_temperature",
		Help: "Temperature measured by BME280 sensor",
	})

	bme280PressureGauge := prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "bme280_pressure",
		Help: "Pressure measured by BME280 sensor",
	})

	bme280HumidityGauge := prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "bme280_humidity",
		Help: "Humidity measured by BME280 sensor",
	})

	registry.MustRegister(bme280TemperatureGauge)
	registry.MustRegister(bme280PressureGauge)
	registry.MustRegister(bme280HumidityGauge)

	device, err := i2c.Open(&i2c.Devfs{Dev: *i2cbus}, *i2caddr)
	if err != nil {
		fmt.Println(err)
		return 1
	}

	b := bme280.New(device)
	if err := b.Init(); err != nil {
		fmt.Println(err)
		return 1
	}

	temperature, pressure, humidity, err := b.EnvData()
	if err != nil {
		fmt.Println(err)
		return 1
	}

	bme280TemperatureGauge.Set(temperature)
	bme280PressureGauge.Set(pressure)
	bme280HumidityGauge.Set(humidity)

	if *mkdirp {
		if err := os.MkdirAll(filepath.Dir(*filename), 0755); err != nil {
			fmt.Println(err)
			return 1
		}
	}

	if err := prometheus.WriteToTextfile(*filename, registry); err != nil {
		fmt.Println(err)
		return 1
	}

	return 0
}

func main() {
	kingpin.Version("0.0.2")
	kingpin.HelpFlag.Short('h')
	kingpin.Parse()
	os.Exit(Run())
}
