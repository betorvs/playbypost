# playbypost

[![on-push](https://github.com/betorvs/playbypost/actions/workflows/on-push.yaml/badge.svg)](https://github.com/betorvs/playbypost/actions/workflows/on-push.yaml)

## Apps

- playbypost: backend that requires PostgreSQL. Divided in two parts, server and worker. Server is a rest JSON api which interacts with plugins, and worker which process all requests hosted in database. Worker cannot run in parallel right now (flags: `-autoplay-worker` and `-stage-worker`). 
- plugin: discord or slack, is a connector with a Chat application.
- CLI: admin-ctl is used to make some administrative tasks and play is used to interact as a player or storyteller with backend. 

## For users

- [Users Guide](./UserGuide.md)


## development instructions

- [Development Guide](./docs/README.md)


## References

[bot image source](https://www.freepik.com/free-vector/floating-robot_82654546.htm#fromView=search&page=1&position=13&uuid=44c37a73-28a9-4b70-8d0d-711903439bc1)
