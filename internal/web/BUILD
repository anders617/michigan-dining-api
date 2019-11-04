load("@io_bazel_rules_go//go:def.bzl", "go_library")

go_library(
    name = "mdiningserver",
    srcs = ["mdiningserver.go"],
    importpath = "github.com/MichiganDiningAPI/internal/web/mdiningserver",
    visibility = ["//visibility:public"],
    deps = [
        "//internal/util:date",
        "//internal/processing:mdiningprocessing",
        "//db:dynamoclient",
        "@com_github_anders617_mdining_proto//proto:mdining_go_proto",
        "@com_github_golang_glog//:go_default_library",
        "@org_golang_google_grpc//codes:go_default_library",
        "@org_golang_google_grpc//status:go_default_library",
    ],
)
