# Integration Test for watcherx/FileWatcher

As kubernetes has a special way to change mounted config map values we want to
make sure our file watcher is compatible with that.

## Perquisites

The versions are the ones that definitely work.

- kind (v0.8.1)
- kubectl (v1.18.5)
- docker (v19.03.12-ce)
- make (v4.3)

## Structure

The `main.go` just logs all events it gets. It is deployed to a kind kubernetes
cluster together with a configmap that gets updated during the test. For details
on the test steps have a look at the `Makefile`.

## Running

To generate the log snapshot run `make snapshot`. That snapshot should be
committed. To check if the FileWatcher works run `make check`. For debugging
purposes single steps of the setup have descriptive make target names and can be
run separately. It is safe to delete the cluster at any point or rerun snapshot
generation.
