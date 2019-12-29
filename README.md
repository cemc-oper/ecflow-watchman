# ecflow-watchman

**WARNING: This program has serious memory leak problem which I am working hard to solve.**

Watch ecflow servers.

## Install

`ecflow-watchman` uses `ecflow-client-go` package which requires ecFlow and boost. 

Set some environment variables before build the library.

```bash
export ECFLOW_BUILD_DIR=/some/path/to/ecflow/build
export ECFLOW_SOURCE_DIR=/some/path/to/ecflow/source
export BOOST_LIB_DIR=/some/path/to/boost/stage/lib
```

Please visit [ecflow-client-go](https://github.com/nwpc-oper/ecflow-client-go) for more information.

Use `Makefile` to build the project and `ecflow_watchman` will be under `bin` directory.

## Getting Started

`ecflow_watchman watch-all` command watches all ecflow servers listed in the config file, 
and sends collected status into a redis server.

```bash
ecflow_watchman watch-all --config-file=/some/config/file/path
```

## Config

The following is an example config file.

```yaml
global:
  scrape_interval: 20s
  scrape_timeout: 10s # not worked

scrape_configs:
  -
    job_name: job name
    owner: owner
    repo: repo
    host: ecflow server host
    port: ecflow server port

sink_config:
  type: redis # only redis is supported
  url: redis url
```

`owner` and `repo` are used in key name for redis.

## Warning

This program has serious memory leak problem caused by getting ecflow status using `ecflow-client-go`.
Strings passed from c++ to goroutine are not released during loops.
I am trying to solve this problem but haven't got any progress.

## License

Copyright 2019, perillaroc at nwpc-oper.

`ecflow-watchman` is licensed under [MIT License](./LICENSE.md).