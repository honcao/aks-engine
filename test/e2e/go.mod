module github.com/Azure/aks-engine/test/e2e

go 1.14

require (
	github.com/Azure/aks-engine v0.54.0
	github.com/Azure/azure-sdk-for-go v43.0.0+incompatible
	github.com/Azure/go-autorest/autorest/to v0.3.0
	github.com/devigned/pub v0.2.0 // indirect
	github.com/go-bindata/go-bindata v3.1.2+incompatible // indirect
	github.com/gogo/protobuf v1.3.1 // indirect
	github.com/google/gofuzz v1.1.0 // indirect
	github.com/influxdata/influxdb v1.7.9
	github.com/kelseyhightower/envconfig v1.4.0
	github.com/mattn/go-isatty v0.0.12 // indirect
	github.com/mitchellh/gox v1.0.1 // indirect
	github.com/onsi/ginkgo v1.12.2
	github.com/onsi/gomega v1.10.1
	github.com/pkg/errors v0.8.1
	golang.org/x/crypto v0.0.0-20191206172530-e9b2fee46413
)

replace github.com/Azure/aks-engine v0.43.0 => ../..
