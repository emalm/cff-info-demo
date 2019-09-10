## CFF info idea

### TODO

- [x] member app, init
- [x] member app manifests for A, C, S
- [x] member app with BAD_LISTENER flag
- [x] member app manifest for C, no listen

- [x] cff info app with assets
- [x] cff info app with call to local member source, request info, report as JSON
- [x] cff info app with call to member app, request info, report as JSON
- [x] cff info app with basic web page showing assets, text, metadata
- [x] cff info app with styling
- [x] cff info app manifest with basic config
- [x] cff info app with parametrized styling
- [x] cff info app manifest with fancy styling
- [x] cff info app retries on errors (including 503s from local envoy)


### cff info app

- `/random` endpoint
  - returns page with a random CFF member from the list
  - asks the members service for a payload
  - gets response with id, full name, title, bio?
    - metadata: host, IP, response time
  - renders page with link to image with id
- `/photos/id.png` endpoint
  - returns PNG image asset at id.png
  - can likely use the assets file router for this

### member app

- / endpoint (or anything)
  - returns JSON response, `{"id": "ID", "name": "Full Name", "title": ""}`

### testing

```

make && ./cff-info-app.darwin

make && MEMBER_URL=https://members.istio.geordi.malm.co ./cff-info-app.darwin

make && MEMBER_ID=foo MEMBER_NAME=Foo PORT=8081 bin/darwin/member-app
make && MEMBER_URL=http://localhost:8081 ./cff-info-app.darwin
```

### initial skeleton

```
cf push -f manifests/cff-info-v1.yml
cf map-route cff-info-v1 istio.geordi.malm.co -n cff-info
cf map-route cff-info-v1 geordi.malm.co -n cff-info-v1

cf push -f manifests/member-abby.yml
cf map-route member-abby istio.apps.internal -n members
cf map-route member-abby istio.geordi.malm.co -n members
cf map-route member-abby istio.geordi.malm.co -n member-abby
cf map-route member-abby geordi.malm.co -n member-abby
cf add-network-policy cff-info-v1 --destination-app member-abby --protocol tcp --port 8080

cf push -f manifests/member-chip.yml
cf map-route member-chip istio.apps.internal -n members
cf map-route member-chip istio.geordi.malm.co -n members
cf map-route member-chip istio.geordi.malm.co -n member-chip
cf add-network-policy cff-info-v1 --destination-app member-chip --protocol tcp --port 8080

cf push -f manifests/member-swarna.yml
cf map-route member-swarna istio.apps.internal -n members
cf map-route member-swarna istio.geordi.malm.co -n members
cf map-route member-swarna istio.geordi.malm.co -n member-swarna
cf add-network-policy cff-info-v1 --destination-app member-swarna --protocol tcp --port 8080

cf remove-network-policy cff-info-v1 --destination-app member-chip --protocol tcp --port 8080
cf add-network-policy cff-info-v1 --destination-app member-chip --protocol tcp --port 8080

cf set-env member-chip BAD_LISTENER true
cf rs member-chip

cf unset-env member-chip BAD_LISTENER
cf rs member-chip

cf push -f manifests/cff-info-v2.yml
cf map-route cff-info-v2 istio.geordi.malm.co -n cff-info

cf add-network-policy cff-info-v2 --destination-app member-abby --protocol tcp --port 8080
cf add-network-policy cff-info-v2 --destination-app member-chip --protocol tcp --port 8080
cf add-network-policy cff-info-v2 --destination-app member-swarna --protocol tcp --port 8080

cf unmap-route cff-info-v2 istio.geordi.malm.co -n cff-info
```
