# cbfullofit

A full-text indexer for Couchbase.

## Prerequisites

There are 2 libraries that must be available prior to building the project.

* icu4c (recommend version > 52) - [http://site.icu-project.org/](http://site.icu-project.org/)
* leveldb (must be version > 1.6) - [https://code.google.com/p/leveldb/](https://code.google.com/p/leveldb/)

If these libraries are installed in the standard system locations the following commands should work.  If the libraries are installed in custom locations, you may need to use additional environment variables: `CGO_CFLAGS="-I/path/to/leveldb/include CGO_LDFLAGS="-L/path/to/leveldb/lib`

## Installation

go get github.com/couchbaselabs/cbfullofit

## Running

cbfullofit requires the use of a Couchbase bucket to coordinate its indexing activities.  By default it will look for a bucket named `cbfullofit` at `localhost:8091`.  This can configured with the `-bucket` and `-couchbase` command-line options respectively.  A data directory is also required to store index data.  By default a directory named `data` is created/used in the program's working directory.  This can be configured with the `-datadir` command-line option.

    cd $GOPATH/src/github.com/couchbaselabs/cbfullofit
    cbfullofit

## Using

Navigate to `http://localhost:8094`.

Here you can:

* Create new index definitions
* Assign cbfullofit nodes to build the indexes you've defined
* Run queries against the indexes

## Drone.io Build Status

[![Build Status](https://drone.io/github.com/couchbaselabs/cbfullofit/status.png)](https://drone.io/github.com/couchbaselabs/cbfullofit/latest)