# confy

confy is a lightweight configuration module for your app. It's quiet simple and easy to use.

#### Usage

```go
// confy supports 6 kinds of type as following:
confy.Define("http.host", "127.0.0.1")               // Text
confy.Define("http.port", 3389)                      // Integer
confy.Define("http.max.memory.rate", 66.6)           // Float
confy.Define("http.auth.enabled", true)              // Boolean
confy.Define("app.terminal.date", time.Now())        // time.Time
confy.Define("app.refresh.interval", time.Second*12) // time.Duration

// start confy deamon for local file and remote refreshing
confy.Ready()

// print confy value
fmt.Printf("http.host:[%s]", confy.Text("http.host"))
fmt.Printf("http.port:[%d]", confy.Int("http.port"))
fmt.Printf("max.memory.rate:[%f]", confy.Float("http.max.memory.rate"))
fmt.Printf("auth.enabled:[%v]", confy.Bool("http.auth.enabled"))

fmt.Printf("terminal.date:[%v]", confy.Time("app.terminal.date"))
fmt.Printf("refresh.interval:[%v]", confy.Duration("app.refresh.interval"))
```
