# SRE-to-ODF-Engineering escalation path

## Purpose

This SOP describes how the SRE Team will escalate issues to the
ODF Engineering Team in case of incidents which would require help
from ODF engineering to resolve them.

## Scope

The scope of this escalation path is to cover the incidents which have
unavailability/data loss associated with the customer and none of the
existing SOPs help to resolve the respective incident at the same time.

## Prerequisites

* `MTSRE` Role permissions over OCM.

## Responsibilities

* SRE to only escalate the issue to engineering on an urgent basis when the
incident is associated with unavailability or data loss for the customer.
* ocm-cli version 0.1.63

## Procedure

### Determine whether the incident leads to unavailability / data loss for the customer

The best way to determine this is to see if the corresponding provider cluster
has any storageconsumers or not.

* Determine the cluster ID of the Provider cluster associated with the customer,
in case it's not directly available from PagerDuty incident: Happens when the
incident's source is an ocs-consumer.

* Backplane access into the provider cluster (Checkout [References](#references)
for the SOP to do so).

* See if there are any storageconsumers or not.

    oc get storageconsumers.ocs.openshift.io -n openshift-storage

If there are no storageconsumers, then wait until the next set of ODF's working
hours to escalate the issue to engineering.

### Escalation path when the incident urgently requires help from ODF Engineering

* Follow the steps in [link](<https://source.redhat.com/groups/public/openshiftplatformsre/wiki/how_to_sre_escalate_towards_engineering>),
starting with submitting the Google Form,
to get help from the ODF Engineering Team.

### If above steps do not work

Follow the manual steps to escalate to ODF Engineering Team

* Fetch certain details of the impacted cluster:
  * Cluster Details:

    ocm describe cluster <internal/external-id-of-cluster>

    For example

    ```bash
    ❯ ocm describe cluster 1s0q9mp65e2f370q9s7ip75g6sj0sifg

    ID:      1s0q9mp65e2f370q9s7ip75g6sj0sifg
    External ID:    7ca71103-f0d8-4e66-a586-35332eaa4b10
    Name:      e2e-ci-prov
    State:      ready
    API URL:    https://<management-domain-address>:6443
    API Listening:    external
    Console URL:    https://<url-address-through-which-we-access-cluster-console-ui>
    Masters:    3
    Infra:      2
    Computes:    3
    Product:    osd
    Provider:    aws
    Version:    4.10.11
    Region:      us-west-2
    Multi-az:    false
    CCS:      true
    Subnet IDs:    [subnet-06a34a7bb480707bf subnet-07b3053eb531c5e70]
    PrivateLink:    false
    STS:      false
    Existing VPC:    true
    Channel Group:    stable
    Cluster Admin:    true
    Organization:    Red Hat-Storage-OCS
    Creator:    <name>-storage-ocs
    Email:      <name>@domain.com
    AccountNumber: <acc-number>
    Created:    2022-05-05T15:53:01Z
    Expiration:    2022-06-16T14:09:19Z
    Shard:      https://domain-name.openshiftapps.com:6443
    ```

* Open a Jira ticket in the [RHOSDFP Jira board](<https://issues.redhat.com/projects/RHOSDFP>)
with the following details:

    ```bash
    Basic Details:

    Addon name:
    Addon version: <can be gotten via
                   `oc get csvs -n <addon-namespace>` on the affected cluster>
    OCP version:
    Business Impact:
    Description:

    Cluster and Customer Details:
        <output of the following command> `ocm describe cluster <internal-id>`

    Pagerduty Incident: <PD URL pointing to the incident>
    ```

    Additional Operators which came with the addon (including the core addon's
    operator): \<output of the following command from inside the affected cluster>
    `oc get csvs -n openshift-storage`

    (Example ticket mentioned in the [References](#references))

* Open a google chat thread in [ODF Managed Services Escalation](https://mail.google.com/chat/u/0/?zx=545uzcc7jlkp#chat/space/AAAAOdTnXXo)
  * @hey odf-cae-team I need engineering assistance with \<RHOSDFP ticket you created>
  * Then in the next message in the thread provide a short description of the issue.
* One of the Engineering escalation contacts is expected to ack the escalation
coming from SRE in under 15 minutes and involve the right set of engineering
team on the issue.
* Get in touch with the associated person of ODF Engineering and get on a bridge
call with them if need be to resolve the issue.
* If there is no acknowledgement by the engineering escalation contacts, a call
should be made to the engineering escalation contact as per [this table](<https://docs.google.com/document/d/1RKvxXnxoIaIPW-tbONnZ1t9Pmec5TH2qYgDzAKSKtO4/edit#bookmark=kix.x1bsgr3rpjx6>).

## References

* [OCM CLI Version 0.1.63](https://github.com/openshift-online/ocm-cli/releases/tag/v0.1.63)
* [RHSTOR Jira board](https://issues.redhat.com/projects/RHSTOR)
* [ODF Managed Services Escalation](https://mail.google.com/chat/u/0/?zx=545uzcc7jlkp#chat/space/AAAAOdTnXXo)
* [Example of RHSTOR ticket](https://issues.redhat.com/projects/RHSTOR/issues/RHSTOR-3329)
* [Steps to backplane access into a cluster](https://gitlab.cee.redhat.com/service/managed-tenants-sops/-/blob/main/MT-SRE/sops/addon-enable-backplane.md)
* [Interim SRE-to-ODF-Engineering escalation path](https://docs.google.com/document/d/1RKvxXnxoIaIPW-tbONnZ1t9Pmec5TH2qYgDzAKSKtO4/edit)
* [RHSTOR ticket tracking the RFE to allow fetching the Provider Cluster ID from any of its associated consumers](https://issues.redhat.com/browse/RHSTOR-3381)
