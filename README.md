<div align="center"><img src="https://github.com/hollerith/pingo/assets/659626/772ebd26-6a22-4cc9-add0-d8a97b322e12" width="50%"></div>

# Pingo

Pingo is a simple command line tool that attempts to guess the operating system of a network host based on TTL (Time To Live) values in ICMP echo replies.

## Usage

Pingo sends an ICMP echo request to a network host and uses the TTL value in the echo reply along with the count of network hops to guess the host's operating system.

The target is the hostname or IP address of the network host to ping. The usage is as follows:

```bash
pingo [target] [--retries N] [--tolerance N] [--verbose]
```

Refer to the [manual](MANPAGE.md) for a detailed explanation of the options and examples.

## Building

To build `pingo`, make sure you have Go installed, then run

```bash
go mod init
go get github.com/spf13/pflag
go build pingo.go
```

Optionally, copy pingo to somewhere like `~/.local/bin`

## Contributing

Pull requests are welcome. For major changes, please open an issue first to discuss what you would like to change.

## Disclaimer

Pingo is provided as-is, without any warranty. Use it at your own risk and always ensure you have the necessary permissions before scanning any network host. Responsible use is strongly encouraged.

## License

MIT
