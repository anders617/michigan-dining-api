load("@io_bazel_rules_go//go:def.bzl", "go_binary", "go_library")

go_library(
    name = "go_default_library",
    srcs = ["main.go"],
    importpath = "github.com/MichiganDiningAPI/cmd/web",
    visibility = ["//visibility:private"],
    deps = [
        "//api/proto:go_default_library",
        "@com_github_golang_protobuf//jsonpb:go_default_library_gen",
        "@com_github_google_logger//:go_default_library",
        "@org_golang_google_grpc//:go_default_library",
    ],
)

go_binary(
    name = "web",
    embed = [":go_default_library"],
    visibility = ["//visibility:public"],
)
