# NodeNetworkInterfaceDown

## Meaning

This alert fires when one or more network interfaces on a node have been down
for more than 5 minutes. The alert excludes virtual ethernet (veth) devices and
bridge tunnels.

## Impact

Network interface failures can lead to:
- Reduced network connectivity for pods on the affected node
- Potential service disruptions if critical network paths are affected
- Degraded cluster communication if management interfaces are impacted

## Diagnosis

1. Identify the affected node and interfaces:
   ```bash
   oc get nodes
   ```

   ```bash
   ssh <node-address>
   ```

   ```bash
   ip link show | grep -i down
   ```

2. Check network interface details:
   ```bash
   ip addr show
   ```

   ```bash
   ethtool <interface-name>
   ```

3. Review system logs for network-related issues:
   ```bash
   journalctl -u NetworkManager
   ```

   ```bash
   dmesg | grep -i eth
   ```

## Mitigation

1. For physical interface issues:
   - Check physical cable connections
   - Verify switch port configuration
   - Test the interface with a different cable/port

2. For software or configuration issues:
   ```bash
   # Restart NetworkManager
   systemctl restart NetworkManager
   ```

   ```bash
   # Bring interface up manually
   ip link set <interface-name> up
   ```

3. If the issue persists:
   - Check network interface configuration files
   - Verify driver compatibility
   - If the failure is on a physical interface, consider hardware replacement

## Additional notes
- Monitor interface status after mitigation
- Document any hardware replacements or configuration changes
- Consider implementing network redundancy for critical interfaces

If you cannot resolve the issue, log in to the
[Customer Portal](https://access.redhat.com) and open a support case,
attaching the artifacts gathered during the diagnosis procedure.