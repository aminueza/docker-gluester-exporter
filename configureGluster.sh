#!/bin/bash
gluster peer probe server2 
gluster peer probe server3
gluster volume create gfs replica 3 server1:/data/brick server2:/data/brick server3:/data/brick force 
gluster volume start gfs