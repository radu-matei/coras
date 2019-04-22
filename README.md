# CORAS
This is an ***experimental*** project for pushing and pulling Cloud Native Application Bundles to OCI registries using the [`oras`](https://github.com/deislabs/oras) project - see https://github.com/radu-matei/coras/issues/3

Do ***NOT*** use it, unless you ***really*** know what you are doing, and are comfortable with breaking changes.

## Building and using:

To build, you need the Go toolchain, `make`, and `dep`:

```
$ make bootstrap
$ make build
```

This will compiled the binary in `bin/`.

```
$ ./bin/coras push testdata/test.json --target cnabregistry.azurecr.io/coras:latest
WARN[0000] encountered unknown type application/vnd.cnab.bundle.thin.v1-wd+json; children may not be fetched
WARN[0000] reference for unknown type: application/vnd.cnab.bundle.thin.v1-wd+json  digest="sha256:4a74a6a6b9e16b63da724718a7622b8db44a1dde0dec0eae0ab12bb072df4090" mediatype=application/vnd.cnab.bundle.thin.v1-wd+json size=1501

$ ./bin/coras pull <output-directory> --target cnabregistry.azurecr.io/coras:latest
WARN[0000] reference for unknown type: application/vnd.cnab.bundle.thin.v1-wd+json  digest="sha256:4a74a6a6b9e16b63da724718a7622b8db44a1dde0dec0eae0ab12bb072df4090" mediatype=application/vnd.cnab.bundle.thin.v1-wd+json size=1501
WARN[0000] unknown type: application/vnd.oci.image.config.v1+json
WARN[0001] encountered unknown type application/vnd.cnab.bundle.thin.v1-wd+json; children may not be fetched
descriptor: {application/vnd.oci.image.manifest.v1+json sha256:9fab77a3b9015346e4bf8c194d9849aeef581ee683471ea24b7928bba77687c6 413 [] map[] <nil>}

 layers: [{application/vnd.cnab.bundle.thin.v1-wd+json sha256:4a74a6a6b9e16b63da724718a7622b8db44a1dde0dec0eae0ab12bb072df4090 1501 [] map[org.opencontainers.image.title:testdata/test.json] <nil>}]
```

> The `<output-directory>` argument passed to the pull operation is where bundle directory will be created. Don't rely on this, as it's just the default way `oras` works, and will be changed here - see https://github.com/radu-matei/coras/issues/5

> Note that this doesn't currently deal with pulling or pushing container images. This will change in the future - see https://github.com/radu-matei/coras/issues/2


The manifest generated in the registry:

```
{
  "schemaVersion": 2,
  "config": {
    "mediaType": "application/vnd.oci.image.config.v1+json",
    "digest": "sha256:44136fa355b3678a1146ad16f7e8649e94fb4fc21fe77e8310c060f61caaff8a",
    "size": 2
  },
  "layers": [
    {
      "mediaType": "application/vnd.cnab.bundle.thin.v1-wd+json",
      "digest": "sha256:4a74a6a6b9e16b63da724718a7622b8db44a1dde0dec0eae0ab12bb072df4090",
      "size": 1501,
      "annotations": {
        "org.opencontainers.image.title": "testdata/test.json"
      }
    }
  ]
}
```