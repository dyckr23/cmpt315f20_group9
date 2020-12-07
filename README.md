# cmpt315f20_group9
CMPT315 Fall 2020 - Group 9 - Project Codenames
Ben Ha, Ron Dyck

Codenames game implemented with a Go backend, and Typescript frontend. We chose
redis for the datastore. Compiling and adding the RedisJSON module to the redis
datastore gave us a native JSON type. This made saving and retrieving JSON
objects very easy, and since the game state and player moves are JSON objects
already, the result is fast keystore peformance with very little overhead or
complexity. We used the Go-ReJSON client to "add" those capabilities to the
Redigo client.

We chose WebSockets over Long Polling. We found it complicated to set up server
side, but well worth it in terms of responsiveness and performance.

The frontend is styled using Bootstrap CSS, and uses some jQuery and Popper
library functionality. doT is used for frontend templating.

Module list:

github.com/gomodule/redigo          redis client for Go
github.com/nitishm/go-rejson        Go-ReJSON support for Go redis clients
github.com/RedisJSON/RedisJSON      redis module, native JSON 
github.com/gorilla/mux              Gorilla router
github.com/gorilla/websocket        Gorilla websockets

Additional notes, server config:

The project obtained a free trial cloud VM on the Azure platform to use as a 
dev and demo platform. All software was installed from Ubuntu packages, or
compiled from scratch in the case of both redis and RedisJSON. Other software
that was installed: Go, nodejs & npm, typescript, emacs & lsp-mode, Caddy. We
used vsCode with the Remote - SSH plugin as our IDE.