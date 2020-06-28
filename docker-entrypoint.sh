#!/bin/bash
#/usr/sbin/init &

/usr/sbin/glusterd -p /var/run/glusterd.pid --log-level INFO &

sleep 2s

echo "Run docker-gluster-exporter"

/usr/bin/docker-gluester-exporter

sleep inf