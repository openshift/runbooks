# KubeletHealthState

## Meaning

The Machine Config Daemon (MCD) pods will run a
health check on kubelet every 30 seconds by running
`curl http://localhost:10248/healthz` on the nodes.
If the kubelet is reachable then it will
return `Ok`, if the endpoint does not return this value
then the MCD will increment the `mcd_kubelet_state`
by 1. The counter will continue to increment for each
failed health check, and fire the alert if the value
of `mcd_kubelet_state` > 2.

## Impact

The `kubelet` is a primary node agent that is
required by each node in a cluster to operate.
It facilitates communication, container management,
and is responsible for node administration. If
the kubelet is not running or reachable on a node,
the node will go into `NotReady` status and can become
inoperable.

## Diagnosis

The health check is initiated and ran by the
machine-config-daemon pod scheduled on each node.
Each pod will log the following when starting up.

```console
I0428 02:23:13.498282   11418 daemon.go:1250] Enabling Kubelet Healthz Monitor
```

The health check will run every 30 seconds,
this check also has a timeout value of 30 seconds

```console
kubeletHealthzPollingInterval = 30 * time.Second
kubeletHealthzTimeout         = 30 * time.Second
```
If the curl does not return `Ok` within that
30 seconds then it's considered
a failure. The machine-config-daemon pod will
log any failure to `stdout`

```console
W0917 11:21:46.203938 1688435 daemon.go:662] Failed kubelet health check: Get http://localhost:10248/healthz: dial tcp: lookup localhost on <node-ip>:53: no such host
W0917 11:21:46.204070 1688435 daemon.go:596] Got an error from auxiliary tools: kubelet health failure threshold reached
```

Depending on the state of the kubelet the node
may be completely inaccessible via the console
or `oc` commands. You may need to login via
external console into the node via root
or SSH directly into the node.

Once in the node you should check the status of
the kubelet service

```console
# sudo systemctl status kubelet.service 
```

Depending on the state of the kubelet it
could be `active`, `inactive` or `failed`.
The command will also display the latest
journal logs from that service. If you
want to see the full journal logs for the
kubelet service you can run:

```console
# journalctl -u kubelet 
```

If the kubelet is failing or inactive the
journal logs are the best place to start
troubleshooting. There could be multiple
reasons the kubelet is failing. The journal
logs will give you the exact reasons why.

If the kubelet service is active but the
health check is still failing then, you'll
need to deep dive into why it can't resolve
the curl command.

Check that kubelet is listening on port 10248.

```console
# netstat -antpu | grep 10248
```
You can try restarting kubelet.service or the
underlying node itself.

```console
# systemctl restart kubelet.service
```
```console
# systemctl reboot
```
You should also check the `/etc/hosts` file
and make sure it has the `localhost` and
`localhost.localdomain` entries for both
ipv4 and ipv6

```console
sh-5.1# cat etc/hosts
127.0.0.1   localhost localhost.localdomain localhost4 localhost4.localdomain4
::1         localhost localhost.localdomain localhost6 localhost6.localdomain6
```

## Mitigation

The general state of the node
affects the state of the kubelet service.
You should make sure the nodes storage disks
are sufficiently sized for your planned cluster size.
That the nodes are not hitting
any network issues and that the network speed is sufficient.
That the nodes have sufficient CPU and memory allocated
for your planned workloads on this cluster.

OpenShift clusters are designed to have high
availability. Control plane nodes can maintain
the ETCD database so long as two control
plane nodes are working. Worker nodes can take on
more workloads if another worker node becomes
unavailable. You [can add, remove, and replace nodes](https://docs.redhat.com/en/documentation/openshift_container_platform/4.20/html/nodes/working-with-nodes)
to add stability to a cluster. OpenShift clusters
should be planned with this in mind in order to take
full advantage of OpenShift features.


