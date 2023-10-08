# ogpk (opengraph peek)

`ogpk` is a simple CLI tool written in Go that fetches OpenGraph data from a given URL. If the optional dependency `timg` is installed, `ogpk` can also display the `og:image` directly in the terminal.

<img src="./dist/screen.png"/>

### Installation

On macOS:

```bash
brew tap almonk/ogpk
brew install ogpk
```

On linux:

* Go to the [releases page](https://github.com/almonk/ogpk/releases) and download the latest release for your platform.
* Extract the archive and move the executable to a directory in your `PATH` (e.g. `/usr/local/bin`)
* Make the executable executable, e.g.:

```bash
chmod +x /usr/local/bin/ogpk
```

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

Output data as JSON:
```bash
ogpk [URL] --json
```

### Building from source

Clone this repository:
```bash
git clone https://github.com/almonk/ogpk.git
```
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
