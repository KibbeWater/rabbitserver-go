# RabbitServer GO

Simple WebSocket wrapper to communicate with RabbitHole through regular WebSockets

## Getting Started

Begin by copying the build script

```bash
cat build.sh.example > build.sh
```

Then, edit the `build.sh` file and set the required environment variables.

```sh
# Set ENV vars
export APP_VERSION=<App-Version>
export OS_VERSION=<OS-Version>
```

Finally, run the build script and the binary will be created in the `bin` directory.

```bash
./build.sh
cd bin
# Run the binary, different depending on the OS
./rabbit
```

## Usage

The server will start on port `8080` by default. You can change this by setting the `PORT` environment variable.

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
  "data": "success"
}
```

You can now start sending messages to the server. The server then responds with the same format as the request

```json
{
  "type": "message",
  "data": "<Message>"
}
```

## API

You can find the API documentation in the [API.md](API.md) file.
