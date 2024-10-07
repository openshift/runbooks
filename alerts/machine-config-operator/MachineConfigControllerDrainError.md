# MCCDrainError

## Meaning 

Alerts the user to when the Machine Config Operator (MCO) 
fails to drain a node, which prevents the MCO from restarting.

>Alerts the user to a failed node drain. Always triggers when the failure
>happens one or more times.

This alert will fire as a warning when: 

- The MCO is unable to drain the node for an hour
- The failure occurs one or more times in a node 

## Impact 

If the MCO fails to drain a node, it will be unable to reboot the node, which prevents any changes to the cluster through a MachineConfig and prevents a cluster upgrade. If the MCO fails to drain a node during an upgrade, the upgrade will not be able to progress/complete.
## Diagnosis 

If a node fails to drain, first check the `machine-config-controller` pod inside the `openshift-machine-config-operator` namespace by using the following command. 
The `machine-config-controller` pod is the central point of management for incoming updates to 
machines.

For the following command, replace the $CONTROLLERPOD variable with the name of your own `machine-config-controller` pod name. 

```console
oc -n openshift-machine-config-operator logs $CONTROLLERPOD -c machine-config-controller
```
When the MCO starts draining a node, the Machine Config Controller (MCC) records the following log entry:

```console
1 drain_controller.go:173] node xxxxxxx-xxxxxx-xxxxxx-xxxxx: initiating drain
```
The MCC then logs the name of each pod that is drained from the node.

```console
  1 drain_controller.go:173] node xxxxxxx-xxxxxx-xxxxxx-xxxxx: Evicted pod <namespace>/<pod-name>
```

If the MCO/MCC is unable to drain a pod after 1m30s, the MCC logs the following error message:

```console
1 drain_controller.go:173] node xxxxxxx-xxxxxx-xxxxxx-xxxxx: Drain failed. Waiting 1 minute then retrying. Error message from drain: error when waiting for pod "xxxx-xxxx-xx" in namespace "xxxxxxx" to terminate: global timeout reached: 1m30s
```

If the drain continues to fail, the MCC logs the a second error message:

```console
 1 drain_controller.go:173] node xxxxxxx-xxxxxx-xxxxxx-xxxxx: Drain failed. Drain has been failing for more than 10 minutes. Waiting 5 minutes then retrying. Error message from drain: error when waiting for pod "xxxx-xxxx-xx" in namespace "xxxxxxx" to terminate: global timeout reached: 1m30s
```

After one hour has passed and the drain is still failing, the MCC logs the following error messages and the node is marked degraded.

```console
1 drain_controller.go:352] node xxxxxxx-xxxxxx-xxxxxx-xxxxx: drain exceeded timeout: 1h0m0s. Will continue to retry.
```

```console
1 status.go:126] Degraded Machine: xxxxxxx-xxxxxx-xxxxxx-xxxxx: and Degraded Reason: failed to drain node: xxxxxxx-xxxxxx-xxxxxx-xxxxx after 1 hour. Please see machine-config-controller logs for more information
```

The MCC logs explicitly list which pods are failing to drain. You should start examining the listed pods for the problem that is causing the drain to fail. Common reasons why a node cannot drain a pod include the following conditions:

- The pod has a PodDisruptionBudget (PBD) that prevents the MCO/MCC from deleting the pod.
- The pod has storage attached and the kubelet is unable to unmount the storage.
- The pod has a webhook that is configured to target UPDATE operations or the webhook is not being called by the kube-apiserver.
- The pod has finalizers set in the pod that are preventing it from terminating.

## Mitigation 

- If a PDB is causing the failure you can temporarily patch it to set the `minAvailable` to 0 so it can scale the pod down successfully. Then patch it to the previous value after the upgrade completes.

```console
oc patch pdb $PDB -n $NS --type=merge -p '{"spec":{"minAvailable":0}}'
```
- If a webhook is preventing the pod deletion:

Check the `openshift-kube-apiserver` logs to see exactly what webhook is preventing deletion.

```console
Failed calling webhook, failing open vault.hashicorp.com: failed calling webhook "webhook_name": failed to call webhook: Post "https://hashi-vault-agent-injector-svc.vault-injector.svc:443/mutate?timeout=30s": context canceled2022-10-11T19:13:53.345381984Z E1011 19:13:53.345348      16 dispatcher.go:184] failed calling webhook "webhook_name": failed to call webhook: Post "https://name-injector-svc.vault-injector.svc:443/mutate?timeout=30s": context canceled
```
Check if it's a `mutating` or `validating` webhook.
```console
$ oc get validatingwebhookconfiguration
$ oc get mutatingwebhookconfiguration
```
Backup and delete the webhook 
```console
$ oc get validatingwebhookconfiguration/<webhook_name> -o yaml > webhook.yaml
$ oc delete <webhook_type> <webhook_name>         
```
This should allow the drain to continue and if the webhook does not come back automatically you can recreate it via the backup.

- If a pod cannot unmount storage, troubleshoot why it's failing. For example, if you are using NFS storage, the problem could be a network issue with the storage server. 

Otherwise, if you are comfortable with possible data loss, you can force delete the pod and immediately remove resources from the API and bypass graceful deletion: 

```console
$ oc delete pod <pod-name> --force=true --grace-period=0
```
- If the finalizers are causing the pod to be stuck in terminating status you can try patching the finalizers to null:

```console
oc patch pod <pod_name> -p '{"metadata":{"finalizers":null}}' -n <namespace_name>
```
Or, you can force delete the pod using the previous command. 


