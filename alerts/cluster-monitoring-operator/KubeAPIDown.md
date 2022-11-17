# KubeAPIDown

## Meaning

The `KubeAPIDown` alert is triggered when all Kubernetes API servers have not
been reachable by the monitoring system for more than 15 minutes.

## Impact

This is a critical alert. It indicates that the Kubernetes API is not
responding, and the cluster might be partially or fully non-functional.

## Diagnosis

1. Verify the status of the API server targets in Prometheus in the OpenShift
web console.

1. Confirm whether the API is also unresponsive:

    ```console
    $ oc cluster-info
    ```

1. If you can still reach the API server, a network issue might exist between
the Prometheus instances and the API server pods. Review the status of the API
server pods:

    ```console
    $ oc -n openshift-kube-apiserver get pods
    $ oc -n openshift-kube-apiserver logs -l 'app=openshift-kube-apiserver'
    ```
## Mitigation

If you can still reach the API server intermittently, you might be able to
troubleshoot this issue as you would for any other failing deployment.

If the API server is not reachable at all, refer to the disaster recovery
documentation for your version of OpenShift.
