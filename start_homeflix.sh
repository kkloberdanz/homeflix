#!/bin/bash
sshfs kyle@192.168.0.189:/home/kyle/data/movies /mnt/hdd/movies-rpi4/
sudo iptables -t nat -A PREROUTING -p tcp --dport 80 -j REDIRECT --to-port 8080
homeflix /mnt/hdd/movies /mnt/hdd/movies-rpi4 > /dev/null 2> /dev/null &
