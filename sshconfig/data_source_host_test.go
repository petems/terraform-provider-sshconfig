package sshconfig

import (
	"fmt"
	"io/ioutil"
	"regexp"
	"strings"
	"testing"

	"github.com/andreyvit/diff"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

const testDataSourceUnknownKeyError = `
data "sshconfig_host" "fail_compilation_unknown_key" {
  host              = "example.com"
  this_doesnt_exist = "foo"
}
`

const testDataSourceValidSSHConfig = `
data "sshconfig_host" "valid_config" {
  host = "example.com"
  path = "%s"
}

output "host_render" {
  value = "${data.sshconfig_host.valid_config.rendered}"
}

output "host_map" {
  value = "${data.sshconfig_host.valid_config.host_map}"
}
`

const exampleHostFileContents = `
Host example.com
  ControlMaster auto
  ControlPersist yes
  User git

Host vaultstack-servers-0.centralus.cloudapp.azure.com
  StrictHostKeyChecking no
`

const justExampleHost = `
Host example.com
  ControlMaster auto
  ControlPersist yes
  User git
`

func createTmpHostFile(t *testing.T, hostfileContents string) string {
	tmpfile, err := ioutil.TempFile("", "sshconfigfile")
	if err != nil {
		t.Errorf("Unable to open a temporary file.")
	}

	if _, err := tmpfile.Write([]byte(exampleHostFileContents)); err != nil {
		t.Errorf("Unable to write to temporary file %q: %v", tmpfile.Name(), err)
	}
	if err := tmpfile.Close(); err != nil {
		t.Errorf("Unable to close temporary file %q: %v", tmpfile.Name(), err)
	}

	return tmpfile.Name()
}

func TestDataSource_compileUnknownKeyError(t *testing.T) {
	resource.UnitTest(t, resource.TestCase{
		Providers: testProviders,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config:      testDataSourceUnknownKeyError,
				ExpectError: regexp.MustCompile("An argument named \"this_doesnt_exist\" is not expected here."),
			},
		},
	})
}

func TestDataSource_validConfigFileRender(t *testing.T) {

	hostFilePath := createTmpHostFile(t, exampleHostFileContents)

	resource.UnitTest(t, resource.TestCase{
		Providers: testProviders,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: fmt.Sprintf(testDataSourceValidSSHConfig, hostFilePath),
				Check: func(s *terraform.State) error {
					_, ok := s.RootModule().Resources["data.sshconfig_host.valid_config"]
					if !ok {
						return fmt.Errorf("missing data resource")
					}

					outputs := s.RootModule().Outputs
					expected := justExampleHost
					actual := outputs["host_render"].Value.(string)

					if a, e := strings.TrimSpace(actual), strings.TrimSpace(expected); a != e {
						t.Errorf("Result not as expected:\n%v", diff.LineDiff(e, a))
					}

					return nil
				},
			},
		},
	})
}

func TestDataSource_validConfigFileConfigMap(t *testing.T) {

	hostFilePath := createTmpHostFile(t, exampleHostFileContents)

	resource.UnitTest(t, resource.TestCase{
		Providers: testProviders,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: fmt.Sprintf(testDataSourceValidSSHConfig, hostFilePath),
				Check: func(s *terraform.State) error {
					resourceState := s.Modules[0].Resources["data.sshconfig_host.valid_config"]
					if resourceState == nil {
						return fmt.Errorf("resource not found in state %v", s.Modules[0].Resources)
					}

					iState := resourceState.Primary
					if iState == nil {
						return fmt.Errorf("resource has no primary instance")
					}

					if got, want := iState.Attributes["host_map.ControlMaster"], "auto"; got != want {
						return fmt.Errorf("host_map.ControlMaster contains %s; want %s", got, want)
					}

					return nil
				},
			},
		},
	})
}
