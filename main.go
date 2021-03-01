package main

import (
	"log"

	"github.com/AliyunContainerService/grpc-transcoder/grpc_transcoder"
	"github.com/spf13/cobra"
)

var (
	version            string
	serviceName        string
	servicePort        int
	packages           []string
	services           []string
	descriptorFilePath string
	headers            []string
)

func main() {
	grpcTranscoderEnvoyFilterCmd := &cobra.Command{
		Short: "grpc-transcoder",
		Example: "grpc-transcoder [--service_port 80] [--service_name foo] " +
			"[--proto_pkg acme.example] [--proto_svc 'http.*,echo.*'] [--version 1.8] [--headers x=a,y=b] " +
			"--descriptor /path/to/descriptor",
		RunE: func(cmd *cobra.Command, args []string) error {
			s, err := grpc_transcoder.BuildGrpcTranscoder(descriptorFilePath, packages, services, version, serviceName, servicePort)
			if err != nil {
				return err
			} else {
				err = grpc_transcoder.BuildHeaderToMetadata(headers, version, serviceName, servicePort)
				if err != nil {
					return err
				} else {
					log.Printf("DONE.\nPlease apply the below yaml files:\n%s\n%s",
						grpc_transcoder.TranscodeFile,
						grpc_transcoder.H2MFile)
					log.Printf("EnvoyFilter:%s", *s)
				}
			}
			return nil
		},
	}

	grpcTranscoderEnvoyFilterCmd.PersistentFlags().IntVarP(&servicePort, "service_port", "p", 80, "")
	grpcTranscoderEnvoyFilterCmd.PersistentFlags().StringVarP(&serviceName, "service_name", "s", "grpc-transcoder", "")
	grpcTranscoderEnvoyFilterCmd.PersistentFlags().StringVarP(&version, "version", "v", "1.8", "The version of proxy")
	grpcTranscoderEnvoyFilterCmd.PersistentFlags().StringSliceVar(&packages, "proto_pkg", []string{}, "")
	grpcTranscoderEnvoyFilterCmd.PersistentFlags().StringSliceVar(&services, "proto_svc", []string{}, "")
	grpcTranscoderEnvoyFilterCmd.PersistentFlags().StringSliceVar(&headers, "header", []string{}, "headers(x,y) to metadata(a,b)")
	grpcTranscoderEnvoyFilterCmd.PersistentFlags().StringVarP(&descriptorFilePath, "descriptor", "d", "", "")

	if err := grpcTranscoderEnvoyFilterCmd.Execute(); err != nil {
		log.Fatal(err)
	}
}
