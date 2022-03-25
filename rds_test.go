package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	// "io/ioutil"
	"net/http"
	"os"
	"testing"

	"github.com/opentelekomcloud/gophertelekomcloud"
	"github.com/opentelekomcloud/gophertelekomcloud/acceptance/tools"
	"github.com/opentelekomcloud/gophertelekomcloud/openstack/networking/v1/subnets"
	"github.com/opentelekomcloud/gophertelekomcloud/openstack/networking/v1/vpcs"
	"github.com/opentelekomcloud/gophertelekomcloud/openstack/networking/v2/extensions/security/groups"
	"github.com/opentelekomcloud/gophertelekomcloud/openstack/rds/v3/instances"
	th "github.com/opentelekomcloud/gophertelekomcloud/testhelper"
	fake "github.com/opentelekomcloud/gophertelekomcloud/testhelper/client"
)

// test service function with testhelper gophertelekomcloud and
// ProviderClient and main() with the MockMuxer on "http://127.0.0.1:50000/v3")

func Test_secgroupGet(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()

	MockRdsResponse(t)

	sg, err := secgroupGet(fake.ServiceClient(), &groups.ListOpts{Name: "golang"})
	th.AssertNoErr(t, err)

	th.AssertEquals(t, "golang", sg.Name)
	th.AssertEquals(t, "85cc3048-abc3-43cc-89b3-377341426ac5", sg.ID)
}

func Test_subnetGet(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()

	MockRdsResponse(t)

	sg, err := subnetGet(fake.ServiceClient(), &subnets.ListOpts{Name: "golang"})
	th.AssertNoErr(t, err)

	th.AssertEquals(t, "golang", sg.Name)
	th.AssertEquals(t, "011fc878-5521-4654-a1ad-f5b0b5820302", sg.ID)
}

func Test_vpcGet(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()

	MockRdsResponse(t)

	sg, err := vpcGet(fake.ServiceClient(), &vpcs.ListOpts{Name: "golang"})
	th.AssertNoErr(t, err)

	th.AssertEquals(t, "golang", sg.Name)
	th.AssertEquals(t, "13551d6b-755d-4757-b956-536f674975c0", sg.ID)
}

func Test_rdsGet(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()

	MockRdsResponse(t)

	sg, err := rdsGet(fake.ServiceClient(), "ed7cc6166ec24360a5ed5c5c9c2ed726in01")
	th.AssertNoErr(t, err)

	th.AssertEquals(t, "default", sg.Name)
	th.AssertEquals(t, "ed7cc6166ec24360a5ed5c5c9c2ed726in01", sg.Id)
	tools.PrintResource(t, sg)
}

func MockMuxer() {
	mux := http.NewServeMux()

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case "GET":
			w.Header().Add("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			_, _ = fmt.Fprint(w, ProviderGetResponse)
		case "POST":
			w.Header().Add("X-Subject-Token", "dG9rZW46IDEyMzQK")
			w.Header().Add("Content-Type", "application/json")
			w.WriteHeader(http.StatusCreated)
			_, _ = fmt.Fprint(w, ProviderPostResponse)
		}
	})

	mux.HandleFunc("/v2.0/security-groups", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case "GET":
			w.Header().Add("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)

			_, _ = fmt.Fprint(w, SecurityGroupListResponse)
			/* Debug output of the request
			uri := r.URL.String()
			fmt.Printf("Uri: %s\n", string(uri))
			body, err := ioutil.ReadAll(r.Body)
			if err != nil {
				panic(err)
			}
			fmt.Printf("Body: %s\n", body)
			*/
		}
	})

	mux.HandleFunc("/v1/bbfe8c41dd034a07bebd592bf03b4b0c/subnets", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case "GET":
			w.Header().Add("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)

			_, _ = fmt.Fprint(w, SubnetListResponse)
		}
	})

	mux.HandleFunc("/v1/bbfe8c41dd034a07bebd592bf03b4b0c/vpcs", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case "GET":
			w.Header().Add("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)

			_, _ = fmt.Fprint(w, VpcListResponse)
		}
	})

	mux.HandleFunc("/instances", func(w http.ResponseWriter, r *http.Request) {
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

	mux.HandleFunc("/jobs", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case "GET":
			w.Header().Add("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)

			_, _ = fmt.Fprint(w, RdsJobResponse)
		}
	})

	fmt.Println("Listening...")

	var retries int = 3

	for retries > 0 {
		err := http.ListenAndServe("127.0.0.1:50000", mux)
		if err != nil {
			fmt.Println("Restart http server ... ", err)
			retries -= 1
		} else {
			break
		}
	}

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
		/* Example output http response
		uri := r.URL.String()
		fmt.Println(string(uri))
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			panic(err)
		}
		fmt.Printf("%s", body)
		*/
	})
}

func Test_rdsCreate(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()

	MockRdsResponse(t)

	RdsCreateOpts := &instances.CreateRdsOpts{
		Name:     "default",
		Port:     "3306",
		Password: "acc-test-password1!",
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

}

func Test_flags(t *testing.T) {
	testCases := []struct {
		name     string
		flags    []string
		expected int
	}{
		{"help_exit_0", []string{"-help"}, 0},
		{"version_exit_0", []string{"-version"}, 0},
	}
	for _, tc := range testCases {
		t.Run("test "+tc.name, func(t *testing.T) {

			// prevent os.Exit exit here
			// https://stackoverflow.com/questions/40615641/testing-os-exit-scenarios-in-go-with-coverage-information-coveralls-io-goverall
			oldOsExit := osExit
			defer func() {
				osExit = oldOsExit
			}()

			var got int
			tmpExit := func(code int) {
				got = code
			}

			osExit = tmpExit

			// flag.CommandLine = flag.NewFlagSet(tc.name, flag.ExitOnError)
			os.Args = append([]string{tc.name}, tc.flags...)

			err := os.Setenv("OS_AUTH_URL", "")
			th.AssertNoErr(t, err)

			getFlags(os.Args[1])

			if got != tc.expected {
				t.Errorf("Expected exit code: %d, got: %d", tc.expected, got)
			}
		})
	}
}

func Test_getProvider(t *testing.T) {
	go MockMuxer()

	err := os.Setenv("OS_USERNAME", "test")
	th.AssertNoErr(t, err)
	err = os.Setenv("OS_USER_DOMAIN_NAME", "test")
	th.AssertNoErr(t, err)
	err = os.Setenv("OS_PASSWORD", "test")
	th.AssertNoErr(t, err)
	err = os.Setenv("OS_IDENTITY_API_VERSION", "3")
	th.AssertNoErr(t, err)
	err = os.Setenv("OS_AUTH_URL", "http://127.0.0.1:50000/v3")
	th.AssertNoErr(t, err)

	provider := getProvider()
	defer getProvider()

	p := &golangsdk.ProviderClient{
		UserID: "91dca41cc55e4614aaca83b78af8ddc5",
	}
	th.CheckDeepEquals(t, p.UserID, provider.UserID)
	fmt.Println("IdentityEndpoint: ", provider.IdentityEndpoint)
	return
}

func PrettyJson(str string) (string, error) {
	var prettyJSON bytes.Buffer
	if err := json.Indent(&prettyJSON, []byte(str), "", "    "); err != nil {
		return "", err
	}
	return prettyJSON.String(), nil
}

func Test_main(t *testing.T) {
	oldOsExit := osExit
	defer func() {
		osExit = oldOsExit
	}()

	var got int
	tmpExit := func(code int) {
		got = code
	}

	osExit = tmpExit

	err := os.Setenv("OS_AUTH_URL", "")
	th.AssertNoErr(t, err)

	main()

	if got != 0 {
		t.Errorf("Expected exit code: %d", got)
	}
}

func Test_Create(t *testing.T) {
	go MockMuxer()

	err := os.Setenv("OS_USERNAME", "test")
	th.AssertNoErr(t, err)
	err = os.Setenv("OS_USER_DOMAIN_NAME", "test")
	th.AssertNoErr(t, err)
	err = os.Setenv("OS_PASSWORD", "test")
	th.AssertNoErr(t, err)
	err = os.Setenv("OS_IDENTITY_API_VERSION", "3")
	th.AssertNoErr(t, err)
	err = os.Setenv("OS_AUTH_URL", "http://127.0.0.1:50000/v3")
	th.AssertNoErr(t, err)

	defer Create()
}
