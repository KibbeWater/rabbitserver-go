# RabbitServer API

Every request to the server must be in JSON format. The server will respond with a JSON object as well.

Every request must contain a `type` field, this identifies the type of request.

## Logon

When a client connects you will not be able to perform anything until you send a `logon` message. This message should contain the following fields:

```json
{
  "type": "logon",
  "data": {
    "imei": "<IMEI>",
    "accountKey": "<Account-Key>" // The key can be seen when registering the device
  }
}
```

After sending this message, you will receive a response with the following format:

```json
{
  "type": "logon",
  "data": "success" // or "failure"
}
```

## Message

You can now start sending messages to the server. The server then responds with the same format as the request

```json
{
  "type": "message",
  "data": "<Message>"
}
```
