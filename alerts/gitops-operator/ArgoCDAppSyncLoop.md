# ArgoCDAppSyncLoop

## Meaning

An Argo CD Application has a sustained high sync operation rate
(`argocd_app_sync_total`). Warning and critical share this alert name and
differ by the `severity` label:

* **warning** — above ~0.01 syncs/sec for 20m (~one sync every ~100s)
* **critical** — above ~0.1 syncs/sec for 10m (~one sync every ~10s)

This usually means the app is thrashing (for example auto-sync + selfHeal
fighting another controller such as an HPA), not a one-time OutOfSync drift.

## Impact

* Extra load on the application controller, repo-server, and Kubernetes API
* Application may stay OutOfSync / Progressing and never settle
* Can affect other apps on the same Argo CD instance under heavy load

## Diagnosis

1. Identify the Application from the alert labels (`name`, `namespace`).
2. Inspect sync / operation status:

   ```console
   $ oc get application <name> -n <namespace> -o yaml
   $ oc get application <name> -n <namespace> -o jsonpath='{.status.operationState}{"\n"}'
   ```

3. Check Diff in the Argo CD UI or CLI for fields that keep changing.
4. Look for conflicting controllers on the same resources (HPA, other
   operators, mutating webhooks, LimitRanges).
5. Confirm the metric rate:

   ```promql
   gitops:argocd_app_sync:rate10m{name="<name>",namespace="<namespace>"}
   ```

## Mitigation

* Stop declaring contested fields in Git (for example `replicas` when HPA
  owns them).
* Add `ignoreDifferences` for fields rewritten by admission or defaulting.
* Temporarily disable `selfHeal` / automated sync while fixing the conflict.
* Fix failing manifests or webhook policies if syncs are also failing
  (see also `ArgoCDAppSyncFailureLoop`).
