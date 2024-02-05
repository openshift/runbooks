# Find connected ODF Client cluster

List all the storageconsumers that are onboarded onto the provider cluster
```bash
  oc get storageconsumer -nopenshift-storage
```

Extract the interested client cluster id from storageconsumer name
```bash
  CONSUMERNAME="storageconsumer-<ID>"
  CLIENTCLUSTERID=${CONSUMERNAME/storageconsumer-/}
  CONSUMERUID=$(oc get storageconsumer -nopenshift-storage ${CONSUMERNAME} \
    -ojsonpath='{.metadata.uid}')
```

Connected ODF Client cluster id will be stored in `CLIENTCLUSTERID` and UID in
`CONSUMERUID`
