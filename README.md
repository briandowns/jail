# Jail

Jail provides native FreeBSD Jail syscalls in Go.  Other implementations require CGO however we've reimplemented some of the native kernel code in Go and provided wrappers to avoid this need.

As of now, FreeBSD 11.1 defines the Jail API at version 2. This is the only version supported at this time.

## Examples

Lock a process into a jail.

```
jo := &JailOpts{
    Path: "/path/to/jail/dir",
    Name: "jailName",
    Hostname: "jailHostname",
}
jail, err := jail.New(jo)
if err != nil {
    log.Fatalln(err)
}
```
