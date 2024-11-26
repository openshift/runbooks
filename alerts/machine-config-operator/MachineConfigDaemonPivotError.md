# MCDPivotError

## Meaning

This alert is triggered when the Machine Config Daemon (MCD)
detects an error when trying to pivot the
operating system (OS) image to a new version or a kernel
change. If the MCD is unable to complete
the pivot or change within 2 minutes the alert
will fire.

## Impact

If the MCD is unable to update the OS and
finish pivoting then this can prevent an
OpenShift upgrade from completing.
This can leave the cluster in a
unstable state and further affect the operation of the cluster.

## Diagnosis

When a node pivots to update the OS image,
the `rpm-ostree` service logs any
actions taken to the `machine-config-daemon-*` pod.

If the MCD fails to
update or pivot a node's OS or kernel,
the first logs that you should
check are the `machine-config-daemon-*`
pod logs for the cluster.

For the following command, replace the $DAEMONPOD variable
with the name of your own machine-config-daemon-* pod name.
That is scheduled on the node expriencing the error.

```console
oc logs -f -n openshift-machine-config-operator $DAEMONPOD -c machine-config-daemon
```
When a pivot is occuring the following will be logged.

```console
I1126 17:15:38.991090    3069 rpm-ostree.go:243] Executing rebase to quay.io/my-registry/custom-image@blah
```
The MCD will log its attempt to pivot.

```console
I1126 17:15:38.991115    3069 update.go:2618] Running: rpm-ostree rebase --experimental ostree-unverified-registry:quay.io/my-registry/custom-image@blah
Pulling manifest: ostree-unverified-registry:quay.io/my-registry/custom-image@blah
```
If the MCD fails to update the OS it will rollback.

```console
E1126 17:16:07.549890    3069 writer.go:226] Marking Degraded due to: failed to update OS to quay.io/my-registry/custom-image@blah : error running rpm-ostree rebase --experimental ostree-unverified-registry:quay.io/my-registry/custom-image@blah: error: Creating importer: Failed to invoke skopeo proxy method OpenImage: remote error: invalid reference format
: exit status 1
```

The MCD will then mark the node degraded.

```console
2024-06-20T17:56:33.930222523Z E0620 17:56:33.930211 4168959 writer.go:135] Marking Degraded due to: failed to update OS to quay.io/openshift-release-dev/ocp-v4.0-art-dev@sha256:368d9b871acb9fc29eea6a4f66e42894677594e91834958c015ed15c03ebe79e : error running rpm-ostree rebase --experimental /run/mco-machine-os-content/os-content-501945936/srv/repo:afdc646803e2d9d774fbf3429cf91de6222e45a85ceabbafe4ee78aca74c2d7b --custom-origin-url pivot://quay.io/openshift-release-dev/ocp-v4.0-art-dev@sha256:368d9b871acb9fc29eea6a4f66e42894677594e91834958c015ed15c03ebe79e --custom-origin-description Managed by machine-config-operator: error: Timeout was reached
```

The node will continue trying to pivot despite the error.
You should start troubleshooting
by examining the main error message and the
stated reason it gives for not being able to pivot. The following are
common reasons a pivot can fail.

- The rpm-ostree service is unable to
pull the image from quay succesfully.
- There are issues with the rpm-ostree service itself such as
being unable to start, or unable to build the OsImage folder,
unable to pivot from the current configuration.
- The rpm-ostree service gets stuck and client is unable to finish
the transaction and gets stuck in a loop as a result.
- The ostree-finalize-staged.service
  (responsible for rolling out pending updates) is having issues.
- Networking issues on the node prevent quay from
being reachable, such as an interface being
down, a firewall blocking the connection, the proxy, etc.

## Mitigation

- If the machine-config-daemon pod is
  logging the following errors from rpm-ostree,
  it is likely that rpm-ostree is stuck in a
  loop and cannot finish the transaction
  or there are issues related to the rpm-ostree service itself.

  ```console
   error: Transaction in progress: (null)
  ```
  You should restart the `rpm-ostreed`
  service on the node that are failing.

  ```console
  $ oc debug node <Node-name>
  # chroot /host
  # systemctl restart rpm-ostreed
  ```
- If the machine-config-daemon pod is
  logging the following errors, it is likely due to networking errors related to
  timeouts or Internal Server errors in the daemon or controller pods.

     ```console
    received unexpected HTTP status: 500 Internal Server Error
     ```

  - Make sure that `quay.io` and
    it's subdomains are whitelisted by your firewall and proxy.
    You can test manual pulls with `podman pull`.

    ```console
    $ podman pull quay.io/startx/couchbase:ubi8 --log-level debug
    ```
    - You should also validate that the internal
     network on the node and the external network
     are healthy and there are no blockers
     such as a firewall or switch issues.

- You can also try to force a manual upgrade to the new image
  if the pivot is stuck.

  ```console
  -delete the currentconfig(rm /etc/machine-config-daemon/currentconfig)
  -create a forcefile(touch /run/machine-config-daemon-force) to retry the OS upgrade.
  ```

  - If you want to dive deeper you can also use the
    rpm-ostree command line tool to attempt to force a pivot

    ```console
     rpm-ostree rebase --experimental ostree-unverified-registry:quay.io/openshift-release-dev/ocp-v4.0-art-dev@sha256:xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx
    ```

- If these troubleshooting steps do not fix the problem, gather
  the following logs, which will be helpful
  for identifying the issue.

    ```console
    # rpm-ostree status
    # journalctl -u ostree-finalize-staged
    # journalctl -b -1 -u rpm-ostreed.service
    ```