package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"golang.org/x/exp/slog"
)

var (
	co2Level = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Name: "co2_level",
			Help: "Current CO2",
		},
	)
	HumiditiyLevel = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Name: "humidity_level",
			Help: "Current Humidity",
		},
	)
	TemperatureLevel = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Name: "temperature_level",
			Help: "Current temperature",
		},
	)

	co2Histogram = prometheus.NewHistogram(
		prometheus.HistogramOpts{
			Name: "co2_level_histogram",
			Help: "CO2 level histogram",
			// Define custom buckets for CO2 levels
			Buckets: []float64{300, 400, 500, 600, 700, 800, 900, 1000},
		},
	)

	temperatureHistogram = prometheus.NewHistogram(
		prometheus.HistogramOpts{
			Name: "temperature_level_histogram",
			Help: "temperature level histogram",
			// Define custom buckets for CO2 levels
			Buckets: []float64{300, 400, 500, 600, 700, 800, 900, 1000},
		},
	)

	humidityHistogram = prometheus.NewHistogram(
		prometheus.HistogramOpts{
			Name: "humidity_level_histogram",
			Help: "humidity level histogram",
			// Define custom buckets for CO2 levels
			Buckets: []float64{300, 400, 500, 600, 700, 800, 900, 1000},
		},
	)

	humidityCounter = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "humidity_value_counter",
			Help: "Number of humidity values received",
		},
	)

	temperatureCounter = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "temperature_value_counter",
			Help: "Number of temperature values received",
		},
	)

	co2ValueCounter = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "co2_value_counter",
			Help: "Number of CO2 values received",
		},
	)

	mqttServer = os.Getenv("MQTT_SERVER")
)

func main() {
	prometheus.MustRegister(co2Level, co2Histogram, co2ValueCounter, temperatureCounter, temperatureHistogram, temperatureCounter, HumiditiyLevel, humidityHistogram, humidityCounter)

	// MQTT connection configuration
	var broker = "mqtt-svc"
	var port = 1883
	opts := mqtt.NewClientOptions()
	opts.AddBroker(fmt.Sprintf("tcp://%s:%d", broker, port))

	client := mqtt.NewClient(opts)

	if token := client.Connect(); token.Wait() && token.Error() != nil {
		panic(token.Error())
	}
	defer client.Disconnect(250) // Disconnect gracefully after 250 milliseconds

	// Subscribe to the MQTT topic
	topic := "sensors/co2"
	token := client.Subscribe(topic, 0, func(client mqtt.Client, msg mqtt.Message) {
		// Assuming the message contains CO2 level value
		co2Value := string(msg.Payload())
		slog.Info("Received CO2 level: " + co2Value)

		co2Val, _ := strconv.Atoi(co2Value)

		co2Histogram.Observe(float64(co2Val))
		co2Level.Set(float64(co2Val))
		co2ValueCounter.Inc()
	})

	// Subscribe to the MQTT topic
	humiditytopic := "sensors/humidity"
	_ = client.Subscribe(humiditytopic, 0, func(client mqtt.Client, msg mqtt.Message) {
		// Assuming the message contains humidity level value
		humidityValue := string(msg.Payload())
		slog.Info("Received humidty level: " + humidityValue)

		humidityVal, _ := strconv.Atoi(humidityValue)

		humidityHistogram.Observe(float64(humidityVal))
		HumiditiyLevel.Set(float64(humidityVal))
		humidityCounter.Inc()
	})

	// Subscribe to the MQTT topic
	temperaturetopic := "sensors/temperature"
	_ = client.Subscribe(temperaturetopic, 0, func(client mqtt.Client, msg mqtt.Message) {
		// Assuming the message contains temperature level value
		temperatureValue := string(msg.Payload())
		slog.Info("Received temprature level: " + temperatureValue)

		temperatureVal, _ := strconv.Atoi(temperatureValue)

		temperatureHistogram.Observe(float64(temperatureVal))
		HumiditiyLevel.Set(float64(temperatureVal))
		temperatureCounter.Inc()
	})

	if token.Wait() && token.Error() != nil {
		log.Fatal(token.Error())
	}

	// Start HTTP server to expose Prometheus metrics
	slog.Info("starting exporter")
	http.Handle("/metrics", promhttp.Handler())
	log.Fatal(http.ListenAndServe(":9093", nil))
}
