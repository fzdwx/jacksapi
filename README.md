# Jack's ChatGpt api client

It's a simple client for the https://chat-one.emmmm.dev/ .

![](.github/ask.gif)

## Usage

### lib

```shell
go get github.com/fzdwx/jacksapi@latest
```

code:

```go
var (
    content    = strings.Join(os.Args[1:], " ")
    accessCode = os.Getenv("EMM_API_KEY")
)

jacksapi.NewClient(accessCode).
    ChatStream(
        []jacksapi.ChatMessage{
            {Role: "user", Content: content},
        }).
    DoWithCallback(jacksapi.Output)
```


### cli

```shell
go install github.com/fzdwx/jacksapi/cmd/ask@latest
```

command:

```shell
# ask chatgpt how are you
ask how are you

# start api server
ask server 1333 
curl -X POST -H "Content-Type: application/json" \ 
     -d '{"messages":[{"role":"user","content":"how are you"}]}' \
     http://localhost:1333 

# simulate the API of ChatGPT.
curl 'http://localhost:1333/v1/chat/completions' \
  -H 'Accept: */*' \
  -H 'Connection: keep-alive' \
  -H 'Content-Type: application/json' \
  -H 'User-Agent: Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/113.0.0.0 Safari/537.36' \
  --data-raw '{"model":"gpt-3.5-turbo","temperature":0,"max_tokens":1000,"top_p":1,"frequency_penalty":1,"presence_penalty":1,"stream":true,"messages":[{"role":"system","content":"你是一个翻译引擎，请翻译给出的文本，只需要翻译不需要解释。当且仅当文本只有一个单词时，请给出单词原始形态（如果有）、单词的语种、对应的音标或转写、所有含义（含词性）、双语示例，至少三条例句，请严格按照下面格式给到翻译结果：\n                <单词>\n                [<语种>] · / <Pinyin>\n                [<词性缩写>] <中文含义>]\n                例句：\n                <序号><例句>(例句翻译)\n                词源：\n                <词源>"},{"role":"user","content":"好的，我明白了，请给我这个单词。"},{"role":"user","content":"单词是：hello"}]}' \
  --compressed
```
