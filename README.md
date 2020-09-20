# matrix-bot

Based on [mautrix-go](https://github.com/tulir/mautrix-go/).

A bot that listens for keywords (e.g. `pma!27`) and sends the full URL to the GitLab issue / merge request back (e.g. `https://gitlab.com/postmarketOS/pmaports/merge_requests/27`).

## Building

```sh
go build
```

## Usage

```sh
./pmos-bot -homeserver https://my.homeserver -username botusername -password botpassword -stateStoragePath /tmp/pmos-bot.json
```

For production usage you can also add the parameter `-deviceid` to set the Matrix device id that will be used (default is an auto-generated one). Additionally you should set `-stateStoragePath` to a persistent directory (e.g. `/var/lib/pmos-matrix-bot/state.json`)
