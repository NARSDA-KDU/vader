#include <PubSubClient.h>
#include <WiFi.h>



// Wifi configuration
const char* ssid = "TeDroid";
const char* password = "01010100";

const char* server = "192.168.43.131";
const char* topic = "sensors/co2";
const char* clientName = "vader.esp";



String payload;

WiFiClient wifiClient;
PubSubClient client(wifiClient);

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
  Serial.print("WiFi connected.");
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

void mqttEmit(String topic, String value)
{
  client.publish((char*) topic.c_str(), (char*) value.c_str());
}


int COpin = A0; // analog input pin of MQ7 sensor
int COvalue = 0; // variable to store CO concentration value

void setup() {
 Serial.begin(115200);
  wifiConnect();
  client.setServer(server, 1883);

   delay(1500);
}


void loop() {

  if (!client.connected()) {
    mqttReConnect();
  }

  COvalue = analogRead(COpin); // read analog value from MQ7 sensor
  Serial.print("CO concentration: "); // print CO concentration label
  Serial.print(COvalue); // print CO concentration value
  Serial.println(" ppm"); // print ppm unit

  mqttEmit(topic, (String) COvalue);
  delay(500);

}
