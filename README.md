# OpenPayment

I have created this as an example of using Event sourced application state. 

Key principles: -
- Application state can be completely constructed from the stream
- Application state changes triggered only by a received event
- Any change received from an external system produces an event

## Account

The account is simply to have an entity to demonstrate simple actions a service
might need to do. The related items (Customer, Card, etc) are referenced by a URL
which shows how you might relate data from other services into this service.

## Event stream

For the example and tests there is an in memory stream. For an actual service,
a service designed for event streams like [kafka](https://kafka.apache.org/) would be 
best. You could use a fast, fault-tolerant store like [cassandra](https://cassandra.apache.org/) 
if that was already in the system. The key aspects of the event stream system are: -

- consistent ordering
- entire stream persisted

Message queue systems initially look attractive for this but typically have a retention 
period which only allows a limited time coverage.

To reload the state from the stream simply requires clearing the state and processing
the events sequentially. The In memory stream just tracks the index of the next event
to send. So this is set back to 0.

Snapshots can be taken to store the entire application state at a point in time. This 
means that only events from that point forward will be required to rebuild the current 
state.

## Typical log showing event handling

```text
go run ./cmd/simpleAccountWebservice/main.go
{"severity":"debug","eventId":"cbc769ce-719b-4890-bc4f-1137049254ad","eventType":"account add","streamLen":1,"time":"2020-11-01T20:11:52Z","message":"event published"}
{"severity":"debug","time":"2020-11-01T20:11:52Z","message":"processed account add request"}
{"severity":"debug","eventId":"cbc769ce-719b-4890-bc4f-1137049254ad","eventType":"account add","time":"2020-11-01T20:11:52Z","message":"event received"}
{"severity":"debug","eventId":"cbc769ce-719b-4890-bc4f-1137049254ad","time":"2020-11-01T20:11:52Z","message":"received event"}
{"severity":"debug","eventId":"cbc769ce-719b-4890-bc4f-1137049254ad","time":"2020-11-01T20:11:52Z","message":"event processed by all subscribers"}
{"severity":"debug","eventId":"cbc769ce-719b-4890-bc4f-1137049254ad","eventType":"account add","time":"2020-11-01T20:11:52Z","message":"received account event"}
{"severity":"debug","eventId":"cbc769ce-719b-4890-bc4f-1137049254ad","eventType":"account add","time":"2020-11-01T20:11:52Z","message":"processed event"}
{"severity":"debug","eventId":"cbc769ce-719b-4890-bc4f-1137049254ad","time":"2020-11-01T20:11:52Z","message":"event processed by simple account"}
{"severity":"debug","accountNumber":"123456789","time":"2020-11-01T20:11:56Z","message":"processed account get request"}
{"severity":"debug","eventId":"ca6b12e9-e3ae-4d59-a7f0-f637df226b13","eventType":"account update","streamLen":2,"time":"2020-11-01T20:11:59Z","message":"event published"}
{"severity":"debug","time":"2020-11-01T20:11:59Z","message":"processed account update request"}
{"severity":"debug","eventId":"ca6b12e9-e3ae-4d59-a7f0-f637df226b13","eventType":"account update","time":"2020-11-01T20:12:00Z","message":"event received"}
{"severity":"debug","eventId":"ca6b12e9-e3ae-4d59-a7f0-f637df226b13","time":"2020-11-01T20:12:00Z","message":"received event"}
{"severity":"debug","eventId":"ca6b12e9-e3ae-4d59-a7f0-f637df226b13","time":"2020-11-01T20:12:00Z","message":"event processed by all subscribers"}
{"severity":"debug","eventId":"ca6b12e9-e3ae-4d59-a7f0-f637df226b13","eventType":"account update","time":"2020-11-01T20:12:00Z","message":"received account event"}
{"severity":"debug","eventId":"ca6b12e9-e3ae-4d59-a7f0-f637df226b13","eventType":"account update","time":"2020-11-01T20:12:00Z","message":"processed event"}
{"severity":"debug","accountNumber":"123456789","time":"2020-11-01T20:12:05Z","message":"processed account get request"}
```
