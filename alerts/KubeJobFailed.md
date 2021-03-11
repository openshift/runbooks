# KubeJobFailed

## Meaning

[This alert][KubeJobFailed] is triggered for the case that the number of job execution
attempts exceeds the `backoffLimit`. A job can therefore create one or many pods
for its tasks.

## Impact

A task has not finished correctly. Depending on the task, this has an different
impact.

## Diagnosis

The alert should contain the job name and the namespace where that job failed.
Follow the particular mitigation steps according to that. For example:

```text
 - alertname = KubeJobFailed
... 
 - message = Job openshift-logging/elasticsearch-delete-app-1600855200 failed to complete.
```

## Mitigation

### Failed `elasticsearch-*` jobs in `openshift-logging` namespace

Make sure you are in the `openshift-logging` project

```console
$ oc project openshift-logging
```

Look at the elasticsearch pods. One easy way to filter these is to use the
`component` label with the value `indexManagement`:

```console
$ oc get pod -l component=indexManagement
```

If, as in this case, you see an `Error` pod with a subsequent `Completed` pod of
the same base name, the error was transient, and the `Error` pod can be deleted
safely.

Have a look at the jobs themselves:

```console
$ oc get jobs -n openshift-logging
```

If, as in this case, there is a healthy job for every failed one, it is safe to
delete the failed jobs. The alert should resolve itself after a few minutes.

[KubeJobFailed]: https://github.com/openshift/cluster-monitoring-operator/blob/aefc8fc5fc61c943dc1ca24b8c151940ae5f8f1c/assets/control-plane/prometheus-rule.yaml#L186-L195