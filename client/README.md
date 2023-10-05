# Bark Client

## IMPORTANT: It is not Yet Built

## What does the client do? 
The Bark client (just _client_ henceforth) is the _library_ side of the bark. It is the piece that takes in the logs from any golang program and sends it to the server which is configured against the client. It is supposed to have the utility functions 

What has to be done finally:
```
User -> Request -> (Parse -> Channel -> Make a single or a batch Network Call)
```

But we can start with: 

```
User -> Request -> Make a single Network Call
```

```
fmt.Println("A#1L1ZYG - Something")
```

```go
package main

import (
	"fmt"
	"github.com/techrail/bark/client"
)

func main() {
	log := client.NewClient("http://bark.example.com", "INFO", "auth_servive", "auth_pod_abcd-xyz")
	
	log.Printf("Some log message")
	log.Info("Some info log content")
	log.Error("Some error occurrec")
}
```

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