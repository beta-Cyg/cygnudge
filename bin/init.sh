#!/bin/sh
mkdir /var/lib/cygnudge
mkdir /var/lib/cygnudge/archive
mkdir /tmp/cygnudge
mkdir /bin/cygnudge
mkdir /etc/cygnudge
mkdir /etc/cygnudge/server
cp ./bin/cygpack.py /bin/cygnudge/
cp ./config/server/compile.json /etc/cygnudge/server/compile.json
cp ./config/server/server.json /etc/cygnudge/server/server.json
# cp ./config/server/server.json /etc/cygnudge_server.json
