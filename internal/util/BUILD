load("@io_bazel_rules_go//go:def.bzl", "go_library")

go_library(
    name = "io",
    srcs = ["io.go"],
    importpath = "github.com/MichiganDiningAPI/internal/util/io",
    visibility = ["//visibility:public"],
    deps = [
        "@com_github_golang_glog//:go_default_library",
        "@com_github_golang_protobuf//proto:go_default_library",
    ],
)

go_library(
    name = "containers",
    srcs = ["containers.go"],
    importpath = "github.com/MichiganDiningAPI/internal/util/containers",
    visibility = ["//visibility:public"],
    deps = [
        "@com_github_golang_glog//:go_default_library",
        "@com_github_golang_protobuf//proto:go_default_library",
    ],
)

go_library(
    name = "date",
    srcs = ["date.go"],
    importpath = "github.com/MichiganDiningAPI/internal/util/date",
    visibility = ["//visibility:public"],
    deps = [
        "@com_github_golang_glog//:go_default_library",
        "@com_github_golang_protobuf//proto:go_default_library",
    ],
)
