package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/gophercloud/utils/client"
	"github.com/opentelekomcloud/gophertelekomcloud"
	"github.com/opentelekomcloud/gophertelekomcloud/openstack"
	"github.com/opentelekomcloud/gophertelekomcloud/openstack/networking/v1/subnets"
	"github.com/opentelekomcloud/gophertelekomcloud/openstack/networking/v1/vpcs"
	"github.com/opentelekomcloud/gophertelekomcloud/openstack/networking/v2/extensions/security/groups"
	"github.com/opentelekomcloud/gophertelekomcloud/openstack/rds/v3/instances"
	"gopkg.in/yaml.v3"
	"k8s.io/klog/v2"
)

var osExit = os.Exit

//var getenv = os.Getenv

const (
	AppVersion = "0.0.4"
	RdsYaml    = "rds.yaml"
)

type conf struct {
	Name             string          `yaml:"name"`
	Datastore        *Datastore      `yaml:"datastore"`
	Ha               *Ha             `yaml:"ha"`
	Port             string          `yaml:"port"`
	Password         string          `yaml:"password"`
	BackupStrategy   *BackupStrategy `yaml:"backupstrategy"`
	FlavorRef        string          `yaml:"flavorref"`
	Volume           *Volume         `yaml:"volume"`
	Region           string          `yaml:"region"`
	AvailabilityZone string          `yaml:"availabilityzone"`
	Vpc              string          `yaml:"vpc"`
	Subnet           string          `yaml:"subnet"`
	SecurityGroup    string          `yaml:"securitygroup"`
}

type Datastore struct {
	Type    string `json:"type" required:"true"`
	Version string `json:"version" required:"true"`
}

type Ha struct {
	Mode            string `json:"mode" required:"true"`
	ReplicationMode string `json:"replicationmode,omitempty"`
}

type BackupStrategy struct {
	StartTime string `json:"starttime" required:"true"`
	KeepDays  int    `json:"keepdays,omitempty"`
}

type Volume struct {
	Type string `json:"type" required:"true"`
	Size int    `json:"size" required:"true"`
}

func secgroupGet(client *golangsdk.ServiceClient, opts *groups.ListOpts) (*groups.SecGroup, error) {

	pages, err := groups.List(client, *opts).AllPages()
	if err != nil {
		return nil, err
	}
	n, err := groups.ExtractGroups(pages)
	if len(n) == 0 {
		klog.Exitf("no secgroup found")
	}

	return &n[0], nil
}

func subnetGet(client *golangsdk.ServiceClient, opts *subnets.ListOpts) (*subnets.Subnet, error) {

	n, err := subnets.List(client, *opts)
	if err != nil {
		return nil, err
	}
	if len(n) == 0 {
		klog.Exitf("no subnet found")
	}

	return &n[0], nil
}

func vpcGet(client *golangsdk.ServiceClient, opts *vpcs.ListOpts) (*vpcs.Vpc, error) {

	n, err := vpcs.List(client, *opts)
	if err != nil {
		return nil, err
	}

	if len(n) == 0 {
		klog.Exitf("no vpc found")
	}

	return &n[0], nil
}

func rdsGet(client *golangsdk.ServiceClient, rdsId string) (*instances.RdsInstanceResponse, error) {

	listOpts := instances.ListRdsInstanceOpts{
		Id: rdsId,
	}
	allPages, err := instances.List(client, listOpts).AllPages()
	if err != nil {
		return nil, err
	}

	n, err := instances.ExtractRdsInstances(allPages)
	if err != nil {
		return nil, err
	}
	if len(n.Instances) == 0 {
		return nil, nil
	}
	return &n.Instances[0], nil
}

func rdsCreate(netclient1 *golangsdk.ServiceClient, netclient2 *golangsdk.ServiceClient, client *golangsdk.ServiceClient, opts *instances.CreateRdsOpts) error {

	var c conf
	c.getConf()

	g, err := secgroupGet(netclient2, &groups.ListOpts{Name: c.SecurityGroup})
	if err != nil {
		klog.Exitf("error getting secgroup state: %v", err)
	}

	s, err := subnetGet(netclient1, &subnets.ListOpts{Name: c.Subnet})
	if err != nil {
		klog.Exitf("error getting subnet state: %v", err)
	}

	v, err := vpcGet(netclient1, &vpcs.ListOpts{Name: c.Vpc})
	if err != nil {
		klog.Exitf("error getting vpc state: %v", err)
	}

	createOpts := instances.CreateRdsOpts{
		Name: c.Name,
		Datastore: &instances.Datastore{
			Type:    c.Datastore.Type,
			Version: c.Datastore.Version,
		},
		Ha: &instances.Ha{
			Mode:            c.Ha.Mode,
			ReplicationMode: c.Ha.ReplicationMode,
		},
		Port:     c.Port,
		Password: c.Password,
		BackupStrategy: &instances.BackupStrategy{
			StartTime: c.BackupStrategy.StartTime,
			KeepDays:  c.BackupStrategy.KeepDays,
		},
		FlavorRef: c.FlavorRef,
		Volume: &instances.Volume{
			Type: c.Volume.Type,
			Size: c.Volume.Size,
		},
		Region:           c.Region,
		AvailabilityZone: c.AvailabilityZone,
		VpcId:            v.ID,
		SubnetId:         s.ID,
		SecurityGroupId:  g.ID,
	}

	createResult := instances.Create(client, createOpts)
	r, err := createResult.Extract()
	if err != nil {
		klog.Exitf("error creating rds instance: %v", err)
	}
	jobResponse, err := createResult.ExtractJobResponse()
	if err != nil {
		klog.Exitf("error creating rds job: %v", err)
	}

	if err := instances.WaitForJobCompleted(client, int(1800), jobResponse.JobID); err != nil {
		klog.Exitf("error getting rds job: %v", err)
	}

	rdsInstance, err := rdsGet(client, r.Instance.Id)

	fmt.Println(rdsInstance.PrivateIps[0])
	if err != nil {
		klog.Exitf("error getting rds state: %v", err)
	}

	return nil
}

func (c *conf) getConf() *conf {

	yfile, err := ioutil.ReadFile(RdsYaml)
	if err != nil {
		klog.Exitf("error reading yaml file: %v", err)
	}

	err = yaml.Unmarshal(yfile, c)
	if err != nil {
		klog.Exitf("error unmarshal yaml file: %v", err)
	}

	return c
}

// func (p *golangsdk.ProviderClient) getProvider() *golangsdk.ProviderClient {
func getProvider() *golangsdk.ProviderClient {

	if os.Getenv("OS_AUTH_URL") == "" {
		os.Setenv("OS_AUTH_URL", "https://iam.eu-de.otc.t-systems.com:443/v3")
	}

	if os.Getenv("OS_IDENTITY_API_VERSION") == "" {
		os.Setenv("OS_IDENTITY_API_VERSION", "3")
	}

	if os.Getenv("OS_REGION_NAME") == "" {
		os.Setenv("OS_REGION_NAME", "eu-de")
	}

	if os.Getenv("OS_PROJECT_NAME") == "" {
		os.Setenv("OS_PROJECT_NAME", "eu-de")
	}

	opts, err := openstack.AuthOptionsFromEnv()
	if err != nil {
		klog.Exitf("error getting auth from env: %v", err)
	}

	provider, err := openstack.AuthenticatedClient(opts)
	if err != nil {
		klog.Exitf("unable to initialize openstack client: %v", err)
	}

	if os.Getenv("OS_DEBUG") != "" {
		provider.HTTPClient = http.Client{
			Transport: &client.RoundTripper{
				Rt:     &http.Transport{},
				Logger: &client.DefaultLogger{},
			},
		}
	}
	return provider
}

func main() {

	version := flag.Bool("version", false, "app version")
	help := flag.Bool("help", false, "print out the help")

	flag.Parse()

	if *help {
		fmt.Println("Provide ENV variable to connect OTC: OS_PROJECT_NAME, OS_REGION_NAME, OS_AUTH_URL, OS_IDENTITY_API_VERSION, OS_USER_DOMAIN_NAME, OS_USERNAME, OS_PASSWORD")
		osExit(0)
		return
	}

	if *version {
		fmt.Println("version", AppVersion)
		osExit(0)
		return
	}

	provider := getProvider()
	network1, err := openstack.NewNetworkV1(provider, golangsdk.EndpointOpts{})
	network2, err := openstack.NewNetworkV2(provider, golangsdk.EndpointOpts{})
	rds, err := openstack.NewRDSV3(provider, golangsdk.EndpointOpts{})
	if err != nil {
		klog.Exitf("unable to initialize rds client: %v", err)
	}

	rdsCreate(network1, network2, rds, &instances.CreateRdsOpts{})
	if err != nil {
		klog.Exitf("rds creating failed: %v", err)
	}
}
