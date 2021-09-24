<p align="center"><a href="#readme"><img src="https://gh.kaos.st/{{SHORT_NAME}}.svg"/></a></p>

<p align="center">
  <a href="https://kaos.sh/w/{{SHORT_NAME}}/ci"><img src="https://kaos.sh/w/{{SHORT_NAME}}/ci.svg" alt="GitHub Actions CI Status" /></a>
  <a href="https://kaos.sh/r/{{SHORT_NAME}}"><img src="https://kaos.sh/r/{{SHORT_NAME}}.svg" alt="GoReportCard" /></a>
  <a href="https://kaos.sh/b/{{SHORT_NAME}}"><img src="https://kaos.sh/b/{{CODEBEAT_UUID}}.svg" alt="Codebeat badge" /></a>
  <a href="https://kaos.sh/w/{{SHORT_NAME}}/codeql"><img src="https://kaos.sh/w/{{SHORT_NAME}}/codeql.svg" alt="GitHub Actions CodeQL Status" /></a>
  <a href="#license"><img src="https://gh.kaos.st/apache2.svg"></a>
</p>

<p align="center"><a href="#installation">Installation</a> • <a href="#command-line-completion">Command-line completion</a> • <a href="#man-documentation">Man documentation</a> • <a href="#usage">Usage</a> • <a href="#build-status">Build Status</a> • <a href="#contributing">Contributing</a> • <a href="#license">License</a></p>

<br/>

`{{SHORT_NAME}}` is a {{DESC_README}}.

### Installation

#### From source

To build the `{{SHORT_NAME}}` from scratch, make sure you have a working Go 1.16+ workspace (_[instructions](https://golang.org/doc/install)_), then:

```
go install github.com/essentialkaos/{{SHORT_NAME}}
```

#### Prebuilt binaries

You can download prebuilt binaries for Linux and OS X from [EK Apps Repository](https://apps.kaos.st/{{SHORT_NAME}}/latest):

```bash
bash <(curl -fsSL https://apps.kaos.st/get) {{SHORT_NAME}}
```

### Command-line completion

You can generate completion for `bash`, `zsh` or `fish` shell.

Bash:
```bash
sudo {{SHORT_NAME}} --completion=bash 1> /etc/bash_completion.d/{{SHORT_NAME}}
```

ZSH:
```bash
sudo {{SHORT_NAME}} --completion=zsh 1> /usr/share/zsh/site-functions/{{SHORT_NAME}}
```

Fish:
```bash
sudo {{SHORT_NAME}} --completion=fish 1> /usr/share/fish/vendor_completions.d/{{SHORT_NAME}}.fish
```

### Man documentation

You can generate man page using next command:

```bash
{{SHORT_NAME}} --generate-man | sudo gzip > /usr/share/man/man1/{{SHORT_NAME}}.1.gz
```

### Usage

```

```

### Build Status

| Branch | Status |
|--------|----------|
| `master` | [![CI](https://kaos.sh/w/{{SHORT_NAME}}/ci.svg?branch=master)](https://kaos.sh/w/{{SHORT_NAME}}/ci?query=branch:master) |
| `develop` | [![CI](https://kaos.sh/w/{{SHORT_NAME}}/ci.svg?branch=develop)](https://kaos.sh/w/{{SHORT_NAME}}/ci?query=branch:develop) |

### Contributing

Before contributing to this project please read our [Contributing Guidelines](https://github.com/essentialkaos/contributing-guidelines#contributing-guidelines).

### License

[Apache License, Version 2.0](http://www.apache.org/licenses/LICENSE-2.0)

<p align="center"><a href="https://essentialkaos.com"><img src="https://gh.kaos.st/ekgh.svg"/></a></p>
