#!/usr/bin/env bats

load helpers

function setup() {
  teardown_busybox
  setup_busybox
}

function teardown() {
  teardown_busybox
}

@test "runc pause and resume" {
  # run busybox detached
  runc run -d --console-socket $CONSOLE_SOCKET test_busybox
  [ "$status" -eq 0 ]

  wait_for_container 15 1 test_busybox

  # pause busybox
  runc pause test_busybox
  [ "$status" -eq 0 ]

  # test state of busybox is paused
  testcontainer test_busybox paused

  # resume busybox
  runc resume test_busybox
  [ "$status" -eq 0 ]

  # test state of busybox is back to running
  testcontainer test_busybox running
}

@test "runc pause and resume with multi-container" {
  # run test_busybox1 detached
  runc run -d --console-socket $CONSOLE_SOCKET test_busybox1
  [ "$status" -eq 0 ]

  wait_for_container 15 1 test_busybox1

  # run test_busybox2 detached
  runc run -d --console-socket $CONSOLE_SOCKET test_busybox2
  [ "$status" -eq 0 ]

  wait_for_container 15 1 test_busybox2

  # pause test_busybox1 and test_busybox2
  runc pause test_busybox1 test_busybox2
  [ "$status" -eq 0 ]

  # test state of test_busybox1 and test_busybox2 is paused
  testcontainer test_busybox1 paused
  testcontainer test_busybox2 paused

  # resume test_busybox1 and test_busybox2
  runc resume test_busybox1 test_busybox2
  [ "$status" -eq 0 ]

  # test state of two containers is back to running
  testcontainer test_busybox1 running
  testcontainer test_busybox2 running

  # delete test_busybox1 and test_busybox2
  runc delete --force test_busybox1 test_busybox2

  runc state test_busybox1
  [ "$status" -ne 0 ]

  runc state test_busybox2
  [ "$status" -ne 0 ]
}

@test "runc pause and resume with nonexist container" {
  # run test_busybox1 detached
  runc run -d --console-socket $CONSOLE_SOCKET test_busybox1
  [ "$status" -eq 0 ]

  wait_for_container 15 1 test_busybox1

  # run test_busybox2 detached
  runc run -d --console-socket $CONSOLE_SOCKET test_busybox2
  [ "$status" -eq 0 ]

  wait_for_container 15 1 test_busybox2

  # pause test_busybox1, test_busybox2 and nonexistent container
  runc pause test_busybox1 test_busybox2 nonexistent
  [ "$status" -ne 0 ]

  # test state of test_busybox1 and test_busybox2 is paused
  testcontainer test_busybox1 paused
  testcontainer test_busybox2 paused

  # resume test_busybox1, test_busybox2 and nonexistent container
  runc resume test_busybox1 test_busybox2 nonexistent
  [ "$status" -ne 0 ]

  # test state of two containers is back to running
  testcontainer test_busybox1 running
  testcontainer test_busybox2 running

  # delete test_busybox1 and test_busybox2
  runc delete --force test_busybox1 test_busybox2

  runc state test_busybox1
  [ "$status" -ne 0 ]

  runc state test_busybox2
  [ "$status" -ne 0 ]
}
