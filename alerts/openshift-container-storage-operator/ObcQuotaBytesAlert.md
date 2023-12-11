# ObcQuotaBytesAlert

## Meaning

The 'ObcQuotaBytesAlert' is triggered when an ObjectBucketClaim (OBC) has
crossed 80% of its quota in bytes.

## Impact

The OBC has reached 80% of its quota, and it will become read-only on reaching
the quota limit. This may impact write operations to the ObjectBucketClaim.

## Diagnosis

Check the usage of the ObjectBucketClaim and its quota to determine the extent
of the breach:

```bash
ocs_objectbucketclaim_info * on (namespace, objectbucket) group_left() (ocs_objectbucket_used_bytes/ocs_objectbucket_max_bytes)
```

## Mitigation

- Review and Increase Quota: Update the quota options on the OBC using the
  following YAML configuration in the CRD:

```yaml
apiVersion: objectbucket.io/v1alpha1
kind: ObjectBucketClaim  
metadata:  
  name: <obc name>  
  namespace: <namespace>  
spec:  
  bucketName:  <name of the backend bucket>  
  storageClassName: <name of storage class>  
  additionalConfig:   
    maxObjects: "1000" # sets limit on the number of objects this OBC can hold  
    maxSize: "2G" # sets the max limit for the size of data this OBC can hold
```
