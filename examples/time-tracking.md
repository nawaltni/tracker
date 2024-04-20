# Examples of Time Tracking

## Check In

```sh
grpcurl -d @ -plaintext localhost:50051 \
nawalt.tracker.v1.TrackingService.CheckIn <<EOM
{
    "user_id": "018bdb59-6ad8-7980-a347-0fa6c27ae9ea",
    "location": {
        "latitude": 12.10193458760576,
        "longitude": -86.26166582107545
    }
}
EOM
```

## Check Out

```sh
grpcurl -d @ -plaintext localhost:50051 \
nawalt.tracker.v1.TrackingService.CheckOut <<EOM
{
    "user_id": "018bdb59-6ad8-7980-a347-0fa6c27ae9ea",
    "session_id": "018d5d33-dd00-79d3-bda0-2439255ee235",
    "location": {
        "latitude": 12.10193458760576,
        "longitude": -86.26166582107547
    }
}
EOM
```

## Start Break

```sh
grpcurl -d @ -plaintext localhost:50051 \
nawalt.tracker.v1.TrackingService.StartBreak <<EOM
{
    "user_id": "018bdb59-6ad8-7980-a347-0fa6c27ae9ea",
    "session_id": "018d7bbc-4388-70dd-b997-785a4a0028f3",
    "location": {
        "latitude": 12.10193458760576,
        "longitude": -86.26166582107547
    }
}
EOM
```

## StartBreak

```sh
grpcurl -d @ -plaintext localhost:50051 \
nawalt.tracker.v1.TrackingService.EndBreak <<EOM
{
    "user_id": "018bdb59-6ad8-7980-a347-0fa6c27ae9ea",
    "session_id": "018d7bbc-4388-70dd-b997-785a4a0028f3",
    "location": {
        "latitude": 12.10193458760576,
        "longitude": -86.26166582107547
    }
}
EOM
```