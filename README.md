# Tracker Service

## Install

Spin up a database:

```sh
$ docker-compose up -d db
```

Start Application

```sh
$ make start [--build] tracker
```

Seed Application

In a new terminal type:

```sh
$ make seed 
```

We can now test the service by running:

```sh
# must have installed grpcurl https://github.com/fullstorydev/grpcurl

$  grpcurl -d @ -plaintext localhost:50051 nawalt.tracker.v1.TrackingService.RecordPosition <<EOM
{
    "user_id": "018bdb59-6ad8-7980-a347-0fa6c27ae9ea",
    "location": {
        "latitude": 12.10139,
        "longitude": -86.25856
    },
    "timestamp": "2023-11-17T03:39:48+00:00",
    "client_id": "018bdb5e-743d-74e6-a6cd-abfef24bc260",
    "metadata": {
        "device_id": "352032066517282"
    }
}
EOM
```bash

$  grpcurl -d @ -plaintext localhost:50051 nawalt.tracker.v1.TrackingService.RecordPosition <<EOM
{
    "user_id": "018bdb59-6ad8-7980-a347-0fa6c27ae9ea",
    "location": {
        "latitude": 12.10193458760576,
        "longitude": -86.26166582107545
    },
    "timestamp": "2023-11-17T03:39:48+00:00",
    "client_id": "018bdb5e-743d-74e6-a6cd-abfef24bc260",
    "metadata": {
        "device_id": "352032066517282"
    }
}
EOM
```
