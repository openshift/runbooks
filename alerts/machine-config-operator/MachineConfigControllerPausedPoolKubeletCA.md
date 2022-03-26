# MachineConfigControllerPausedPoolKubeletCA

## Meaning

The apiserver has performed a certificate rotation, but the specified
MachineConfigPool is paused, preventing deployment of the MachineConfig
containing the rotated  `kublet-ca.crt` bundle to the pool's nodes.

>You will  need to unpause the specified pool to allow the new certificate
>through before the existing one expires.

This alert fires as a warning for a pool when:

- A MachineConfigPool is paused AND
- There is a new `kubelet-ca.crt` specified in the MachineConfigPool's spec AND
- The pool has been in this state for more than `1 hour`

This alert becomes critical for a pool when:

- The above conditions are met, and there are only two weeks ( `14 days` )
  remaining until the expiry date of the pool's most recent
  `kube-apiserver-to-kubelet-signer` certificate.

For clarity, the pool has a MachineConfig specified in
`status.Configuration.Name`. That MachineConfig has a file inside it with a path
of `/etc/kubernetes/kubelet-ca.crt`.  Inside that file is a certificate
`kube-apiserver-to-kubelet-signer` which has to be valid for your nodes to work.
The MachineConfigPool's `pause` feature is preventing that configuration from
being updated.

## Impact

If the pool remains paused, the nodes in the specified pool will stop working
when the certificates in the pool's existing `kubelet-ca.crt`  bundle expire.

### Short term: warning

> NOTE: This is not a desirable state to be in unless you know what you're
> doing.

After (at most) 12 hours following a certificate rotation, you will experience
the following negative symptoms if your pool is still paused:

- Pod logs for nodes in the specified pool will not be viewable in the web
  console
- The commands `oc logs`  ,  `oc debug` ,  `oc attach`,  `oc exec`  will not be
  usable.

This happens because:

- These features depend on having a  `kubelet-client` certificate in the cluster
  that matches the node's `kube-apiserver-to-kubelet-signer`
- The `kubelet-client` certificate in the cluster rotates every 12 hours
- The `kubelet-client` is signed by the *most recent*
  `kube-apiserver-to-kubelet-signer`
- When a certificate rotation happens, nodes in the paused pool no longer have
  the most recent `kube-apiserver-to-kubelet-signer` (it is in `kubelet-ca.crt`,
  which is stuck behind pause)
- So once `kubelet-client` rotates and gets signed with the most recent signer,
  nodes in the paused pool cannot verify `kubelet-client` and therefore do not
  trust it.

Other than these symptoms, your nodes should work normally until the
kube-apiserver-to-kubelet-signer expires.

### Long term: critical

The `kube-apiserver-to-kubelet-signer` certificate in the `kubelet-ca.crt`
bundle is what enables trust between your node's kubelet and your cluster. The
apiserver schedules a rotation when this certificate reaches 80% of its
(currently 365 day) lifetime. The new one *must* be deployed before all of the
previous ones expire.

The kubelets on the nodes in the specified pool will stop communicating with the
cluster apiserver if the `kube-apiserver-to-kubelet-signer`  is allowed to
expire. The nodes will no longer be able to participate in the cluster. *This is
very bad*.

To avoid this, you will need to unpause the specified pool to allow the new
`kubelet-ca.crt` bundle to be deployed.

## Diagnosis

If a pool is paused, it was paused by an administrator; pools do not pause
themselves automatically.

>NOTE: there are some operators (like SR-IOV) that may briefly pause a pool to
>do some work, but they will not leave the pool paused long-term.

For the commands below, replace `$MCP_NAME` with the name of your pool.

You can see the pool's paused status by looking at the `spec.paused` field for
the pool:

```console
oc -n openshift-machine-config-operator get mcp $MCP_NAME -o jsonpath='{.spec.paused}'
```

For the `kube-apiserver-to-kubelet-signer` certificate in the cluster, you can
check its annotations to see when it was rotated:

```console
oc -n openshift-kube-apiserver-operator describe secret kube-apiserver-to-kubelet-signer
```

For the `kubelet-client` cert (the one that is responsible for `oc logs`, etc
working) in the cluster, you can check which signer it was signed with:

```console
oc -n openshift-kube-apiserver describe secret/kubelet-client
```

You can also check the certificates that are present in the `kublet-ca.crt`
bundle on one of your nodes to see when they expire:

```console
openssl crl2pkcs7 -nocrl -certfile /etc/kubernetes/kubelet-ca.crt | openssl pkcs7 -print_certs -text -noout
``````

>NOTE: There may be multiple `kube-apiserver-to-kubelet-signer` certificates in
>the certificate bundles, kube-apiserver does not pull them out until they
>expire. You want to look at the *newest* `kube-apiserver-to-kubelet-signer`.

You can find which MachineConfig a pool is currently using by looking at the
pool status:

```console
oc -n openshift-machine-config-operator get mcp $MCP_NAME -o jsonpath='{.status.configuration.name}'
```

Use the following commands to check the expiry dates, based on how the bundle is
encoded.

If the bundle is URL-encoded, use the following command with the desired
MachineConfig to decode it:

```console
oc get mc rendered-worker-bc1470f2331a3136999e0b49d85e1e21 -o jsonpath='{.spec.config.storage.files[?(@.path=="/etc/kubernetes/kubelet-ca.crt")].contents.source}' | python3 -c 'import sys, urllib.parse; print(urllib.parse.unquote(sys.stdin.read()))' | openssl x509 -text -noout
```

If the bundle is base64-encoded and gzipped, use the following command with the
desired MachineConfig to decode it:

```console
ENCODEDCERT=$(oc get mc rendered-worker-bc1470f2331a3136999e0b49d85e1e21 -o jsonpath='{.spec.config.storage.files[?(@.path=="/etc/kubernetes/kubelet-ca.crt")].contents.source}') CHOMPED=${ENCODEDCERT#"data:;base64,"} echo $CHOMPED | base64 -d | gzip -d | openssl x509 -text -noout
```

## Mitigation

You must unpause the pool.

>NOTE: Unpausing the pool will result in **ALL** pending configuration that was
>"waiting behind pause" being applied, and not just the certificate bundle.

Unpause the pool (substitute the pool name for $MCP_NAME):

```console
oc patch mcp $MCP_NAME --type='json' -p='[{"op": "replace", "path": "/spec/paused", "value":false}]'
```

You can also unpause manually by:

```console
oc edit mcp $MCP_NAME
```

and changing `spec.paused` to `false`.
