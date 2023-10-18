# Zap Parser for Bark

🚀 **Zap Parser** is a utility that simplifies the process of transferring logs generated using Uber's Zap library to a [Bark](https://github.com/techrail/bark) server for easy monitoring and analysis.

## Features

- Import logs from Uber's Zap library into Bark.
- Customizable options for specifying log file path, service name, and session name.
- Hassle-free integration with Bark for log analysis.

## Table of Contents

- [Getting Started](#getting-started)
- [Usage](#usage)
- [Contributing](#contributing)
- [License](#license)

## Getting Started

### Manual Build

1. Clone this repository to your local machine:

   ```bash
   git clone https://github.com/techrail/bark
   ```

2. Navigate to the project directory:

   ```bash
   go build cmd/zapparser/zapparser.go
   ```

3. Execute the Zap Parser with the desired parameters, as explained in the [Usage](#usage) section.

## Usage

You can execute the Zap Parser by running the following command, customizing the options to your needs:

```shell
./zapparser.exe -file="log.txt" -db="postgres://auth_user:pass@localhost:5432/bark?sslmode=disable" -service="example service" -session="example session"
```

- `-file` (optional): The path to the Zap log file you want to process. Defaults to `log.txt`.
- `-db` (or `-server`) (required): The URL of the Bark server or the database connection URL.
- `-service` (optional): The name of the service associated with the log entries. Defaults to "No service name."
- `-session` (optional): The name of the session for the log entries. Defaults to "No session name."
- `-format` (optional): Specifies a custom time format for timestamps in log entries (e.g., "2006-01-02 15:04:05").

By executing this command, the Zap Parser will process the specified Zap log file and send the log entries to the Bark server.

### Timestamp Handling

The Zap Parser is designed to handle timestamps in log entries gracefully. If the input log file contains a Unix timestamp in the log entries, it will be used for timestamp processing. However, if log entries do not include Unix timestamps, the parser will default to using the RFC3339 format for timestamps.

You can further customize the timestamp format using the `-format` option. This allows you to specify a time format other than RFC3339, giving you flexibility in how timestamps are interpreted. The `-format` option accepts standard Go time layout formats.

By providing this flexibility, the Zap Parser ensures that log entries with different timestamp formats can be effectively processed and sent to the Bark server.

> [!WARNING]  
> The Zap Parser exclusively processes logs generated in production mode and captures data only for the fields "level," "timestamp," and "msg."
## Contributing

If you'd like to contribute to this project, please read our [Contribution Guidelines](../../CONTRIBUTING.md).

## License

This project is licensed under the [MIT License](../../LICENSE).

---

Happy logging with Bark! 🌟

**Zap Parser for Bark** is not affiliated with or endorsed by Uber Technologies, Inc.