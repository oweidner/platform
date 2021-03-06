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

# Recompile and Run - Oct. 2016

```
go install -v  github.com/oweidner/tutorbox.backend && ./bin/tutorbox.backend --config=etc/sample_local.cfg
```


## Develop Platform

Create a new workspace:

```
mkdir -p ~/workspace/platform
cd ~/workspace/platform
```

Get a copy of the platform sources:

```
export $GOPATH=`pwd`
go get github.com/oweidner/platform
```

The Platform sources are in ` ~/workspace/platform/src/github.com/oweidner/platform`.

### Recompile

???

```
go install -a -v github.com/oweidner/platform && go build && ./example -config=sample.cfg
```

### Run the tests

```
 go test github.com/oweidner/platform -test.v --config=tests/test.cfg
```

### Run the Example

Build the example server:

```
cd  ~/workspace/platform/src/github.com/oweidner/platform/example
go install
```

### Generate JWT Keypairs

Generate an RSA keypair for JWT encryption.

```
openssl genrsa -out jwt_sample.rsa 4096
openssl rsa -in jwt_sample.rsa -pubout > jwt_sample.rsa.pub
```

### Generate TLS Certificate

Generate a self-signed X.509 certificate tp run the _Platform_ server in TLS mode.

```
go run /usr/local/go/src/crypto/tls/generate_cert.go --host="localhost"
```

Run the example server:

```
$GOPATH/bin/example --config=sample.cfg
```

Make OS X accept the certificate: http://stackoverflow.com/questions/7580508/getting-chrome-to-accept-self-signed-localhost-certificate


## Notes

* Testing: http://dennissuratna.com/testing-in-go/
