# CephNodeDown

## Meaning

This alert indicates that one of the Ceph storage node went down.

## Impact

A node running Ceph pods is down. While storage operations will continue to
function as Ceph is designed to deal with a node failure, it is recommended
to resolve the issue to minimise risk of another node going down and affecting
storage functions.

## Diagnosis

The alert message will clearly indicate which node is down

    Storage node <nod-name> went down. Please check the node immediately.

## Mitigation

Document the current OCS pods (running and failing):

    oc -n openshift-storage get pods

The OCS resource requirements must be met in order for the osd pods to be
scheduled on the new node. This may take a few minutes as the ceph cluster
recovers data for the failing but now recovering osd.

To watch this recovery in action ensure the osd pods were actually placed on the
new worker node.

Check if the previous failing osd pods are now running:

    oc -n openshift-storage get pods

If the previously failing osd pods have not been scheduled, use describe and
check events for reasons the pods were not rescheduled.

Describe events for failing osd pod:

    oc -n openshift-storage get pods | grep osd

Find a failing osd pod(s):

    oc -n openshift-storage describe pods/<osd podname from previous step>

In the event section look for failure reasons, such as resources not being met.

In addition, you may use the rook-ceph-toolbox to watch the recovery. This step
is optional but can be helpful for large Ceph clusters.

**Determine failed OCS Node** [determine_failed_ocs_node](helpers/determineFailedOcsNode.md)

[access toolbox](helpers/accessToolbox.md)

From the rsh command prompt, run the following and watch for "recovery" under
the io section:

    ceph status

[gather logs](helpers/gatherLogs.md)

