load("@bazel_tools//tools/build_defs/repo:http.bzl", "http_archive")

http_archive(
    name = "io_bazel_rules_docker",
    sha256 = "aed1c249d4ec8f703edddf35cbe9dfaca0b5f5ea6e4cd9e83e99f3b0d1136c3d",
    strip_prefix = "rules_docker-0.7.0",
    urls = ["https://github.com/bazelbuild/rules_docker/archive/v0.7.0.tar.gz"],
)

http_archive(
    name = "io_bazel_rules_go",
    urls = [
        "https://storage.googleapis.com/bazel-mirror/github.com/bazelbuild/rules_go/releases/download/0.18.6/rules_go-0.18.6.tar.gz",
        "https://github.com/bazelbuild/rules_go/releases/download/0.18.6/rules_go-0.18.6.tar.gz",
    ],
    sha256 = "f04d2373bcaf8aa09bccb08a98a57e721306c8f6043a2a0ee610fd6853dcde3d",
)

http_archive(
    name = "bazel_gazelle",
    sha256 = "3c681998538231a2d24d0c07ed5a7658cb72bfb5fd4bf9911157c0e9ac6a2687",
    urls = ["https://github.com/bazelbuild/bazel-gazelle/releases/download/0.17.0/bazel-gazelle-0.17.0.tar.gz"],
)

http_archive(
    name = "bazel_skylib",
    sha256 = "eb5c57e4c12e68c0c20bc774bfbc60a568e800d025557bc4ea022c6479acc867",
    strip_prefix = "bazel-skylib-0.6.0",
    urls = ["https://github.com/bazelbuild/bazel-skylib/archive/0.6.0.tar.gz"],
)

load(
    "@io_bazel_rules_docker//repositories:repositories.bzl",
    container_repositories = "repositories",
)

container_repositories()

load(
    "@io_bazel_rules_docker//container:container.bzl",
    "container_pull",
)
load("@io_bazel_rules_go//go:deps.bzl", "go_register_toolchains", "go_rules_dependencies")

go_rules_dependencies()

go_register_toolchains()

load("@bazel_gazelle//:deps.bzl", "gazelle_dependencies", "go_repository")

gazelle_dependencies()

go_repository(
    name = "com_github_gorilla_mux",
    commit = "ed099d42384823742bba0bf9a72b53b55c9e2e38",
    importpath = "github.com/gorilla/mux",
)

go_repository(
    name = "com_github_docker_distribution",
    build_extra_args = ["-exclude=vendor"],
    commit = "79f6bcbe169ad8aa34433c9b1e3b4e38e26ba789",
    importpath = "github.com/docker/distribution",
)

go_repository(
    name = "com_github_jinzhu_gorm",
    commit = "01b66011427614f01e84a473b0303c917179f2a0",
    importpath = "github.com/jinzhu/gorm",
)

go_repository(
    name = "com_github_aws_aws_sdk_go",
    commit = "b4ef9e4c9898fde0de7068049129aca9f749c7ab",
    importpath = "github.com/aws/aws-sdk-go",
)

go_repository(
    name = "com_github_opencontainers_go_digest",
    commit = "ac19fd6e7483ff933754af248d80be865e543d22",
    importpath = "github.com/opencontainers/go-digest",
)

go_repository(
    name = "com_github_prometheus_client_golang",
    commit = "3d8379da8fc2309dd563f593f15a739054c098bc",
    importpath = "github.com/prometheus/client_golang",
)

go_repository(
    name = "com_github_docker_docker",
    build_extra_args = ["-exclude=vendor"],
    commit = "384c782721c7d0865b7d40ce7ca402022c690058",
    importpath = "github.com/docker/docker",
)

go_repository(
    name = "com_github_jinzhu_inflection",
    commit = "f5c5f50e6090ae76a29240b61ae2a90dd810112e",
    importpath = "github.com/jinzhu/inflection",
)

go_repository(
    name = "com_github_prometheus_client_model",
    commit = "fd36f4220a901265f90734c3183c5f0c91daa0b8",
    importpath = "github.com/prometheus/client_model",
)

go_repository(
    name = "com_github_prometheus_common",
    commit = "31bed53e4047fd6c510e43a941f90cb31be0972a",
    importpath = "github.com/prometheus/common",
)

go_repository(
    name = "com_github_lib_pq",
    commit = "2ff3cb3adc01768e0a552b3a02575a6df38a9bea",
    importpath = "github.com/lib/pq",
)

go_repository(
    name = "com_github_matttproud_golang_protobuf_extensions",
    commit = "c182affec369e30f25d3eb8cd8a478dee585ae7d",
    importpath = "github.com/matttproud/golang_protobuf_extensions",
)

go_repository(
    name = "com_github_prometheus_procfs",
    commit = "f959769dfe3ee5212b9fb905319b8811f12bf7a2",
    importpath = "github.com/prometheus/procfs",
)

go_repository(
    name = "com_github_beorn7_perks",
    commit = "4b2b341e8d7715fae06375aa633dbb6e91b3fb46",
    importpath = "github.com/beorn7/perks",
)

go_repository(
    name = "com_github_gorilla_context",
    commit = "51ce91d2eaddeca0ef29a71d766bb3634dadf729",
    importpath = "github.com/gorilla/context",
)

go_repository(
    name = "com_github_sirupsen_logrus",
    commit = "2a22dbedbad1fd454910cd1f44f210ef90c28464",
    importpath = "github.com/sirupsen/logrus",
)

go_repository(
    name = "com_github_docker_go_connections",
    commit = "fd1b1942c4d55f7f210a8387e612dc6ffee78ff6",
    importpath = "github.com/docker/go-connections",
)

go_repository(
    name = "com_github_pkg_errors",
    commit = "27936f6d90f9c8e1145f11ed52ffffbfdb9e0af7",
    importpath = "github.com/pkg/errors",
)

go_repository(
    name = "com_github_opencontainers_image_spec",
    commit = "da296dcb1e473a9b4e2d148941d7faa9ac8fea3f",
    importpath = "github.com/opencontainers/image-spec",
)

go_repository(
    name = "com_github_nvveen_gotty",
    commit = "cd527374f1e5bff4938207604a14f2e38a9cf512",
    importpath = "github.com/Nvveen/Gotty",
)

go_repository(
    name = "org_golang_x_crypto",
    commit = "ea8f1a30c4438cc8b13f05538385ad8dc6049b43",
    importpath = "golang.org/x/crypto",
)

go_repository(
    name = "com_github_docker_go_units",
    commit = "519db1ee28dcc9fd2474ae59fca29a810482bfb1",
    importpath = "github.com/docker/go-units",
)

go_repository(
    name = "com_github_docker_libtrust",
    commit = "aabc10ec26b754e797f9028f4589c5b7bd90dc20",
    importpath = "github.com/docker/libtrust",
)

go_repository(
    name = "com_github_docker_go_metrics",
    commit = "b84716841b82eab644a0c64fc8b42d480e49add5",
    importpath = "github.com/docker/go-metrics",
)

go_repository(
    name = "com_github_morikuni_aec",
    commit = "39771216ff4c63d11f5e604076f9c45e8be1067b",
    importpath = "github.com/morikuni/aec",
)
