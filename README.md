# gloeg-ws
Websocket binding for gloegg

Feed it a ws/wss url and an optional token and it will forward logs in json format and local flag updates over websockets. It will listen for remote flag updates and update local flags accordingly.

## Encoded behaviors

To get an idea of how the lib behaves, play around with the server and client in the example folder.

An example of encoded behavior is that when no websocket connection is established, logs are dumped in raw json to stdout. Local flag updates and request for flags are silently droped.

Each reconnect attempt is also logged to stdout.