# pushproxy-light

A primitive proxy implementation to transverse NAT for [Prometheus](https://prometheus.io/) like [PushProx](https://github.com/prometheus-community/PushProx) and [Pushgateway](https://github.com/prometheus/pushgateway).

# Difference between PushProx / Pushgateway

## PushProx
 - Push metrics from client exporting program directly.


## Pushgateway
- Metrics data has lifetime.

# Motivation

I tried to use PushProx to monitor a Raspberry Pi Model B that was doing environmental observations of the room and executing dump1090/rtl-sdr. However I could not execute PushProx client due to limitation of computing resouces(dump1090 is so heavy). ;(.

# Usage
 `pushprox-light --preSharedKey=hogehoge ...`

- `--preSharedKey` (mandatory)
  - Pre-Shared Key to sign metrics data.
- `--metricsLifetime` default=5
  - Metrics Lifetime in minutes.
- `--httpListener`  default=`:8080`
  - Listening Address/Port
- `--proxyMetricsPath` default=`/metrics`
  - Path URI of Metrics where Prometheus scapes.

# Push Protocol
PushProx-light accepts JSON text as metrics data(body).

Send signature with pre-shared key for tamper detection and authentication.

## HTTP Methods
- GET
  - `URL?signature=SIGNATURE&body=BODY`
- POST as Form data (`application/x-www-form-urlencoded`).
  - `signature=SIGNATURE&body=BODY`
- POST JSON (`application/json`)
  - `X-Signature` Header shows signature.
  - body is POST body.

## Push URL Path

`/push/{{hostname}}`
{{hostname}} is the same as the hostname in the Prometheus configuration.

## body
Metrics data expressed as JSON.

```json
{
	"name of metrics": {
		"type": "gauge/counter",
		"help": "description of metrics",
		"metrics": [{
			"value": "11.1",
			"labels": {
				"key1": "value1",
				"key2": "value2"
			}
		}]
	}
}
```

## signature

`alg := '1'(md5),'5'(sha256)','6'(sha512)`

`signature :=  fmt.Sprintf("%c%x",alg, md5/sha256/sha512.Sum(preSharedKey+body))`

# Dependeicies
- [github.com/gorilla/mux](https://github.com/gorilla/mux)
- [github.com/patrickmn/go-cache](https://github.com/patrickmn/go-cache)

# Author
walkure 