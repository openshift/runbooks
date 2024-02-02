# StorageClientIncompatibleOperatorVersion

## Meaning

OCS Client operator version of connected ODF Client is not same as OCS
operator of ODF Provider

## Impact

At Warning level, ODF Client is lagging ODF Provider by one minor version. This
stops ODF Provider from getting upgraded to next minor/patch version.

At Critical level, ODF Client is ahead or lagging by two minor versions of ODF
Provider. This reduces the supportability of connected ODF Client.

## Diagnosis

From OCP Console on ODF Provider cluster:

1. Take note of ODF Version by following __Operators -> Installed Operators ->
Project: All Projects -> Search for "OpenShift Data Foundation"__
2. Take note of connected clients by following __Storage -> Storage Clients ->
Data Foundation version column corresponding to the storageconsumer from the
alert"__

Observe the difference in minor versions in info gathered from above process.

NOTE: You might want to run below command for enabling ODF Client console if
you don't see __Storage Clients__ UI

```bash
 oc patch console.v1.operator.openshift.io cluster --type=json \
 -p="[{'op': 'add', 'path': '/spec/plugins', 'value':[odf-client-console]}]"
```

## Mitigation

Update ocs-client-operator on ODF Client cluster to be on same major and minor
version as odf-operator on ODF Provider cluster.
