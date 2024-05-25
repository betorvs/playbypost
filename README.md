# playbypost

## Directories

- cmd: command line tools to run it
- core: main packages, not designed to be shared with others (not using pkg pattern because of it)
- docs: documentation about this project

## migrations

```bash
cd core/sys/db/migrate
migrate create -ext sql -dir migrations/ -seq create_base_tables
migrate create -ext sql -dir migrations/ -seq create_encounter_and_participants_table
```

## References

http://go-database-sql.org/errors.html
https://go.dev/blog/routing-enhancements
https://dev.to/mokiat/proper-http-shutdown-in-go-3fji
https://eli.thegreenplace.net/2023/better-http-server-routing-in-go-122

[bot image source](https://www.freepik.com/free-vector/floating-robot_82654546.htm#fromView=search&page=1&position=13&uuid=44c37a73-28a9-4b70-8d0d-711903439bc1)
