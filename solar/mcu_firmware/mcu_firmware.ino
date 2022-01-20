#include <avr/sleep.h>

int ledPin = 13;
int batteryPin = A0;
int regulatorPin = 9;
int batteryValue = 0;
int batteryPowerOnThreshold = 700;
int batteryPowerOffThreshold = 400;

void setup() {

  //Serial.begin(9600);
  //pinMode(batteryPin, INPUT);

}

void loop() {

  // Check battery voltage using pin A0
  batteryValue = analogRead(batteryPin);

  // Log battery voltage value to serial.
  //Serial.println(batteryValue);

  // TODO: Would probably be better to test the average of the last n measurements
  // to smooth-out spikes when the SBC load changes.
  if(batteryValue > batteryPowerOnThreshold){
    digitalWrite(regulatorPin, HIGH);
  }

  if(batteryValue < batteryPowerOffThreshold){
    digitalWrite(regulatorPin, LOW);
  }

  //  check again in one second (TODO: Change to one minute if not testing).
  delay(60000);

}
