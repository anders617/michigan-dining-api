# michigan-dining-api

A system for scraping and serving information from the University of Michigan dining API.

## Setup
Clone this repo
```shell
git clone https://github.com/anders617/michigan-dining-api.git
```

Install the [Bazel](https://docs.bazel.build/versions/master/install.html) build system

## Executables

Run the web server:
```shell
bazel run //cmd/web:web
```

Run the fetch executable:
```shell
bazel run //cmd/fetch:fetch
```
