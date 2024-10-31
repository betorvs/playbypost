# Development documentation

## Directories

- app: applications main files.
- core: main packages, not designed to be shared with others (not using pkg pattern because of it).
- ui: react with boostrap application used as Web interface.
- docs: documentation about this project.
- library: RPG related json files.

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


### Requirements

- Go > 1.23 and [golangci-lint](https://github.com/golangci/golangci-lint)
- PostgreSQL (docker or local installation or remote service)
- [Taskfile](https://taskfile.dev/)
- [Tilt](tilt.dev) OR [watchexec](https://github.com/watchexec/watchexec)
- [Zellij](https://zellij.dev/) (optional)

### Create both files with all environment variables
```bash
.env
.env.task
```

### Create playbypost database:

```bash
task migrate_up
```

### test and generate all binaries

```bash
task test
task build_all
```

### Run tilt with task

```bash
task dev
```

#### Change Tiltfile

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

### Run zellij with watchexec

```bash
task zterm
```


## Diagrams

Go to [diagrams](./diagrams/README.md) 


## References

http://go-database-sql.org/errors.html  
https://go.dev/blog/routing-enhancements  
https://dev.to/mokiat/proper-http-shutdown-in-go-3fji  
https://eli.thegreenplace.net/2023/better-http-server-routing-in-go-122  

