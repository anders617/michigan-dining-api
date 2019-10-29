load("@io_bazel_rules_go//go:def.bzl", "go_binary", "go_library")

go_library(
    name = "go_default_library",
    srcs = ["main.go"],
    importpath = "github.com/MichiganDiningAPI/cmd/web",
    visibility = ["//visibility:private"],
    deps = [
        "//api/proto:go_default_library",
        "//internal/util:io",
        "//db:dynamoclient",
        "@com_github_golang_glog//:go_default_library",
        "@com_github_golang_protobuf//proto:go_default_library",
        "@grpc_ecosystem_grpc_gateway//runtime:go_default_library",
        "@org_golang_google_grpc//:go_default_library",
        "@com_github_soheilhy_cmux//:go_default_library"
    ],
)

go_binary(
    name = "web",
    data = [
        "//api/proto/sample:proto_sample_data",
    ],
    embed = [":go_default_library"],
    visibility = ["//visibility:public"],
)
