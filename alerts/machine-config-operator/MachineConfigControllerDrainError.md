# MCCDrainError

## Meaning 

Alerts the user to when the Machine Config Operator 
failed to drain a node. Preventing the MCO from restarting 
the node. 

>Alerts the user to a failed node drain. Always triggers when the failure
>happens one or more times.

This alert will fire as a warning when: 

- If the mco is unable to drain the node for an hour
- When the failure occurs one or more times in a node 

## Impact 

If the mco fails to drain a node it will be unable to reboot the node. 
Preventing any changes via MachineConfig or a cluster upgrade to be applied. 
This will prevent a cluster upgrade from progressing/completing.

## Diagnosis 

If a node is failing to draining the first place the alert itself recommends to check the machine-config-controller pod inside the openshift-machine-config-operator namespace. 
The machine-config-controller pod is the central point of management for incoming updates to 
machines.

For the commands below replace the $CONTROLLERPOD variable with the name of your own machine-config-controller pod name. 

```console
oc -n openshift-machine-config-operator logs $CONTROLLERPOD -c machine-config-controller
```

When a drain is initiated the following log will be recorded.

```console
1 drain_controller.go:173] node xxxxxxx-xxxxxx-xxxxxx-xxxxx: initiating drain
```
The controller will then log each and every pod drained from the node. 

```console
  1 drain_controller.go:173] node xxxxxxx-xxxxxx-xxxxxx-xxxxx: Evicted pod <namespace>/<pod-name>
```

If the controller is unable to drain a pod after 1m30s then it will fire an error log.

```console
1 drain_controller.go:173] node xxxxxxx-xxxxxx-xxxxxx-xxxxx: Drain failed. Waiting 1 minute then retrying. Error message from drain: error when waiting for pod "xxxx-xxxx-xx" in namespace "xxxxxxx" to terminate: global timeout reached: 1m30s
```

If the drain continues to fail then another failure log will fire.

```console
 1 drain_controller.go:173] node xxxxxxx-xxxxxx-xxxxxx-xxxxx: Drain failed. Drain has been failing for more than 10 minutes. Waiting 5 minutes then retrying. Error message from drain: error when waiting for pod "xxxx-xxxx-xx" in namespace "xxxxxxx" to terminate: global timeout reached: 1m30s
```

Finally if 1 hour has passed and the drain is still failing the following logs will be fired and the node will be marked degraded.

```console
1 drain_controller.go:352] node xxxxxxx-xxxxxx-xxxxxx-xxxxx: drain exceeded timeout: 1h0m0s. Will continue to retry.
```

```console
1 status.go:126] Degraded Machine: xxxxxxx-xxxxxx-xxxxxx-xxxxx: and Degraded Reason: failed to drain node: xxxxxxx-xxxxxx-xxxxxx-xxxxx after 1 hour. Please see machine-config-controller logs for more information
```

The machine-controller logs will explicitly tell you what pods are failing to drain succesfully. Looking at the pods named in the logs is where you should start. 

Common reasons a node cannot drain a pod successfully are:

- The pod has a PodDisruptionBudget (PBD) set preventing the controller from deleting the pod.
- The pod has storage attached and the kubelet is unable to unmount the storage succesfully.
- The pod deletion may be prevented by a webhook if it's configured to target UPDATE operations or the webhook is failing to be called by the kube-apiserver.
- The pod deletion could be stuck because the finalizers set in the pod are preventing it from terminating succesfully.

## Mitigation 

- If a PDB is causing the failure you can temporarily patch it to set the `minAvailable` to 0 so it can scale the pod down successfully. Then patch it to the old value after the upgrade completes.

```console
oc patch pdb $PDB -n $NS --type=merge -p '{"spec":{"minAvailable":0}}'
```
- If a webhook is preventing the pod deletion you can check the `openshift-kube-apiserver` logs to see exactly what webhook is preventing deletion.

```console
Failed calling webhook, failing open vault.hashicorp.com: failed calling webhook "webhook_name": failed to call webhook: Post "https://hashi-vault-agent-injector-svc.vault-injector.svc:443/mutate?timeout=30s": context canceled2022-10-11T19:13:53.345381984Z E1011 19:13:53.345348      16 dispatcher.go:184] failed calling webhook "webhook_name": failed to call webhook: Post "https://name-injector-svc.vault-injector.svc:443/mutate?timeout=30s": context canceled
```
Then check if it's a `mutating` or `validating` webhook.
```console
$ oc get validatingwebhookconfiguration
$ oc get mutatingwebhookconfiguration
```
Then backup and delete the webhook 
```console
$ oc get validatingwebhookconfiguration/<webhook_name> -o yaml > webhook.yaml
$$ oc delete <webhook_type> <webhook_name>         
```
This should allow the drain to continue and if the webhook does not come back automatically you can recreate it via the backup.

- If a pod cannot amount storage you either need to troubleshoot why it's failing such as if it's a NFS storage then it could be a network issue with the storage server. If you are okay with possible data loss you can force delete the pod and immediately remove resources from the API and bypass graceful deletion. 

```console
$ oc delete pod <pod-name> --force=true --grace-period=0
```
- If the pod is stuck in terminating status you can try patching the finalizers to null.

```console
oc patch pod <pod_name> -p '{"metadata":{"finalizers":null}}' -n <namespace_name>
```
or you can try force deleting the pod as described above. 



