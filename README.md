# Gipity

A CLI application for sending queries to ChatGPT.


## Installation

Clone the repo, install dependencies and build.

```bash
  git clone https://github.com/CaveScraps/Gipity.git
  cd Gipity
  go mod tidy
  go build -o gpt
```
You will need to source your own API key from https://platform.openai.com/api-keys.  
Then you will need to make a .env file named ".gipityenv" in the root directory of the project and add the following line:
```bash
  OPENAI_API_KEY="your-api-key-here"
```

## Usage:

```bash
  ./gpt "Your question goes here"
```


## License

[The Unlicense](./LICENSE)
