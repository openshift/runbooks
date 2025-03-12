# Node Clock Issues (NodeClockSkewDetected and NodeClockNotSynchronising)

## Meaning

These alerts indicate problems with system time synchronization on an OpenShift worker node (RHCOS):
- NodeClockSkewDetected: System clock is offset by more than 300ms
- NodeClockNotSynchronising: System is unable to synchronize its clock

## Impact

- etcd cluster instability and potential data inconsistency
- Certificate validation failures for OpenShift API communications
- Authentication issues between OpenShift components
- Kubernetes scheduler timing inconsistencies
- Inaccurate logging timestamps affecting troubleshooting
- Operator reconciliation timing issues

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

4. Check the chrony sources and status:
```shell
oc debug node/<NODE_NAME> -- chroot /host chronyc sources
oc debug node/<NODE_NAME> -- chroot /host chronyc tracking
```

5. Review chrony configuration:
```shell
oc debug node/<NODE_NAME> -- chroot /host cat /etc/chrony.conf
```

6. Check for NTP traffic restrictions:
```shell
oc debug node/<NODE_NAME> -- chroot /host firewall-cmd --list-all
```

## Mitigation

1. Chronyd service issues:
   - Restart the chronyd service:
   ```shell
   oc debug node/<NODE_NAME> -- chroot /host systemctl restart chronyd
   ```
   
   - Force time synchronization:
   ```shell
   oc debug node/<NODE_NAME> -- chroot /host chronyc makestep
   ```

2. Configuration fixes through MachineConfig:
   - Create a MachineConfig CR to update chrony configuration:
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
             source: data:,server%20time.example.com%20iburst%0Adriftfile%20/var/lib/chrony/drift%0Amakestep%201.0%203%0Artcsync%0Alogdir%20/var/log/chrony
           mode: 0644
           path: /etc/chrony.conf
           overwrite: true
   ```

3. For consistent cluster-wide configuration:
   - Use the same NTP servers across all nodes
   - Consider local NTP servers if running air-gapped
   - Ensure proper network connectivity to time sources

4. For serious clock drift:
   - If the node continues to have problems, consider reboot:
   ```shell
   oc debug node/<NODE_NAME> -- chroot /host systemctl reboot
   ```
   
   - For persistent issues, the node may need replacement:
   ```shell
   oc adm cordon <NODE_NAME>
   oc adm drain <NODE_NAME> --ignore-daemonsets
   ```

## Additional Notes

- Time synchronization is absolutely critical for OpenShift cluster health
- RHCOS uses chronyd as the default time synchronization daemon
- Chrony configuration can be managed via MachineConfig operator
- Time skew most severely impacts etcd and authentication subsystems
- Check for virtualization issues if running on virtual infrastructure
- Hardware clock issues may require investigation with your infrastructure provider 