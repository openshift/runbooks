# CephOSDVersionMismatch

## Meaning

There are different versions of Ceph OSD components running. Typically this
alert triggers during an upgrade that is taking a long time.

## Impact

It will impact cluster performance.

## Diagnosis

Verify Ceph version in OSD Pods:

```bash
    oc describe pods -n openshift-storage --selector app=rook-ceph-osd | grep CONTAINER_IMAGE
```

All the pods must have the same image

## Mitigation

Usually the problem is solved (all OSDs daemons running same Ceph version) when
 the upgrade has finished.

If the alert persists after upgrade:

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

If no issues found, [gather_logs](helpers/gatherLogs.md) to provide more
information to support teams.
