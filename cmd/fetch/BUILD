load("@io_bazel_rules_go//go:def.bzl", "go_binary", "go_library")

go_library(
    name = "go_default_library",
    srcs = ["main.go"],
    importpath = "github.com/MichiganDiningAPI/cmd/fetch",
    visibility = ["//visibility:private"],
    deps = [
        "//internal/processing:mdiningprocessing",
        "//db:dynamoclient",
        "//api/mdining:mdiningclient",
        "//internal/util:io",
        "//internal/util:containers",
        "@com_github_anders617_mdining_proto//proto:mdining_go_proto",
        "@com_github_aws_aws_sdk_go_v2//aws:go_default_library",
        "@com_github_aws_aws_sdk_go_v2//aws/endpoints:go_default_library",
        "@com_github_aws_aws_sdk_go_v2//aws/external:go_default_library",
        "@com_github_aws_aws_sdk_go_v2//service/dynamodb:go_default_library",
        "@com_github_aws_aws_sdk_go_v2//service/dynamodb/dynamodbattribute:go_default_library",
        "@com_github_golang_glog//:go_default_library",
        "@com_github_golang_protobuf//proto:go_default_library",
        "@com_github_golang_protobuf//jsonpb:go_default_library",
    ],
)

go_binary(
    name = "fetch",
    data = ["@com_github_anders617_mdining_proto//proto/sample:proto_sample_data",],
    embed = [":go_default_library"],
    visibility = ["//visibility:public"],
)
