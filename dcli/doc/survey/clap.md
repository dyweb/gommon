# Clap

- [Command line apps in Rust](https://rust-cli.github.io/book/index.html)
- TBH, compared with cobra it's hard to use ... requires to run the dispatch logic by doing pattern matching by yourself
  - maybe I not using it in the right way

## StructOpt

- https://github.com/TeXitoi/structopt
- https://clap.rs/2019/03/08/clap-v3-update-structopt/

```rust
use structopt::StructOpt;

/// Search for a pattern in a file and display the lines that contain it.
#[derive(StructOpt)]
struct Cli {
    /// The pattern to look for
    pattern: String,
    /// The path to the file to read
    #[structopt(parse(from_os_str))]
    path: std::path::PathBuf,
}

```