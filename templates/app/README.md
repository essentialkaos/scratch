<p align="center"><a href="#readme"><img src="https://gh.kaos.st/{{SHORT_NAME}}.svg"/></a></p>

<p align="center">
  <a href="https://kaos.sh/w/{{SHORT_NAME}}/ci"><img src="https://kaos.sh/w/{{SHORT_NAME}}/ci.svg" alt="GitHub Actions CI Status" /></a>
  <a href="https://kaos.sh/l/{{SHORT_NAME}}"><img src="https://kaos.sh/l/{{CODECLIMATE_ID}}.svg" alt="Code Climate Maintainability" /></a>
  <a href="https://kaos.sh/b/{{SHORT_NAME}}"><img src="https://kaos.sh/b/{{CODEBEAT_UUID}}.svg" alt="Codebeat badge" /></a>
  <a href="https://kaos.sh/w/{{SHORT_NAME}}/codeql"><img src="https://kaos.sh/w/{{SHORT_NAME}}/codeql.svg" alt="GitHub Actions CodeQL Status" /></a>
  <a href="#license"><img src="https://gh.kaos.st/apache2.svg"></a>
</p>

<p align="center"><a href="#installation">Installation</a> • <a href="#usage">Usage</a> • <a href="#ci-status">CI Status</a> • <a href="#contributing">Contributing</a> • <a href="#license">License</a></p>

<br/>

`{{SHORT_NAME}}` is {{DESC_README}}.

### Installation

#### From source

To build the `{{SHORT_NAME}}` from scratch, make sure you have a working Go 1.20+ workspace (_[instructions](https://go.dev/doc/install)_), then:

```bash
go install github.com/essentialkaos/{{SHORT_NAME}}@latest
```

#### From [ESSENTIAL KAOS Public Repository](https://pkgs.kaos.st)

```bash
sudo yum install -y https://pkgs.kaos.st/kaos-repo-latest.el$(grep 'CPE_NAME' /etc/os-release | tr -d '"' | cut -d':' -f5).noarch.rpm
sudo yum install {{SHORT_NAME}}
```

### Usage

```

```

### CI Status

| Branch | Status |
|--------|----------|
| `master` | [![CI](https://kaos.sh/w/{{SHORT_NAME}}/ci.svg?branch=master)](https://kaos.sh/w/{{SHORT_NAME}}/ci?query=branch:master) |
| `develop` | [![CI](https://kaos.sh/w/{{SHORT_NAME}}/ci.svg?branch=develop)](https://kaos.sh/w/{{SHORT_NAME}}/ci?query=branch:develop) |

### Contributing

Before contributing to this project please read our [Contributing Guidelines](https://github.com/essentialkaos/contributing-guidelines#contributing-guidelines).

### License

[Apache License, Version 2.0](http://www.apache.org/licenses/LICENSE-2.0)

<p align="center"><a href="https://essentialkaos.com"><img src="https://gh.kaos.st/ekgh.svg"/></a></p>
