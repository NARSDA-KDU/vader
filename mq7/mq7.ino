#include <PubSubClient.h>
#include <WiFi.h>
#include "DHT.h"

// WiFi configuration
const char* ssid = "Mustafar";
const char* password = "vader1234";

// MQTT configuration
const char* server = "192.168.4.2";
const int port = 1883;
const char* co2Topic = "sensors/co2";
const char* clientName = "vader.esp";
const char* tempTopic = "sensors/temperature";
const char* humTopic = "sensors/humidity";

// Variables to hold sensor readings
float temp;
float hum;

// Digital pin connected to the DHT sensor
#define DHTPIN 5
#define DHTTYPE DHT11 

// Initialize DHT sensor
DHT dht(DHTPIN, DHTTYPE);

WiFiClient wifiClient;
PubSubClient client(wifiClient);

// Time management
unsigned long previousMillis = 0;  // Store the last time a new reading was published
const long interval = 10000;       // Interval at which to publish sensor readings (milliseconds)

void wifiConnect() {
  Serial.println();
  Serial.print("Connecting to ");
  Serial.println(ssid);

  WiFi.begin(ssid, password);

  while (WiFi.status() != WL_CONNECTED) {
    delay(500);
    Serial.print(".");
  }
  Serial.println("");
  Serial.println("WiFi connected.");
  Serial.print("IP address: ");
  Serial.println(WiFi.localIP());
}

void mqttReConnect() {
  while (!client.connected()) {
    Serial.print("Attempting MQTT connection...");
    if (client.connect(clientName)) {
      Serial.println("connected");
    } else {
      Serial.print("failed, rc=");
      Serial.print(client.state());
      Serial.println(" try again in 5 seconds");
      delay(5000);
    }
  }
}

void mqttEmit(const char* topic, const char* payload) {
  client.publish(topic, payload);
}

void mqttEmit(String topic, String value) {
  client.publish((char*)topic.c_str(), (char*)value.c_str());
}

int COpin = A0; // analog input pin of MQ7 sensor
int COvalue = 0; // variable to store CO concentration value

void setup() {
  Serial.begin(115200);
  Serial.println();
  dht.begin();
  wifiConnect();
  client.setServer(server, port);
  delay(1500);
}

void loop() {
  if (!client.connected()) {
    mqttReConnect();
  }
  client.loop();

  unsigned long currentMillis = millis();

  // Every X number of seconds (interval = 10 seconds) 
  // it publishes a new MQTT message
  if (currentMillis - previousMillis >= interval) {
    // Save the last time a new reading was published
    previousMillis = currentMillis;

    int hum = 0;
    int temp = 0;
    // New DHT sensor readings
    hum = dht.readHumidity();
    temp = dht.readTemperature(); // Read temperature as Celsius (the default)

    // Publish an MQTT message on topic esp/dht/temperature
    mqttEmit(tempTopic, String(temp).c_str()); 
    Serial.print("Temperature: "); // Print Temperature label
    Serial.print(temp); // Print Temperature value
    Serial.println(" C"); // Print ppm unit                           

    // Publish an MQTT message on topic esp/dht/humidity
    mqttEmit(humTopic, String(hum).c_str());                            
    Serial.printf("Publishing on topic %s at QoS 1: ", humTopic);
    Serial.printf("Message: %.2f \n", hum);
  }

  COvalue = analogRead(COpin); // Read analog value from MQ7 sensor
  Serial.print("CO concentration: "); // Print CO concentration label
  Serial.print(COvalue); // Print CO concentration value
  Serial.println(" ppm"); // Print ppm unit

  mqttEmit(co2Topic, String(COvalue).c_str());
  delay(500);
}
