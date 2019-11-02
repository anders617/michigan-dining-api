load("@io_bazel_rules_go//go:def.bzl", "go_library")

go_library(
    name = "mdiningclient",
    srcs = ["mdiningclient.go"],
    importpath = "github.com/MichiganDiningAPI/api/mdining/mdiningclient",
    visibility = ["//visibility:public"],
    deps = [
        "//api/proto:go_default_library",
        "//api/proto:mdiningapi_proto_library",
        "//internal/util:date",
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
