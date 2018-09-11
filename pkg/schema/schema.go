package schema

type ContainerImage struct {
	// UUID that identifies the container image internally.
	Id string `gorm:"primary_key"`

	// UUID of the repository containing the container image.
	RepositoryId string

	// Digest of the container image. Typically of the form
	// "sha256:...".
	Digest string

	// MIME type of the manifest. Only set if the image is present.
	ManifestMediatype *string

	// Manifest contents of the container image. Only set if the image is
	// present.
	Manifest *[]byte
}

type ContainerRegistry struct {
	// UUID that identifies the container registry internally.
	Id string `gorm:"primary_key"`

	// URI of the container registry (e.g.,
	// "https://index.docker.io/").
	Uri string
}

type ContainerRepository struct {
	// UUID that identifies the container repository internally.
	Id string `gorm:"primary_key"`

	// UUID of the registry containing the repository.
	RegistryId string

	// Name of the repository within the registry (e.g.,
	// "library/mysql", "grafana/grafana").
	RepositoryName string
}

// File holds information of a single-file object that needs to be
// stored by the distfile mirroring service.
type File struct {
	// UUID that identifies the file internally.
	Id string `gorm:"primary_key"`

	// URI of the file.
	Uri string

	// SHA-256 checksum of the file. May be empty if the file has
	// not yet been downloaded.
	Sha256 *string

	// Size of the file in bytes. May be empty if the file has not
	// yet been downloaded.
	Size *uint64

	// Whether the file has already been downloaded successfully.
	Present bool
}
