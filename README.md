# DOP - Day One Parser

I built [Day One](http://dayoneapp.com/) Parser to host my Food Log online and to share it with my nutritionist.  However it can be used for just about anything.  Here is what it does:

* displays Day One Tags as a list of color coded slugs (tag colors are defined in `CssLookup` see conf/dop.json warning=yellow, danger=red, success=green)
* places Day One Image at the top of the entry to which it was attached
* counts the number of "Search-String" in Day One entry and displays it in a badge (Search-String is defined in `Count` see conf/dop.json)
* decodes a starred entry and displayes it as a YES/no slug

![dop screen shot](https://s3.amazonaws.com/mve-shared/dop1.png)



## Installation

    GOPATH=/data/app/dev/golang; export GOPATH
    GO=/usr/local/go/bin/go; export GO

    $GO get github.com/rigingo/dop
    $GO install github.com/rigingo/dop

## Configuration

Modify conf/dop.json and define `Dir` to point to your DayOne journal directories where `entries` subdirectory is located (do not add trailing slash).

Modify conf/dop.env and define HTTP Host variables (see example configuration)

## Startup

Source the conf/dop.env:

    . conf/dop.env

Start it:

    ./start.sh

## Stop

    ./stop.sh

## License

Apache License Version 2.0
