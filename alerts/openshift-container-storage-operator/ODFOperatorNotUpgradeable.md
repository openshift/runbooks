# ODFOperatorNotUpgradeable

## Meaning

The ODF Operator has been marked as not upgradeable. This indicates that there are conditions preventing a safe upgrade to the next ODF version.

## Impact

The ODF cluster cannot be upgraded to the next version until the underlying issue is resolved or the condition is overridden. This may delay critical updates, security patches, or new features.

## Diagnosis

### Step 1: Check the OperatorCondition CR

Get the OperatorCondition CR for the ODF operator to identify the specific issue:

```bash
oc get operatorconditions -n openshift-storage
```

Get detailed information about the odf operator condition:

```bash
oc get operatorconditions -n openshift-storage <odf-operator-condition-name> -o yaml
```

Look for the `Upgradeable` condition in `spec.conditions` or `status.conditions`. The condition will contain:
- 
- `status`: "False" indicates not upgradeable
- `reason`: The reason why the operator is not upgradeable
- `message`: Detailed description of the issue

Example output:
```yaml
spec:
  conditions:
  - lastTransitionTime: "2025-11-03T13:02:13Z"
    message: 'CephCluster health is HEALTH_WARN. Details: [MON_DOWN: HEALTH_WARN - 1/3 mons down, quorum b,c]'
    reason: CephClusterNotHealthOK
    status: "False"
    type: Upgradeable
```

### Step 2: Check the underlying component

Based on the `reason` and `message` from the OperatorCondition, investigate the specific component that is causing the issue. Common reasons include:
- OCP Version
- OCP Upgrade status
- Ceph cluster health issues (HEALTH_WARN or HEALTH_ERR)
- Storage capacity issues
- Component availability issues

Check the specific component's status and logs to understand the root cause.

## Mitigation

### Option 1: Resolve the underlying issue (Recommended)

Address the specific issue mentioned in the OperatorCondition message. For example:
- If Ceph health is HEALTH_WARN, resolve the Ceph issues first
- If a component is down, restore the component
- If there are capacity issues, add more storage or clean up unused data

Once the issue is resolved, the operator will automatically mark itself as upgradeable again.

### Option 2: Override the condition (Use with caution)

If you have verified that the issue is non-critical or you need to proceed with the upgrade despite the condition, you can override the upgradeable status.

**Warning**: Overriding the upgradeable condition bypasses safety checks and may lead to upgrade failures or data loss. Only override if you fully understand the risks.

To override:

```bash
oc patch operatorconditions -n openshift-storage <operator-condition-name> --type=merge -p '
spec:
  overrides:
  - type: Upgradeable
    status: "True"
    reason: ManualOverride
    message: "Manually overriding upgradeable condition"
'
```

Note: `spec.overrides` takes precedence over `spec.conditions` according to the OLM behavior.

## Additional Resources

- [OLM Operator Conditions Documentation](https://olm.operatorframework.io/docs/advanced-tasks/communicating-operator-conditions-to-olm/)
