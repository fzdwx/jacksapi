# Jack's ChatGpt api client

It's a simple client for the https://chat-one.emmmm.dev/ .

```go
var (
    content    = strings.Join(os.Args[1:], " ")
    accessCode = os.Getenv("EMM_API_KEY")
)

api.NewClient(accessCode).
    ChatStream(
        []api.ChatMessage{
            {Role: "user", Content: content},
        }).
    DoWithCallback(cb.CopyToStdio)
```
