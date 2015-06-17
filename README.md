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

### Run the Example

Build the example server:

```
cd  ~/workspace/platform/src/github.com/codewerft/platform/example
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