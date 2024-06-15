# RabbitServer GO

Simple WebSocket wrapper to communicate with RabbitHole through regular WebSockets

## Getting Started

This guide will cover how to build and run the server. If you need docker instructions, please refer to the [Docker](#docker) section.

Begin by copying the build script

```bash
cat start.sh.example > start.sh
```

Then, edit the `start.sh` file and set the required environment variables.

```sh
# Set ENV vars
export APP_VERSION=<App-Version>
export OS_VERSION=<OS-Version>
```

Before the next step, we should build the project to generate the `bin` directory

```bash
./build.sh
```

The executable will also look for a key.pub file in the same directory as the executable. This file should contain the public RSA key used to sign the Device-Health messages.

Finally, we can run the start script and the server will start

```bash
./start.sh
```

## Docker

You can run the server using Docker. First, build the Docker image:

```bash
docker build -t rabbitserver .
```

Then, run the Docker container:

```bash
docker run -p 8080:8080 -e APP_VERSION=<AppVer> -e OS_VERSION=<OSVer> rabbitserver
```

For a simpler setup, use Docker Compose. Here's an example `docker-compose.yml`:

```yaml
version: '3.8'

services:
  rabbitserver:
    image: rabbitserver
    build: .
    ports:
      - "8080:8080"
    environment:
      - APP_VERSION=<AppVer>
      - OS_VERSION=<OSVer>
```

Replace `<AppVer>` and `<OSVer>` with your application and OS versions, respectively.

Finally, run the Docker container with Docker Compose:

```bash
docker-compose up
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
