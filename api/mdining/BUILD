load("@io_bazel_rules_go//go:def.bzl", "go_library")

go_library(
    name = "mdiningclient",
    srcs = ["mdiningclient.go"],
    importpath = "github.com/MichiganDiningAPI/api/mdining/mdiningclient",
    visibility = ["//visibility:public"],
    deps = [
        "//internal/util:date",
        "@com_github_anders617_mdining_proto//proto:mdining_go_proto",
        "@com_github_anders617_mdining_proto//proto:mdiningapi_go_proto",
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

go_library(
    name = "mdiningclient2",
    srcs = ["mdiningclient2.go"],
    importpath = "github.com/MichiganDiningAPI/api/mdining/mdiningclient2",
    visibility = ["//visibility:public"],
    deps = [
        "//internal/util:date",
        "@com_github_anders617_mdining_proto//proto:mdining_go_proto",
        "@com_github_anders617_mdining_proto//proto:mdiningapi2_go_proto",
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
