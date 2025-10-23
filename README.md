<p align="center">
  <a href="jail"><img src="https://blogger.googleusercontent.com/img/b/R29vZ2xl/AVvXsEi39c9ab6rTHulzqrvy45M_omMN8cUyRxfaAph0UhlpubhMxgnJVyOEarYGmHNZgt1uUZmO8cobmrloSiAfxUjgjNOVvRZrF9n9b5tO0S-sG7e9DHfalqyYQZm6aY1jV55IzPbGPA/s1600/freebsd_jail.png" width="200" height="200" border="0" alt="jail"></a>
</p>
<p align="center">
  <a href="https://godoc.org/github.com/briandowns/jail"><img src="https://godoc.org/github.com/briandowns/jail?status.svg" alt="GoDoc"></a>
  <a href="https://opensource.org/licenses/BSD-3-Clause"><img src="https://img.shields.io/badge/License-BSD%203--Clause-orange.svg?" alt="License"></a>
  <a href="https://github.com/briandowns/jail/releases"><img src="https://img.shields.io/badge/version-0.1.0-green.svg?" alt="Version"></a>
</p>

# Jail

Jail provides native FreeBSD Jail syscalls in Go.  As of now, FreeBSD 14.3 defines the Jail API at version 2.  This is the only version supported at this time.  The syscalls supported are:

* jail(2)
* jail_set(2)
* jail_get(2)
* jail_attach(2)
* jail_remove(2)

To get specifics on the syscalls themselves can be referenced [here](https://www.freebsd.org/cgi/man.cgi?query=jail_set&apropos=0&sektion=2&manpath=FreeBSD+14.3-RELEASE&arch=default&format=html).

For examples, please reference the `examples` directory.

## Contributing

Please feel free to open a PR!

## License

Jail source code is available under the BSD 2 clause [License](/LICENSE).

## Contact

[@bdowns328](http://twitter.com/bdowns328)

## Image Credit

www.debarbora.com
