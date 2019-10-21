load("@io_bazel_rules_go//go:def.bzl", "go_binary", "go_library")

go_library(
    name = "go_default_library",
    srcs = ["main.go"],
    importpath = "github.com/MichiganDiningAPI/cmd/db",
    visibility = ["//visibility:private"],
    deps = [
        "//db:dynamoclient",
        "@com_github_golang_glog//:go_default_library",
    ],
)

go_binary(
    name = "db",
    data = [],
    embed = [":go_default_library"],
    visibility = ["//visibility:public"],
)
