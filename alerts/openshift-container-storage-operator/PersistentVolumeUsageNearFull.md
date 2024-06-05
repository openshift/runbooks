# PersistentVolumeUsageNearFull

## Meaning

Persistent Volume Claim (PVC) usage is alarmingly very high,
indicating an imminent risk of reaching full capacity.
Please see the alert documentation text for an exact threshold limit.

## Impact

**Severity:** Warning
**Potential Customer Impact:** High

This alert signifies that a PVC is nearing its full capacity, potentially
leading to data loss if not addressed promptly.

## Diagnosis

The alert triggers when a Persistent Volume Claim (PVC) approaches or surpasses
very high capacity limit. It indicates the need to expand the PVC size to
accommodate more data or to remove unnecessary data to free up space.
Please see the alert documentation text for an exact threshold limit.

**Prerequisites:** [Prerequisites](helpers/diagnosis.md)

## Mitigation

### Recommended Actions

- **Expand the PVC Size:** Increase the capacity of the PVC to prevent data loss.
  ![Expand PVC Dropdown](helpers/screenshots/expand-pvc-dropdown.png)
  ![Expand PVC Dialog](helpers/screenshots/expand-pvc-dialog.png)
  
- **Delete Unnecessary Data:** Remove unnecessary data occupying space in the PVC.
