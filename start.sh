#!/bin/bash

err() {
    echo "ERROR: ${1} - exiting"
    exit 1
}

if [ -f "$GOPATH/bin/dop" ]; then
    echo "OK: found $GOPATH/bin/dop"
else
    err "dop binary is missing in $GOBIN/dop"
fi

if [ -f "$DOPROOT/conf/dop.json" ]; then
    echo "OK: found dop conf file in $DOPROOT/conf/dop.json"
else
    err "dop conf file is missing in $DOPROOT/conf/dop.json"
fi

if [ -d "$JDIR/entries" ]; then
    echo "OK: found entries in $JDIR/entries"
else
    err "Day One Journal entries are in $JDIR/entries"
fi


nohup $GOPATH/bin/dop \
    -jDir="${JDIR}" \
    -jTemplate="${JTEMPLATE}" \
    -dopRoot="${DOPROOT}" \
    -httpHost="${HTTPHOST}" \
    -httpMount="${HTTPMOUNT}" \
    -httpPort="${HTTPPORT}" \
    -httpHostExt="${HTTPHOSTEXT}" >> ${DOPROOT}/server.log 2>&1 </dev/null &


