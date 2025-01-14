# Threadviewer

Threadviewer is a Go command-line application designed to load a given assistant thread by ID and print it to the stdout with markdown formatting suitable for terminal rendering.

## Key Features

- Loads threads using the OpenAI API.
- Formats and prints threads with markdown rendering.
- Uses Cobra for command-line parsing.
- Uses Viper for configuration management, including automatic environment variable binding.

## Getting Started

### Prerequisites

- OpenAI API Key.

### Installation

#### Using Eget:

```sh
eget harnyk/threadviewer
```

#### Manually:

Download a binary from the [releases](https://github.com/harnyk/threadviewer/releases) page and place it in a directory on your PATH.


### Configuration

The configuration file is located **in the current directory** or at `<config>/config/.threadviewer.yaml` and is in YAML format. You can specify the config file path using the `--config` flag.

Example config file:

```yaml
API_KEY: your_openai_api_key
```

### Usage

```sh
./threadviewer --threadID <thread_id> --apiKey <api_key>
```

### Flags

- `--threadID`, `-t`: Thread ID to retrieve.
- `--apiKey`: OpenAI API Key.
- `--config`: Config file path.

### Environment Variables

- `API_KEY`: OpenAI API Key.

## License

This project is licensed under the WTFPL License.

## Contributing

1. Fork the repository.
2. Create a new branch: `git checkout -b feature-branch`
3. Make your changes and commit them: `git commit -m 'Add new feature'`
4. Push to the branch: `git push origin feature-branch`
5. Submit a pull request.
