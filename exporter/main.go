package main

import (
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/yosssi/gmq/mqtt"
	"github.com/yosssi/gmq/mqtt/client"
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
)

func main() {
	prometheus.MustRegister(co2Level, co2Histogram, co2ValueCounter)

	// MQTT connection configuration
	// Connect to the MQTT Server.
	// Create an MQTT Client.
	c := client.New(&client.Options{
		// Define the processing of the error handler.
		ErrorHandler: func(err error) {
			fmt.Println(err)
		},
	})

	// Terminate the Client.

	defer c.Terminate()
	err := c.Connect(&client.ConnectOptions{
		Network:  "tcp",
		Address:  "localhost:1883",
		ClientID: []byte("vader.exporter"),
	})
	if err != nil {
		panic(err)
	}
	defer c.Disconnect()

	// Subscribe to the MQTT topic
	topic := "sensors/co2"
	err = c.Subscribe(&client.SubscribeOptions{SubReqs: []*client.SubReq{
		&client.SubReq{
			TopicFilter: []byte(topic),
			QoS:         mqtt.QoS0,
			Handler: func(topicName, message []byte) {
				// Assuming the message contains CO2 level value
				co2Value := string(message)
				slog.Info("Received CO2 level: " + co2Value)

				co2Val, _ := strconv.Atoi(co2Value)

				co2Histogram.Observe(float64(co2Val))
				co2Level.Set(float64(co2Val))
				co2ValueCounter.Inc()
			},
		},
	}})
	if err != nil {
		log.Fatal(err)
	}

	// Start HTTP server to expose Prometheus metrics
	slog.Info("starting exporter")
	http.Handle("/metrics", promhttp.Handler())
	log.Fatal(http.ListenAndServe(":8080", nil))
}
