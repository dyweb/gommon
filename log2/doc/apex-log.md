# Apex-Log

https://github.com/apex/log

- use handler instead of formatter + writer
  - like `net/http` you can use a function as handler
  - [ ] TODO: why it is `func (f HandlerFunc)` instead of `func (f *HandlerFunc)`, can we have a `type HandlerFuncPtr`?

````go
func AuthWrapper(h http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        // some auth logic ...
        h.ServeHTTP(w, r)
    })
}
````

entry -> formatter -> bytes -> writer
entry -> handler (convert entry to format accepted by destination)