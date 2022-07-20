package provider

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func execGit(t *testing.T, arg ...string) string {
	t.Helper()

	output, err := exec.Command("git", arg...).Output()
	if err != nil {
		t.Fatal(err)
	}
	return strings.TrimSpace(string(output))
}

func TestDataSourceRepository(t *testing.T) {
	dir, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	dir = filepath.Join(dir, "../..")
	dir = filepath.ToSlash(dir)

	expectedBranch := execGit(t, "rev-parse", "--abbrev-ref", "HEAD")
	expectedCommit := execGit(t, "rev-parse", "HEAD")
	expectedMessage := execGit(t, "show", "-s", "--format='%s'")

	resource.UnitTest(t, resource.TestCase{
		ProviderFactories: map[string]func() (*schema.Provider, error){
			"git": func() (*schema.Provider, error) {
				return New(), nil
			},
		},
		Steps: []resource.TestStep{
			{
				Config: testDataSourceRepositoryConfig(dir),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.git_repository.test", "branch", strings.TrimSpace(string(expectedBranch))),
					resource.TestCheckResourceAttr("data.git_repository.test", "commit_hash", strings.TrimSpace(string(expectedCommit))),
					resource.TestCheckResourceAttr("data.git_repository.test", "commit_message", strings.TrimSpace(string(expectedMessage))),
				),
			},
		},
	})
}

func testDataSourceRepositoryConfig(path string) string {
	return fmt.Sprintf(`
data "git_repository" "test" {
	path = "%s"
}
`, path)
}
