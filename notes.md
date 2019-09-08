## CFF info idea

### TODO

- [x] member app, init
- [x] member app manifests for A, C, S
- [ ] member app with BAD_LISTENER flag
- [x] member app manifest for C, no listen

- [x] cff info app with assets
- [x] cff info app with call to local member source, request info, report as JSON
- [x] cff info app with call to member app, request info, report as JSON
- [ ] cff info app with basic web page showing assets, text, metadata
- [ ] cff info app with styling
- [x] cff info app manifest with basic config
- [ ] cff info app with parametrized styling
- [ ] cff info app manifest with basic config



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

```

### initial skeleton

```
cf push -f manifests/cff-info-v1.yml
cf map-route cff-info-v1 istio.geordi.malm.co -n cff-info

cf push -f manifests/member-abby.yml
cf map-route member-abby istio.apps.internal -n members
cf add-network-policy cff-info-v1 --destination-app member-abby --protocol tcp --port 8080

cf push -f manifests/member-chip.yml
cf map-route member-chip istio.apps.internal -n members
cf add-network-policy cff-info-v1 --destination-app member-chip --protocol tcp --port 8080


