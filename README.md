```
  ██████╗ ██████╗ ██████╗ ███████╗██╗    ██╗███████╗██████╗ ███████╗████████╗
 ██╔════╝██╔═══██╗██╔══██╗██╔════╝██║    ██║██╔════╝██╔══██╗██╔════╝╚══██╔══╝
 ██║     ██║   ██║██║  ██║█████╗  ██║ █╗ ██║█████╗  ██████╔╝█████╗     ██║
 ██║     ██║   ██║██║  ██║██╔══╝  ██║███╗██║██╔══╝  ██╔══██╗██╔══╝     ██║
 ╚██████╗╚██████╔╝██████╔╝███████╗╚███╔███╔╝███████╗██║  ██║██║        ██║
  ╚═════╝ ╚═════╝ ╚═════╝ ╚══════╝ ╚══╝╚══╝ ╚══════╝╚═╝  ╚═╝╚═╝        ╚═╝

```

# Platform
Codewerft Platform is the core of all Codewerft API server applications.


## Develop Platform

Create a new workspace:

```
mkdir -p ~/workspace/platform
cd ~/workspace/platform
```

Get a copy of the platform sources:

```
go get github.com/codewerft/platform
```

The Platform sources are in ` ~/workspace/platform/src/github.com/codewerft/platform`.

### Run the example

Build the example server:

```
cd  ~/workspace/platform/src/github.com/codewerft/platform/example
go install
```

Run the example server:

```
$GOPATH/bin/example --config=sample.cfg
```