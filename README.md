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
    DoWithCallback(cb.Output)
```


## Usage

```shell
# ask chatgpt how are you
ask how are you

# start api server
ask server 1333 
curl -X POST -H "Content-Type: application/json" \ 
     -d '{"messages":[{"role":"user","content":"how are you"}]}' \
     http://localhost:1333 
```
