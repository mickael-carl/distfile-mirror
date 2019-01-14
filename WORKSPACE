load("@bazel_tools//tools/build_defs/repo:http.bzl", "http_archive")

http_archive(
    name = "io_bazel_rules_docker",
    sha256 = "29d109605e0d6f9c892584f07275b8c9260803bf0c6fcb7de2623b2bedc910bd",
    strip_prefix = "rules_docker-0.5.1",
    urls = ["https://github.com/bazelbuild/rules_docker/archive/v0.5.1.tar.gz"],
)

http_archive(
    name = "io_bazel_rules_go",
    sha256 = "7be7dc01f1e0afdba6c8eb2b43d2fa01c743be1b9273ab1eaf6c233df078d705",
    urls = ["https://github.com/bazelbuild/rules_go/releases/download/0.16.5/rules_go-0.16.5.tar.gz"],
)

http_archive(
    name = "bazel_gazelle",
    sha256 = "7949fc6cc17b5b191103e97481cf8889217263acf52e00b560683413af204fcb",
    urls = ["https://github.com/bazelbuild/bazel-gazelle/releases/download/0.16.0/bazel-gazelle-0.16.0.tar.gz"],
)

http_archive(
    name = "bazel_skylib",
    sha256 = "eb5c57e4c12e68c0c20bc774bfbc60a568e800d025557bc4ea022c6479acc867",
    strip_prefix = "bazel-skylib-0.6.0",
    urls = ["https://github.com/bazelbuild/bazel-skylib/archive/0.6.0.tar.gz"],
)

load(
    "@io_bazel_rules_docker//container:container.bzl",
    container_repositories = "repositories",
)

container_repositories()

load("@io_bazel_rules_go//go:def.bzl", "go_register_toolchains", "go_rules_dependencies")

go_rules_dependencies()

go_register_toolchains()

load("@bazel_gazelle//:deps.bzl", "gazelle_dependencies", "go_repository")

gazelle_dependencies()

go_repository(
    name = "com_github_gorilla_mux",
    commit = "08e7f807d38d6a870193019bb439056118661505",
    importpath = "github.com/gorilla/mux",
)

go_repository(
    name = "com_github_docker_distribution",
    commit = "91b0f0559eb6710c574451049f8c71eb9c3952a7",
    importpath = "github.com/docker/distribution",
    build_extra_args = ["-exclude=vendor"],
)

go_repository(
    name = "com_github_jinzhu_gorm",
    commit = "9f1a7f53511168c0567b4b4b4f10ab7d21265174",
    importpath = "github.com/jinzhu/gorm",
)

go_repository(
    name = "com_github_aws_aws_sdk_go",
    commit = "9b56194bccedb23d7e753822b44f6b2e69487ff8",
    importpath = "github.com/aws/aws-sdk-go",
)

go_repository(
    name = "com_github_opencontainers_go_digest",
    commit = "c9281466c8b2f606084ac71339773efd177436e7",
    importpath = "github.com/opencontainers/go-digest",
)

go_repository(
    name = "com_github_prometheus_client_golang",
    commit = "d2ead25884778582e740573999f7b07f47e171b4",
    importpath = "github.com/prometheus/client_golang",
)

go_repository(
    name = "com_github_docker_docker",
    commit = "ebc0750e9fa657ebf5ea08fdeae3242144e784aa",
    importpath = "github.com/docker/docker",
    build_extra_args = ["-exclude=vendor"],
)

go_repository(
    name = "com_github_jinzhu_inflection",
    commit = "04140366298a54a039076d798123ffa108fff46c",
    importpath = "github.com/jinzhu/inflection",
)

go_repository(
    name = "com_github_prometheus_client_model",
    commit = "f287a105a20ec685d797f65cd0ce8fbeaef42da1",
    importpath = "github.com/prometheus/client_model",
)

go_repository(
    name = "com_github_prometheus_common",
    commit = "2998b132700a7d019ff618c06a234b47c1f3f681",
    importpath = "github.com/prometheus/common",
)

go_repository(
    name = "com_github_lib_pq",
    commit = "9eb73efc1fcc404148b56765b0d3f61d9a5ef8ee",
    importpath = "github.com/lib/pq",
)

go_repository(
    name = "com_github_matttproud_golang_protobuf_extensions",
    commit = "c182affec369e30f25d3eb8cd8a478dee585ae7d",
    importpath = "github.com/matttproud/golang_protobuf_extensions",
)

go_repository(
    name = "com_github_prometheus_procfs",
    commit = "b1a0a9a36d7453ba0f62578b99712f3a6c5f82d1",
    importpath = "github.com/prometheus/procfs",
)

go_repository(
    name = "com_github_beorn7_perks",
    commit = "3a771d992973f24aa725d07868b467d1ddfceafb",
    importpath = "github.com/beorn7/perks",
)

go_repository(
    name = "com_github_gorilla_context",
    commit = "51ce91d2eaddeca0ef29a71d766bb3634dadf729",
    importpath = "github.com/gorilla/context",
)

go_repository(
    name = "com_github_sirupsen_logrus",
    commit = "78fb3852d92683dc28da6cc3d5f965100677c27d",
    importpath = "github.com/sirupsen/logrus",
)

go_repository(
    name = "com_github_docker_go_connections",
    commit = "97c2040d34dfae1d1b1275fa3a78dbdd2f41cf7e",
    importpath = "github.com/docker/go-connections",
)

go_repository(
    name = "com_github_pkg_errors",
    commit = "ffb6e22f01932bf7ac35e0bad9be11f01d1c8685",
    importpath = "github.com/pkg/errors",
)

go_repository(
    name = "com_github_opencontainers_image_spec",
    commit = "e74bb5aed1c876be723e181895359d8dca9b56a6",
    importpath = "github.com/opencontainers/image-spec",
)

go_repository(
    name = "com_github_nvveen_gotty",
    commit = "cd527374f1e5bff4938207604a14f2e38a9cf512",
    importpath = "github.com/Nvveen/Gotty",
)

go_repository(
    name = "org_golang_x_crypto",
    commit = "ff983b9c42bc9fbf91556e191cc8efb585c16908",
    importpath = "golang.org/x/crypto",
)

go_repository(
    name = "com_github_docker_go_units",
    commit = "2fb04c6466a548a03cb009c5569ee1ab1e35398e",
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
