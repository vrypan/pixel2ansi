# pixel2ansi

`pixel2ansi` is a command-line tool for analyzing and converting pixel art
images. It can detect the smallest repeating pixel unit, handle color
variations with adjustable tolerance, and render the image as ANSI art in
the terminal. Features include transparent color detection, optional
cropping, and parallel processing for speed.

`pixel2ansi` is able to read PNG, BMP, and JPEG images. However, if you use
a lossy format like JPEG, the "pixel unit" is often blurred and you may
need to adjust the tolerance value (start with `--tolerance=150` or so).

# Installation

- From source: Clone the repo and run `make`
- Binaries: Download from [releases](https://github.com/vrypan/pixel2ansi/releases)
- Homebrew: Install with
```
brew install vrypan/pixel2ansi/pixel2ansi
```

# Usage
To analyse an image, use `pixel2ansi inspect`:

```
$ pixel2ansi inspect bcbc14b662011b1df244857d8bedb0ce.png

Block size: 10x10 pixels
Grid size: 32x32 pixels
```

To output an image as ANSI art, use `pixel2ansi print`:

```
$ pixel2ansi print bcbc14b662011b1df244857d8bedb0ce.png
```

For more info:

```
$ pixel2ansi

Usage:
  pixel2ansi [command]

Available Commands:
  about       Show about info
  completion  Generate the autocompletion script for the specified shell
  help        Help about any command
  inspect     Get unit pixel and grid dimensions
  print       Print image as ANSI blocks
  version     Get the current version number

Flags:
  -h, --help   help for pixel2ansi

Use "pixel2ansi [command] --help" for more information about a command.
```
