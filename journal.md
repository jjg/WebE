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

For now I've exposed an experimental website at [solar.jasongullickson.com](http://solar.jasongullickson.com/index.html) (which reminds me, it would be nice if jsfs handled default documents...).


### References

* https://wiki.dfrobot.com/Solar_Power_Manager_5V_SKU__DFR0559
* https://wiki.pine64.org/wiki/ROCK64#Expansion_Ports
* https://armbian.tnahosting.net/dl/rock64/archive/
* https://github.com/jjg/jsfs
* https://medium.com/gowombat/tutorial-how-to-use-ssh-tunnel-to-expose-a-local-server-to-the-internet-4e975e1965e5
* https://www.digitalocean.com/community/tutorials/iptables-essentials-common-firewall-rules-and-commands
