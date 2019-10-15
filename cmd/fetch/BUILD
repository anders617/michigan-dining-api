load("@io_bazel_rules_go//go:def.bzl", "go_binary", "go_library")

go_library(
    name = "dynamoclient",
    srcs = ["dynamoclient.go"],
    importpath = "github.com/MichiganDiningAPI/cmd/fetch/dynamoclient",
    visibility = ["//visibility:private"],
    deps = [
        "//api/proto:go_default_library",
        "@com_github_aws_aws_sdk_go_v2//aws:go_default_library",
        "@com_github_aws_aws_sdk_go_v2//aws/endpoints:go_default_library",
        "@com_github_aws_aws_sdk_go_v2//aws/external:go_default_library",
        "@com_github_aws_aws_sdk_go_v2//service/dynamodb:go_default_library",
        "@com_github_aws_aws_sdk_go_v2//service/dynamodb/dynamodbattribute:go_default_library",
        "@com_github_golang_glog//:go_default_library",
        "@com_github_golang_protobuf//proto:go_default_library",
    ]
)

go_library(
    name = "mdiningclient",
    srcs = ["mdiningclient.go"],
    importpath = "github.com/MichiganDiningAPI/cmd/fetch/mdiningclient",
    visibility = ["//visibility:private"],
    deps = [
        "//api/proto:go_default_library",
        "//api/proto:mdiningapi_proto_library",
        "@com_github_aws_aws_sdk_go_v2//aws:go_default_library",
        "@com_github_aws_aws_sdk_go_v2//aws/endpoints:go_default_library",
        "@com_github_aws_aws_sdk_go_v2//aws/external:go_default_library",
        "@com_github_aws_aws_sdk_go_v2//service/dynamodb:go_default_library",
        "@com_github_aws_aws_sdk_go_v2//service/dynamodb/dynamodbattribute:go_default_library",
        "@com_github_golang_glog//:go_default_library",
        "@com_github_golang_protobuf//proto:go_default_library",
        "@com_github_golang_protobuf//jsonpb:go_default_library",
    ]
)

go_library(
    name = "go_default_library",
    srcs = ["main.go"],
    importpath = "github.com/MichiganDiningAPI/cmd/fetch",
    visibility = ["//visibility:private"],
    deps = [
        ":dynamoclient",
        ":mdiningclient",
        "//api/proto:go_default_library",
        "//util:io",
        "@com_github_aws_aws_sdk_go_v2//aws:go_default_library",
        "@com_github_aws_aws_sdk_go_v2//aws/endpoints:go_default_library",
        "@com_github_aws_aws_sdk_go_v2//aws/external:go_default_library",
        "@com_github_aws_aws_sdk_go_v2//service/dynamodb:go_default_library",
        "@com_github_aws_aws_sdk_go_v2//service/dynamodb/dynamodbattribute:go_default_library",
        "@com_github_golang_glog//:go_default_library",
        "@com_github_golang_protobuf//proto:go_default_library",
    ],
)

go_binary(
    name = "fetch",
    data = ["//api/proto/sample:proto_sample_data"],
    embed = [":go_default_library"],
    visibility = ["//visibility:public"],
)
