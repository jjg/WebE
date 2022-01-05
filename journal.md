# WebE Project Journal

## 01032022

Beginning work on the [solar](./solar)-powered project.

### Hardware

5VDC [Solar Panel](https://www.dfrobot.com/product-1775.html)l -> DFRobot [Solar Power Manager 5V](https://wiki.dfrobot.com/Solar_Power_Manager_5V_SKU__DFR0559) -> [Rock64](https://wiki.pine64.org/wiki/ROCK64)

### Software 

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

My initial thought is to "hijack" the line connected to the `CHG` and `DONE` LEDs, and feed this to the GPIO on the SBC.  There would need to be some logic to determine if the sun has gone away and to initiate shutdown while power remains in the battery.  I think it would work like this:

| CHG | DONE | OUT  | 
|-----|------|------|
| ON  | OFF  | HIGH |
| ON  | ON   | HIGH |
| OFF | ON   | HIGH |
| OFF | OFF  | LOW  |

I'm not sure if I want to implement this logic externally or bring both signals back to the SBC and do the logic there.  The former seems like the right thing to do, but the latter seems simpler (less electronics work, more software work).  I guess it will come down to how much hacking I have to do to get at the signals (if I have to a lot of soldering to just get the signal it might not be much more work to implement the logic at the same time).

Software-wise the next thing we can do is open-up access to the node to the public Internet.  The most reliable way to do this is to use an ssh tunnel, so let's give that a try.

1. Create a new ssh key on webe-solar: `ssh-keygen`
2. Copy the key to some server on the public Internet
3. Open a tunnel on webe-solar: `ssh -N -i /home/jason/.ssh/id_rsa_tunnel jason@theneuromantics.net -R 2022:localhost:7302 -C`
4. Test the connection from the Internet server: `curl -v http://localhost:2022/`
5. Make sure Gateway Ports: yes in /etc/ssh/sshd_config
6. Make sure the port isn't blocked by the firewall
  + `iptables -I INPUT -p tcp --dport 2022 -m conntrack --ctstate NEW,ESTABLISHED -j ACCEPT`
    + TODO: Make this persistent once a persistent host is chosen
7. Test connection from the public internet: `curl -v http://68.183.206.69:2022/`
8. Automate!
  + `sudo apt install autossh`
  + Connect with `autossh` from root once to establish key validity, etc: `sudo autossh -M 20000 -N -i /home/jason/.ssh/id_rsa_tunnel jason@theneuromantics.net -R 2022:localhost:7302 -C`
  + Create systemd service file
  + `sudo systemctl daemon-reload`
  + `sudo systemctl start jsfs-tunnel`
  + Test connecting from the public Internet
  + `sudo systemctl enable jsfs-tunnel`
  + TODO: Create service file to start `jsfs` and *then* start `jsfs-tunnel`

OK, at this point I have both JSFS and the tunnel starting automatically at boot.  Once I hack the electronics to tell the SBC the sun is gone (and write a script to read this and shut the SBC down) we should be really close to making this usable for some testing.

...well there is also the matter of booting things back up when the sun comes back.  Hmm...

For now I've exposed an experimental website at [solar.jasongullickson.com](http://solar.jasongullickson.com:2022/index.html) (which reminds me, it would be nice if jsfs handled default documents...).


### References

* https://wiki.dfrobot.com/Solar_Power_Manager_5V_SKU__DFR0559
* https://wiki.pine64.org/wiki/ROCK64#Expansion_Ports
* https://armbian.tnahosting.net/dl/rock64/archive/
* https://github.com/jjg/jsfs
* https://medium.com/gowombat/tutorial-how-to-use-ssh-tunnel-to-expose-a-local-server-to-the-internet-4e975e1965e5
* https://www.digitalocean.com/community/tutorials/iptables-essentials-common-firewall-rules-and-commands


## 01042022

Last night I was studying the charge controller's schematics and learned that I can actually read the battery voltage (the `BAT` pin next to the `EN` pin) and use it as a proxy for [state-of-charge](https://en.wikipedia.org/wiki/State_of_charge).  This is exciting because I don't have to try and solder new leads to the surface mount LEDs or come up with complex logic to decide when to power-down the SBC.  The downside is that the signal is analog, and the SBC I'm using doesn't have any analog inputs (most Linux SBC's don't) so if I want to use this I'm going to have to either add analog input to the SBC (probably using a [MCP3008](https://www.microchip.com/en-us/product/MCP3008)) or adding another device that has an analog input that can either control the power supply, signal the SBC to control it or maybe both.

I'm leaning toward adding the MCP3008 because there may be other analog signals we'll want to access in the future (temp, light level, etc.) and I really don't want to have yet another codebase (even if it's small), toolchain, etc. to program a microcontroller to do this.  Of course the counterargument is that if a mcu is used, it could do things independent of the SBC (and do so using much less power) so I'm still undecided.  In the long run I can see a custom charge controller that includes an MCU integrated into the board to provide this sort of "supervisor" processor but we're not there yet.

If I don't want to wait for parts I might be able to cobble-together a basic comparitor circuit that could be read like a digital signal (basically an [open collector](https://en.wikipedia.org/wiki/Open_collector) connected to a GPIO pin on the SBC) which would be good enough to tell the SBC when to shutdown.  This might also work to determine when to power the SBC back up, because it could control the 5v out pins on the charge controller (via the `EN` pin).  It might require two comparitors, one to send the "time to shut down" signal and another to actually cut the juice, but I'm pretty sure I have parts on-hand to make something like this to keep moving forward until I can get some MCP3008's.

```
SPM     SBC       comparator

5V  ->  5V
5V
5V
GND ->  GND
GND
GND

GND ------------> gnd
EN
BAT ------------> base 

        GPIO  <-  collector
        GND -----> gnd

```

Hmm... looking at the various DIY comparator options I'm leaning toward just ordering some MCP3008's and going that way, but if something else comes up in the meantime I'll experiment.

Preston suggested the idea of just letting it die when it runs out of power.  I thought this might be more complicated than doing a gracefull shutdown (given how often I've seen filesystem/SD card corruption with SBC's that get turned off like this) but given how much trouble detecting the power supply state has turned out to be, this might be worth looking into.

I know that Armbian already stores its logs in a "RAM disk", so at least those won't be writing to disk when we pull the plug, but I'm going to have to poke-around a bit to see what other measures could be taken to allow one of these boards to simply be unplugged and then plugged back in without messing up the O/S.

After digging a bit it turns out that the solution was obvious: run Linux in RAM like a USB boot drive.  Turns out this is the default mode of operation for [Alpine Linux on Raspberry Pi](https://wiki.alpinelinux.org/wiki/Raspberry_Pi) so this might just do the trick (assuming you're using a pi).  In my case I'm using a different SBC at the moment, but I'll give it a try and see what happens.

If it doesn't work I might just have to wait until I have a Pi to try.

OK, looks like Alpine's [Generic ARM](https://dl-cdn.alpinelinux.org/alpine/v3.15/releases/aarch64/alpine-uboot-3.15.0-aarch64.tar.gz) package might work with [one of the boards I have](https://wiki.alpinelinux.org/wiki/Pine64_A64_LTS).  This would be a lot easier if I could just use a Raspberry Pi, but apparently there are none for sale in the world right now...


## 01052022

Thought of a fairly simple way to provide access to the nodes as they come and go with the sunshine.

![](solar/docs/public_proxy.png)

By adding [HA Proxy](http://www.haproxy.org/) to the existing public host (the place the SSH tunnel terminates) we can setup each node to open an ssh tunnel to the host and then loadbalance requests from a single public address across all of the solar nodes that happen to be up at the time.  HA Proxy can be configured to perform a health check against each configured node so when nodes go off/online, incoming requests are automatically routed to whatever nodes are up at the time.

Ideally this would be 100% dynamic (nodes join using auto-generated ports with no up-front configuration of HA Proxy) but I don't know how to do that yet so it will have to be configured manually for now.

Since I don't want to break the host I'm using for solar.jasongullickson.com by installimg HA Proxy (it's already running other things), I think its time to setup a dedicated public host for this project.  That will also let me put the public interface on port 80 which most clients expect.


### Next Steps

1. Create a Raspberry Pi Alpine SD card
2. Boot a Raspberry Pi 3 Model A (+?) with the card and see if we can set it up over serial
3. Setup JSFS and test JSFS API locally, then over WiFi
4. Setup ssh tunnel to existing public host (new port)
5. Test public access
6. Test pulling the plug
  + Maybe setup a scheduled curl to poll the system as it goes up and down?
7. If power failure recovery looks good, swap the boards
  + Keep the Rock64 online, we'll use it to test the load balancer & federation later

Once this works we can start setting up the new public host with HA proxy, etc.
