# fpTracking_server

This repository contains a golang web server library which provides an interface to launch the algorithms which can track a fingerprint over the time, and to see the results on graphics. These algorithms were developped by the Spirals Team of the INRIA in Lille.

## Getting Started
### Prerequisites

To use this library, you'll have to install some packages :

```
 $ go get github.com/gorilla/mux
 $ go get github.com/gorilla/sessions
 $ go get github.com/satori/go.uuid
```

You'll also need *github.com/clementmaerten/fpTracking* package. To install it, just follow the instructions written in the readme of the library.

### Installing

For now, as the repository is private, **go get** function to get this package won't work. So here are the steps in order to install it :

 * Go into the github.com/clementmaerten directory 
```
 $ cd $(go env GOPATH)/src/github.com/clementmaerten
```

 * Then clone the repository
```
 $ git clone https://github.com/clementmaerten/fpTracking_server.git
 or
 $ git clone git@github.com:clementmaerten/fpTracking_server.git
```

 * Go in the directory created
```
 $ cd fpTracking_server
```

 * Then **build** the application inside this directory
```
 $ go build
```

### Running

To run the web server, just launch the created app. By default, it will be hosted on http://localhost:8080.

If you want to stop the server, just send signal SIGINT to it (kill -2 or Ctrl+C).