# ogpk (opengraph peek)

`ogpk` is a simple CLI tool written in Go that fetches OpenGraph data from a given URL. If the optional dependency `timg` is installed, `ogpk` can also display the `og:image` directly in the terminal.


### Usage

To fetch OpenGraph data from a website:
```bash
ogpk [URL]
```

For example:
```bash
ogpk https://example.com
```

To display the `og:image` in the terminal (requires `timg`):
```bash
ogpk [URL] --p
```

### Building from source

Clone this repository:
```bash
git clone https://github.com/YOUR_USERNAME/ogpk.git
```
(Replace `YOUR_USERNAME` with your actual GitHub username)

Navigate to the cloned directory:
```bash
cd ogpk
```

Build the tool:
```bash
go build -o ogpk
```

This will produce an executable named `ogpk` in the current directory.


### Optional Dependency on `timg`

ogpk has an optional dependency on `timg`, a terminal image viewer. If `timg` is installed and available in the `PATH`, ogpk can display the `og:image` directly in the terminal when the `--p` flag is used.

To install `timg`, refer to its [official documentation](https://github.com/hzeller/timg).