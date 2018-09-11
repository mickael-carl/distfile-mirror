# Distfile Mirror

Distfile Mirror is a set of applications that make it easier to mirror
artifacts from the Internet within an organisation. These artifacts may
be downloaded from the mirror service through a HTTP proxy that spoofs
SSL certificates on demand. This allows these original artifacts to be
fetched under their original URLs, under the condition that the user
installs a custom CA certificate.

Unlike a general purpose caching proxy, this system does not download
artifacts from the Internet by simply fetching them on demand. Instead,
a web UI is offered that can be used to explicitly add URIs of artifacts
that need to be present. The reason for this is twofold:

- Once an artifact has been declared through the web UI and is
  downloaded, it will never be purged from storage. There is no size
  limited caching policy. This makes Distfile Mirror a useful tool for
  being able to reliably do builds of software whose build processes
  download files (e.g., [Bazel](https://bazel.build/),
  [Buildroot](https://buildroot.org/)). Even if the upstream file is
  deleted or modified.

- By requiring users to add artifacts manually, the occasion could be
  used to document what the artifact is being used for. Especially in
  larger organisations it is useful to have some bookkeeping of the use
  of open-source software.

Right now the service is capable of storing two types of resources:

- Files, downloaded over HTTP or HTTPS. Files are identified by URI.

- Docker container images. Container images are identified by registry
  URI, repository name and image digest (SHA-256). Mirroring by tag is
  not supported, as experience has shown that suppliers of container
  images often overwrite tags to point to newer versions of an image.
  This is bad for reproducibility of work.

TODO(edsch): Continue README.
