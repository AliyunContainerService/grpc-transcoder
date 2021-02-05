proto_path=../hello-servicemesh-grpc/proto
proto_dep_path=~/cooding/github/grpc-gateway/third_party/googleapis

protoc \
    --proto_path=${proto_path} \
    --proto_path=${proto_dep_path} \
    --include_imports \
    --include_source_info \
    --descriptor_set_out=landing.proto-descriptor \
    "${proto_path}"/landing2.proto

make build

./grpc-transcoder \
--version 1.7 \
--service_port 9996 \
--service_name grpc-server-svc \
--proto_pkg org.feuyeux.grpc \
--proto_svc LandingService \
--header x=a,y=b \
--descriptor landing.proto-descriptor
