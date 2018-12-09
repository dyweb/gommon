# Slice on stack

- slice must have static length, i.e. `make(0, i)` where `i` is variable will cause it to escape
- no `parameter to indirect call`
  - I am not sure I understand it correctly, but is seems when use methods from an interface, it's parameter to indirect 
  call because it's opaque to the compiler so it can't analysis what happens to the bytes slice