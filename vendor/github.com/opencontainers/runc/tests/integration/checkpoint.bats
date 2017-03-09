#!/usr/bin/env bats

load helpers

function setup() {
  teardown_busybox
  setup_busybox
}

function teardown() {
  teardown_busybox
}

@test "checkpoint and restore" {
  requires criu

  # criu does not work with external terminals so..
  # setting terminal and root:readonly: to false
  sed -i 's;"terminal": true;"terminal": false;' config.json
  sed -i 's;"readonly": true;"readonly": false;' config.json
  sed -i 's/"sh"/"sh","-c","while :; do date; sleep 1; done"/' config.json

  (
    # run busybox (not detached)
    runc run test_busybox
    [ "$status" -eq 0 ]
  ) &

  # check state
  wait_for_container 15 1 test_busybox

  runc state test_busybox
  [ "$status" -eq 0 ]
  [[ "${output}" == *"running"* ]]

  # checkpoint the running container
  runc --criu "$CRIU" checkpoint test_busybox
  # if you are having problems getting criu to work uncomment the following dump:
  #cat /run/opencontainer/containers/test_busybox/criu.work/dump.log
  [ "$status" -eq 0 ]

  # after checkpoint busybox is no longer running
  runc state test_busybox
  [ "$status" -ne 0 ]

  # restore from checkpoint
  (
    runc --criu "$CRIU" restore test_busybox
    [ "$status" -eq 0 ]
  ) &

  # check state
  wait_for_container 15 1 test_busybox

  # busybox should be back up and running
  runc state test_busybox
  [ "$status" -eq 0 ]
  [[ "${output}" == *"running"* ]]
}
