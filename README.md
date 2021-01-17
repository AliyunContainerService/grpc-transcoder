# grpc-transcoder

- support istio1.6+

```bash
# https://github.com/AliyunContainerService/hello-servicemesh-grpc
proto_path={path/to/hello-servicemesh-grpc}/proto
# https://github.com/grpc-ecosystem/grpc-gateway/tree/master/third_party/
proto_dep_path={path/to/third_party}
protoc \
    --proto_path=${proto_path} \
    --proto_path=${proto_dep_path} \
    --include_imports \
    --include_source_info \
    --descriptor_set_out=landing.proto-descriptor \
    "${proto_path}"/landing.proto
```