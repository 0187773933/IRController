# IR Controller

- https://irdroid.com/irdroid-usb-ir-transceiver
- https://irdroid.com/irdroid-wifi-version-3-0
- https://www.lirc.org/software.html
- https://github.com/Irdroid/Irdroid-USB/blob/master/app/src/main/java/com/microcontrollerbg/usbirtoy/IrToy.java
- https://github.com/ww24/lirc-web-api/blob/master/lirc/client.go
- https://github.com/ww24/lirc-web-api/blob/master/lirc/test/lircd.conf
- https://lirc-remotes.sourceforge.net/remotes-table.html
- https://github.com/Irdroid/USB_IR_Toy
- https://irdroid.com/2016/10/how-to-turn-your-raspberry-pi-into-a-fully-functional-infrared-remote-control/
- https://anibit.com/product/ptt08001
- https://www.lirc.org/html/irtoy.html
- https://github.com/irdroid/usb_infrared_transceiver
- https://community.home-assistant.io/t/how-to-add-ir-tranmitter-to-enhance-your-audio-video-experince/446033


### Observe Protocol and Key Codes - V1

```
sudo ir-keytable -v -t -p rc-5,rc-5-sz,jvc,samsung,sony,nec,sanyo,mce_kbd,rc-6,sharp,xmp -s rc0
```

### Replay IR Codes - V1

```
ir-ctl -S necx:0x70702
```

---

### Observe Protocol and Key Codes - V2

```
ir-ctl -d /dev/lirc0 --receive=samnsung_power.key
```

### Intermediate convert_ir_ctl_recieve_format.sh

```
#!/bin/bash

if [ "$#" -ne 1 ]; then
	echo "Usage: $0 <filename>"
	exit 1
fi
output_file="samnsung_power_fixed.key"
> "$output_file"
content=$(cat $1)
for value in $content; do
	if [[ $value == +* ]]; then
		echo "pulse ${value#+}" >> "$output_file"
	elif [[ $value == -* ]]; then
		echo "space ${value#-}" >> "$output_file"
	fi
done
```

### Replay IR Codes - V2

```
ir-ctl -d /dev/lirc0 --send=samnsung_power_fixed.key
```