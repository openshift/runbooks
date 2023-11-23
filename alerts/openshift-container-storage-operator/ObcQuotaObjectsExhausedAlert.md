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
For quota help please visit: https://access.redhat.com/articles/6541861

## Mitigation

Immediately increase the quota for the OBC, specified in the alert details.

