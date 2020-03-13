module github.com/honcao/aks-engine

go 1.12

require (
	github.com/Azure/aks-engine v0.47.0
	github.com/Azure/aks-engine/test/e2e v0.0.0-20200313194836-e30bd50eb644 // indirect
	github.com/Azure/azure-sdk-for-go v36.2.0+incompatible
	github.com/Azure/go-autorest/autorest v0.9.2
	github.com/Azure/go-autorest/autorest/adal v0.8.0
	github.com/Azure/go-autorest/autorest/azure/cli v0.3.0
	github.com/Azure/go-autorest/autorest/date v0.2.0
	github.com/Azure/go-autorest/autorest/to v0.3.0
	github.com/Jeffail/gabs v1.4.0
	github.com/blang/semver v3.5.1+incompatible
	github.com/davecgh/go-spew v1.1.1
	github.com/devigned/pub v0.2.0 // indirect
	github.com/fatih/structs v1.1.0
	github.com/ghodss/yaml v1.0.0
	github.com/go-bindata/go-bindata v3.1.2+incompatible // indirect
	github.com/gogo/protobuf v1.3.1 // indirect
	github.com/google/go-cmp v0.3.0
	github.com/google/gofuzz v1.1.0 // indirect
	github.com/google/uuid v1.1.1
	github.com/imdario/mergo v0.3.6
	github.com/jarcoal/httpmock v1.0.1
	github.com/leonelquinteros/gotext v1.4.0
	github.com/mattn/go-colorable v0.0.9
	github.com/mattn/go-isatty v0.0.12 // indirect
	github.com/mitchellh/go-homedir v1.1.0
	github.com/mitchellh/gox v1.0.1 // indirect
	github.com/onsi/ginkgo v1.10.3
	github.com/onsi/gomega v1.4.3
	github.com/pkg/errors v0.8.1
	github.com/sirupsen/logrus v1.4.2
	github.com/spf13/cobra v0.0.5
	github.com/spf13/pflag v1.0.5
	github.com/x-cray/logrus-prefixed-formatter v0.5.2
	golang.org/x/crypto v0.0.0-20191105034135-c7e5f84aec59
	golang.org/x/sync v0.0.0-20190423024810-112230192c58
	gopkg.in/go-playground/validator.v9 v9.25.0
	gopkg.in/ini.v1 v1.41.0
	k8s.io/api v0.0.0-20190222213804-5cb15d344471
	k8s.io/apimachinery v0.0.0-20190221213512-86fb29eff628
	k8s.io/client-go v10.0.0+incompatible
)
