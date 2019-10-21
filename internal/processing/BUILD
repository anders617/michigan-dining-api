load("@io_bazel_rules_go//go:def.bzl", "go_library")

go_library(
    name = "mdiningprocessing",
    srcs = ["mdiningprocessing.go"],
    importpath = "github.com/MichiganDiningAPI/internal/processing/mdiningprocessing",
    visibility = ["//visibility:public"],
    deps = [
        "//api/proto:go_default_library",
        "//internal/util:containers",
        "@com_github_golang_glog//:go_default_library",
        "@com_github_golang_protobuf//proto:go_default_library",
    ]
)
