# VirtPlatformAutopilotThrashingDetected

## Meaning

This alert fires when virt-platform-autopilot detects an "edit war" on a
managed resource and pauses its reconciliation to protect the API server.

An edit war happens when another controller or a user repeatedly modifies a
resource that virt-platform-autopilot manages, so each apply is immediately
overwritten. When this is detected, the operator sets the
`platform.kubevirt.io/reconcile-paused=true` annotation on the resource and
stops reconciling it.

The alert is driven by the `kubevirt_autopilot_paused_resources` gauge,
labelled by `kind`, `name`, and `namespace`:

- `1` = reconciliation paused for the resource
- `0` = reconciliation active

The alert triggers immediately on `kubevirt_autopilot_paused_resources > 0`
(there is no `for` clause). The related `kubevirt_autopilot_thrashing_total`
counter tracks the anti-thrashing gate hits that lead up to a pause.

## Impact

This is a `warning` alert. While a resource is paused, its drift is no longer
corrected and it may diverge from the desired state. This is a deliberate
safety mechanism that prevents a reconciliation loop from overloading the API
server; the operator itself is healthy.

## Diagnosis

1. Export the `NAMESPACE` environment variable:

   ```bash
   $ export NAMESPACE="$(oc get deployment virt-platform-autopilot -A \
       -o custom-columns="":.metadata.namespace | tr -d '[:space:]')"
   ```

2. Read the `kind`, `name`, and `namespace` of the paused resource from the
   alert labels. Replace `<kind>`, `<name>`, and `<namespace>` with the
   corresponding alert-label values in all commands below. Confirm the pause
   annotation is set:

   ```bash
   $ oc get <kind> <name> -n <namespace> \
       -o jsonpath='{.metadata.annotations.platform\.kubevirt\.io/reconcile-paused}'
   ```

3. Identify the competing writer. As a best-effort check, inspect the
   resource's `managedFields`: the operator applies with the field manager
   `virt-platform-autopilot`, so any other manager is a likely source of the
   conflict:

   ```bash
   $ oc get <kind> <name> -n <namespace> \
       -o jsonpath='{.metadata.managedFields[*].manager}'
   ```

   For an authoritative record of who wrote the resource, consult the API
   server audit log (where audit logging is enabled).

## Mitigation

Choose the option that matches the intent:

- **The external change is unintentional:** stop or fix the conflicting
  controller or workflow, then resume reconciliation by removing the pause
  annotation:

  ```bash
  $ oc annotate <kind> <name> -n <namespace> \
      platform.kubevirt.io/reconcile-paused-
  ```

- **The external change is intentional** (another tool should own the
  resource): transfer ownership so the operator stops managing it:

  ```bash
  $ oc annotate <kind> <name> -n <namespace> \
      platform.kubevirt.io/mode=unmanaged --overwrite
  ```

- **Only specific fields conflict:** keep the resource managed but carve out
  the contested fields, using either a JSON Patch override or an ignore list:

  ```bash
  $ oc annotate <kind> <name> -n <namespace> \
      platform.kubevirt.io/patch='[{"op":"replace","path":"/spec/yourField","value":"yourValue"}]' \
      --overwrite

  $ oc annotate <kind> <name> -n <namespace> \
      platform.kubevirt.io/ignore-fields='/spec/replicas' --overwrite
  ```

  The `patch` annotation uses RFC 6902 JSON Patch array format. The
  `ignore-fields` annotation uses comma-separated RFC 6901 JSON Pointer paths
  (e.g., `/spec/replicas,/metadata/labels/app`).

The alert resolves when `kubevirt_autopilot_paused_resources` returns to `0`,
which happens once the pause annotation is removed and no further conflict is
detected.

If you cannot resolve the issue, log in to the
[Red Hat Customer Portal](https://access.redhat.com) and open a support case,
attaching the artifacts gathered during the diagnosis procedure.