# otc-rds-client
![Coverage](https://img.shields.io/badge/Coverage-%7B%7D%00type..eq.%5B4%5Dk8s.io%2Fklog%2Fv2.flushSyncWriter%00type..eq.k8s.io%2Fklog%2Fv2.traceLocation%00type..eq.k8s.io%2Fklog%2Fv2.syncBuffer%00k8s.io%2Fklog%2Fv2.syncBuffer.Flush%00k8s.io%2Fklog%2Fv2.syncBuffer.ReadFrom%00k8s.io%2Fklog%2Fv2.syncBuffer.Size%00k8s.io%2Fklog%2Fv2.syncBuffer.WriteByte%00k8s.io%2Fklog%2Fv2.syncBuffer.WriteString%00net%2Fhttp%2Fhttptest.newLocalListener%00net%2Fhttp%2Fhttptest.init.0%00net%2Fhttp%2Fhttptest.strSliceContainsPrefix%00net%2Fhttp%2Fhttptest.NewServer%00net%2Fhttp%2Fhttptest.%28%2AServer%29.Start%00net%2Fhttp%2Fhttptest.%28%2AServer%29.Close%00net%2Fhttp%2Fhttptest.%28%2AServer%29.logCloseHangDebugInfo%00net%2Fhttp%2Fhttptest.%28%2AServer%29.goServe%00net%2Fhttp%2Fhttptest.%28%2AServer%29.wrap%00net%2Fhttp%2Fhttptest.%28%2AServer%29.forgetConn%00net%2Fhttp%2Fhttptest.%28%2AServer%29.goServe.func1%00net%2Fhttp%2Fhttptest.%28%2AServer%29.wrap.func1%00net%2Fhttp%2Fhttptest.%28%2AServer%29.logCloseHangDebugInfo-fm%00github.com%2Fopentelekomcloud%2Fgophertelekomcloud%2Ftesthelper.prefix%00github.com%2Fopentelekomcloud%2Fgophertelekomcloud%2Ftesthelper.logFatal%00github.com%2Fopentelekomcloud%2Fgophertelekomcloud%2Ftesthelper.AssertEquals%00github.com%2Fopentelekomcloud%2Fgophertelekomcloud%2Ftesthelper.AssertNoErr%00github.com%2Fopentelekomcloud%2Fgophertelekomcloud%2Ftesthelper.TeardownHTTP%00github.com%2Fopentelekomcloud%2Fgophertelekomcloud%2Ftesthelper.TestMethod%00github.com%2Fopentelekomcloud%2Fgophertelekomcloud%2Ftesthelper.TestHeader%00rds%252ego.secgroupGet%00rds%252ego.Test_secgroupGet%00rds%252ego.Test_main%00rds%252ego.Test_secgroupGet.func1%00main.init.0%00main.init.1%00main.coverRegisterFile%00main.main%00main.init%00-red)

creates rds ha instances in OTC tenant based on `rds.yaml`

[![CodeQL](https://github.com/eumel8/otc-rds-client/actions/workflows/codeql-analysis.yml/badge.svg)](https://github.com/eumel8/otc-rds-client/actions/workflows/codeql-analysis.yml)


## example

```yaml
name: "mydb"
datastore:
  type: "MySQL"
  version: "5.7"
volume:
  type: "COMMON"
  size: 40
ha:
  mode: "Ha"
  replicationmode:  "semisync"
port: "3306"
password: "A12345678+"
backupstrategy:
  starttime: "01:00-02:00"
  keepdays: 10
flavorref: "rds.mysql.c2.xlarge.ha"
region: "eu-de"
availabilityzone: "eu-de-01,eu-de-02"
vpc: "golang"
subnet: "golang"
securitygroup: "golang"
```

## usage


```bash
$ ./rds
192.168.0.159
```

verify result

```bash
$ openstack rds instance show mydb
```

note: will wait until instance is in state `ACTIVE` (timeout: 1800)

## find valid values

```bash
$ openstack rds datastore type list
+------------+
| Name       |
+------------+
| MySQL      |
| PostgreSQL |
| SQLServer  |
+------------+
```

```bash
$ openstack rds datastore version list mysql
+--------------------------------------+------+
| ID                                   | Name |
+--------------------------------------+------+
| bf5a9a94-dbb1-3a0f-b37b-e257595441fa | 5.6  |
| b5863f8e-8081-3066-8221-7b3760218bc3 | 5.7  |
| c4f55bf1-0f4b-32ab-aa98-9becf6bdfef8 | 8.0  |
+--------------------------------------+------+
```

```bash
$ openstack rds flavor list mysql 8.0
+----------------------------+---------------+-------+-----+
| name                       | instance_mode | vcpus | ram |
+----------------------------+---------------+-------+-----+
| rds.mysql.c2.medium.rr     | replica       | 1     |   2 |
| rds.mysql.m1.large.rr      | replica       | 2     |  16 |
| rds.mysql.m1.xlarge.rr     | replica       | 4     |  32 |
| rds.mysql.m1.2xlarge.rr    | replica       | 8     |  64 |
| rds.mysql.s1.medium.ha     | ha            | 1     |   4 |
| rds.mysql.c2.medium.ha     | ha            | 1     |   2 |
| rds.mysql.m1.large.ha      | ha            | 2     |  16 |
| rds.mysql.m1.xlarge.ha     | ha            | 4     |  32 |
| rds.mysql.m1.2xlarge.ha    | ha            | 8     |  64 |
| rds.mysql.s1.medium        | single        | 1     |   4 |
| rds.mysql.c2.medium        | single        | 1     |   2 |
| rds.mysql.m1.large         | single        | 2     |  16 |
| rds.mysql.m1.xlarge        | single        | 4     |  32 |
| rds.mysql.m1.2xlarge       | single        | 8     |  64 |
| rds.mysql.s1.medium.rr     | replica       | 1     |   4 |
| rds.mysql.m3.15xlarge.8.rr | replica       | 60    | 512 |
| rds.mysql.m1.4xlarge       | single        | 16    | 128 |
| rds.mysql.m1.4xlarge.ha    | ha            | 16    | 128 |
| rds.mysql.m1.4xlarge.rr    | replica       | 16    | 128 |
| rds.mysql.m1.8xlarge.rr    | replica       | 32    | 256 |
| rds.mysql.m1.8xlarge       | single        | 32    | 256 |
| rds.mysql.m1.8xlarge.ha    | ha            | 32    | 256 |
| rds.mysql.m3.15xlarge.8    | single        | 60    | 512 |
| rds.mysql.m3.15xlarge.8.ha | ha            | 60    | 512 |
| rds.mysql.s1.large.rr      | replica       | 2     |   8 |
| rds.mysql.s1.xlarge.rr     | replica       | 4     |  16 |
| rds.mysql.s1.2xlarge.rr    | replica       | 8     |  32 |
| rds.mysql.c2.large.rr      | replica       | 2     |   4 |
| rds.mysql.c2.xlarge.rr     | replica       | 4     |   8 |
| rds.mysql.s1.large.ha      | ha            | 2     |   8 |
| rds.mysql.s1.xlarge.ha     | ha            | 4     |  16 |
| rds.mysql.s1.2xlarge.ha    | ha            | 8     |  32 |
| rds.mysql.c2.large.ha      | ha            | 2     |   4 |
| rds.mysql.c2.xlarge.ha     | ha            | 4     |   8 |
| rds.mysql.s1.large         | single        | 2     |   8 |
| rds.mysql.s1.xlarge        | single        | 4     |  16 |
| rds.mysql.s1.2xlarge       | single        | 8     |  32 |
| rds.mysql.c2.large         | single        | 2     |   4 |
| rds.mysql.c2.xlarge        | single        | 4     |   8 |
| rds.mysql.c3.15xlarge.4    | single        | 60    | 256 |
| rds.mysql.c3.15xlarge.4.ha | ha            | 60    | 256 |
| rds.mysql.c3.15xlarge.4.rr | replica       | 60    | 256 |
| rds.mysql.c2.2xlarge       | single        | 8     |  16 |
| rds.mysql.c2.2xlarge.ha    | ha            | 8     |  16 |
| rds.mysql.c2.2xlarge.rr    | replica       | 8     |  16 |
| rds.mysql.c2.4xlarge       | single        | 16    |  32 |
| rds.mysql.c2.4xlarge.ha    | ha            | 16    |  32 |
| rds.mysql.c2.4xlarge.rr    | replica       | 16    |  32 |
| rds.mysql.s1.4xlarge       | single        | 16    |  64 |
| rds.mysql.s1.4xlarge.ha    | ha            | 16    |  64 |
| rds.mysql.s1.4xlarge.rr    | replica       | 16    |  64 |
| rds.mysql.c2.8xlarge       | single        | 32    |  64 |
| rds.mysql.c2.8xlarge.ha    | ha            | 32    |  64 |
| rds.mysql.c2.8xlarge.rr    | replica       | 32    |  64 |
| rds.mysql.s1.8xlarge       | single        | 32    | 128 |
| rds.mysql.s1.8xlarge.ha    | ha            | 32    | 128 |
| rds.mysql.s1.8xlarge.rr    | replica       | 32    | 128 |
| rds.mysql.c3.15xlarge.2.ha | ha            | 60    | 128 |
| rds.mysql.c3.15xlarge.2.rr | replica       | 60    | 128 |
| rds.mysql.c3.15xlarge.2    | single        | 60    | 128 |
+----------------------------+---------------+-------+-----+
```

list default parameter group or create your own

```bash
openstack rds configuration list
```

volume types: `COMMON|ULTRAHIGH`

ha mode: `async|semisync` (MySQL) `async|sync` (PostgreSQL) hint: no single instance supported

refer [OTC API DOC](https://docs.otc.t-systems.com/api/rds/rds_01_0002.html)


## Credits

Frank Kloeker f.kloeker@telekom.de

Life is for sharing. If you have an issue with the code or want to improve it, feel free to open an issue or an pull request.
