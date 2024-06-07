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

## PTT

Updates your current PTT status, the payload also has a `image` field which accepts a base64 encoded image which you can use as "vision".

```json
{
  "type": "ptt",
  "data": {
    "active": true, // or false
    "image": "<Base64-Image>" // optional
  }
}
```

When speech is recognized by the server, you will receive the spoken text in the following format:

```json
{
  "type": "ptt",
  "data": "<Spoken-Text>"
}
```

## Audio

This should be used in tandem with the PTT type. When the PTT is active, you can send wav audio data to the server.

```json
{
  "type": "audio",
  "data": "<Base64-Audio>"
}
```

When you receive a text response, a audio reponse usually gets sent right after it containing the spoken text and information about the audio.

```json
{
  "type": "audio",
  "data": {
    "text": "<Stringified-Json>", // ex: {\"language\":\"en\",\"chars\":[\" \",\"H\"],\"char_start_times_ms\":[0,0],\"char_durations_ms\":[0,93]}
    "audio": "<Base64-Audio>"
  }
}
```

## Register

This type is used to register a new device, it requires a base64 encoded QR code from the activation page.

```json
{
  "type": "register",
  "data": "<Base64-QR-Code>"
}
```

The server will then parse and register the device to your account, returning the following response:

```json
{
  "type": "register",
  "data": {
    "imei": "<IMEI>",
    "accountKey": "<Account-Key>",
    "userName": "<User-Name>",
    "userId": "<User-ID>",
    "actualUserId": "<Actual-User-ID>"
  }
}
```

## Long

This type is sent by the server when user has asked for a knowledge prompt. Theres no request using this type

The server will respond with a JSON object containing the following fields:

```json
{
  "type": "long",
  "data": {
    "text": "<Text>",
    "images": ["<Image-URLs>"]
  }
}
```
