# OpenShift Runbooks
A collection of runbooks for alert rules that are shipped with OpenShift.

## How To Use
OpenShift Container Platform is shipping a finely tuned set of alerts to inform
the cluster's owner and/or operator of events and bad conditions in the cluster.
This repository is a consistently growing collection of runbooks for said
alerts, that are intended to guide a cluster owner and/or operator through the
steps of fixing problems on clusters, which are surfaced by alerts.

Going forward alerts in OpenShift will ship a link to the corresponding runbook
in this repository, to make fixing problems even easier.

## Repository Layout
Runbooks in this repository are grouped by the operator that is responsible for
shipping the respective alert. This results in a structure as follows:

`root/alerts/operator_repository_name/some_alert.md`

The root folder contains all information and files required for this repository
in general, while all runbooks live in their respective subfolder.

## Onboarding
If your operator has not shipped any runbooks so far, please create a subfolder
following the repository layout and place an OWNERS file in it. The OWNERS file
should be the same as the one in the operator's main repository. Additionally
add all runbooks in said folder, that you want to contribute.

## Contributing
We are welcoming *any* contributions! While Red Hat is going to keep adding more
runbooks for new alerts as they surface, we would also love to extend existing
runbooks with any real world experiences that cluster owners or operators have
made.

Steps to your contribution are outlined below:
* Fork this repository
* Make your change (We recommend using a branch in your fork.)
* Submit a Pull Request

Our integrations and bots will walk you through any other required steps, like
assigning someone to look at your PR, or inform you if tests are failing.
