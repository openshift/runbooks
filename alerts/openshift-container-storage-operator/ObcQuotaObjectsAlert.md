# ObcQuotaObjectsAlert

## Meaning

An ObjectBucketClaim object has crossed 80% of the size limit set by the quota(objects)
and will become read-only on reaching the quota limit.

## Impact

OBC has reached 80% of it's limit and soon will get exhausted once reaching
the quota limit.

## Diagnosis

Alert message will clearly indicate which OBC is being filled up fast.
Look at the deployments attached to the OBC and
see what all apps are using/filling-up the OBC.
For quota help please visit: https://access.redhat.com/articles/6541861

## Mitigation

Increase the quota limit for the ObjectBucketClaim custom resource.

