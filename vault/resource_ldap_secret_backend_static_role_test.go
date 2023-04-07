package vault

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-provider-vault/internal/consts"
	"github.com/hashicorp/terraform-provider-vault/testutil"
	"testing"
)

func TestAccLDAPSecretBackendStaticRole(t *testing.T) {
	backend := acctest.RandomWithPrefix("tf-test-ldap-static-role")
	bindDN, bindPass, url := testutil.GetTestLDAPCreds(t)

	resourceType := "vault_ldap_secret_backend_static_role"
	resourceName := resourceType + ".role"
	resource.Test(t, resource.TestCase{
		ProviderFactories: providerFactories,
		PreCheck:          func() { testutil.TestAccPreCheck(t) },
		CheckDestroy:      testCheckMountDestroyed(resourceType, consts.MountTypeLDAP, consts.FieldBackend),

		Steps: []resource.TestStep{
			{
				Config: testLDAPSecretBackendStaticRoleConfig(backend, bindDN, bindPass, url, "alice", "cn=alice,ou=users,dc=example,dc=org", "alice", 60),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "dn", "cn=alice,ou=users,dc=example,dc=org"),
					resource.TestCheckResourceAttr(resourceName, "username", "alice"),
					resource.TestCheckResourceAttr(resourceName, "rotation_period", "60"),
				),
			},
			{
				Config: testLDAPSecretBackendStaticRoleConfig(backend, bindDN, bindPass, url, "bob", "cn=bob,ou=users,dc=example,dc=org", "bob", 120),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "dn", "cn=bob,ou=users,dc=example,dc=org"),
					resource.TestCheckResourceAttr(resourceName, "username", "bob"),
					resource.TestCheckResourceAttr(resourceName, "rotation_period", "120"),
				),
			},
			testutil.GetImportTestStep(resourceName, false, nil, "backend", "role", "disable_remount"),
		},
	})
}

func testLDAPSecretBackendStaticRoleConfig(backend, bindDN, bindPass, url, username, dn, role string, rotationPeriod int) string {
	return fmt.Sprintf(`
resource "vault_ldap_secret_backend" "test" {
  backend                   = "%s"
  description               = "test description"
  default_lease_ttl_seconds = "3600"
  max_lease_ttl_seconds     = "7200"
  binddn                    = "%s"
  bindpass                  = "%s"
  url                       = "%s"
  userdn                    = "CN=Users,DC=corp,DC=example,DC=net"
}

resource "vault_ldap_secret_backend_static_role" "role" {
    backend = vault_ldap_secret_backend.test.backend
    username = "%s"
    dn = "%s"
    role = "%s"
    rotation_period = %d
}
`, backend, bindDN, bindPass, url, username, dn, role, rotationPeriod)
}
