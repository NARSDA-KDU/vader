# python 3.11


## for local testing purposes 
import random

import time 

from paho.mqtt import client as mqtt_client


broker = 'mqtt'
port = 1883
topic = "sensors/co2"
# Generate a Client ID with the subscribe prefix.
client_id = f'subscribe-{random.randint(0, 100)}'


def connect_mqtt() -> mqtt_client:
    def on_connect(client, userdata, flags, rc):
        if rc == 0:
            print("Connected to MQTT Broker!")
        else:
            print("Failed to connect, return code %d\n", rc)

    client = mqtt_client.Client(client_id)
    # client.username_pw_set(username, password)
    client.on_connect = on_connect
    client.connect(broker, port)
    return client


def subscribe(client: mqtt_client):
    def on_message(client, userdata, msg):
        print(f"Received `{msg.payload.decode()}` from `{msg.topic}` topic")

    client.subscribe(topic)
    client.on_message = on_message




def publish(client):
    while True:
        co2_value = random.randint(300, 1000)  # Random CO2 value between 300 and 1000 ppm
        result = client.publish(topic, co2_value)
        # result: [0, 1]
        status = result[0]
        if status == 0:
            print(f"Published CO2 value: {co2_value}")
        else:
            print(f"Failed to publish CO2 value with status code {status}")
        time.sleep(3)  # Adjust the delay as needed

def run():
    client = connect_mqtt()
    # subscribe(client)
    publish(client)
    client.loop_forever()


if __name__ == '__main__':
    run()
