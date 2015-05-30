#!/bin/bash

##pidof dop | xargs kill -9
ps -ef | grep "\-dopRoot" | grep -v grep | awk '{print $2}' | xargs kill -9
