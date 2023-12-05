# KMSServerConnectionAlert

## Meaning

Storage Cluster KMS Server is in un-connected state for more than 5s.
Please check KMS config

## Impact

Critical.
Encryption in block and file storage will not be available.
Information cannot be retrieved or stored properly.

## Diagnosis

Connection with external key management service is not working.

## Mitigation

Review configuration values in the ´ocs-kms-connection-details´ confimap.

Verify the connectivity with the external KMS, verifying
[network connectivity](helpers/networkConnectivity.md)
