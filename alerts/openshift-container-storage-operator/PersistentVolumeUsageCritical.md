# PersistentVolumeUsageCritical

## Meaning

A PVC is nearing its full capacity and may lead to data loss if not attended to
timely. The alert will be fired when Persistent Volume Claim usage has exceeded
more than 85% of its capacity.

## Impact

A PVC is about to be exhausted, any Write in PVS using this PVC will be blocked
when PVC will not have available space.

## Diagnosis

Using the Openshift console, go to Storage-PersistentVolumeClaims.
A list of the available PVCs with basic information about space used and
available will be shown.
Enter in the PVC affected to have more details.


## Mitigation

Expand the PVC size to increase the capacity.
In the list of PVCs (Storage-PersistentVolumeClaims), press the "three points"
button shown at the end of the affected PVC row. Select "Expand PVC" and
increase the size of the PVC.

![pvc-dropdown](helpers/screenshots/expand-pvc-dropdown.png)
![pvc-dialog](helpers/screenshots/expand-pvc-dialog.png)

Alternatively, you can also delete unnecessary data in PVs that may be taking
 up space

