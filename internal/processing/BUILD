load("@io_bazel_rules_go//go:def.bzl", "go_library")

go_library(
    name = "mdiningprocessing",
    srcs = ["mdiningprocessing.go"],
    importpath = "github.com/MichiganDiningAPI/internal/processing/mdiningprocessing",
    visibility = ["//visibility:public"],
    deps = [
        "//internal/util:containers",
        "@com_github_anders617_mdining_proto//proto:mdining_go_proto",
        "@com_github_golang_glog//:go_default_library",
        "@com_github_golang_protobuf//proto:go_default_library",
    ]
)
