#!/bin/bash

FILE_CHANGE_LOG_FILE=/tmp/changes.log
SERVICE_ARGS="$@"

echo "*******$SERVICE_ARGS*********"

log() {
  echo "***** $1 *****"
}

init() {
  log "Initializing"
  truncate -s 0 ${FILE_CHANGE_LOG_FILE}
  tail -f ${FILE_CHANGE_LOG_FILE} &
}

build() {
  log "Building ${SERVICE_NAME} binary"
  go env -w GOPROXY="proxy.golang.org,direct"
  go mod download
  go build -gcflags "all=-N -l" -o /${SERVICE_NAME}
}

start() {
  log "Starting Delve"
  # ./entrypoint.sh serve all -c ../hydra/config.yml --dangerous-force-http
  # dlv --listen=:20001 --headless=true --api-version=2 --accept-multiclient exec /${SERVICE_NAME} -- ${SERVICE_ARGS} &
  /${SERVICE_NAME} ${SERVICE_ARGS} &
}

restart() {
  build

  log "Killing old processes"
  killall dlv
  killall ${SERVICE_NAME}

  start
}

watch() {
  log "Watching for changes"
  inotifywait -e "MODIFY,DELETE,MOVED_TO,MOVED_FROM" -m -r ${PWD} | (
    while true; do
      read path action file
      ext=${file: -3}
      if [[ "$ext" == ".go" ]]; then
        echo "$file"
      fi
    done
  ) | (
    WAITING=""
    while true; do
      file=""
      read -t 1 file
      if test -z "$file"; then
        if test ! -z "$WAITING"; then
          echo "CHANGED"
          WAITING=""
        fi
      else
        log "File ${file} changed" >> ${FILE_CHANGE_LOG_FILE}
        WAITING=1
      fi
    done
  ) | (
    while true; do
      read TMP
      restart
    done
  )
}

# main part
init
build
start
# watch
