# KubeJobFailed

## Meaning

[This alert][KubeJobFailed] is triggered for the case that the number of job
execution attempts exceeds the `backoffLimit`. A job can therefore create one or
many pods for its tasks.

## Impact

A task has not finished correctly. Depending on the task, this has a different
impact.

## Diagnosis

The alert should contain the job name and the namespace where that job failed.
Follow the particular mitigation steps according to that. For example:

```text
 - alertname = KubeJobFailed
...
 - job_name = elasticsearch-delete-app-1600903800
 - namespace = openshift-logging
... 
 - message = Job openshift-logging/elasticsearch-delete-app-1600855200 failed to complete.
```

## Mitigation

Find the pods that belong to that job:

```console
$ oc get pod -n $NAMESPACE -l job-name=$JOBNAME
```

If you see an `Error` pod with a subsequent `Completed` pod of
the same base name, the error was transient, and the `Error` pod can safely be deleted.

Have a look at the jobs themselves:

```console
$ oc get jobs -n $NAMESPACE
```

If there is a healthy job for every failed one, it is safe to delete the failed
jobs. The alert should resolve itself after a few minutes.

[KubeJobFailed]: https://github.com/openshift/cluster-monitoring-operator/blob/aefc8fc5fc61c943dc1ca24b8c151940ae5f8f1c/assets/control-plane/prometheus-rule.yaml#L186-L195