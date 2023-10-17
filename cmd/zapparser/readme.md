## Zap Parser
Bark's zap parser is a utility that allows you to save logs generated with uber's zap library into bark server.

### How to use it?
You can use the suitable executable zapparser file, and pass any/all the following parameters:

- `file` - the path of the zap's log file, this defaults to `log.txt` if not specified.
- `server` - the url of the bark server, this defaults to `http://localhost:8080/` if not specified.
- `service` - any service name that you want to specify, this defaults to `No service name` if not specified.
- `session` - any session name that you want to specify, this defaults to `No session name` if not specified.

This is how you would execute the command,

```shell
./zapparser.exe -file="log.txt" -server="http://localhost:8080/" -service="example service" -session="example session"
```

