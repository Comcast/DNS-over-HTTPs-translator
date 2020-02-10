# DNS-over-HTTPs-translator

HTTP proxy server that accepts DNS queries and responds with replies from the configured DNS resolver.

## License

Apache License 2.0

## Working

The DoH-Translator listens on port 80 and translates GET or POST requests containing messages of media type - "application/dns-message". The translator then queries the configured DNS Resolver (default is the public anycast Comcast DNS Resolver) and correlates the DNS response from the resolver with the  HTTP exchange.

## Design

The DoH-Translator uses the following interfaces: Config, Controller and Proxy. This keeps the design modular and allows for more functionality to be added by adding features to the interface implementation. The current Proxy interface is implemented using the net/http go package. Proxy can also be implemented using a different http package if desired.

## Configuration

The translator is configured using the config-doh-translator.yaml file at the path /etc/translator on a linux machine.

Since the translator is designed to be hosted behind an nginx HTTPS service and within a provider's network, the resolver's IP could be appropriately chosen to be the geographically closest internal resolver. For now the public anycast Comcast DNS IP is set as default.

Example:
```yaml
# IP address and port of the resolver.
# Note: Comcast resolver set as default if no resolver
#       is provided here.
resolver: "75.75.75.75:53"

# Other configuration options such as caching options, rate-limiting
# options go below this. (TBD)
```

## Requirements

Please verify the success of the below checks before executing the translator:
- `echo $GOPATH` must be set as a system-wide environment variable.

## Execution
```shell
mkdir -p /etc/translator/
mkdir -p $GOPATH/src/github.com/Comcast/DNS-over-HTTPs-translator
cd $GOPATH/src/github.com/Comcast/DNS-over-HTTPs-translator
git clone git@github.com:Comcast/DNS-over-HTTPs-translator.git
cd translator
cp $GOPATH/src/github.com/Comcast/DNS-over-HTTPs-translator/config-doh-translator.yaml /etc/translator
make build
sudo .build/doh-translator-linux-amd64 start
```
Note: To exit to terminal or stop translator, hit `Ctrl + C` twice.

## Deploy
```shell
# Run script to deploy the translator as a systemd service (sudo commands used)
sh deploy.sh

# To start service...
sudo service doh-translator start

# To check status of service...
sudo service doh-translator status

# To stop service...
sudo service doh-translator stop
```
