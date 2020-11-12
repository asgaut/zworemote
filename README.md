# zworemote
Web interface for ZWO ASI Cameras (tested on ASI1600MM) and EFW mini filter wheel.

# Build and Run

- Install /usr/local/lib/libusb-1.0.0.dylib
  RPi: sudo apt install libusb-1.0-0-dev libudev-dev
- cd src
- go build zworemote.go
- ./zworemote
- Visit http://localhost:8080/zworemote/

Capture three image sets (written to `/tmp`) with exposure 150, gain 150, one image in 
each set prefixed r, g, b respectively, at 16 bit, with filter positions 1, 2, 3 respectively:

```shell
caffeinate curl "http://localhost:8080/zworemote/series?&e=150&g=150&n=1&p=r&d=16&f=1";
caffeinate curl "http://localhost:8080/zworemote/series?&e=150&g=150&n=1&p=g&d=16&f=2";
caffeinate curl "http://localhost:8080/zworemote/series?&e=150&g=150&n=1&p=b&d=16&f=3"
```
