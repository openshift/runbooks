# NodeClockNotSynchronising

## Meaning

The `NodeClockNotSynchronising` alert triggers when an OpenShift worker node (RHCOS) is unable to synchronize its clock with configured time servers. This indicates a critical time synchronization failure that, if left unaddressed, can lead to serious cluster issues.

The alert is specifically triggered when:
- The Prometheus metric `node_timex_sync_status` returns `0` (indicating NTP is not synchronizing)
- The metric `node_timex_maxerror_seconds` is greater than or equal to `16` seconds

This condition indicates that the node's chronyd service is either not running correctly, unable to reach time servers, or experiencing significant synchronization issues.

## Impact

- etcd cluster instability and potential data inconsistency
- Certificate validation failures for OpenShift API communications
- Authentication issues and failed connections due to invalid timestamps
- Kubernetes API server communication failures
- Control plane components may become unavailable
- Pods may fail scheduling or function improperly

Time synchronization failure is a critical issue that can have cascading effects throughout the cluster, potentially impacting all operations.

## Diagnosis

1. Identify the affected node from the alert labels:
```shell
oc get nodes <NODE_NAME> -o wide
```

2. Check time synchronization status on the node:
```shell
oc debug node/<NODE_NAME> -- chroot /host timedatectl status
```

3. Verify chronyd service status (RHCOS uses chronyd):
```shell
oc debug node/<NODE_NAME> -- chroot /host systemctl status chronyd
```

4. Check chrony logs for error messages:
```shell
oc debug node/<NODE_NAME> -- chroot /host journalctl -u chronyd
```

5. Check the chrony sources to see if time servers are reachable:
```shell
oc debug node/<NODE_NAME> -- chroot /host chronyc sources
```

6. Review chrony configuration:
```shell
oc debug node/<NODE_NAME> -- chroot /host cat /etc/chrony.conf
```

7. Check for NTP traffic restrictions:
```shell
oc debug node/<NODE_NAME> -- chroot /host firewall-cmd --list-all
```

8. Verify the Prometheus metrics that triggered the alert:
```shell
oc debug node/<NODE_NAME> -- chroot /host curl -s localhost:9100/metrics | grep node_timex_sync_status
oc debug node/<NODE_NAME> -- chroot /host curl -s localhost:9100/metrics | grep node_timex_maxerror_seconds
```

## Mitigation

1. Chronyd service issues:
   - If the chronyd service is not running, start it:
   ```shell
   oc debug node/<NODE_NAME> -- chroot /host systemctl start chronyd
   ```
   
   - If the chronyd service is running but not synchronizing, restart it:
   ```shell
   oc debug node/<NODE_NAME> -- chroot /host systemctl restart chronyd
   ```
   
   - Force time synchronization:
   ```shell
   oc debug node/<NODE_NAME> -- chroot /host chronyc makestep
   ```

2. Configuration fixes through MachineConfig:
   - Create a MachineConfig CR to update chrony configuration with reliable time servers:
   ```yaml
   apiVersion: machineconfiguration.openshift.io/v1
   kind: MachineConfig
   metadata:
     labels:
       machineconfiguration.openshift.io/role: worker
     name: worker-chrony-configuration
   spec:
     config:
       ignition:
         version: 3.2.0
       storage:
         files:
         - contents:
             source: data:,server%200.rhel.pool.ntp.org%20iburst%0Aserver%201.rhel.pool.ntp.org%20iburst%0Aserver%202.rhel.pool.ntp.org%20iburst%0Aserver%203.rhel.pool.ntp.org%20iburst%0Adriftfile%20/var/lib/chrony/drift%0Amakestep%201.0%203%0Artcsync%0Alogdir%20/var/log/chrony
           mode: 0644
           path: /etc/chrony.conf
           overwrite: true
   ```

3. For connectivity issues:
   - Check if the node can reach the time servers:
   ```shell
   oc debug node/<NODE_NAME> -- chroot /host ping -c 5 0.rhel.pool.ntp.org
   ```
   
   - Verify firewall rules allow NTP traffic (port 123/UDP):
   ```shell
   oc debug node/<NODE_NAME> -- chroot /host firewall-cmd --zone=public --add-service=ntp --permanent
   oc debug node/<NODE_NAME> -- chroot /host firewall-cmd --reload
   ```

4. For severe issues:
   - If the node continues to have synchronization problems, consider reboot:
   ```shell
   oc debug node/<NODE_NAME> -- chroot /host systemctl reboot
   ```
   
   - For persistent synchronization issues, the node may need replacement:
   ```shell
   oc adm cordon <NODE_NAME>
   oc adm drain <NODE_NAME> --ignore-daemonsets
   ```

## Additional Notes

- Time synchronization is absolutely critical for OpenShift cluster health
- This alert differs from NodeClockSkewDetected, which focuses on clock offsets of more than 300ms
- The `NodeClockNotSynchronising` alert indicates a more serious condition where the synchronization process itself is failing
- RHCOS uses chronyd as the default time synchronization daemon
- Chrony configuration can be managed via MachineConfig operator
- Time synchronization failures most severely impact etcd and authentication subsystems
- In virtualized environments, check hypervisor time synchronization settings
- Hardware clock issues may require investigation with your infrastructure provider

## OpenShift Documentation References

For more information on configuring time synchronization in OpenShift, refer to the following documentation:

- [Configuring chrony time service](https://docs.openshift.com/container-platform/latest/post_installation_configuration/machine-configuration-tasks.html#installation-special-config-chrony_post-install-machine-configuration-tasks) - Official OpenShift documentation on configuring the chrony time service
- [OpenShift 4 chrony configuration](https://access.redhat.com/solutions/4510631) - Red Hat solution article with additional details on configuring chrony in OpenShift 4
- [NTP server synchronization for installation](https://docs.openshift.com/container-platform/latest/installing/installing_bare_metal_ipi/ipi-install-installation-workflow.html#checking-ntp-synchronization_ipi-install-installation-workflow) - Information on ensuring proper NTP synchronization during OpenShift installation

For troubleshooting chrony on RHEL/RHCOS systems:
- [RHEL 9 - Configuring time synchronization](https://docs.redhat.com/en/documentation/red_hat_enterprise_linux/9/html/configuring_basic_system_settings/configuring-time-synchronization_configuring-basic-system-settings)
- [RHEL System Roles - timesync](https://docs.redhat.com/en/documentation/red_hat_enterprise_linux/9/html/configuring_basic_system_settings/configuring-time-synchronization_configuring-basic-system-settings#managing-time-synchronization-using-rhel-system-roles_configuring-time-synchronization)
