# CMPE-273-Assignment-3


### POST Request

```
curl -H "Content-Type: application/json" -X POST -d '{"starting_from_location_id":"565199c90ed56134e00c0f92","location_ids":["565199d50ed56134e00c0f93","565199e30ed56134e00c0f94","565199ef0ed56134e00c0f95","565199fc0ed56134e00c0f96"]}' localhost:8080/trips
```

### GET Request

```
curl localhost:8080/trips/{tripid}
```

### PUT Request

```
curl -X PUT localhost:8080/trips/{tripid}
