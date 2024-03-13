# OSDCpuLoadHigh

## Meaning

This alert indicates that the CPU usage in the OSD (Object Storage Daemon)
container on a specific pod has exceeded 80%, potentially affecting the
performance of the OSD.

## Impact

OSD is a critical component in Ceph storage, responsible for managing data
placement and recovery. High CPU usage in the OSD container suggests increased
processing demands, potentially leading to degraded storage performance.

## Diagnosis

To diagnose the alert, follow these steps:

1. Navigate to the Kubernetes dashboard or equivalent.
2. Access the "Workloads" section and select the relevant pod associated with
the OSD alert.
3. Click on the "Metrics" tab to view CPU metrics for the OSD container.
4. Verify that the CPU usage exceeds 80% over a significant period
(as specified in the alert configuration).

## Mitigation

If the OSD CPU usage is consistently high, consider taking the following steps:

1. Evaluate the overall storage cluster performance and identify the OSDs
contributing to high CPU usage.
2. Increase the number of OSDs in the cluster by adding more new storage
devices in the existing nodes or adding new nodes with new storage devices.
Review the Openshift Scaling storage documentation. This would help distribute
the load and improve overall system performance.
