#!/bin/bash
apt install golang xterm git python3-pip upx-ucl -y
go get github.com/atotto/clipboard
go get github.com/bmatcuk/doublestar
go get github.com/dustin/go-humanize
go get github.com/kbinani/screenshot
go get github.com/lxn/win
go get github.com/mattn/go-shellwords
go get github.com/miekg/dns
go get golang.org/x/sys/windows
pip3 install -r requirements.txt
cd impacket
python setup.py install
