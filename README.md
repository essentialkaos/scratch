<p align="center"><a href="#readme"><img src="https://gh.kaos.st/scratch.svg"/></a></p>

<p align="center">
  <a href="https://kaos.sh/w/scratch/ci"><img src="https://kaos.sh/w/scratch/ci.svg" alt="GitHub Actions CI Status" /></a>
  <a href="https://kaos.sh/b/scratch"><img src="https://kaos.sh/b/3b2ed0f2-1e39-4366-93f6-d955ca22ce3a.svg" alt="Codebeat badge" /></a>
  <a href="https://kaos.sh/w/scratch/codeql"><img src="https://kaos.sh/w/scratch/codeql.svg" alt="GitHub Actions CodeQL Status" /></a>
  <a href="#license"><img src="https://gh.kaos.st/apache2.svg"></a>
</p>

<p align="center"><a href="#usage-demo">Usage demo</a> • <a href="#installation">Installation</a> • <a href="#command-line-completion">Command-line completion</a> • <a href="#usage">Usage</a> • <a href="#contributing">Contributing</a> • <a href="#license">License</a></p>

<br/>

`scratch` is a simple utility for generating blank files for Go apps, utilities and packages.

### Installation

#### From sources

To install the `scratch` from sources, make sure you have a working Go 1.19+ workspace (_[instructions](https://go.dev/doc/install)_), then:

```
go install github.com/essentialkaos/scratch
```

### Command-line completion

You can generate completion for `bash`, `zsh` or `fish` shell.

Bash:
```bash
sudo scratch --completion=bash 1> /etc/bash_completion.d/scratch
```


ZSH:
```bash
sudo scratch --completion=zsh 1> /usr/share/zsh/site-functions/scratch
```


Fish:
```bash
sudo scratch --completion=fish 1> /usr/share/fish/vendor_completions.d/scratch.fish
```

### Man documentation

You can generate man page for `scratch` using next command:

```bash
scratch --generate-man | sudo gzip > /usr/share/man/man1/scratch.1.gz
```

### Usage

```
Usage: scratch {options} template dir

Options

  --no-color, -nc    Disable colors in output
  --help, -h         Show this help message
  --version, -v      Show version

Examples

  scratch package
  List files in template "package"

  scratch package .
  Generate files based on template "package" in current directory

  scratch service $GOPATH/src/github.com/essentialkaos/myapp
  Generate files based on template "service" in given directory
```

### Contributing

Before contributing to this project please read our [Contributing Guidelines](https://github.com/essentialkaos/contributing-guidelines#contributing-guidelines).

### License

[Apache License, Version 2.0](https://www.apache.org/licenses/LICENSE-2.0)

<p align="center"><a href="https://essentialkaos.com"><img src="https://gh.kaos.st/ekgh.svg"/></a></p>
