# HashIDGo

HashIDGo is a Go implementation inspired by [hashID](https://github.com/psypanda/hashID).
It identifies the format of a given hash string using pattern definitions from the `prototypes.json` file, which is managed as a submodule under `hashID/`.

## Features

- Detects hash types based on regular expressions and metadata.
- Uses the latest `prototypes.json` from the official hashID repository via Git submodule.
- Fast and portable Go implementation.

## Project Structure

- `main.go` — Main program logic.
- `hashID/` — Git submodule containing the original hashID resources, including `prototypes.json`.

## Getting Started

1. Clone this repository and initialize submodules:
   ```sh
   git clone --recurse-submodules <this-repo-url>
   ```

2. Build the program:
   ```sh
   go build -o hashidgo main.go
   ```

3. Run the program:
   ```sh
   ./hashidgo <Hash>
   ```

## License

This project follows the license of the original [hashID](https://github.com/psypanda/hashID) project.