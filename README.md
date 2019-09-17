# cff-info-demo

## Prerequisites

- Install the `cf` CLI.
- Install the `jq` command-line tool.
- Log into a CFAR installation with the Istio service mesh enabled.
- Change to the root directory of this repository.
- Run `export CF_DOMAIN=<REPLACE_WITH_CF_DOMAIN>` with `<REPLACE_WITH_CF_DOMAIN>` replaced with the base apps domain for your CFAR installation.


## Service-Mesh Load Balancing

Push the front-end CFF info apps:

```
cf push -f manifests/cff-info-v1.yml
cf push -f manifests/cff-info-v2.yml
```

Map an Istio ingress route to the `cff-info-v1` app

```
cf map-route cff-info-v1 "istio.$CF_DOMAIN" -n cff-info
```

Push the back-end CFF member info apps:

```
for n in abby chip swarna;
do
  cf push -f "manifests/member-$n.yml"
done
```

Add network policies from the front-end apps to the back-end member apps:

```
./scripts/set-network-policies
```

Make a request to the front-end app:

```
curl "cff-info.istio.$CF_DOMAIN/random"
```

Make additional requests and observe that the member info changes while the URL and resolved IP in the member request metadata remain the same.


### Fault handling

Change the `member-chip` app to listen on port 9999 instead of on 8080:

```
cf set-env member-chip BAD_LISTENER true
cf restart member-chip
```

Make several requests to the random CFF info page and observe that sometimes the front-end app makes several requests to the back-end internal route before getting a successful response.

Change the `member-chip` app to listening on 8080:

```
cf unset-env member-chip BAD_LISTENER
cf restart member-chip
```


## Weighted Routing

Set the weighted-routing rules to send 80% of traffic to the v1 CFF info app and 20% to the v2 CFF info app:

```
./scripts/set-route-weights 80
```

Make many requests to the random info page and observe differences in styling.


Set the weighted-routing rules to send 10% of traffic to the v1 CFF info app and 90% to the v2 CFF info app:

```
./scripts/set-route-weights 90
```

Make many more requests to the random info page and observe the new style eventually occurs much more frequently.

Return to routing 100% of traffic to the v1 CFF info app:

```
./scripts/clear-route-weights
```
