# chaos-redis [under development]
[CHAOS] Project to simulate connection to redis


We want reproduce when Redis is full in connection and The application can't connect to Redis.


Run script : 
```shell
go run main.go ping --help

Usage:
  chaos-redis ping [flags]

Flags:
  -d, --duration int         chaos duration (default 300)
  -h, --help                 help for ping
      --host string          redis host (default "localhost")
  -c, --max-connection int   Max Active connection, by default unlimited
  -w, --max-parallel int     Max Go routine in parallel (default 10)
      --password string      redis password
  -P, --port int             redis port (default 6379)
```


Run locally redis to test :

```shell
docker run --name redis -p 6379:6379 -d redis
```
