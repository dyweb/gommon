# Library Style

Style guide for writing library using gommon

## Error handling

- DO NOT use `log.Fatal`, `panic`, always return error, if an error is added later, many application won't compile