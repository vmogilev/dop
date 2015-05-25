ps -ef | grep dop | grep -v grep | awk '{print $2}' | xargs kill -9
