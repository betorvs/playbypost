# User Guide

## Download 

Go to [releases](https://github.com/betorvs/playbypost/releases) and download for your Operation System

## Create a Local Database

Could be a remote PostgreSQL database, like [elephantSQL](https://www.elephantsql.com/) or a [local installation](https://www.postgresql.org/download/) or docker. 

You need to have your database credentials as environment variables:

```bash
export PGUSER="postgres"
export PGPASSWORD="mypassword"
export PGHOST="localhost"
export PGPORT=5432
export PGDATABASE="playbypost"
```

Or Windows ([source](https://www.ibm.com/docs/en/informix-servers/12.10?topic=windows-using-command-prompt-change-environment-variables)):
```
set PGUSER="postgres"
set PGPASSWORD="mypassword"
set PGHOST="localhost"
set PGPORT=5432
set PGDATABASE="playbypost"
```

### Changing library files

You can change or add more content to [library files](./library/) with caution. 

If you need to change path, these are environment variable used:

- "D10HM_DEFINITION": default value "./library/definitions-d10HM.json".
- "PFD20_DEFINITION": default value "./library/definitions-pfd20.json".
- "PFD20_ANCESTRIES_DEFINITION": default value "./library/definitions-pfd20-ancestries.json".
- "PFD20_BACKGROUNDS_DEFINITION": default value "./library/definitions-pfd20-backgrounds.json".
- "PFD20_CLASSES_DEFINITION": default value "./library/definitions-pfd20-classes.json".


## Create Bot Credentials

It can only connect to a one Chat. Then choose one and create credentials for it.

### Discord

Follow a tutorial about Discord Bot creation on Discord Developer Portal and give to your Bot the following permission in the end:
- Bor > Privileged Gateway Intents > Message Content Intent.


### Slack

Follow a tutorial about Slack Bot and give him the following permissions:
- Bot Token Scopes:
  - app_mentions:read
  - channels:history
  - chat:write
  - chat:write.customize
  - commands
  - im:history
  - im:write
  - incoming-webhook
  - reactions:write
  - users:read


For Slack it's important to set up Socket Mode with: Interactivity & Shortcuts, Slash Commands and Event Subscriptions. 

Slash command recommendation: `/play-by-post`.

### Add Bot credentials

```bash
export SLACK_AUTH_TOKEN=""
export SLACK_APP_TOKEN=""
export DISCORD_TOKEN=""
export DISCORD_GUILD_ID=""
```

Or Windows:
```
set SLACK_AUTH_TOKEN=""
set SLACK_APP_TOKEN=""
set DISCORD_TOKEN=""
set DISCORD_GUILD_ID=""
```

## Run your applications

After exporting your environment variables, create your database:
```bash
./admin-ctl db ping
./admin-ctl db create
./admin-ctl db up
```

First playbypost
```bash
./playbypost -autoplay-worker -stage-worker
```

In another terminal, your bot:
```bash
./discord-plugin
```

Or

```bash
./slack-plugin
```

### Create your first writer user

You can find admin token in the logs, like `msg="adding admin user one year token" admin=admin token=HERE`. It changes on every restart. 

```bash
./admin-ctl writer create -u YOUR-USERNAME --password YOUR-PASSWORD --token ADMIN-TOKEN
```

### Associate a User with a Writer

Once you have created a writer, you can associate a Discord or Slack user with that writer using the `/iamwriter` command. This allows the bot to identify the user as a specific writer within the game.

**Discord:**

```
/iamwriter writer_username
```

Replace `writer_username` with the username of the writer you created using the `admin-ctl` tool.

**Slack:**

```
/play-by-post iamwriter writer_username
```

Replace `writer_username` with the username of the writer you created using the `admin-ctl` tool.

IMPORTANT: All tokens expires in each restart of this application. Admin token and Writers tokens (For Web interface).

### Access your local Web Interface

[PlaybyPost-Web](http://localhost:3000/)

At home page, you can find more information about all related concepts and how to use it.


## Feedback

Please, share with me what do you think about it and how can we improve it. Thanks and Enjoy!