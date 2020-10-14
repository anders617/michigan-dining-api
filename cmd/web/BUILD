load("@io_bazel_rules_go//go:def.bzl", "go_binary", "go_library")
load("@io_bazel_rules_docker//go:image.bzl", "go_image")
load("@io_bazel_rules_docker//container:container.bzl", "container_push")

go_library(
    name = "go_default_library",
    srcs = ["main.go"],
    importpath = "github.com/MichiganDiningAPI/cmd/web",
    visibility = ["//visibility:private"],
    deps = [
        "//api/analytics:analyticsclient",
        "//db:dynamoclient",
        "//internal/processing:mdiningprocessing",
        "//internal/util:date",
        "//internal/util:io",
        "//internal/web:mdiningserver",
        "@com_github_anders617_mdining_proto//proto:mdining_go_proto",
        "@com_github_golang_glog//:go_default_library",
        "@com_github_golang_protobuf//proto:go_default_library",
        "@com_github_improbable_eng_grpc_web//go/grpcweb:go_default_library",
        "@com_github_soheilhy_cmux//:go_default_library",
        "@grpc_ecosystem_grpc_gateway//runtime:go_default_library",
        "@org_golang_google_grpc//:go_default_library",
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

# Build with --platforms=@io_bazel_rules_go//go/toolchain:linux_amd64
go_image(
    name = "web_image",
    args = [
        "--alsologtostderr",
    ],
    data = [
        "//public:public_files",
    ],
    embed = [":go_default_library"],
    visibility = ["//visibility:public"],
)

# A common pattern when users want to avoid trampling
# on each other's images during development.
container_push(
    name = "web_image_publish",
    format = "Docker",
    image = ":web_image",
    # Any of these components may have variables.
    registry = "gcr.io",
    repository = "michigandiningapi/web",
    tag = "latest",
)
