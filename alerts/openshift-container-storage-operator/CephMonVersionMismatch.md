# CephMonVersionMismatch

## Meaning

There are different versions of Ceph Mon components running.. Typically this
alert is triggered during an upgrade that is taking a long time.

## Impact

It will impact cluster availability if the number of monitors are not enough
 to get quorum. Cluster operations will be blocked until quorum will be
 established again

## Diagnosis

Verify Ceph version in Mon Pods:

```console
    oc describe pods -n openshift-storage --selector app=rook-ceph-mon | grep CONTAINER_IMAGE
```

All the pods must have the same image

## Mitigation

Usually the problem is solved (all Monitor daemons running same Ceph version)
when the upgrade has finished.

If the alert persists after upgrade:

Verify the connectivity between monitors is working properly, verifying
[network connectivity](helpers/networkConnectivity.md)

Verify the ODF operator events and logs in order to find an error and an
explanation about the problem updating the OSD daemon with different version.

```bash
    ocs_operator=$(oc describe deployment -n openshift-storage ocs-operator | grep OPERATOR_CONDITION_NAME: | awk '{ print $2 }')
    oc get events --field-selector involvedObject.name=$ocs_operator --namespace openshift-storage
```

```bash
    oc logs -n openshift-storage --selector name=ocs-operator -c ocs-operator
```

If nothing found, verify the
[ODF operator state](helpers/checkOperator.md)

If the ODF operator does not present any problem,
see [general diagnosis document](helpers/diagnosis.md)

If no issues found, [gather_logs](helpers/gather_logs.md) to provide more
information to support teams.
