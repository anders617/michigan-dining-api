load("@io_bazel_rules_go//go:def.bzl", "go_binary", "go_library")

go_library(
    name = "go_default_library",
    srcs = ["main.go"],
    importpath = "github.com/MichiganDiningAPI/cmd/web",
    visibility = ["//visibility:private"],
    deps = [
        "//internal/util:io",
        "//internal/util:date",
        "//internal/processing:mdiningprocessing",
        "//internal/web:mdiningserver",
        "//db:dynamoclient",
        "@com_github_anders617_mdining_proto//proto:mdining_go_proto",
        "@com_github_golang_glog//:go_default_library",
        "@com_github_golang_protobuf//proto:go_default_library",
        "@grpc_ecosystem_grpc_gateway//runtime:go_default_library",
        "@org_golang_google_grpc//:go_default_library",
        "@com_github_soheilhy_cmux//:go_default_library",
        "@com_github_improbable_eng_grpc_web//go/grpcweb:go_default_library",
    ],
)

go_binary(
    name = "web",
    data = [
        "//public:public_files",
    ],
    embed = [":go_default_library"],
    visibility = ["//visibility:public"],
)
