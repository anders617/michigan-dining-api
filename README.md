# michigan-dining-api
[![Build Status](https://travis-ci.com/anders617/michigan-dining-api.svg?token=cMRcZeh9VAjpBXRsmo8P&branch=master)](https://travis-ci.com/anders617/michigan-dining-api)

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

## Deployment

Currently michigan-dining-api is deployed and hosted on [Heroku](https://www.heroku.com/home) at https://michigan-dining-api.herokuapp.com

In order to deploy your own server:
* Setup the Heroku application to point to this repository
* Add the custom [heroku-buildpack-bazel](https://github.com/anders617/heroku-buildpack-bazel) buildpack to allow building with bazel
* Setup the [HerokuScheduler](https://devcenter.heroku.com/articles/scheduler) add on to run the command `cmd/fetch/fetch` daily in order to fill the tables
* Set the following Heroku config vars:
    * `AWS_ACCESS_KEY_ID` - Access key used for AWS DynamoDB access
    * `AWS_SECRET_ACCESS_KEY` - Secret used for AWS DynamoDB access
    * `BAZEL_BUILD_PATH` - `//cmd:all`
    * `BAZEL_VERSION` - `0.29.1` (or later version)
* Go to the deploy tab and click deploy branch

## Endpoints
[/v1/items](https://michigan-dining-api.herokuapp.com/v1/items) \
[/v1/diningHalls](https://michigan-dining-api.herokuapp.com/v1/diningHalls) \
[/v1/filterableEntries](https://michigan-dining-api.herokuapp.com/v1/filterableEntries) \
[/v1/all](https://michigan-dining-api.herokuapp.com/v1/all)
[/v1/menu](https://michigan-dining-api.herokuapp.com/v1/menu?d)
[/v1/stats](https://michigan-dining-api.herokuapp.com/v1/stats)

