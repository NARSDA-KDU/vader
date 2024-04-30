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
			Help: "Current CO2 level",
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

	co2ValueCounter = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "co2_value_counter",
			Help: "Number of CO2 values received",
		},
	)

	mqttServer = os.Getenv("MQTT_SERVER")
)

func main() {
	prometheus.MustRegister(co2Level, co2Histogram, co2ValueCounter)

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

	if token.Wait() && token.Error() != nil {
		log.Fatal(token.Error())
	}

	// Start HTTP server to expose Prometheus metrics
	slog.Info("starting exporter")
	http.Handle("/metrics", promhttp.Handler())
	log.Fatal(http.ListenAndServe(":9093", nil))
}
