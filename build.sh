#!/bin/sh
scp -r * go@corono-work.org.ua:/home/int/corono/corono-work
ssh go@corono-work.org.ua cd /home/int/corono/corono-work && /usr/local/go/bin/go build