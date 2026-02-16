# OpenShiftUpdateRiskMightApply

The alert `OpenShiftUpdateRiskMightApply` is still a
[Technology Preview feature](https://access.redhat.com/support/offerings/techpreview/)
from OpenShift 4.22. If you see this alert firing, that means the feature set TechPreviewNoUpdate
is enabled on the cluster.

## Meaning
This alert is triggered by the
[cluster-version-operator](https://github.com/openshift/cluster-version-operator)
(CVO) when the Cluster Version Operator detects that
[an update risk](https://docs.redhat.com/en/documentation/openshift_container_platform/4.21/html/updating_clusters/understanding-openshift-updates-1)
has been applied to the cluster for more than 10 minutes.

## Impact
If a risk is associated to an OpenShift version, upgrade to it is not recommended
unless it is accepted by a cluster administrator.

## Diagnosis

### 1. Check the information about the risk

The alert description includes the name of the risk and we can use it to get the
link to the risk

* Check the status of `ClusterVersion/version`:
    ```console
    $ oc get -o jsonpath='{.status.conditionalUpdateRisks[?(@.name == "TO_BE_REPLACED")].url}{"\n"}' clusterversion version
    ```

### 2. Assess if the risk is acceptable
Cluster administrators can accept the risk and update to that particular target
version after evaluating the risk based on the information provided in the linked
URL above.

## Mitigation

### Accept the risk

This is also a Technology Preview feature from OpenShift 4.22.

If a risk is acceptable, then
`OC_ENABLE_CMD_UPGRADE_ACCEPT_RISKS=true oc adm upgrade accept TO_BE_REPLACED`.
When the risks exposed to a version are all accepted, then the upgrade to that
version becomes recommended.

### Wait for the fix

If a risk is not acceptable, we have to wait until a fix to the risk is shipped.
The noise coming from the alert can be
[silenced by a cluster administrator](https://docs.redhat.com/en/documentation/monitoring_stack_for_red_hat_openshift/4.21/html/managing_alerts/managing-alerts-as-an-administrator#silencing-alerts_managing-alerts-as-an-administrator).
