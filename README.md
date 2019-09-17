# cff-info-demo

## Prerequisites

- Install the `cf` CLI.
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

