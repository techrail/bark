# Bark Client

## IMPORTANT: It is not Yet Built

## What does the client do? 
The Bark client (just _client_ henceforth) is the _library_ side of the bark. It is the piece that takes in the logs from any golang program and sends it to the server which is configured against the client. It is supposed to have the utility functions to help users log to bark directly from go code without having to worry about network calls and such.

## Levels of Logs
The client defines 7 levels of logs:

1. **Panic (P)** - The message you emit right before the program crashes
2. **Alert (A)** - The message needs to be sent as an alert to someone who must resolve it ASAP
3. **Error (E)** - The message indicating that there was an error and should be checked whenever possible
4. **Warning (W)** - The message indicating that something wrong could have happened but was handled. Can be overlooked in some cases.
5. **Notice (N)** - Something worth noticing, though it is fine to be ignored.
6. **Info (I)** - Just a log of some data - does not indicate any error
7. **Debug (D)** - used for debugging. It can represent any level of information but is only supposed to indicate a message emitted during a debug session

Any single character in the place of error level in a parsable single log message would indicate the level of **INFO**.

## Simple usecase
The client can be initialized and used as follows (we explain the options below code sample): 

```go
barkClient := client.NewClient("<bark_server_url>", "<default_log_level>", "<default_service_name>", "<session_name>", 
    "<enable_slog", "<enable_bulk_dispatch>")

barkClient.Panic("E#1LPW24 - Panic message")
barkClient.Alert("E#1LPW25 - Alert message", false)
barkClient.Error("E#1LPW26 - Error message")
barkClient.Warn("E#1LPW27 - Warn message")
barkClient.Notice("E#1LPW28 - Notice message")
barkClient.Info("E#1LPW29 - Info message")
barkClient.Debug("E#1LPW30 - Debug message")
barkClient.Println("Println message")
barkClient.Default("Println message")
```

The options that are used for initializing the client are as follows:

- **bark_server_url**: This is the URL of a running bark server. It must end in `/`. For example `http://bark.example.com/` or `http://127.0.0.1:8080/`
- **default_log_level**: When you use `Println` or `Default`, the log message is parsed (rules for prasing are described [here](../_nocode/docs/log-string-parsing-in-bark.md)) and if it does not contain any indication for what the log level is, then the value supplied in this field is used as the log level for sent log message. When using dedicated methods for error levels (e.g. `Panic`, `Error` etc.), the parsed error level is overwritten.
- **default_service_name**: This is the name of the service which is sending the log - so it has to be the name of the program or service which is calling it. In case a blank string is sent, the value against `constants.DefaultLogServiceName` (currently set to `def_svc`) is used.
- **session_name**: This is the name of the calling app's session. This value is supposed to indicate which instance among possibly multiple instances of a service sent a log message. For example, in case of the service being deployed within Kubernetes, it might indicate the service's pod's name. If the value is sent as a blank string, client will try to use the machine's hostname. If it fails to fetch the hostname, a random string will be used instead. 
- **enable_slog**: This enables [slog](https://go.dev/blog/slog) for the client. When this option is enabled, all logs in addition to being sent to the bark server is also printed on STDOUT of the service.
- **enable_bulk_dispatch**: Setting this to true would enable the client to push all the requests being received in a channel and start using it. It improves the overall performance of the client sending log entries to the server.

### Simplest usecase (without any server)
The simplest usecase of any logging library is to print to STDOUT. While the primary usecase of bark is to be able to dispatch log messages to a remote server, when we start off with a new project, we often just want to start logging things to STDOUT. Maybe even later, that is how we want to use logs. For such usecases, the client library offers `NewSloggerClient` which uses the built in [slog](https://go.dev/blog/slog) package in go (version 1.21+) to log your messages to STDOUT with levels. Example: 

```go
log := client.NewSloggerClient("INFO")

log.Panic("Panic message")
log.Alert("Alert message", true)
log.Error("Error message")
log.Warn("Warn message")
log.Notice("Notice message")
log.Info("Info message")
log.Debug("Debug message")
log.Println("Println message")
```
The above piece of code will end up printing something like the following (the dates in the beginning of each line will vary): 

```
2023/10/15 21:57:41 PANIC Panic message
2023/10/15 21:57:41 ALERT Alert message
2023/10/15 21:57:41 ERROR Error message
2023/10/15 21:57:41 WARN Warn message
2023/10/15 21:57:41 NOTICE Notice message
2023/10/15 21:57:41 INFO Info message
2023/10/15 21:57:41 DEBUG Debug message
2023/10/15 21:57:41 INFO Println message
```

## Printing logs to a file
Bark client, as shown above, is capable of sending logs to a server as well as printing them to the standard output as well. It can also do both of those things simultaneously. The architecture in very simple representation looks like this: 

![barkslogger.svg](../_nocode/images/barkslogger.svg)

You can use any of the log methods available in bark client to do this. If you want to print the logs to a different output, such as a file, you can use the `SetCustomOut` method. This method takes an `io.Writer` parameter and sets it as the output writer for the bark client. For example, if you want to print the logs to a file named random.txt, you can do this:

```go
log := client.NewClient("http://127.0.0.1:8080/", "INFO", "BarkClientFileTest", "TestClientSession", true, false)

file, err := os.Create("random.txt")
if err != nil {
	fmt.Println("Error when creating new file: ", err)
	return
}

log.SetCustomOut(file)

log.Info("Some Message that'll be sent to random.txt file")
log.WaitAndEnd()
```
The above code will write the output to `random.txt` file. You can expect the file to contain something like this:

```text
2023/10/18 19:27:51 INFO Some Message that'll be sent to random.txt file
```

### Slog and writing to a file 

Bark client uses [slog](https://go.dev/blog/slog) internally to handle the printing of the logs. Slog is a simple and structured logging library that comes with Go (version 1.21+).

You can customize how slog prints the logs by specifying a [handler](https://pkg.go.dev/log/slog#Handler). A handler is a function that takes a log record and writes it to an output. Slog provides some built-in handlers, such as [JSONHandler](https://pkg.go.dev/log/slog#JSONHandler) and [TextHandler](https://pkg.go.dev/log/slog#TextHandler), or you can write your own.

**_Note:_** Changing the handler will only affect how the logs are printed, not how they are sent to bark server.

To specify a handler for the bark client, you can use the `SetSlogHandler` method. This method takes a `handler` function as a parameter and sets it as the handler for the slog logger. For example, if you want to use the `JSONHandler` and print the logs as JSON objects to a file named `random.txt`, you can do this:

```go
package main

import (
	"fmt"
	"github.com/techrail/bark/client"
	"log/slog"
	"os"
)

func main() {
	log := client.NewClient("http://127.0.0.1:8080/", "INFO", "BarkClientFileTest", "TestClientSession", true, false)

	file, err := os.Create("random.txt")
	if err != nil {
		fmt.Println("E#1M5WRN - Error when creating new file: ", err)
		return
	}

	// We are using JSONHandler here so the line that will be logged will actually be a JSON string
	log.SetSlogHandler(slog.NewJSONHandler(file, client.SlogHandlerOptions()))
	// If you want to log to STDOUT, you can use `os.Stdout` in place of the `file` as writer 
	// Of course in case that you would have  remove the unused code from above.

	log.Info("Some Message that'll be sent to random.txt file")
	log.WaitAndEnd()
}
```
You can expect the `random.txt` file to contain something like this with different time being logged (we are using a JSON handler for slog):

```text
{"time":"2023-10-18T19:30:38.773512+05:30","level":"INFO","msg":"Some Message that'll be sent to random.txt file"}
```

You may have noticed that we are passing some options to the `JSONHandler` using the `client.SlogHandlerOptions()` method. This is because slog has predefined labels for only four log levels: `info, warning, debug, and error`. However, bark client supports three additional log levels: `alert, panic, and notice`. The options returned by `client.SlogHandlerOptions()` define labels for these additional log levels.

If you add a nil options, the log labels will appear as described in the [slog documentation here](https://pkg.go.dev/log/slog#Level.String)

[Slog treats log levels as integers](https://pkg.go.dev/log/slog#Level). The predefined log levels have the following values:

> LevelDebug Level = -4 \
> LevelInfo  Level = 0 \
> LevelWarn  Level = 4 \
> LevelError Level = 8 \

The custom log levels defined by bark client have the following values:

```
Notice = 3
Alert = 9
Panic = 10
```

If you are writing a custom handler for slog, please make sure to handle these log levels appropriately.