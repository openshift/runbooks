# Alerts Hierarchy
`alerts-hierarchy.dot` needs to be updated for each new alert added to
Cluster Network Operator that overlaps with another alert using
graphviz and then execute the following command to update a graphic
visualisation of the hierarchy:

```sh
dot ./alerts-hierarchy.dot -Tsvg -Nshape=rect > alerts-hierarchy.svg
```
