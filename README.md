# gocontainer
gocontainer contains several optimized container implementations.

* Package databox defines the DataBox type, which can be used to store data to reduce the number of references and memory fragmentation.
* Package skiplist implements a skip list.
* Package uskiplist implements a skiplist using unusual operations to minimize memory and references.

Documentation
-------------

- [API Reference](http://godoc.org/github.com/someonegg/gocontainer)

Installation
------------

Install gocontainer using the "go get" command:

    go get github.com/someonegg/gocontainer
