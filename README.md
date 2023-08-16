# Rontgen (rn)

Recursively find matches for a given regex pattern in file names or plain text file content.

## Usage

By default, Rontgen will search through the current working directory. You can specify a path of a file or directory.

```
Usage: rn [flags...] <pattern> <path>
Flags:
  -dc int
        Maximum directory depth (default 10)
  -fc int
        Maximum file count (default 100000)
  -fs int
        Maximum file size in kilobytes (default 20000)
  -mc int
        Maximum matches per file (default 1000)
  -n    No colors
  -v    Show version
  -verbose
        Verbose
  <pattern>
        Pattern to search for
  <path>
        Path to directory or file
```

## Installation

To simply install using Homebrew execute the following commands:

```zsh
brew tap yonedash/homebrew-formulae
```
```zsh
brew install rontgen
```

Otherwise you can build it yourself.

## Building it yourself

To build the binary yourself you need to have Git and the Go language installed on your device.

Run the following commands:
```zsh
git clone https://github.com/YONEDASH/Rontgen.git
```
```zsh
cd Rontgen
```
```zsh
go build -o=rn
```

Once you are done the binary called ``rn`` was created. Note that this binary won't have the current version number.

## API

You can simply call the ``Rontgen`` function. 

Its only parameter is a ``Configuration`` struct. Set ``Path`` as the path of an existing file or directory. Set ``Pattern`` as your desired (and compiled) Regex pattern. The struct looks like this:

```go
type Configuration struct {
	Verbose  bool
	Path     string
	Pattern  *regexp.Regexp
	DepthCap int
	SizeCap  int64
	CountCap int
	MatchCap int
}
```

Now you can finally run the ``Rontgen`` function. It returns an array of ``Match`` structs and an error (``nil`` if there is none).

```go
type Match struct {
	Path      string
	NameMatch bool
	Row       int
	Column    int
	Length    int
	Matched   string
	Line      string
}
```

``Path`` is always the path of the matched file.
``NameMatch`` determines whether the file name was matched (``true``) or the file's content was matched (``false``).

If file content was matched the following fields will be set:
- ``Row`` (line) and ``Column`` of the matched text.
- ``Length`` of the matched text.
- ``Matched`` as the matched text itself.
- ``Line`` as the entire line where the text was matched.

If a file name was matched the following fields will be set:
- ``Column`` (index) of the matched name.
- ``Length`` of the matched name.
- ``Matched`` as the matched name itself.


## Todos

- Replace (example: rn <pattern> -r ...)
- ~~Make the ``Rontgen`` function also return errors~~ (maybe with severity?)
- ~~Maximum directory depth, depth cap (-dc ...)~~
- ~~Max file count to search through, file cap (-fc ...)~~
