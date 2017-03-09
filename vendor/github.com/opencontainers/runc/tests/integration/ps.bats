#!/usr/bin/env bats

load helpers

function setup() {
  teardown_busybox
  setup_busybox
}

function teardown() {
  teardown_busybox
}

@test "ps" {
  # start busybox detached
  runc run -d --console-socket $CONSOLE_SOCKET test_busybox
  [ "$status" -eq 0 ]

  # check state
  wait_for_container 15 1 test_busybox

  testcontainer test_busybox running

  runc ps test_busybox
  [ "$status" -eq 0 ]
  [[ ${lines[0]} =~ UID\ +PID\ +PPID\ +C\ +STIME\ +TTY\ +TIME\ +CMD+ ]]
  [[ "${lines[1]}" == *"root"*[0-9]* ]]
}

@test "ps -f json" {
  # start busybox detached
  runc run -d --console-socket $CONSOLE_SOCKET test_busybox
  [ "$status" -eq 0 ]

  # check state
  wait_for_container 15 1 test_busybox

  testcontainer test_busybox running

  runc ps -f json test_busybox
  [ "$status" -eq 0 ]
  [[ ${lines[0]} =~ [0-9]+ ]]
}

@test "ps -e -x" {
  # start busybox detached
  runc run -d --console-socket $CONSOLE_SOCKET test_busybox
  [ "$status" -eq 0 ]

  # check state
  wait_for_container 15 1 test_busybox

  testcontainer test_busybox running

  runc ps test_busybox -e -x
  [ "$status" -eq 0 ]
  [[ ${lines[0]} =~ \ +PID\ +TTY\ +STAT\ +TIME\ +COMMAND+ ]]
  [[ "${lines[1]}" =~ [0-9]+ ]]
}
