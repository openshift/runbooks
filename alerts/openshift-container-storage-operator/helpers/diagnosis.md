# Diagnosis cluster

## Verify cluster access

Check the output to ensure you are in the correct context for the cluster
mentioned in the alert. If not, please change context and proceed.

List clusters you have permission to access:

```bash
    ocm list clusters
```

From the list above, find the cluster id of the cluster named in the alert.
If you do not see the alerting cluster in the list above please refer
[Effective communication with SRE Platform](https://red.ht/srep-comms)

Create a tunnel through backplane by providing SSH key passphrase:

```bash
    ocm backplane tunnel <cluster_id>
```

In a new tab, login to target cluster using backplane by providing 2FA:

```bash
    ocm backplane login <cluster_id>
```

## Check Alerts

Set port-forwarding for alertmanager:

```bash
    oc port-forward alertmanager-managed-ocs-alertmanager-0 9093 -n
    openshift-storage
```

Check all alerts

```bash
    curl http://localhost:9093/api/v1/alerts | jq '.data[] | select( .labels.alertname) | { ALERT: .labels.alertname, STATE: .status.state}'
```

## Check OCS Ceph Cluster Health

You may directly check OCS Ceph Cluster health by using the rook-ceph toolbox.

Step 1: [Check and document ceph cluster health](cephCLI.md):

Step 2: From the bash command prompt, run the following and capture the output.

```bash
    ceph status
    ceph osd status
    exit
```

If `ceph status` is not in **HEALTH\_OK**, please look at the Troubleshooting
 section to resolve issue.

## Check Worker Node Status

If `ceph status` is not **HEALTH\_OK** and all unhealthy components are related
 to a particular node, then the following steps will help identify if the
 underlying infrastructure is at fault.

Step 1: Check Node Health:

```bash
    $ oc get nodes -o wide

    NAME                                         STATUS   ROLES          AGE   VERSION           INTERNAL-IP    EXTERNAL-IP   OS-IMAGE                                                        KERNEL-VERSION                 CONTAINER-RUNTIME
    ip-10-0-128-234.eu-west-2.compute.internal   Ready    master         89m   v1.23.5+3afdacb   10.0.128.234   <none>        Red Hat Enterprise Linux CoreOS 410.84.202206080346-0 (Ootpa)   4.18.0-305.49.1.el8_4.x86_64   cri-o://1.23.3-3.rhaos4.10.git5fe1720.el8
    ip-10-0-129-200.eu-west-2.compute.internal   Ready    infra,worker   66m   v1.23.5+3afdacb   10.0.129.200   <none>        Red Hat Enterprise Linux CoreOS 410.84.202206080346-0 (Ootpa)   4.18.0-305.49.1.el8_4.x86_64   cri-o://1.23.3-3.rhaos4.10.git5fe1720.el8
    ip-10-0-133-54.eu-west-2.compute.internal    Ready    worker         84m   v1.23.5+3afdacb   10.0.133.54    <none>        Red Hat Enterprise Linux CoreOS 410.84.202206080346-0 (Ootpa)   4.18.0-305.49.1.el8_4.x86_64   cri-o://1.23.3-3.rhaos4.10.git5fe1720.el8
    ip-10-0-144-158.eu-west-2.compute.internal   Ready    master         89m   v1.23.5+3afdacb   10.0.144.158   <none>        Red Hat Enterprise Linux CoreOS 410.84.202206080346-0 (Ootpa)   4.18.0-305.49.1.el8_4.x86_64   cri-o://1.23.3-3.rhaos4.10.git5fe1720.el8
    ip-10-0-153-73.eu-west-2.compute.internal    Ready    infra,worker   66m   v1.23.5+3afdacb   10.0.153.73    <none>        Red Hat Enterprise Linux CoreOS 410.84.202206080346-0 (Ootpa)   4.18.0-305.49.1.el8_4.x86_64   cri-o://1.23.3-3.rhaos4.10.git5fe1720.el8
    ip-10-0-159-144.eu-west-2.compute.internal   Ready    worker         82m   v1.23.5+3afdacb   10.0.159.144   <none>        Red Hat Enterprise Linux CoreOS 410.84.202206080346-0 (Ootpa)   4.18.0-305.49.1.el8_4.x86_64   cri-o://1.23.3-3.rhaos4.10.git5fe1720.el8
    ip-10-0-168-106.eu-west-2.compute.internal   Ready    infra,worker   66m   v1.23.5+3afdacb   10.0.168.106   <none>        Red Hat Enterprise Linux CoreOS 410.84.202206080346-0 (Ootpa)   4.18.0-305.49.1.el8_4.x86_64   cri-o://1.23.3-3.rhaos4.10.git5fe1720.el8
    ip-10-0-173-205.eu-west-2.compute.internal   Ready    master         89m   v1.23.5+3afdacb   10.0.173.205   <none>        Red Hat Enterprise Linux CoreOS 410.84.202206080346-0 (Ootpa)   4.18.0-305.49.1.el8_4.x86_64   cri-o://1.23.3-3.rhaos4.10.git5fe1720.el8
    ip-10-0-175-99.eu-west-2.compute.internal    Ready    worker         83m   v1.23.5+3afdacb   10.0.175.99    <none>        Red Hat Enterprise Linux CoreOS 410.84.202206080346-0 (Ootpa)   4.18.0-305.49.1.el8_4.x86_64   cri-o://1.23.3-3.rhaos4.10.git5fe1720.el8
```

If any nodes are not ready/scheduable, then continue to Step 2.

Step 2: Inspect Node Events:

```bash
    oc get events -n default | grep NODE_NAME
```

Example:

```bash
    $ oc get events -n default | grep 10-0-159-144

    57m         Normal    ConfigDriftMonitorStopped                          node/ip-10-0-159-144.eu-west-2.compute.internal      Config Drift Monitor stopped
    57m         Normal    NodeNotSchedulable                                 node/ip-10-0-159-144.eu-west-2.compute.internal      Node ip-10-0-159-144.eu-west-2.compute.internal status is now: NodeNotSchedulable
    57m         Normal    Cordon                                             node/ip-10-0-159-144.eu-west-2.compute.internal      Cordoned node to apply update
    57m         Normal    Drain                                              node/ip-10-0-159-144.eu-west-2.compute.internal      Draining node to update config.
    11m         Normal    OSUpdateStarted                                    node/ip-10-0-159-144.eu-west-2.compute.internal
    11m         Normal    OSUpdateStaged                                     node/ip-10-0-159-144.eu-west-2.compute.internal      Changes to OS staged
    11m         Normal    PendingConfig                                      node/ip-10-0-159-144.eu-west-2.compute.internal      Written pending config rendered-worker-c8a49ffa8d6d6ee43a4e4ae5b5c7f60f
    11m         Normal    Reboot                                             node/ip-10-0-159-144.eu-west-2.compute.internal      Node will reboot into config rendered-worker-c8a49ffa8d6d6ee43a4e4ae5b5c7f60f
    10m         Normal    NodeNotReady                                       node/ip-10-0-159-144.eu-west-2.compute.internal      Node ip-10-0-159-144.eu-west-2.compute.internal status is now: NodeNotReady
    10m         Normal    Starting                                           node/ip-10-0-159-144.eu-west-2.compute.internal      Starting kubelet.
```

Look for any events similar to the above example which may indicate the node
is undergoing maintainence.

## Further info

### OpenShift Data Foundation Dedicated Architecture

Red Hat OpenShift Data Foundation Dedicated (ODF Dedicated) is deployed in
converged mode on OpenShift Dedicated Clusters by the OpenShift Cluster Manager
 add-on infrastructure.

Related Links

* [ODF Dedicated Converged Add-on Architecure](https://docs.google.com/document/d/1ISEY16OfsvEPmlJEjEwPvDvDs0KyNzgl369A-V6-GRA/edit#heading=h.mznotzn8pklp)

* [ODF Product Architecture](https://access.redhat.com/documentation/en-us/red_hat_openshift_container_storage/4.6/html/planning_your_deployment/ocs-architecture_rhocs)


Check the links to identify the errors.

1) [https://access.redhat.com/documentation/en-us/red\_hat\_ceph\_storage/4/html/troubleshooting\_guide/troubleshooting-ceph-osds#common-ceph-osd-error-messages-in-the-ceph-logs\_diag](https://access.redhat.com/documentation/en-us/red_hat_ceph_storage/4/html/troubleshooting_guide/troubleshooting-ceph-osds#common-ceph-osd-error-messages-in-the-ceph-logs_diag)

2) [https://access.redhat.com/documentation/en-us/red\_hat\_ceph\_storage/4/html/troubleshooting\_guide/troubleshooting-ceph-placement-groups#inconsistent-placement-groups\_diag](https://access.redhat.com/documentation/en-us/red_hat_ceph_storage/4/html/troubleshooting_guide/troubleshooting-ceph-placement-groups#inconsistent-placement-groups_diag)
