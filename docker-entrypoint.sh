#!/bin/bash
#/usr/sbin/init &

/usr/sbin/glusterd -p /var/run/glusterd.pid --log-level INFO &

sleep 2s

echo "Run docker-gluster-prometheus"

/usr/bin/docker-gluster-prometheus

sleep inf