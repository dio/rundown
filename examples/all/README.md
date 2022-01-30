# Example

This is an example of using [`github.com/dio/rundown/api/ext_authz`](../../api/ext_authz/) and [`github.com/dio/rundown/api/proxy`]((../../api/proxy/)) packages.

To run the auth service only:

```console
go run main.go --external-auth-service-config path/to/configs/ext_authz.json --disable-proxy
```

To run the proxy only:

```console
go run main.go --proxy-config path/to/configs/proxy.yaml --disable-external-auth-service
```

## Config

Please refer to [authservice/docs](../authservice/docs/README.md) to author a valid configuration for the `auth_server`.

The [ext_authz.json](../configs/ext_authz.json) used in this example is taken from https://github.com/dio/authservice/blob/3f884b8d37b0d754751182fd8b67453f3cf0f4b0/bookinfo-example/config/authservice-configmap-template-for-authn.yaml#L14-L48.
