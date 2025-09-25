# HCP_KubeAPIDown
The  `KubeAPIDown`  alert is triggered in the HCP environment when all Kubernetes API servers are not reachable to the monitoring system for more than 5 minutes.

## Initial troubleshooting steps
1. Start by looking at the dynatrace dashboard

   ```sh
   osdctl dt dash --cluster-id $HCP_CLUSTER_ID
   ```

   Look at the overall situation of this cluster. Check resource usage, request serving node metrics, logs etc.

2. Login to the management cluster mentioned in the alert:

```sh
   ocm backplane login $MC_NAME
  ```
   Set the `HCP_NAMESPACE` variable before running the troubleshooting commands, you will find it in the alert definition
   ```bash
   export HCP_NAMESPACE=<HCP_NAMESPACE>
   ```

3. Check the kube-apiserver pod status and see why it is not coming up.

```sh
oc get po -n $HCP_NAMESPACE -owide | grep kube-apiserver
   kube-apiserver-xxxxxxxxxx-xxxxx           0/5     Pending     0             15d     10.128.202.9     ip-10-0-x-xxx.us-west-2.compute.internal     <none>           <none>
   kube-apiserver-xxxxxxxxxx-xxxxx           5/5     Running     0             14d     10.128.217.7     ip-10-0-xx-xx.us-west-2.compute.internal     <none>           <none>

 ```
   Inspect the details of the kube-apiserver pod:

   ```sh
oc describe po/<kube-api-server pod-name> -n $HCP_NAMESPACE
   ```

   View the logs for the failing kube-apiserver pod to identify errors:

   ```sh
osdctl dt logs <kube-apiserver-pod-name> --namespace $HCP_NAMESPACE
   ```

4. Inspect the Hypershift operator logs for kube-apiserver-related errors:

```sh
   osdctl dt logs --namespace hypershift
```

   _Note_: Refrain yourself from deleting the kube-apiserver pods without finding the proper cause.

   As per the logs find out the error and try to mitigate it. If possible raise a PR to include the troubleshooting steps for the error you've faced for this alert.
