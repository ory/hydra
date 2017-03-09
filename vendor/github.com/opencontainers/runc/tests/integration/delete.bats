#!/usr/bin/env bats

load helpers

function setup() {
  teardown_busybox
  setup_busybox
}

function teardown() {
  teardown_busybox
}

@test "runc delete" {
  # run busybox detached
  runc run -d --console-socket $CONSOLE_SOCKET test_busybox
  [ "$status" -eq 0 ]

  # check state
  wait_for_container 15 1 test_busybox

  testcontainer test_busybox running

  runc kill test_busybox KILL
  # wait for busybox to be in the destroyed state
  retry 10 1 eval "__runc state test_busybox | grep -q 'stopped'"

  # delete test_busybox
  runc delete test_busybox

  runc state test_busybox
  [ "$status" -ne 0 ]
}

@test "runc delete --force" {
  # run busybox detached
  runc run -d --console-socket $CONSOLE_SOCKET test_busybox
  [ "$status" -eq 0 ]

  # check state
  wait_for_container 15 1 test_busybox

  testcontainer test_busybox running

  # force delete test_busybox
  runc delete --force test_busybox

  runc state test_busybox
  [ "$status" -ne 0 ]
}

@test "run delete with multi-containers" {
  # create busybox1 detached
  runc create --console-socket $CONSOLE_SOCKET test_busybox1
  [ "$status" -eq 0 ]

  testcontainer test_busybox1 created

  # run busybox2 detached
  runc run -d --console-socket $CONSOLE_SOCKET test_busybox2
  [ "$status" -eq 0 ]

  wait_for_container 15 1 test_busybox2
  testcontainer test_busybox2 running

  # delete both test_busybox1 and test_busybox2 container
  runc delete test_busybox1 test_busybox2

  runc state test_busybox1
  [ "$status" -ne 0 ]

  runc state test_busybox2
  [ "$status" -eq 0 ]

  runc kill test_busybox2 KILL
  # wait for busybox2 to be in the destroyed state
  retry 10 1 eval "__runc state test_busybox2 | grep -q 'stopped'"

  # delete test_busybox2
  runc delete test_busybox2

  runc state test_busybox2
  [ "$status" -ne 0 ]
}


@test "run delete --force with multi-containers" {
  # create busybox1 detached
  runc create --console-socket $CONSOLE_SOCKET test_busybox1
  [ "$status" -eq 0 ]

  testcontainer test_busybox1 created

  # run busybox2 detached
  runc run -d --console-socket $CONSOLE_SOCKET test_busybox2
  [ "$status" -eq 0 ]

  wait_for_container 15 1 test_busybox2
  testcontainer test_busybox2 running

  # delete both test_busybox1 and test_busybox2 container
  runc delete --force test_busybox1 test_busybox2

  runc state test_busybox1
  [ "$status" -ne 0 ]

  runc state test_busybox2
  [ "$status" -ne 0 ]
}
