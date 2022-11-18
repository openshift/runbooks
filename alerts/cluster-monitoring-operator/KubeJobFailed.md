# KubeJobFailed

## Meaning

The `KubeJobFailed` alert triggers when the number of job execution attempts
exceeds the value defined in `backoffLimit`. If this issue occurs, a job can
create one or many pods for its tasks.

## Impact

A task has not finished correctly. Depending on the task, the severity of the
impact differs.

## Diagnosis

The alert message contains the job name and the namespace in which that job
failed. The following message provides an example:

```text
 - alertname = KubeJobFailed
...
 - job_name = elasticsearch-delete-app-1600903800
 - namespace = openshift-logging
... 
 - message = Job openshift-logging/elasticsearch-delete-app-1600855200 failed to complete.
```

This information is required for you to follow the mitigation steps.

## Mitigation

* Find the pods that belong to the failed job shown in the alert message:

    ```console
    $ oc get pod -n $NAMESPACE -l job-name=$JOBNAME
    ```

    If you see an `Error` pod together with a subsequent `Completed` pod of the
    same base name, the error was transient, and you can safely delete the
    `Error` pod.

* Review the status of the jobs:

    ```console
    $ oc get jobs -n $NAMESPACE
    ```

    If a healthy job exists for every failed job, you can safely delete the
    failed jobs, and the alert will resolve itself after a few minutes.
