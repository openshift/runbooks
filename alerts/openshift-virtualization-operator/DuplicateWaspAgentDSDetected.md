# DuplicateWaspAgentDSDetected [Deprecated]

This alert is deprecated. You can safely ignore or silence it.

## Meaning
wasp-agent is a node-local agent that enables swap for burstable QoS pods.
It mimics the behavior of kubelet swap feature.
wasp-agent deployment consists of Daemonset, serivce account,
role binding, privileged SCC. Wasp-agent currently deployed automatically
by the HCO operator when the memory overcommit percentage is set to a
value higher than 100%. In the past wasp-agent was deployed manually
according to the user-guide. It may happen that the previous wasp-agent
deployment still exists. This alert will be fired when such scenario
detected, so the user would be able to remove the old deployment
according to the steps in the mitigation section.

## Impact
When memory overcommit percentage is set to a value higher than 100% an
automatic deployment of wasp-agent will start. This may lead to two
wasp-agent deployments run at the same time on the cluster, which is undesired.

## Diagnosis
Identify the duplicate wasp-agent daemonset :
   ```bash
   oc get ds wasp-agent -n wasp
   ```

## Mitigation

1. Delete the wasp-agent Daemonset
   ```bash
   oc delete ds wasp-agent -n wasp
   ```
2. Delete the wasp-agent service account
   ```bash
   oc delete sa wasp -n wasp
   ```
3. Delete the wasp-agent cluster role binding
   ```bash
   oc delete clusterrolebinding wasp
   ```
4. Remove the `wasp` service account from privileged SCC
   ```bash
   oc adm policy remove-scc-from-user -n wasp privileged -z wasp
   ```
5. Remove the `wasp` project
   ```bash
   oc delete project wasp
   ```


## Additional notes
* [wasp-agent user-guide](https://docs.redhat.com/en/documentation/openshift_container_platform/4.20/html-single/virtualization/index#virt-configuring-higher-vm-workload-density)

If you cannot resolve the issue, log in to the
[Customer Portal](https://access.redhat.com) and open a support case,
attaching the artifacts gathered during the diagnosis procedure.
