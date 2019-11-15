# michigan-dining-api
[![Build Status](https://travis-ci.org/anders617/michigan-dining-api.svg?branch=master)](https://travis-ci.org/anders617/michigan-dining-api)

[michigan-dining-api.herokuapp.com](https://michigan-dining-api.herokuapp.com/) \
[michigan-dining-api.tendiesti.me](https://michigan-dining-api.tendiesti.me/)

A system for scraping and serving information from the University of Michigan dining API. This repository contains the code for fetching, analyzing and serving dining information. For how to make use of this service, see the [usage](#Usage) section or visit the [mdining-proto](https://github.com/anders617/mdining-proto) repository for protobuf service definitions.

**[Setup](#Setup)** \
**[Executables](#Executables)** \
**[Deployment](#Deployment)** \
**[Usage](#Usage)** 

## Setup
Clone this repo
```shell
git clone https://github.com/anders617/michigan-dining-api.git
```

Install the [Bazel](https://docs.bazel.build/versions/master/install.html) build system

## Executables
This project uses the [glog](https://github.com/golang/glog) library for logging. The `--alsologtostderr` flag can be specified to send log output to stderr.

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

Run the testing client executable to connect to a instance of the web server:
```shell
bazel run //cmd:client -- --alsologtostderr --address=michigan-dining-api.tendiesti.me:443 --use_credentials
```

## Deployment
**[Containers](#Containers)** \
**[AWS](#AWS)** \
**[Heroku](#Heroku)**
### Containers
The `//cmd:web`, `//cmd:fetch`, and `//cmd:analyze` executables all have rules for creating [distroless](https://github.com/GoogleContainerTools/distroless) docker images:
* `//cmd/web:web_image` 
* `//cmd:fetch:fetch_image` 
* `//cmd:analyze:analyze_image`

There are also rules for pushing these container images to container registries:
* `//cmd/web:web_image_publish`
* `//cmd/fetch:fetch_image_publish`
* `//cmd/analyze:analyze_image_publish`

Note that each target above needs to be run with the `--platforms=@io_bazel_rules_go//go/toolchain:linux_amd64` flag set to ensure the binaries are built for running in a linux container. Alternatively, you can specify `--config=container` to use the config set in the `.bazelrc` to avoid having to remember the long platform name.

Currently these rules are configures to push the images to gcr.io/michigandiningapi but can be easily configured to publish to other container registries by editing the rules in the BUILD files.

This means that the latest container image builds for each executable are available at:
* gcr.io/michigandiningapi/web:latest
* gcr.io/michigandiningapi/fetch:latest
* gcr.io/michigandiningapi/analyze:latest

Note that these container images need to have the `AWS_ACCESS_KEY_ID` and `AWS_SECRET_ACCESS_KEY` set when run for the AWS account which will host the dynamodb data tables.

Note that since these are distroless docker images, only the bare minimum for running the executable is included so no shells or other standard Linux programs are included. This means that traditional Docker healthchecks that depend on shell commands will not work and should not be used for determining container health.
### AWS
Currently michigan-dining-api is deployed and hosted on [AWS](https://aws.amazon.com/) using the [Elastic Container Service](https://aws.amazon.com/ecs/) at [michigan-dining-api.tendiesti.me](https://michigan-dining-api.tendiesti.me).

There is a [task definition](https://docs.aws.amazon.com/AmazonECS/latest/developerguide/create-task-definition.html) for each executable container image (listed above). Within each task definition, the `AWS_SECRET_ACCESS_KEY` and `AWS_SECRET_ACCESS_KEY` environment variables must be specified for the AWS account hosting the dynamodb tables. The analyze and web tasks have 0.5GiB memory and 0.25vCPU allocated. The fetch task requires more memory and is allocated 1.0GiB memory and 0.25vCPU.

There is a [Service](https://docs.aws.amazon.com/AmazonECS/latest/developerguide/ecs_services.html) defined for the web server using the web task definition. This service is deployed on a [Fargate](https://docs.aws.amazon.com/AmazonECS/latest/developerguide/AWS_Fargate.html) cluster. The web service is configured to include [load balancing](https://docs.aws.amazon.com/AmazonECS/latest/developerguide/service-load-balancing.html) using a [network load balancer](https://docs.aws.amazon.com/elasticloadbalancing/latest/network/network-load-balancers.html). The network load balancer is configured with an SSL/TLS certificate on its :443 listener and decrypts HTTPS traffic before it is forwarded to the web server. It is important this is a network load balancer instead of an application load balancer since AWS application load balancers do not handle grpc style HTTP/2 traffic correctly.

There are [scheduled tasks](https://docs.aws.amazon.com/AmazonECS/latest/developerguide/scheduled_tasks.html) for the fetch and analyze tasks to run once daily in order to update the dynamodb tables.

### Heroku
Currently michigan-dining-api is deployed and hosted on [Heroku](https://www.heroku.com/home) at https://michigan-dining-api.herokuapp.com.

Heroku is not optimal for hosting grpc servers since it does not support HTTP/2. Therefore, if you plan to take advantage of grpc, I recommend you use a different provider such as AWS.

In order to deploy your own server:
* Setup the Heroku application to point to this repository
* Add the custom [heroku-buildpack-bazel](https://github.com/anders617/heroku-buildpack-bazel) buildpack to allow building with bazel
* Setup the [HerokuScheduler](https://devcenter.heroku.com/articles/scheduler) add on to run the command `cmd/fetch/fetch`  and `cmd/analyze/analyze` daily in order to fill the data tables
* Set the following Heroku config vars:
    * `AWS_ACCESS_KEY_ID` - Access key used for AWS DynamoDB access
    * `AWS_SECRET_ACCESS_KEY` - Secret used for AWS DynamoDB access
    * `BAZEL_BUILD_PATH` - `//cmd:all`
    * `BAZEL_VERSION` - `1.1.0` (or later version)
    * `BUILD_CACHE_LOCATION` - Address of a bazel remote cache server (optional)
* Go to the deploy tab and click deploy branch

## Usage
There are examples of grpc usage and client libraries in the [mdining-proto](https://github.com/anders617/mdining-proto) library. This library also contains the proto definitions of messages and services provided by this service.
### REST Endpoints
[/v1/items](https://michigan-dining-api.herokuapp.com/v1/items) \
[/v1/diningHalls](https://michigan-dining-api.herokuapp.com/v1/diningHalls) \
[/v1/filterableEntries](https://michigan-dining-api.herokuapp.com/v1/filterableEntries) \
[/v1/all](https://michigan-dining-api.herokuapp.com/v1/all) \
[/v1/menus?date={yyyy-MM-dd}&diningHall={DINING_HALL}&meal={MEAL}](https://michigan-dining-api.herokuapp.com/v1/menus?date=2019-11-04&diningHall=Bursley%20Dining%20Hall&meal=LUNCH) \
[/v1/foods?name={LOWERCASE_FOOD_NAME}&date={yyyy-MM-dd}&meal={MEAL}](https://michigan-dining-api/herokuapp.com/v1/foods?name=chicken%20tenders&date=2019-11-08&meal=DINNER) \
[/v1/stats](https://michigan-dining-api.herokuapp.com/v1/stats) \
[/v1/hearts](https://michigan-dining-api.herokuapp.com/v1/hearts?keys=chicken%20tenders)

