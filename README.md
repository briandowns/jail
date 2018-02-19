# Jail

Jail provides native FreeBSD Jail syscalls in Go.  Other implementations require CGO however we've reimplemented some of the native kernel code in Go and provided wrappers to avoid this need.

As of now, FreeBSD 11.1 defines the Jail API at version 2. This is the only version supported.

## Contributing

Please feel free to open a PR!