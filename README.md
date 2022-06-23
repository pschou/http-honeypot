# A very simple HTTP HoneyPot

This is to make bots go crazy for you.  :)

# Usage

To run the http honeypot, you can just run the binary and it will begin listening on port 8090.

```
$ ./http-honeypot -h
http-honeypot - make http scanner bots go crazy, Written by Paul Schou (github.com/pschou/http-honeypot), Version: 0.1.VERSION

Usage: http-honeypot [options...]
Options:
  --header STRING  What header to print for blank requests  (Default: "This is my server.")
  --listen PORT  Listen address  (Default: ":8090")
  --show_header  Should we show headers?
```


A simple systemd script can be written like this /usr/lib/systemd/system/http-honeypot.service
```
[Unit]
Description=HTTP HoneyPot
After=network.target

[Service]
Type=notify
ExecStart=/usr/bin/http-honeypot --header "My nifty server"
KillMode=process
Restart=on-failure
RestartSec=42s

[Install]
WantedBy=multi-user.target
```
