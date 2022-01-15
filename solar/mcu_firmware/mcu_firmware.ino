#include <avr/sleep.h>

int ledPin = 13;
int batteryPin = A0;
int regulatorPin = 9;
int batteryVoltage = 0;
int batteryChargedVoltage = 128;  // TODO: Figure out how this relates to actual voltage.

void setup() {

  pinMode(batteryPin, OUTPUT);

}

void loop() {

  // Check battery voltage using pin A0
  batteryVoltage = analogRead(batteryPin);

  if(batteryVoltage > batteryChargedVoltage){
    // If battery is charged enough, turn-on pin D9.
    digitalWrite(regulatorPin, HIGH);

    // Turn-off the LED to indicate that we think the battery is ready.
    digitalWrite(ledPin, LOW);
    
    // TODO: Reduce power as much as possible w/o letting pin 9 go low.
    /*
    set_sleep_mode(SLEEP_MODE_PWR_DOWN);
    sleep_enable();
    ADCSRA &= ~(1 << ADEN);
    //PRR = 0xFF;
    sleep_mode();
    */
  } else {
    // If battery is not charged enough, turn-on the LED to indicate that we're still waiting for the battery to charge.
    digitalWrite(ledPin, HIGH);
    
    //  check again in one second (TODO: Change to one minute if not testing).
    delay(1000);
  }

}
