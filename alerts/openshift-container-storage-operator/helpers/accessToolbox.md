# Access the Toolbox

Run the following to rsh to the toolbox pod:

    TOOLS_POD=$(oc get pods -n openshift-storage -l app=rook-ceph-tools -o name)
    oc rsh -n openshift-storage $TOOLS_POD
