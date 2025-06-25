# BARISTA Protocol

Inspired from MQTT.

Tradoffs:

```
+ extremely simple
+ extremely fast
- low reliability
- minimal features
```

UDP: strict payload size limit of 1472 to avoid fragmentation.

### Client subscribes to a topic

```
Subscriber --SUBSCRIBE----> Broker
Broker     --SUBACK-------> Client
```

### Client publishes a message

```
Pubilsher  --PUBLISH------> Broker
Broker     --PUBACK-------> Pubilsher
```

### Subscriber receives a message

```
Broker     --PUBLISH------> Subscribers
Subscriber --PUBACK-------> Broker
```

### Client unsubscribes from a topic

```
Subscriber --UNSUBSCRIBE--> Broker
Broker     --UNSUBACK-----> Subscriber
```
