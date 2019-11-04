# michigan-dining-api
[![Build Status](https://travis-ci.com/anders617/michigan-dining-api.svg?token=cMRcZeh9VAjpBXRsmo8P&branch=master)](https://travis-ci.com/anders617/michigan-dining-api)

## [michigan-dining-api.herokuapp.com](http://michigan-dining-api.herokuapp.com/)

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
bazel run //cmd:web -- --alsologtostderr
```

Run the fetch executable to fill the DiningHalls/Foods/Menus tables:
```shell
bazel run //cmd:fetch -- --alsologtostderr
```

Run the analyze executable to fill the FoodStats table (depends on data from running `//cmd:fetch` above):
```shell
bazel run //cmd:analyze -- --alsologtostderr
```

Run the db executable to create tables:
```shell
bazel run //cmd:db -- --alsologtostderr --create
```

Run the db executable to delete tables:
```shell
bazel run //cmd:db -- --alsologtostderr --delete
```

## Deployment

Currently michigan-dining-api is deployed and hosted on [Heroku](https://www.heroku.com/home) at https://michigan-dining-api.herokuapp.com

In order to deploy your own server:
* Setup the Heroku application to point to this repository
* Add the custom [heroku-buildpack-bazel](https://github.com/anders617/heroku-buildpack-bazel) buildpack to allow building with bazel
* Setup the [HerokuScheduler](https://devcenter.heroku.com/articles/scheduler) add on to run the command `cmd/fetch/fetch`  and `cmd/analyze/analyze` daily in order to fill the data tables
* Set the following Heroku config vars:
    * `AWS_ACCESS_KEY_ID` - Access key used for AWS DynamoDB access
    * `AWS_SECRET_ACCESS_KEY` - Secret used for AWS DynamoDB access
    * `BAZEL_BUILD_PATH` - `//cmd:all`
    * `BAZEL_VERSION` - `1.1.0` (or later version)
* Go to the deploy tab and click deploy branch

## Endpoints
You can click each link for an example \
[/v1/items](https://michigan-dining-api.herokuapp.com/v1/items) \
[/v1/diningHalls](https://michigan-dining-api.herokuapp.com/v1/diningHalls) \
[/v1/filterableEntries](https://michigan-dining-api.herokuapp.com/v1/filterableEntries) \
[/v1/all](https://michigan-dining-api.herokuapp.com/v1/all) \
[/v1/menus?date={yyyy-MM-dd}&diningHall={DINING_HALL}&meal={MEAL}](https://michigan-dining-api.herokuapp.com/v1/menus?date=2019-11-04&diningHall=Bursley%20Dining%20Hall&meal=LUNCH) \
[/v1/stats](https://michigan-dining-api.herokuapp.com/v1/stats)

