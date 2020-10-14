load("@io_bazel_rules_go//go:def.bzl", "go_binary", "go_library")
load("@io_bazel_rules_docker//go:image.bzl", "go_image")
load("@io_bazel_rules_docker//container:container.bzl", "container_push")

go_library(
    name = "go_default_library",
    srcs = ["main.go"],
    importpath = "github.com/MichiganDiningAPI/cmd/fetch",
    visibility = ["//visibility:private"],
    deps = [
        "//api/mdining:mdiningclient",
        "//db:dynamoclient",
        "//internal/processing:mdiningprocessing",
        "//internal/util:containers",
        "//internal/util:io",
        "@com_github_anders617_mdining_proto//proto:mdining_go_proto",
        "@com_github_aws_aws_sdk_go_v2//aws:go_default_library",
        "@com_github_aws_aws_sdk_go_v2//aws/endpoints:go_default_library",
        "@com_github_aws_aws_sdk_go_v2//aws/external:go_default_library",
        "@com_github_aws_aws_sdk_go_v2//service/dynamodb:go_default_library",
        "@com_github_aws_aws_sdk_go_v2//service/dynamodb/dynamodbattribute:go_default_library",
        "@com_github_golang_glog//:go_default_library",
        "@com_github_golang_protobuf//jsonpb:go_default_library",
        "@com_github_golang_protobuf//proto:go_default_library",
    ],
)

go_binary(
    name = "fetch",
    data = ["@com_github_anders617_mdining_proto//proto/sample:proto_sample_data"],
    embed = [":go_default_library"],
    visibility = ["//visibility:public"],
)

# Build with --platforms=@io_bazel_rules_go//go/toolchain:linux_amd64
go_image(
    name = "fetch_image",
    args = [
        "--alsologtostderr",
    ],
    data = ["@com_github_anders617_mdining_proto//proto/sample:proto_sample_data"],
    embed = [":go_default_library"],
    visibility = ["//visibility:public"],
)
# A common pattern when users want to avoid trampling
# on each other's images during development.
container_push(
    name = "fetch_image_publish",
    format = "Docker",
    image = ":fetch_image",
    # Any of these components may have variables.
    registry = "gcr.io",
    repository = "michigandiningapi/fetch",
    tag = "latest",
)

go_library(
    name = "fetch2_library",
    srcs = ["main2.go"],
    importpath = "github.com/MichiganDiningAPI/cmd/scrape",
    visibility = ["//visibility:private"],
    deps = [
        "//api/mdining:mdiningclient2",
        "//db:dynamoclient",
        "//internal/util:containers",
        "//internal/processing:mdiningprocessing",
        "@com_github_golang_protobuf//proto:go_default_library",
        "@com_github_golang_glog//:go_default_library",
    ],
)

go_binary(
    name = "fetch2",
    embed = [":fetch2_library"],
    visibility = ["//visibility:public"],
)

# Build with --platforms=@io_bazel_rules_go//go/toolchain:linux_amd64
go_image(
    name = "fetch2_image",
    args = [
        "--alsologtostderr",
    ],
    data = ["@com_github_anders617_mdining_proto//proto/sample:proto_sample_data"],
    embed = [":fetch2_library"],
    visibility = ["//visibility:public"],
)

# A common pattern when users want to avoid trampling
# on each other's images during development.
container_push(
    name = "fetch2_image_publish",
    format = "Docker",
    image = ":fetch2_image",
    # Any of these components may have variables.
    registry = "gcr.io",
    repository = "michigandiningapi/fetch2",
    tag = "latest",
)