



int COpin = A0; // analog input pin of MQ7 sensor
int COvalue = 0; // variable to store CO concentration value

void setup() {
  Serial.begin(9600); // initialize serial communication
}

void loop() {
  COvalue = analogRead(COpin); // read analog value from MQ7 sensor
  Serial.print("CO concentration: "); // print CO concentration label
  Serial.print(COvalue); // print CO concentration value
  Serial.println(" ppm"); // print ppm unit
  delay(1000); // wait for 1 second
}

