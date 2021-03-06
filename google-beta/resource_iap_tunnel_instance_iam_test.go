package google

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccIapTunnelInstanceIamBinding(t *testing.T) {
	t.Parallel()

	project := getTestProjectFromEnv()
	zone := getTestZoneFromEnv()
	instanceName := fmt.Sprintf("tf-test-instance-%s", acctest.RandString(10))

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccIapTunnelInstanceIamBinding_basic(zone, instanceName),
			},
			{
				ResourceName:      "google_iap_tunnel_instance_iam_binding.foo",
				ImportStateId:     fmt.Sprintf("%s/%s/%s roles/iap.tunnelResourceAccessor", project, zone, instanceName),
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				// Test Iam Binding update
				Config: testAccIapTunnelInstanceIamBinding_update(zone, instanceName),
			},
			{
				ResourceName:      "google_iap_tunnel_instance_iam_binding.foo",
				ImportStateId:     fmt.Sprintf("%s/%s/%s roles/iap.tunnelResourceAccessor", project, zone, instanceName),
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccIapTunnelInstanceIamMember(t *testing.T) {
	t.Parallel()

	project := getTestProjectFromEnv()
	zone := getTestZoneFromEnv()
	instanceName := fmt.Sprintf("tf-test-instance-%s", acctest.RandString(10))

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				// Test Iam Member creation (no update for member, no need to test)
				Config: testAccIapTunnelInstanceIamMember_basic(zone, instanceName),
			},
			{
				ResourceName:      "google_iap_tunnel_instance_iam_member.foo",
				ImportStateId:     fmt.Sprintf("%s/%s/%s roles/iap.tunnelResourceAccessor user:admin@hashicorptest.com", project, zone, instanceName),
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccIapTunnelInstanceIamPolicy(t *testing.T) {
	t.Parallel()

	project := getTestProjectFromEnv()
	zone := getTestZoneFromEnv()
	instanceName := fmt.Sprintf("tf-test-instance-%s", acctest.RandString(10))

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccIapTunnelInstanceIamPolicy_basic(zone, instanceName),
			},
			// Test a few import formats
			{
				ResourceName:      "google_iap_tunnel_instance_iam_policy.foo",
				ImportStateId:     fmt.Sprintf("projects/%s/zones/%s/instances/%s", project, zone, instanceName),
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				ResourceName:      "google_iap_tunnel_instance_iam_policy.foo",
				ImportStateId:     fmt.Sprintf("%s/%s/%s", project, zone, instanceName),
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				ResourceName:      "google_iap_tunnel_instance_iam_policy.foo",
				ImportStateId:     fmt.Sprintf("%s/%s", zone, instanceName),
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccIapTunnelInstanceIamBinding_basic(zone, instanceName string) string {
	return fmt.Sprintf(`
resource "google_compute_instance" "test_vm" {
  zone         = "%s"
  name         = "%s"
  machine_type = "n1-standard-1"
  boot_disk {
    initialize_params {
      image = "debian-cloud/debian-9"
    }
  }
  network_interface {
    network = "default"
  }
}

resource "google_iap_tunnel_instance_iam_binding" "foo" {
  project = "${google_compute_instance.test_vm.project}"
  zone = "${google_compute_instance.test_vm.zone}"
  instance = "${google_compute_instance.test_vm.name}"
  role        = "roles/iap.tunnelResourceAccessor"
  members      = ["user:admin@hashicorptest.com"]
}
`, zone, instanceName)
}

func testAccIapTunnelInstanceIamBinding_update(zone, instanceName string) string {
	return fmt.Sprintf(`
resource "google_compute_instance" "test_vm" {
  zone         = "%s"
  name         = "%s"
  machine_type = "n1-standard-1"
  boot_disk {
    initialize_params {
      image = "debian-cloud/debian-9"
    }
  }
  network_interface {
    network = "default"
  }
}

resource "google_iap_tunnel_instance_iam_binding" "foo" {
  project = "${google_compute_instance.test_vm.project}"
  zone = "${google_compute_instance.test_vm.zone}"
  instance = "${google_compute_instance.test_vm.name}"
  role        = "roles/iap.tunnelResourceAccessor"
  members      = [
    "user:admin@hashicorptest.com",
    "user:paddy@hashicorp.com"
  ]
}
`, zone, instanceName)
}

func testAccIapTunnelInstanceIamMember_basic(zone, instanceName string) string {
	return fmt.Sprintf(`
resource "google_compute_instance" "test_vm" {
  zone         = "%s"
  name         = "%s"
  machine_type = "n1-standard-1"
  boot_disk {
    initialize_params {
      image = "debian-cloud/debian-9"
    }
  }
  network_interface {
    network = "default"
  }
}

resource "google_iap_tunnel_instance_iam_member" "foo" {
  project = "${google_compute_instance.test_vm.project}"
  zone = "${google_compute_instance.test_vm.zone}"
  instance = "${google_compute_instance.test_vm.name}"
  role        = "roles/iap.tunnelResourceAccessor"
  member      = "user:admin@hashicorptest.com"
}
`, zone, instanceName)
}

func testAccIapTunnelInstanceIamPolicy_basic(zone, instanceName string) string {
	return fmt.Sprintf(`
resource "google_compute_instance" "test_vm" {
  zone         = "%s"
  name         = "%s"
  machine_type = "n1-standard-1"
  boot_disk {
    initialize_params {
      image = "debian-cloud/debian-9"
    }
  }
  network_interface {
    network = "default"
  }
}

data "google_iam_policy" "foo" {
  binding {
    role = "roles/iap.tunnelResourceAccessor"
    members = ["user:admin@hashicorptest.com"]
  }
}

resource "google_iap_tunnel_instance_iam_policy" "foo" {
  project = "${google_compute_instance.test_vm.project}"
  zone = "${google_compute_instance.test_vm.zone}"
  instance = "${google_compute_instance.test_vm.name}"
  policy_data = "${data.google_iam_policy.foo.policy_data}"
}
`, zone, instanceName)
}
