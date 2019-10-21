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
bazel run //cmd/web:web -- --alsologtostderr
```

Run the fetch executable:
```shell
bazel run //cmd/fetch:fetch -- --alsologtostderr
```

Run the db executable to create tables:
```shell
bazel run //cmd/db:db -- --alsologtostderr --create
```

Run the db executable to delete tables:
```shell
bazel run //cmd/db:db -- --alsologtostderr --delete
```
