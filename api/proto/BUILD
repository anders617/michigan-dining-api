load("@io_bazel_rules_go//go:def.bzl", "go_library")
load("@io_bazel_rules_go//proto:def.bzl", "go_proto_library")

proto_library(
    name = "mdining_proto",
    srcs = ["mdining.proto"],
    visibility = ["//visibility:public"],
)

go_proto_library(
    name = "mdining_go_proto",
    compilers = ["@io_bazel_rules_go//proto:go_grpc"],
    importpath = "github.com/MichiganDiningAPI/api/proto",
    proto = ":mdining_proto",
    visibility = ["//visibility:public"],
)

go_library(
    name = "go_default_library",
    embed = [":mdining_go_proto"],
    importpath = "github.com/MichiganDiningAPI/api/proto",
    visibility = ["//visibility:public"],
)
