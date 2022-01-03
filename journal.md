# WebE Project Journal

## 01032022

Beginning work on the [solar](./solar)-powered project.

### Hardware

5VDC [Solar Pane](https://www.dfrobot.com/product-1775.html)l -> DFRobot [Solar Power Manager 5V](https://wiki.dfrobot.com/Solar_Power_Manager_5V_SKU__DFR0559) -> [Rock64](https://wiki.pine64.org/wiki/ROCK64)

### O/S

1. Burn Armbian to SD card
2. Connect serial console cable
  + `picoterm /dev/ttyUSB0 -b1500000`
3. Boot & run-through initial setup
  + Personal -> Hostname -> `webe-solar.local`
  + System -> CPU -> set min and max speed to slowest option, conservative governor
4. `reboot`
5. Install JSFS
  + `sudo apt install git, vim, nodejs, npm`
  + `git clone https://github.com/jjg/jsfs.git`
  + Walk-through the [configuration steps](https://github.com/jjg/jsfs#configuration) in the README

At this point we can sucessfully make JSFS requests.

There was just enough sunlight to perform this initial test on actual solar power.  What's needed next to run this safely is a way to shut-down cleanly when the sun goes away.  Sadly the charge controller board doesn't have any output to indicate the charging state or amount of energy left in the battery, so this will require some hacking.

My initial thought is to "hijack" the line connected to the `CHG` and `DONE` LEDs, and feed this to the GPIO on the SBC running JSFS.  There would need to be some logic to determine if the sun has gone away and to initiate shutdown while power remains in the battery.  I think it would work like this:

| CHG | DONE | OUT  | 
|-----|------|------|
| ON  | OFF  | HIGH |
| ON  | ON   | HIGH |
| OFF | ON   | HIGH |
| OFF | OFF  | LOW  |

I'm not sure if I want to implement this logic externally or bring both signals back to the SBC and do the logic there.  The former seems like the right thing to do, but the latter seems simpler (less electronics work, more software work).  I guess it will come down to how much hacking I have to do to get at the signals (if I have to a lot of soldering to just get the signal it might not be much more work to implement the logic at the same time).


### References

* https://wiki.dfrobot.com/Solar_Power_Manager_5V_SKU__DFR0559
* https://wiki.pine64.org/wiki/ROCK64#Expansion_Ports
* https://armbian.tnahosting.net/dl/rock64/archive/
* https://github.com/jjg/jsfs
