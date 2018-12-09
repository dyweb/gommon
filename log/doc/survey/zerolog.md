# Zerolog

https://github.com/rs/zerolog

- mainly focus on JSON logging (not using `encoding/json` as I remember)
- also support CBOR http://cbor.io/ RFC 7049 Concise Binary Object Representation Like JSON but in binary, 123456 is no longer encoded as '123456' and bytes are no longer base64 encoded

````go
zerolog.TimestampFieldName = "t"
zerolog.LevelFieldName = "l"
zerolog.MessageFieldName = "m"
````