# ACMThanosCompactHalted

## Meaning

This alert is triggered by the ACM Observability service when the Thanos
Compactor component has entered a "halted" state for more than 5 minutes.

This "halted" state means the `observability-thanos-compact-0` pod is running,
but it has encountered an internal, non-retriable error and has stopped
performing its duties. The alert fires when the metric
`acm_thanos_compact_halted{job="observability-thanos-compact"}` has a value of
1.

The Thanos Compactor is a critical background component responsible for
organizing metric data in object storage, applying retention policies (deleting
old data), and down-sampling data to improve query performance.

## Impact

If the Thanos Compactor is halted, the metric data in the S3 object store is no
longer being managed. This has significant consequences:

* **Data Retention policies are not applied**: Old metric data will not be
  deleted, causing object storage (S3) costs to increase indefinitely.

* **Data Down-sampling fails**: Metric data is not down-sampled into 5-minute
  and 1-hour aggregates. This will cause Grafana dashboards and queries to
  become progressively slower especially when querying for longer rangers.

* **Increased storage usage**: The storage use of the Compactor and Store
  Gateway will keep increasing.

## Diagnosis

The first step is to identify why the compactor has halted by inspecting its
logs for fatal, non-retriable errors.

Confirm the alert is firing by querying the `acm_thanos_compact_halted` metric
in the OpenShift console (Observe -> Metrics) or Grafana. The result should be
1.

```console
acm_thanos_compact_halted{job="observability-thanos-compact"}
```

Inspect the logs of the `observability-thanos-compact-0` pod to find the fatal
error. This pod is in the `open-cluster-management-observability` namespace.

```console
oc logs -n open-cluster-management-observability observability-thanos-compact-0
```

Look for one of the following common failure patterns in the logs:

* **Disk Full**: The pod is in a `CrashLoopBackOff` state, and logs show `no
  space left on device`.

* **Data Corruption**: The logs show `err="compaction: ... get file ...
  meta.json: EOF"`, `unexpected end of JSON input` or `segment index 0 out of
  range`. This indicates a corrupted block in the S3 bucket.

* **S3 Connection Failure**: The logs show S3 authentication or connection
  errors, such as `The request signature we calculated does not match...` (bad
  secret key) or `no such host` (bad endpoint).

### 1. Diagnosis for 'no space left on device'

If the pod is in a `CrashLoopBackOff` state, you cannot rsh into it. You can
verify the disk is full by attempting to attach the `PVC` to a debug pod or by
checking the pod's last-known logs.

### 2. Diagnosis for Data Corruption

This error will appear in the logs while the pod is `Running`, as it attempts
to scan a malformed `meta.json` file from the S3 bucket.

### 3. Diagnosis for S3 Connection Failure

This error will appear in the logs while the pod is `Running`. This indicates
the `thanos-object-storage` secret is misconfigured.

Get the `thanos-object-storage` secret and decode it to verify its contents:

```console
oc get secret thanos-object-storage -n open-cluster-management-observability -o jsonpath='{.data.thanos\.yaml}' | base64 --decode
```

Check the `endpoint`, `access_key`, and `secret_key` values in the output to
ensure they are correct.

**Verify the ObjectBucketClaim status:**

If you suspect S3 connection issues, verify the OBC is healthy:

```console
oc get obc -n open-cluster-management-observability
oc get configmap <OBC-NAME> -n open-cluster-management-observability -o yaml
```

Check that the ConfigMap contains valid endpoint and bucket information.

## Mitigation

Resolution depends on the fatal error found during the Diagnosis.

### 1. Halted due to 'no space left on device'

The pod's persistent volume is full.

* Identify the `StorageClass` used by the `PVC`:

    ```console
    oc get pvc data-observability-thanos-compact-0 -n open-cluster-management-observability -o jsonpath='{.spec.storageClassName}'
    ```

* Verify the `StorageClass` has `allowVolumeExpansion: true`:

    ```console
    oc get sc <YOUR_STORAGE_CLASS> -o yaml | grep allowVolumeExpansion
    ```

* If expansion is allowed, edit the `PVC` to request more storage. This can be
  done while the pod is crashing.

    ```console
    oc edit pvc data-observability-thanos-compact-0 -n open-cluster-management-observability
    ```

* Change `spec.resources.requests.storage` from its current value to a larger
  one (e.g., 98G to 120G).

* Kubernetes and ODF will automatically expand the volume. The pod will stop
  crash-looping and will start successfully, clearing the alert.

### 2. Halted due to Data Corruption

A data block in the S3 object store is malformed.

* Identify the corrupted block ID from the pod logs (e.g., `download block
  01K8JQ8...`).

* Launch a debug pod with the `aws-cli` and the S3 credentials.

* From the debug pod, delete the entire corrupted block directory from the S3
  bucket:

    ```console
    aws --endpoint-url <S3_ENDPOINT> --no-verify-ssl s3 rm s3://<BUCKET_NAME>/<CORRUPTED_BLOCK_ID>/ --recursive
    ```

* Restart the `observability-thanos-compact-0` pod to clear the halted state:

    ```console
    oc delete pod -n open-cluster-management-observability observability-thanos-compact-0
    ```

### 3. Halted due to S3 Connection Failure

The `thanos-object-storage` secret in the
`open-cluster-management-observability` namespace is incorrect.

* Retrieve the correct S3 credentials (endpoint, bucket name, access key, and
  secret key) from your ODF `ObjectBucketClaim` (OBC).

* Create a correct `thanos.yaml` file on your local machine.

* Update the `thanos-object-storage` secret with the correct configuration:

    ```console
    $ oc create secret generic thanos-object-storage \
      -n open-cluster-management-observability \
      --from-file=thanos.yaml=./thanos.yaml \
      --dry-run=client -o yaml | oc apply -f -
    ```

* Restart the `observability-thanos-compact-0` pod to apply the new secret:

    ```console
    oc delete pod -n open-cluster-management-observability observability-thanos-compact-0
    ```

### 4. Verify Resolution

After applying any mitigation, verify the compactor has recovered:

* Check the halted metric has returned to 0:

    ```console
    oc exec -n open-cluster-management-observability observability-thanos-compact-0 -- wget -qO- http://localhost:10902/metrics | grep acm_thanos_compact_halted
    ```

* Monitor the pod logs for successful compaction activity:

    ```console
    oc logs -n open-cluster-management-observability observability-thanos-compact-0 --tail=50 -f
    ```

* The alert should clear within 5-10 minutes of successful compaction.

* If the compactor has not been working for a long period the
  `acm_thanos_compact_todo_compactions` metric is expected to be high. After
  restarting the compactor keep an eye on this metric to ensure that the number
  of todo compactions is decreasing. It might take several weeks to work
  through the full compaction backlog and for the todo compactions to fall down
  to below 10.

**Expected resolution time**: Alert should clear within 5-10 minutes after the
compactor resumes normal operation. The compaction backlog might take several
weeks to get through depending on how long the compactor was non-functioning
for.

For more detailed troubleshooting steps and additional scenarios, refer to the
Red Hat Knowledge Base article: [Thanos compactor
halted](https://access.redhat.com/solutions/7080672)
