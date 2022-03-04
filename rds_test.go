package main

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/opentelekomcloud/gophertelekomcloud/acceptance/tools"
	fake "github.com/opentelekomcloud/gophertelekomcloud/testhelper/client"
	"github.com/opentelekomcloud/gophertelekomcloud/openstack/networking/v2/extensions/security/groups"
	"github.com/opentelekomcloud/gophertelekomcloud/openstack/networking/v1/subnets"
	"github.com/opentelekomcloud/gophertelekomcloud/openstack/networking/v1/vpcs"
	"github.com/opentelekomcloud/gophertelekomcloud/openstack/rds/v3/instances"
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
            "name": "golang",
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
            "name": "golang",
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
const RdsGetResponse = `
{
	"instances": [{
		"id": "ed7cc6166ec24360a5ed5c5c9c2ed726in01",
		"status": "ACTIVE",
		"name": "default",
		"port": 3306,
		"type": "Single",
		"region": "eu-de",
		"datastore": {
			"type": "MySQL",
			"version": "5.7"
		},
		"created": "2018-08-20T02:33:49+0800",
		"updated": "2018-08-20T02:33:50+0800",
		"volume": {
			"type": "ULTRAHIGH",
			"size": 100
		},
		"nodes": [{
			"id": "06f1c2ad57604ae89e153e4d27f4e4b8no01",
			"name": "mysql-0820-022709-01_node0",
			"role": "master",
			"status": "ACTIVE",
			"availability_zone": "eu-de-01"
		}],
		"private_ips": ["192.168.0.142"],
		"public_ips": ["10.154.219.187", "10.154.219.186"],
		"db_user_name": "root",
		"vpc_id": "b21630c1-e7d3-450d-907d-39ef5f445ae7",
		"subnet_id": "45557a98-9e17-4600-8aec-999150bc4eef",
		"security_group_id": "38815c5c-482b-450a-80b6-0a301f2afd97",
		"flavor_ref": "rds.mysql.s1.large",
		"switch_strategy": "",
		"backup_strategy": {
			"start_time": "19:00-20:00",
			"keep_days": 7
		},
		"maintenance_window": "02:00-06:00",
		"related_instance": [],
		"disk_encryption_id": "",
		"time_zone": ""
	}, {
		"id": "ed7cc6166ec24360a5ed5c5c9c2ed726in02",
		"status": "ACTIVE",
		"name": "mysql-0820-022709-02",
		"port": 3306,
		"type": "Single",
		"region": "eu-de",
		"datastore": {
			"type": "MySQL",
			"version": "5.6"
		},
		"created": "2019-08-20T02:33:49+0800",
		"updated": "2019-08-20T02:33:50+0800",
		"volume": {
			"type": "ULTRAHIGH",
			"size": 100
		},
		"nodes": [{
			"id": "06f1c2ad57604ae89e153e4d27f4e4b8no01",
			"name": "mysql-0820-022709-01_node0",
			"role": "master",
			"status": "ACTIVE",
			"availability_zone": "eu-de-01"
		}],
		"private_ips": ["192.168.0.142"],
		"public_ips": ["10.154.219.187", "10.154.219.186"],
		"db_user_name": "root",
		"vpc_id": "b21630c1-e7d3-450d-907d-39ef5f445ae7",
		"subnet_id": "45557a98-9e17-4600-8aec-999150bc4eef",
		"security_group_id": "38815c5c-482b-450a-80b6-0a301f2afd97",
		"flavor_ref": "rds.mysql.s1.large",
		"switch_strategy": "",
		"backup_strategy": {
			"start_time": "19:00-20:00",
			"keep_days": 7
		},
		"maintenance_window": "02:00-06:00",
		"related_instance": [],
		"disk_encryption_id": "",
		"time_zone": ""
	}],
	"total_count": 2
}
`

const RdsGetResponseSingle = `
{
	"instances": [{
		"id": "ed7cc6166ec24360a5ed5c5c9c2ed726in01",
		"status": "ACTIVE",
		"name": "default",
		"port": 3306,
		"type": "Single",
		"region": "eu-de",
		"datastore": {
			"type": "MySQL",
			"version": "5.7"
		},
		"created": "2018-08-20T02:33:49+0800",
		"updated": "2018-08-20T02:33:50+0800",
		"volume": {
			"type": "ULTRAHIGH",
			"size": 100
		},
		"nodes": [{
			"id": "06f1c2ad57604ae89e153e4d27f4e4b8no01",
			"name": "mysql-0820-022709-01_node0",
			"role": "master",
			"status": "ACTIVE",
			"availability_zone": "eu-de-01"
		}],
		"private_ips": ["192.168.0.142"],
		"public_ips": ["10.154.219.187", "10.154.219.186"],
		"db_user_name": "root",
		"vpc_id": "b21630c1-e7d3-450d-907d-39ef5f445ae7",
		"subnet_id": "45557a98-9e17-4600-8aec-999150bc4eef",
		"security_group_id": "38815c5c-482b-450a-80b6-0a301f2afd97",
		"flavor_ref": "rds.mysql.s1.large",
		"switch_strategy": "",
		"backup_strategy": {
			"start_time": "19:00-20:00",
			"keep_days": 7
		},
		"maintenance_window": "02:00-06:00",
		"related_instance": [],
		"disk_encryption_id": "",
		"time_zone": ""
	}],
	"total_count": 1
}
`
const RdsCreateResponse = `
{
  "instance":{
           "id": "dsfae23fsfdsae3435in01",
           "name": "default",
           "datastore": {
             "type": "MySQL",
             "version": "5.6"
           },
           "ha": {
             "mode": "ha",
             "replication_mode": "semisync"
           },
           "flavor_ref": "rds.mysql.s1.large.ha",
           "volume": {
               "type": "ULTRAHIGH",
               "size": 100
           },
           "disk_encryption_id":  "2gfdsh-844a-4023-a776-fc5c5fb71fb4",
           "region": "eu-de",
           "availability_zone": "eu-de-01,en-de-02",
           "vpc_id": "490a4a08-ef4b-44c5-94be-3051ef9e4fce",
           "subnet_id": "0e2eda62-1d42-4d64-a9d1-4e9aa9cd994f",
           "security_group_id": "2a1f7fc8-3307-42a7-aa6f-42c8b9b8f8c5",
           "port": "3306",
           "backup_strategy": {
             "start_time": "08:15-09:15",
             "keep_days": 3
           },
           "configuration_id": "452408-44c5-94be-305145fg",
           "charge_info": {
                   "charge_mode": "postPaid"
           }
  },
  "job_id": "dff1d289-4d03-4942-8b9f-463ea07c000d"
}
`
const RdsJobResponse = `
{
  "job": {
    "created": "2022-03-04T21:38:47+0000",
    "entities": {
      "instance": {
        "datastore": {
          "type": "mysql",
          "version": "5.6"
        },
        "type": "Ha"
      }
    },
    "id": "dff1d289-4d03-4942-8b9f-463ea07c000d",
    "instance": {
      "id": "16fa025c33b444a6bb2f04e705767adbin01",
      "name": "default"
    },
    "name": "CreateMysqlSingleHAInstance",
    "process": "100%",
    "status": "Completed"
  }
}
`
func Test_secgroupGet(t *testing.T) {

	th.SetupHTTP()
	defer th.TeardownHTTP()

	th.Mux.HandleFunc("/security-groups", func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "GET")
		th.TestHeader(t, r, "X-Auth-Token", fake.TokenID)

		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		_, _ = fmt.Fprint(w, SecurityGroupListResponse)
	})

	sg, err := secgroupGet(fake.ServiceClient(),  &groups.ListOpts{Name: "default"})
	th.AssertNoErr(t, err)

	th.AssertEquals(t, "default", sg.Name)
	th.AssertEquals(t, "85cc3048-abc3-43cc-89b3-377341426ac5", sg.ID)

}

func Test_subnetGet(t *testing.T) {

	th.SetupHTTP()
	defer th.TeardownHTTP()

	th.Mux.HandleFunc("/subnets", func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "GET")
		th.TestHeader(t, r, "X-Auth-Token", fake.TokenID)

		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		_, _ = fmt.Fprint(w, SubnetListResponse)
	})

	sg, err := subnetGet(fake.ServiceClient(),  &subnets.ListOpts{Name: "golang"})
	th.AssertNoErr(t, err)

	th.AssertEquals(t, "golang", sg.Name)
	th.AssertEquals(t, "011fc878-5521-4654-a1ad-f5b0b5820302", sg.ID)

}

func Test_vpcGet(t *testing.T) {

	th.SetupHTTP()
	defer th.TeardownHTTP()

	th.Mux.HandleFunc("/vpcs", func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "GET")
		th.TestHeader(t, r, "X-Auth-Token", fake.TokenID)

		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		_, _ = fmt.Fprint(w, VpcListResponse)
	})

	sg, err := vpcGet(fake.ServiceClient(),  &vpcs.ListOpts{Name: "golang"})
	th.AssertNoErr(t, err)

	th.AssertEquals(t, "default", sg.Name)
	th.AssertEquals(t, "13551d6b-755d-4757-b956-536f674975c0", sg.ID)

}

func Test_rdsGet(t *testing.T) {

	th.SetupHTTP()
	defer th.TeardownHTTP()

	th.Mux.HandleFunc("/instances", func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "GET")
		th.TestHeader(t, r, "X-Auth-Token", fake.TokenID)

		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		_, _ = fmt.Fprint(w, RdsGetResponse)
	})

	sg, err := rdsGet(fake.ServiceClient(), "ed7cc6166ec24360a5ed5c5c9c2ed726in01")
	th.AssertNoErr(t, err)

	th.AssertEquals(t, "default", sg.Name)
	th.AssertEquals(t, "ed7cc6166ec24360a5ed5c5c9c2ed726in01", sg.Id)
	tools.PrintResource(t, sg)

}

func MockRdsResponse(t *testing.T) {

	th.Mux.HandleFunc("/instances", func(w http.ResponseWriter, r *http.Request) {
		th.TestHeader(t, r, "X-Auth-Token", fake.TokenID)
		switch r.Method {
		case "GET":
			w.Header().Add("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			_, _ = fmt.Fprint(w, RdsGetResponse)
		case "POST":
			w.Header().Add("Content-Type", "application/json")
			w.WriteHeader(http.StatusAccepted)
			_, _ = fmt.Fprint(w, RdsCreateResponse)
		}
	})
}

func Test_rdsCreate(t *testing.T) {

	th.SetupHTTP()
	defer th.TeardownHTTP()

	MockRdsResponse(t)

	th.Mux.HandleFunc("/security-groups", func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "GET")
		th.TestHeader(t, r, "X-Auth-Token", fake.TokenID)

		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		_, _ = fmt.Fprint(w, SecurityGroupListResponse)
	})

	th.Mux.HandleFunc("/subnets", func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "GET")
		th.TestHeader(t, r, "X-Auth-Token", fake.TokenID)

		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		_, _ = fmt.Fprint(w, SubnetListResponse)
	})

	th.Mux.HandleFunc("/vpcs", func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "GET")
		th.TestHeader(t, r, "X-Auth-Token", fake.TokenID)

		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		_, _ = fmt.Fprint(w, VpcListResponse)
	})

	th.Mux.HandleFunc("/jobs", func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "GET")
		th.TestHeader(t, r, "X-Auth-Token", fake.TokenID)

		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		_, _ = fmt.Fprint(w, RdsJobResponse)
	})

	RdsCreateOpts := &instances.CreateRdsOpts{
		Name:             "default",
		Port:             "3306",
		Password:         "acc-test-password1!",
		BackupStrategy: &instances.BackupStrategy{
			StartTime: "08:15-09:15",
			KeepDays:  12,
		},
		FlavorRef:        "rds.mysql.s1.large.ha",
		Region:           "eu-de",
		AvailabilityZone: "eu-de-01,eu-de-02",
		VpcId:            "490a4a08-ef4b-44c5-94be-3051ef9e4fce",
		SubnetId:         "0e2eda62-1d42-4d64-a9d1-4e9aa9cd994f",
		SecurityGroupId:  "2a1f7fc8-3307-42a7-aa6f-42c8b9b8f8c5",
		Volume: &instances.Volume{
			Type: "ULTRAHIGH",
			Size: 100,
		},
		Datastore: &instances.Datastore{
			Type:    "MySQL",
			Version: "5.6",
		},
		Ha: &instances.Ha{
			Mode:            "ha",
			ReplicationMode: "semisync",
		},
	}
	err := rdsCreate(fake.ServiceClient(), fake.ServiceClient(), fake.ServiceClient(), RdsCreateOpts)
	th.AssertNoErr(t, err)

	//th.AssertEquals(t, "default", sg.Name)
	//th.AssertEquals(t, "ed7cc6166ec24360a5ed5c5c9c2ed726in01", sg.Id)
	//tools.PrintResource(t, sg)

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
