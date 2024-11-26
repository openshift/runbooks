# MCDPivotError

## Meaning

Alerts the User when the Machine Config Daemon
is either unable to or detecs an error when
trying to pivot OS image to a new version.

>Alerts the user when an error is detected upon pivot.
This triggers if the pivot errors are above zero for 2 minutes.

This alert will occur if

- The mcd_pivot_errors_total is
greater then 0 for 2 minutes.
- mcd_pivot_errors_total is a guage counter
that can increment if there are errors
during an update.
- These errors can occur during an OS update
or kernel change.

## Impact

If the MCD is unable to update the OS and
finish pivoting then this can prevent an
OpenShift upgrade from completing.
This can leave the cluster an a
unstable state and further affect the operation of the cluster.

## Diagnosis

When a node pivots to update the OS image.
This is done via the `rpm-ostree` service
that logs any actions taken to the `machine-config-daemon-*` pod

If the Machine Config Daemon is failing to
update or pivot the nodes OS or Kernel
then the first logs that should
be checked are the `machine-config-daemon-*`
pod logs for the cluster

For the following command, replace the $DAEMONPOD variable
with the name of your own machine-config-daemon-* pod name.
To the erroring node.

```console
oc logs -f -n openshift-machine-config-operator $DAEMONPOD -c machine-config-daemon
```
When a pivot is occuring the following will be logged.

```console
I1126 17:15:38.991090    3069 rpm-ostree.go:243] Executing rebase to quay.io/my-registry/custom-image@blah
```
It will log it's attempt to pivot  

```console
I1126 17:15:38.991115    3069 update.go:2618] Running: rpm-ostree rebase --experimental ostree-unverified-registry:quay.io/my-registry/custom-image@blah
Pulling manifest: ostree-unverified-registry:quay.io/my-registry/custom-image@blah
```
If it fails to update the OS it will rollback

```console
E1126 17:16:07.549890    3069 writer.go:226] Marking Degraded due to: failed to update OS to quay.io/my-registry/custom-image@blah : error running rpm-ostree rebase --experimental ostree-unverified-registry:quay.io/my-registry/custom-image@blah: error: Creating importer: Failed to invoke skopeo proxy method OpenImage: remote error: invalid reference format
: exit status 1
```

It will then mark the node degraded

```console
2024-06-20T17:56:33.930222523Z E0620 17:56:33.930211 4168959 writer.go:135] Marking Degraded due to: failed to update OS to quay.io/openshift-release-dev/ocp-v4.0-art-dev@sha256:368d9b871acb9fc29eea6a4f66e42894677594e91834958c015ed15c03ebe79e : error running rpm-ostree rebase --experimental /run/mco-machine-os-content/os-content-501945936/srv/repo:afdc646803e2d9d774fbf3429cf91de6222e45a85ceabbafe4ee78aca74c2d7b --custom-origin-url pivot://quay.io/openshift-release-dev/ocp-v4.0-art-dev@sha256:368d9b871acb9fc29eea6a4f66e42894677594e91834958c015ed15c03ebe79e --custom-origin-description Managed by machine-config-operator: error: Timeout was reached
```

The node will continue trying to pivot despite this.
You should start by examining the main error message and the
stated reason it gives for not being able to pivot. Common reasons a pivot
can fail.

- rpm-ostree is unable to
pull the image from quay succesfully.
- Issues with the rpm-ostree service itself such as
being unable to start, or unable to build the OsImage folder,
unable to pivot from the current configuration.
- rpm-ostree service gets stuck and client is unable to finish
transaction and gets stuck in a loop as a result.
- The ostree-finalize-staged.service
  (responsible for rolling out pending updates) is having issues.
- Networking issues on the node prevent quay from
being reachable such as an interface being
down, a firewall blocking the connection, the proxy, etc.

## Mitigation

- If the machine-config-daemon pod is
  logging the following errors from rpm-ostree
  that relate to transactions being stuck or
  possible issues with `rpm-ostreed` such as.

  ```console
   error: Transaction in progress: (null)
  ```

  Then likely rpm-ostree is stuck in a
  loop and cannot finish the transaction
  or there's issues related to the service itself.
  You should restart the `rpm-ostreed`
  service on the node that are failing.

  ```console
  $ oc debug node <Node-name>
  # chroot /host
  # systemctl restart rpm-ostreed
  ```
- If you see the networking errors related to
  timeoutes or Internal Server errors in the daemon or controller pods.
  For example.


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
     as well are healthy and there are no blockers.
     Such as a firewall or switch issues.

- You can also try to force a manual upgrade to the new image.
  If the pivot is stuck

  ```console
  -delete the currentconfig(rm /etc/machine-config-daemon/currentconfig)
  -create a forcefile(touch /run/machine-config-daemon-force) to retry the OS upgrade.
  ```

  - If you want to dive deeper you can also use the
    rpm-ostree command line tool to try and force a pivot

    ```console
     rpm-ostree rebase --experimental ostree-unverified-registry:quay.io/openshift-release-dev/ocp-v4.0-art-dev@sha256:xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx
    ```

- If all else is still failing then gathering
  the following logs will be helpful for a
  Red Hat support engineer to identifying the issue.

    ```console
    # rpm-ostree status
    # journalctl -u ostree-finalize-staged
    # journalctl -b -1 -u rpm-ostreed.service
    ```