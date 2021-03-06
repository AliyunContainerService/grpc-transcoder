package grpc_transcoder

import "text/template"

var grpcTranscoderTmpl = template.Must(template.New("grpc json transcoder filter").Parse(
	`#Generated by ASM(http://servicemesh.console.aliyun.com)
#GRPC Transcoder EnvoyFilter[{{ .VERSION }}]
apiVersion: networking.istio.io/v1alpha3
kind: EnvoyFilter
metadata:
  name: grpc-transcoder-{{ .FILTER_SERVICE }}
spec:
  workloadSelector:
    labels:
      app: istio-ingressgateway
  configPatches:
    - applyTo: HTTP_FILTER
      match:
        context: GATEWAY
        listener:
          portNumber: {{ .FILTER_PORT }}
          filterChain:
            filter:
              name: "envoy.filters.network.http_connection_manager"
              subFilter:
                name: "envoy.filters.http.router"
        proxy:
          proxyVersion: {{ .RE_VERSION }}
      patch:
        operation: INSERT_BEFORE
        value:
          name: envoy.grpc_json_transcoder
          typed_config:
            '@type': type.googleapis.com/envoy.extensions.filters.http.grpc_json_transcoder.v3.GrpcJsonTranscoder
            proto_descriptor_bin: {{ .FILTER_PB }}
            services: {{ range .PROTO_SERVICE }}
            - {{ . }}{{end}}
            print_options:
              add_whitespace: true
              always_print_primitive_fields: true
              always_print_enums_as_ints: false
              preserve_proto_field_names: false
`))

func GetGrpcTranscoderTmpl() *template.Template {
	return grpcTranscoderTmpl
}

var headerToMetadataTmpl = template.Must(template.New("http header to grpc metadata filter").Parse(
	`#Generated by ASM(http://servicemesh.console.aliyun.com)
#GRPC Transcoder EnvoyFilter[{{ .VERSION }}]
apiVersion: networking.istio.io/v1alpha3
kind: EnvoyFilter
metadata:
  name: header-to-metadata-{{ .FILTER_SERVICE }}
spec:
  workloadSelector:
    labels:
      app: istio-ingressgateway
  configPatches:
    - applyTo: HTTP_FILTER
      match:
        context: GATEWAY
        listener:
          portNumber: {{ .FILTER_PORT }}
          filterChain:
            filter:
              name: "envoy.filters.network.http_connection_manager"
              subFilter:
                name: "envoy.filters.http.router"
        proxy:
          proxyVersion: {{ .RE_VERSION }}
      patch:
        operation: INSERT_FIRST
        value:
          name: envoy.filters.http.header_to_metadata
          typed_config:
            "@type": type.googleapis.com/envoy.extensions.filters.http.header_to_metadata.v3.Config
            request_rules:{{ range $key, $value := .HEADERS }}
              - header: {{ $key }}
                on_header_present:
                  metadata_namespace: envoy.lb
                  key: {{ $value }}
                  type: STRING{{end}}
                remove: false
`))

func GetHeaderToMetadataTmplTmpl() *template.Template {
	return headerToMetadataTmpl
}
