# playbypost

[![on-push](https://github.com/betorvs/playbypost/actions/workflows/on-push.yaml/badge.svg)](https://github.com/betorvs/playbypost/actions/workflows/on-push.yaml)

## Directories

- cmd: command line tools to run it
- core: main packages, not designed to be shared with others (not using pkg pattern because of it)
- docs: documentation about this project

## Environment variables

- `PGUSER` : string variable. 
- `PGPASSWORD` : string variable. 
- `PGHOST` : string variable. 
- `PGPORT` : string variable. 
- `PGDATABASE` : string variable. 
- `SLACK_AUTH_TOKEN` : string variable. 
- `SLACK_APP_TOKEN` : string variable. 
- `SLACK_CHANNEL_ID` : string variable. 
- `DISCORD_TOKEN` : string variable. 
- `DISCORD_GUILD_ID` : string variable. 

## development

### Create both files with all environment variables
```bash
.env
.env.task
```

### Change Tiltfile

tiltfile for slack
```
local_resource(
  name='playbypost-server',
  cmd='task tidy build_assets build_local',
  serve_cmd='./playbypost',
  deps=["app/", "go.mod", "go.sum", "core/"]
)
local_resource(
  name='slack-plugin',
  serve_cmd='./slack-plugin',
  resource_deps=['playbypost-server'],
  deps=["./slack-plugin"]
)
```

tiltfile for discord
```
local_resource(
  name='playbypost-server',
  cmd='task tidy build_assets build_local',
  serve_cmd='./playbypost',
  deps=["app/", "go.mod", "go.sum", "core/"]
)
local_resource(
  name='discord-plugin',
  serve_cmd='./discord-plugin',
  resource_deps=['playbypost-server'],
  deps=["./discord-plugin"]
)
```

### Create playbypost database:
```
. .env
./admin-ctl db create
./admin-ctl db up
```

### Run tilt

```bash
tilt dev
```

## References

http://go-database-sql.org/errors.html
https://go.dev/blog/routing-enhancements
https://dev.to/mokiat/proper-http-shutdown-in-go-3fji
https://eli.thegreenplace.net/2023/better-http-server-routing-in-go-122

[bot image source](https://www.freepik.com/free-vector/floating-robot_82654546.htm#fromView=search&page=1&position=13&uuid=44c37a73-28a9-4b70-8d0d-711903439bc1)
