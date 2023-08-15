# Rontgen (rn)

Recursively find matches for a given regex pattern in file names or plain text file content.

## Usage

By default, Rontgen will search through the current working directory. You can specify a path of a file or directory.

```
Usage: rn [-v] [-i] <pattern> <path>
Flags:
 -v Show version
 -verbose Verbose
Pattern:
 <path> Path to directory or file
 <pattern> Pattern to search for
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

To build the project yourself you need to have the Go language and Git installed on your device.

Run the following commands:
```zsh
git clone https://github.com/YONEDASH/Rontgen.git
```
```zsh
cd Rontgen
```
```zsh
chmod +x build.sh && ./build.sh
```

Once you are done the binary called ``rn`` was created.

## API

You can simply call the ``Rontgen`` function. 

Its only parameter is a ``Configuration`` struct. Set ``Path`` as the path of an existing file or directory. Set ``Pattern`` as your desired (and compiled) Regex pattern. The struct looks like this:

```go
type Configuration struct {
	Verbose bool
	Path    string
	Pattern *regexp.Regexp
}
```

Now you can finally run the ``Rontgen`` function. It returns an array of Match structs:

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

- Make the ``Rontgen`` function also return errors (maybe with severity?)
- Replace (example: rn <pattern> -r ...)
- Maximum directory depth, depth cap (-dc ...)
- Max file count to search through, file cap (-fc ...)
