# DOP - Day One Parser

I built [Day One](http://dayoneapp.com/) Parser to host my Food Log online and to share it with my nutritionist.  However it can be also used to run a blog powered by Day One Journal.  Here is what it does:

* displays Day One Tags as a list of color coded slugs (tag colors are defined in `CssLookup` see conf/dop.json warning=yellow, danger=red, success=green)
* places Day One Image at the top of the entry to which it was attached
* counts the number of "Search-String" in Day One entry and displays it in a badge (Search-String is defined in `Count` see conf/dop.json)
* decodes a starred entry and displayes it as a YES/no slug

DOP also supports custom Title, URL and Description for each entry as follows:
* if First line of the entry begins with "# " what follows is used for Title
* if Second line of the entry begins with "//dop:desc" what follows  is used for Description
* if Third line of the entry begins with "//dop:link" what follws is used for url link

Example:

    # Custom Entry Title
    //dd: Longer Description of the entry
    //dl: my-entry-url.html



![dop screen shot](https://s3.amazonaws.com/mve-shared/dop1.png)



## Installation

    GOPATH=/data/app/dev/golang; export GOPATH
    GO=/usr/local/go/bin/go; export GO

    $GO get github.com/rigingo/dop
    $GO install github.com/rigingo/dop

## Configuration

Modify conf/dop.json and define the following variables:
* `Title`: Site's Top Level Title
* `Desc`: Site's description - maps to `<meta name="description" content="">` on top level page/root
* `PubStarred`: `true|false` - if `true` only Starred entries are published
* `Count`: Only used with `dop_food` template - counts occurances of string and sets result in a badge
* `CssLookup`: map of `tag` -> `css-tag` - allowed values [`success` `warning` `danger`]

Modify conf/dop.env and define HTTP Host variables and JDIR - path to your DayOne journal directory where `entries` subdirectory is located (see example configuration)

## Startup

Source the conf/dop.env:

    . conf/dop.env

Start it:

    ./start.sh

## Stop

    ./stop.sh

## License

Apache License Version 2.0

*rev:   203 *

