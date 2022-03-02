# Stock Ticker

This stock ticker app queries the alpha vantage API for a given stock symbol and returns a configurable number of days data along with the average closing price. The data is cached in memory for 1 hour.

## Building

The container image builds from a UBI8 build image, which downloads Go 1.17.7 and builds the stock ticker. This is then copied into a fresh UBI8 image.

To build, use buildah or podman (or docker):

`buildah bud -t <your tag here> .` or `podman build -t <your tag here> .`

To build the binary on it's own:

```
CGO_ENABLED=0 #optional if you want a statically compiled binary
go build -o stockticker main.go
```

You can also build a from scratch version (be sure to mount in your ssl certs directory containing your trusted ca bundle into /etc/ssl/certs when running):

`buildah bud -f Containerfile.scratch-t <your tag here> .` or `podman build -f Containerfile.scratch -t <your tag here> .`

## Running with podman

To run this with podman (or docker), you'll need an environment file that looks like this:
```
NDAYS=7
SYMBOL=MSFT
APIKEY=yourkeyhere
```
Alternatively, you can put these on the podman command line.

With an environment file, run it like this:

`podman run -d --env-file=<envfile> -p 8080:8080 <your tag here>`

## Running on Kubernetes

An example Kubernetes deployment manifest can be found in deployment/. This will create a stockticker namespace, a config map, a secret, the deployment, a service and an ingress. You will need to modify the secret to contain your alpha vantage API key, and the symbol to find, along with number of days of data to return, can be set in the config map.

Apply to your cluster with `kubectl apply -f deployment/stockticker.yaml`.
