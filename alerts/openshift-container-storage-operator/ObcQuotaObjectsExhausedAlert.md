# ObcQuotaObjectsExhausedAlert

## Meaning

ObjectBucketClaim has crossed the limit set by the quota(objects) and
will be read-only now.

## Impact

Application won't be able to do any transaction through the OBC and will be stalled.

## Diagnosis

Alert message will indicate which OBC has reached the object quota limit.
Look at the deployments attached to the OBC and
see what all apps are using/filling-up the OBC.

## Mitigation

Immediately increase the quota for the OBC, specified in the alert details.
We can increase quota option on OBC by using the `maxObjects` and
`maxSize` options in the ObjectBucketClaim CRD

```yaml
apiVersion: objectbucket.io/v1alpha1
kind: ObjectBucketClaim
metadata:
  name: <obc_name>
  namespace: <namespace>
spec:
  bucketName: <name_of_the_backend_bucket>
  storageClassName: <name_of_storage_class>
  additionalConfig:
    maxObjects: "1000" # sets limit on no of objects this obc can hold
    maxSize: "2G" # sets max limit for the size of data this obc can hold
```

