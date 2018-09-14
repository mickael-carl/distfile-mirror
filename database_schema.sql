CREATE TABLE container_registries (
	id UUID NOT NULL DEFAULT gen_random_uuid(),
	uri STRING NOT NULL,
	CONSTRAINT "primary" PRIMARY KEY (id ASC),
	UNIQUE INDEX container_registries_uri_key (uri ASC),
	FAMILY "primary" (id, uri)
);

CREATE TABLE container_repositories (
	id UUID NOT NULL DEFAULT gen_random_uuid(),
	registry_id UUID NOT NULL,
	repository_name STRING NOT NULL,
	CONSTRAINT "primary" PRIMARY KEY (id ASC),
	CONSTRAINT fk_registry_id_ref_container_registries FOREIGN KEY (registry_id) REFERENCES container_registries (id),
	UNIQUE INDEX container_repositories_registry_id_repository_name_key (registry_id ASC, repository_name ASC),
	FAMILY "primary" (id, registry_id, repository_name)
);

CREATE TABLE container_images (
	id UUID NOT NULL DEFAULT gen_random_uuid(),
	repository_id UUID NOT NULL,
	digest STRING NOT NULL,
	manifest_mediatype STRING NULL,
	manifest BYTES NULL,
	CONSTRAINT "primary" PRIMARY KEY (id ASC),
	CONSTRAINT fk_repository_id_ref_container_repositories FOREIGN KEY (repository_id) REFERENCES container_repositories (id),
	UNIQUE INDEX container_images_repository_id_digest_key (repository_id ASC, digest ASC),
	FAMILY "primary" (id, repository_id, digest, manifest_mediatype, manifest),
	CONSTRAINT check_manifest_manifest_mediatype CHECK ((manifest IS NULL) = (manifest_mediatype IS NULL))
);

CREATE TABLE files (
	id UUID NOT NULL DEFAULT gen_random_uuid(),
	uri STRING NOT NULL,
	sha256 STRING NULL,
	size INTEGER NULL,
	present BOOL NOT NULL DEFAULT false,
	CONSTRAINT "primary" PRIMARY KEY (id ASC),
	UNIQUE INDEX files_uri_key (uri ASC),
	FAMILY "primary" (id, uri, sha256, size, present),
	CONSTRAINT check_sha256 CHECK (sha256 ~ '^[0-9a-f]{64}$'),
	CONSTRAINT check_present_sha256 CHECK ((NOT present) OR (sha256 IS NOT NULL)),
	CONSTRAINT check_present_size CHECK ((NOT present) OR (size IS NOT NULL))
);
