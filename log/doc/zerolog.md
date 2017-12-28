# Zerolog

https://github.com/rs/zerolog

- mainly focus on JSON logging (not using `encoding/json` as I remember)

````go
zerolog.TimestampFieldName = "t"
zerolog.LevelFieldName = "l"
zerolog.MessageFieldName = "m"
````