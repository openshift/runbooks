# ObcQuotaBytesExhausedAlert

## Meaning

This is the next stage once we have reached [ObcQuotaObjectsAlert](ObcQuotaObjectsAlert.md).
ObjectBucketClaim has crossed the limit set by the quota(bytes) and will be
read-only now. Increase the quota in the OBC custom resource immediately.

## Impact

OBC has exhausted and reached it's limit.

## Diagnosis

Alert message will clearly indicate which OBC has reached the quota bytes limit.
Look at the deployments attached to the OBC and see what all apps are
using/filling-up the OBC.
For quota help please visit: https://access.redhat.com/articles/6541861

## Mitigation

Need to increase the quota limit immediately for the ObjectBucketClaim custom resource.

