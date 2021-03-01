package grpc_transcoder

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"regexp"
	"sort"
	"strings"

	"github.com/golang/protobuf/proto"
	"github.com/golang/protobuf/protoc-gen-go/descriptor"
	"github.com/hashicorp/go-multierror"
	"github.com/pkg/errors"
)

const (
	TranscodeFile = "grpc-transcoder-envoyfilter.yaml"
	H2MFile       = "header2metadata-envoyfilter.yaml"
)

var (
	V16 = `^1\.6.*`
	V17 = `^1\.7.*`
	V18 = `^1\.8.*`
)

func BuildHeaderToMetadata(headers []string, version string, serviceName string, servicePort int) error {
	maps := make(map[string]string)
	for _, kv := range headers {
		kvs := strings.Split(kv, "=")
		maps[kvs[0]] = kvs[1]
	}
	params := map[string]interface{}{
		"VERSION":        version,
		"RE_VERSION":     getReVersion(version),
		"FILTER_SERVICE": serviceName,
		"FILTER_PORT":    servicePort,
		"HEADERS":        maps,
	}
	f, err := os.Create(H2MFile)
	if err != nil {
		log.Fatal(err)
	}

	return GetHeaderToMetadataTmplTmpl().Execute(f, params)
}

func BuildGrpcTranscoder(descriptorFilePath string, packages []string, services []string, version string, serviceName string, servicePort int) (*string, error) {
	if _, err := os.Stat(descriptorFilePath); os.IsNotExist(err) {
		log.Printf("error opening descriptor file %q\n", descriptorFilePath)
		return nil, err
	}

	descriptorBytes, err := ioutil.ReadFile(descriptorFilePath)
	if err != nil {
		log.Printf("error reading descriptor file %q\n", descriptorFilePath)
		return nil, err
	}

	return buildGrpcTranscoder(descriptorBytes, packages, services, version, serviceName, servicePort)
}

func buildGrpcTranscoder(descriptorBytes []byte, packages []string, services []string, version string, serviceName string, servicePort int) (*string, error) {
	protoServices, err := getServices(&descriptorBytes, packages, services)
	if err != nil {
		log.Printf("error extracting services from descriptor: %v\n", err)
		return nil, err
	}
	sort.Strings(protoServices)

	descriptorBinary := base64.StdEncoding.EncodeToString(descriptorBytes)

	params := map[string]interface{}{
		"VERSION":        version,
		"RE_VERSION":     getReVersion(version),
		"FILTER_SERVICE": serviceName,
		"FILTER_PORT":    servicePort,
		"FILTER_PB":      descriptorBinary,
		"PROTO_SERVICE":  protoServices,
	}

	buf := new(bytes.Buffer)
	err = GetGrpcTranscoderTmpl().Execute(buf, params)
	s := buf.String()
	return &s, err
}

func getReVersion(version string) string {
	switch version {
	case "1.6":
		return V16
	case "1.7":
		return V17
	case "1.8":
		return V18
	default:
		return V17
	}
}
func getServices(b *[]byte, packages []string, services []string) ([]string, error) {
	var (
		fds  descriptor.FileDescriptorSet
		out  []string
		rexp []*regexp.Regexp
		errs error
	)
	if err := proto.Unmarshal(*b, &fds); err != nil {
		return out, errors.Wrapf(err, "error proto unmarshall to FileDescriptorSet")
	}
	rexp = make([]*regexp.Regexp, 0)
	for _, r := range services {
		re, err := regexp.Compile(r)
		if err != nil {
			errs = multierror.Append(errs, err)
		} else {
			rexp = append(rexp, re)
		}
	}

	// package
	findPkg := func(name string) bool {
		for _, p := range packages {
			if strings.HasPrefix(name, p) {
				return true
			}
		}
		return len(packages) == 0
	}

	// service
	findSvc := func(s string) bool {
		for _, r := range rexp {
			if r.MatchString(s) {
				return true
			}
		}
		return len(rexp) == 0
	}

	for _, f := range fds.GetFile() {
		if !findPkg(f.GetPackage()) {
			continue
		}
		for _, s := range f.GetService() {
			if findSvc(s.GetName()) {
				out = append(out, fmt.Sprintf("%s.%s", f.GetPackage(), s.GetName()))
			}
		}
	}
	return out, errs
}
