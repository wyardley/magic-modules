package accesscontextmanager_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"

	"github.com/hashicorp/terraform-provider-google/google/acctest"
	"github.com/hashicorp/terraform-provider-google/google/envvar"
	"github.com/hashicorp/terraform-provider-google/google/tpgresource"
	transport_tpg "github.com/hashicorp/terraform-provider-google/google/transport"
)

// Since each test here is acting on the same organization and only one AccessPolicy
// can exist, they need to be run serially. See AccessPolicy for the test runner.

func testAccAccessContextManagerGcpUserAccessBinding_basicTest(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"org_id":        envvar.GetTestOrgFromEnv(t),
		"org_domain":    envvar.GetTestOrgDomainFromEnv(t),
		"cust_id":       envvar.GetTestCustIdFromEnv(t),
		"random_suffix": acctest.RandString(t, 10),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckAccessContextManagerGcpUserAccessBindingDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccAccessContextManagerGcpUserAccessBinding_accessContextManagerGcpUserAccessBindingBasicExample(context),
			},
			{
				ResourceName:            "google_access_context_manager_gcp_user_access_binding.gcp_user_access_binding",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"organization_id"},
			},
			{
				Config: testAccAccessContextManagerGcpUserAccessBinding_accessContextManagerGcpUserAccessBindingNamedExample(context),
			},
			{
				ResourceName:            "google_access_context_manager_gcp_user_access_binding.gcp_user_access_binding",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"organization_id"},
			},
		},
	})
}

func testAccAccessContextManagerGcpUserAccessBinding_accessContextManagerGcpUserAccessBindingBasicExample(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_cloud_identity_group" "group" {
  display_name = "tf-test-my-identity-group%{random_suffix}"

  parent = "customers/%{cust_id}"

  group_key {
    id = "tf-test-my-identity-group%{random_suffix}@%{org_domain}"
  }

  labels = {
    "cloudidentity.googleapis.com/groups.discussion_forum" = ""
  }
}

resource "google_access_context_manager_access_level" "tf_test_access_level_id_for_user_access_binding%{random_suffix}" {
  parent = "accessPolicies/${google_access_context_manager_access_policy.access-policy.name}"
  name   = "accessPolicies/${google_access_context_manager_access_policy.access-policy.name}/accessLevels/tf_test_chromeos_no_lock%{random_suffix}"
  title  = "tf_test_chromeos_no_lock%{random_suffix}"
  basic {
    conditions {
      device_policy {
        require_screen_lock = true
        os_constraints {
          os_type = "DESKTOP_CHROME_OS"
        }
      }
      regions = [
  "US",
      ]
    }
  }
}

resource "google_access_context_manager_access_policy" "access-policy" {
  parent = "organizations/%{org_id}"
  title  = "my policy"
}

resource "google_access_context_manager_gcp_user_access_binding" "gcp_user_access_binding" {
  organization_id = "%{org_id}"
  group_key       = trimprefix(google_cloud_identity_group.group.id, "groups/")
  access_levels   = [
    google_access_context_manager_access_level.tf_test_access_level_id_for_user_access_binding%{random_suffix}.name,
  ]
  session_settings {
    session_length = "1800s"
    session_length_enabled = true
    session_reauth_method = "LOGIN"
    use_oidc_max_age = false
  }
  scoped_access_settings {
    scope {
      client_scope {
        restricted_client_application {
	        client_id = "TEST_APPLICATION"
         }
      }
    }
    active_settings {
  	  access_levels = [
  		  google_access_context_manager_access_level.tf_test_access_level_id_for_user_access_binding%{random_suffix}.name,
  	  ]
  	  session_settings {
  		  session_length = "1800s"
  		  session_length_enabled = true
  		  session_reauth_method = "LOGIN"
  		  use_oidc_max_age = false
  	  }
  	}
  	dry_run_settings {
  	  access_levels = [
  		  google_access_context_manager_access_level.tf_test_access_level_id_for_user_access_binding%{random_suffix}.name,
  	  ]
  	}
  }
}
`, context)
}

func testAccAccessContextManagerGcpUserAccessBinding_accessContextManagerGcpUserAccessBindingNamedExample(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_cloud_identity_group" "group" {
  display_name = "tf-test-my-identity-group%{random_suffix}"

  parent = "customers/%{cust_id}"

  group_key {
    id = "tf-test-my-identity-group%{random_suffix}@%{org_domain}"
  }

  labels = {
    "cloudidentity.googleapis.com/groups.discussion_forum" = ""
  }
}

resource "google_access_context_manager_access_level" "tf_test_access_level_id_for_user_access_binding%{random_suffix}" {
  parent = "accessPolicies/${google_access_context_manager_access_policy.access-policy.name}"
  name   = "accessPolicies/${google_access_context_manager_access_policy.access-policy.name}/accessLevels/tf_test_chromeos_no_lock%{random_suffix}"
  title  = "tf_test_chromeos_no_lock%{random_suffix}"
  basic {
    conditions {
      device_policy {
        require_screen_lock = true
        os_constraints {
          os_type = "DESKTOP_CHROME_OS"
        }
      }
      regions = [
  "US",
      ]
    }
  }
}

resource "google_access_context_manager_access_policy" "access-policy" {
  parent = "organizations/%{org_id}"
  title  = "my policy"
}

resource "google_access_context_manager_gcp_user_access_binding" "gcp_user_access_binding" {
  organization_id = "%{org_id}"
  group_key       = trimprefix(google_cloud_identity_group.group.id, "groups/")
  access_levels   = [
    google_access_context_manager_access_level.tf_test_access_level_id_for_user_access_binding%{random_suffix}.name,
  ]
  session_settings {
    session_length = "1800s"
    session_length_enabled = true
    session_reauth_method = "LOGIN"
    use_oidc_max_age = false
  }
  scoped_access_settings {
    scope {
      client_scope {
        restricted_client_application {
	        name = "Cloud Console"
         }
      }
    }
    active_settings {
  	  access_levels = [
  		  google_access_context_manager_access_level.tf_test_access_level_id_for_user_access_binding%{random_suffix}.name,
  	  ]
  	  session_settings {
  		  max_inactivity = "400s"
  		  session_length = "1800s"
  		  session_length_enabled = true
  		  session_reauth_method = "LOGIN"
  		  use_oidc_max_age = false
  	  }
  	}
  	dry_run_settings {
  	  access_levels = [
  		  google_access_context_manager_access_level.tf_test_access_level_id_for_user_access_binding%{random_suffix}.name,
  	  ]
  	}
  }
}
`, context)
}

func testAccCheckAccessContextManagerGcpUserAccessBindingDestroyProducer(t *testing.T) func(s *terraform.State) error {
	return func(s *terraform.State) error {
		for name, rs := range s.RootModule().Resources {
			if rs.Type != "google_access_context_manager_gcp_user_access_binding" {
				continue
			}
			if strings.HasPrefix(name, "data.") {
				continue
			}

			config := acctest.GoogleProviderConfig(t)

			url, err := tpgresource.ReplaceVarsForTest(config, rs, "{{AccessContextManagerBasePath}}organizations/{{organization_id}}/gcpUserAccessBindings/{{name}}")
			if err != nil {
				return err
			}

			_, err = transport_tpg.SendRequest(transport_tpg.SendRequestOptions{
				Config:    config,
				Method:    "GET",
				RawURL:    url,
				UserAgent: config.UserAgent,
			})
			if err == nil {
				return fmt.Errorf("AccessContextManagerGcpUserAccessBinding still exists at %s", url)
			}
		}

		return nil
	}
}
