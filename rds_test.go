package main

import (
	"fmt"
	"net/http"
	"testing"

	fake1 "github.com/opentelekomcloud/gophertelekomcloud/openstack/networking/v1/common"
	fake2 "github.com/opentelekomcloud/gophertelekomcloud/openstack/networking/v2/common"
	"github.com/opentelekomcloud/gophertelekomcloud/openstack/networking/v2/extensions/security/groups"
	"github.com/opentelekomcloud/gophertelekomcloud/openstack/networking/v1/subnets"
	"github.com/opentelekomcloud/gophertelekomcloud/openstack/networking/v1/vpcs"
	th "github.com/opentelekomcloud/gophertelekomcloud/testhelper"
)

const SecurityGroupListResponse = `
{
    "security_groups": [
        {
            "description": "default",
            "id": "85cc3048-abc3-43cc-89b3-377341426ac5",
            "name": "default",
            "security_group_rules": [],
            "tenant_id": "e4f50856753b4dc6afee5fa6b9b6c550"
        }
    ]
}
`

const SecurityGroupGetResponse = `
{
    "security_group": {
        "description": "default",
        "id": "85cc3048-abc3-43cc-89b3-377341426ac5",
        "name": "default",
        "security_group_rules": [
            {
                "direction": "egress",
                "ethertype": "IPv6",
                "id": "3c0e45ff-adaf-4124-b083-bf390e5482ff",
                "port_range_max": null,
                "port_range_min": null,
                "protocol": null,
                "remote_group_id": null,
                "remote_ip_prefix": null,
                "security_group_id": "85cc3048-abc3-43cc-89b3-377341426ac5",
                "tenant_id": "e4f50856753b4dc6afee5fa6b9b6c550"
            },
            {
                "direction": "egress",
                "ethertype": "IPv4",
                "id": "93aa42e5-80db-4581-9391-3a608bd0e448",
                "port_range_max": null,
                "port_range_min": null,
                "protocol": null,
                "remote_group_id": null,
                "remote_ip_prefix": null,
                "security_group_id": "85cc3048-abc3-43cc-89b3-377341426ac5",
                "tenant_id": "e4f50856753b4dc6afee5fa6b9b6c550"
            }
        ],
        "tenant_id": "e4f50856753b4dc6afee5fa6b9b6c550"
    }
}`

const SubnetListResponse = `
{
    "subnets": [
        {
            "name": "default",
            "cidr": "172.16.236.0/24",
            "id": "011fc878-5521-4654-a1ad-f5b0b5820302",
            "enable_dhcp": true,
            "network_id": "48efad0c-079d-4cc8-ace0-dce35d584124",
            "tenant_id": "bbfe8c41dd034a07bebd592bf03b4b0c",
            "project_id": "bbfe8c41dd034a07bebd592bf03b4b0c",
            "dns_nameservers": [],
            "allocation_pools": [
                {
                    "start": "172.16.236.2",
                    "end": "172.16.236.251"
                }
            ],
            "host_routes": [],
            "ip_version": 4,
            "gateway_ip": "172.16.236.1",
            "created_at": "2018-03-26T08:23:43",
            "updated_at": "2018-03-26T08:23:44"
        }
    ]
}`

const VpcListResponse = `
{
    "vpcs": [
        {
            "id": "13551d6b-755d-4757-b956-536f674975c0",
            "name": "default",
            "description": "test",
            "cidr": "172.16.0.0/16",
            "status": "OK",
            "enterprise_project_id": "0",
            "routes": [],
            "enable_shared_snat": false
        },
        {
            "id": "3ec3b33f-ac1c-4630-ad1c-7dba1ed79d85",
            "name": "222",
            "description": "test",
            "cidr": "192.168.0.0/16",
            "status": "OK",
            "enterprise_project_id": "0635d733-c12d-4308-ba5a-4dc27ec21038",
            "routes": [],
            "enable_shared_snat": false
        },
        {
            "id": "99d9d709-8478-4b46-9f3f-2206b1023fd3",
            "name": "vpc",
            "description": "test",
            "cidr": "192.168.0.0/16",
            "status": "OK",
            "enterprise_project_id": "0",
            "routes": [],
            "enable_shared_snat": false
        }
    ]
}
`

func Test_secgroupGet(t *testing.T) {

	th.SetupHTTP()
	defer th.TeardownHTTP()

	th.Mux.HandleFunc("/v2.0/security-groups", func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "GET")
		th.TestHeader(t, r, "X-Auth-Token", fake2.TokenID)

		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		_, _ = fmt.Fprint(w, SecurityGroupListResponse)
	})

	sg, err := secgroupGet(fake2.ServiceClient(),  &groups.ListOpts{Name: "default"})
	th.AssertNoErr(t, err)

	th.AssertEquals(t, "default", sg.Name)
	th.AssertEquals(t, "85cc3048-abc3-43cc-89b3-377341426ac5", sg.ID)

}

func Test_subnetGet(t *testing.T) {

	th.SetupHTTP()
	defer th.TeardownHTTP()

	th.Mux.HandleFunc("/v1/85636478b0bd8e67e89469c7749d4127/subnets", func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "GET")
		th.TestHeader(t, r, "X-Auth-Token", fake1.TokenID)

		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		_, _ = fmt.Fprint(w, SubnetListResponse)
	})

	sg, err := subnetGet(fake1.ServiceClient(),  &subnets.ListOpts{Name: "default"})
	th.AssertNoErr(t, err)

	th.AssertEquals(t, "default", sg.Name)
	th.AssertEquals(t, "011fc878-5521-4654-a1ad-f5b0b5820302", sg.ID)

}

func Test_vpcGet(t *testing.T) {

	th.SetupHTTP()
	defer th.TeardownHTTP()

	th.Mux.HandleFunc("/v1/85636478b0bd8e67e89469c7749d4127/vpcs", func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "GET")
		th.TestHeader(t, r, "X-Auth-Token", fake1.TokenID)

		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		_, _ = fmt.Fprint(w, VpcListResponse)
	})

	sg, err := vpcGet(fake1.ServiceClient(),  &vpcs.ListOpts{Name: "default"})
	th.AssertNoErr(t, err)

	th.AssertEquals(t, "default", sg.Name)
	th.AssertEquals(t, "13551d6b-755d-4757-b956-536f674975c0", sg.ID)

}

func Test_main(t *testing.T) {
	tests := []struct {
		name string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			main()
		})
	}
}
