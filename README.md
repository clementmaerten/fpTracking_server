# fpTracking_server

This repository contains a golang web server library which provides an interface to launch the algorithms which can track a fingerprint over the time, and to see the results on graphics. These algorithms were developped by the Spirals Team of the INRIA in Lille.

## Getting Started
### Prerequisites

To use this library, you'll have to install some libraries :

```
 $ github.com/gorilla/mux
 $ github.com/gorilla/sessions
 $ github.com/satori/go.uuid
```

There is a Makefile in this package in order to install these libraries.

You'll also need *github.com/clementmaerten/fpTracking* package. To install it, just follow the instructions written in the readme of the library.

### Installing

To download and use this package, follow the instructions below :

 * Get the repository
```
 $ go get github.com/clementmaerten/fpTracking_server
```

 * Go inside this directory
```
 $ cd $(go env GOPATH)/src/github.com/clementmaerten/fpTracking_server
```

 * Then execute the Makefile (it will download the required libraries and build the package)
```
 $ make
```

### Configure database information

Configure the conf/conf.json file with your database information.

### Running

To run the web server, just launch the created app. By default, it will be hosted on http://localhost:8080.

If you want to stop the server, just send signal SIGINT to it (kill -2 or Ctrl+C).