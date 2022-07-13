# etcdHighNumberOfFailedGRPCRequests

## Meaning

This alert fires when at least 50% of etcd gRPC requests failed in the past 10
minutes and sends a warning at 10%.

## Impact

First establish which gRPC method is failing, this will be visible in the alert.
If it's not part of the alert, the following query will display method and etcd
instance that has failing requests:

```sh
(sum(rate(grpc_server_handled_total{job="etcd", grpc_code=~"Unknown|FailedPrecondition|ResourceExhausted|Internal|Unavailable|DataLoss|DeadlineExceeded"}[5m])) without (grpc_type, grpc_code)
    /
(sum(rate(grpc_server_handled_total{job="etcd"}[5m])) without (grpc_type, grpc_code)
    > 2 and on ()(sum(cluster_infrastructure_provider{type!~"ipi|BareMetal"} == bool 1)))) * 100 > 10
```

## Diagnosis

All the gRPC errors should also be logged in each respective etcd instance logs.
You can get the instance name from the alert that is firing or by running the
query detailed above. Those etcd instance logs should serve as further insight
into what is wrong.

To get logs of etcd containers either check the instance from the alert and
check logs directly or run the following:

```sh
oc logs -n openshift-etcd -lapp=etcd -c etcd
```

### Defrag method errors

If defrag method is failing, this could be due to defrag that is periodically
performed by cluster-etcd-operator pe starting from OpenShift v4.9 onwards. To
verify this check the logs of cluster-etcd-operator.

```sh
oc logs -l app=etcd-operator -n openshift-etcd-operator --tail=-1
```

If you have run defrag manually on older OpenShift versions check the errors of
those manual runs.

### MemberList method errors

Member list is most likely performed by cluster-etcd-operator, so it's also best
to check also logs of cluster-etcd-operator for any errors:

```sh
oc logs -l app=etcd-operator -n openshift-etcd-operator --tail=-1
```

## Mitigation

Depending on the above diagnosis, the issue will most likely be described in the
error log line of either etcd or openshift-etcd-operator. Most likely causes
tend to be networking issues.
