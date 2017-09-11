# gocontainer
gocontainer contains several optimized container implementations.

* Package bufpool is a wrapper over the sync.Pool for buffer objects, it simplifies the use.
* Package databox defines the DataBox type, which can be used to store data to reduce the number of references and memory fragmentation.
* Package queue contains some queue implements.
* Package rbuf implements simple ring buffer.
* Package skiplist implements a skip list.

Documentation
-------------

- [API Reference](http://godoc.org/github.com/someonegg/gocontainer)

Installation
------------

Install gocontainer using the "go get" command:

    go get github.com/someonegg/gocontainer

The Go distribution is gocontainer's only dependency.
