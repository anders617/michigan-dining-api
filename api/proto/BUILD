load("@io_bazel_rules_go//go:def.bzl", "go_library")
load("@io_bazel_rules_go//proto:def.bzl", "go_proto_library")

proto_library(
    name = "dininghalls_proto",
    srcs = ["dininghalls.proto"],
    visibility = ["//visibility:public"],
)

proto_library(
    name = "items_proto",
    srcs = ["items.proto"],
    visibility = ["//visibility:public"],
    deps = [
        ":food_proto",
    ]
)

proto_library(
    name = "filterable_entries_proto",
    srcs = ["filterableentries.proto"],
    visibility = ["//visibility:public"],
)

proto_library(
    name = "menu_proto",
    srcs = ["menu.proto"],
    visibility = ["//visibility:public"],
)

proto_library(
    name = "food_proto",
    srcs = ["food.proto"],
    visibility = ["//visibility:public"],
    deps = [
        ":menu_proto"
    ]
)

proto_library(
    name = "foodstat_proto",
    srcs = ["foodstat.proto"],
    visibility = ["//visibility:public"],
)

proto_library(
    name = "mdining_proto",
    srcs = ["mdining.proto"],
    visibility = ["//visibility:public"],
    deps = [
        ":dininghalls_proto",
        ":filterable_entries_proto",
        ":items_proto",
        ":menu_proto",
        ":food_proto",
        ":foodstat_proto",
        "@go_googleapis//google/api:annotations_proto",
    ],
)

proto_library(
    name = "mdiningapi_proto",
    srcs = ["mdiningapi.proto"],
    visibility = ["//visibility:public"],
    deps = [
        ":dininghalls_proto",
        ":filterable_entries_proto",
        ":items_proto",
        ":menu_proto",
        ":food_proto",
    ],
)

go_proto_library(
    name = "mdiningapi_proto_library",
    compilers = [
        "@io_bazel_rules_go//proto:go_grpc",
        "@grpc_ecosystem_grpc_gateway//protoc-gen-grpc-gateway:go_gen_grpc_gateway",
    ],
    importpath = "github.com/MichiganDiningAPI/api/proto/mdiningapi",
    protos = [
        ":mdiningapi_proto",
    ],
    visibility = ["//visibility:public"],
    deps = [
        "//api/proto:go_default_library",
        "@com_github_golang_protobuf//descriptor:go_default_library",
        "@com_github_golang_protobuf//proto:go_default_library",
        "@grpc_ecosystem_grpc_gateway//runtime:go_default_library",
        "@grpc_ecosystem_grpc_gateway//utilities:go_default_library",
        "@org_golang_google_genproto//googleapis/api/annotations:go_default_library",
        "@org_golang_google_grpc//:go_default_library",
        "@org_golang_google_grpc//codes:go_default_library",
        "@org_golang_google_grpc//grpclog:go_default_library",
        "@org_golang_google_grpc//status:go_default_library",
    ],
)

go_proto_library(
    name = "go_default_library",
    compilers = [
        "@io_bazel_rules_go//proto:go_grpc",
        "@grpc_ecosystem_grpc_gateway//protoc-gen-grpc-gateway:go_gen_grpc_gateway",
    ],
    importpath = "github.com/MichiganDiningAPI/api/proto",
    protos = [
        ":dininghalls_proto",
        ":filterable_entries_proto",
        ":items_proto",
        ":mdining_proto",
        ":menu_proto",
        ":food_proto",
        ":foodstat_proto",
    ],
    visibility = ["//visibility:public"],
    deps = [
        "@com_github_golang_protobuf//descriptor:go_default_library",
        "@com_github_golang_protobuf//proto:go_default_library",
        "@grpc_ecosystem_grpc_gateway//runtime:go_default_library",
        "@grpc_ecosystem_grpc_gateway//utilities:go_default_library",
        "@org_golang_google_genproto//googleapis/api/annotations:go_default_library",
        "@org_golang_google_grpc//:go_default_library",
        "@org_golang_google_grpc//codes:go_default_library",
        "@org_golang_google_grpc//grpclog:go_default_library",
        "@org_golang_google_grpc//status:go_default_library",
    ],
)
