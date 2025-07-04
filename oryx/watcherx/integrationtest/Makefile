SHELL=/bin/bash -euo pipefail

CLUSTER_NAME=watcherx-integration-test
SNAPSHOT_FILE=eventlog_snapshot

define generate_snapshot
	sleep 5
	make update
	sleep 1
	kubectl logs eventlogger --context kind-${CLUSTER_NAME} >> $(1)
	make apply
	sleep 1
	kubectl logs eventlogger --context kind-${CLUSTER_NAME} >> $(1)
	make update
	sleep 1
	kubectl logs eventlogger --context kind-${CLUSTER_NAME} >> $(1)
endef

.PHONY: build
build:
	docker build -f Dockerfile -t eventlogger:latest ../..

.PHONY: create
create:
	kind create cluster --name ${CLUSTER_NAME} --wait 1m || true

.PHONY: load
load:
	kind load docker-image eventlogger:latest --name ${CLUSTER_NAME}

.PHONY: apply
apply:
	kubectl apply -f configmap.yml -f event_logger.yml --context kind-${CLUSTER_NAME}

.PHONY: delete
delete:
	kind delete cluster --name ${CLUSTER_NAME}

.PHONY: setup
setup: build create load apply

.PHONY: snapshot
snapshot: setup container-restart
	rm ${SNAPSHOT_FILE}
	${call generate_snapshot,$(SNAPSHOT_FILE)}

.PHONY: check
check: setup container-restart
	rm tmp_snapshot || true
	${call generate_snapshot,tmp_snapshot}
	diff tmp_snapshot ${SNAPSHOT_FILE}

.PHONY: logs
logs:
	kubectl logs eventlogger --context kind-${CLUSTER_NAME}

.PHONY: container-restart
container-restart:
	kubectl delete -f event_logger.yml --context kind-${CLUSTER_NAME}
	kubectl apply -f event_logger.yml --context kind-${CLUSTER_NAME}

.PHONY: update
update:
	cat configmap.yml | sed 's/somevalue/othervalue/' | kubectl apply -f - --context kind-${CLUSTER_NAME}
	cat event_logger.yml | sed 's/somevalue/othervalue/' | kubectl apply -f - --context kind-${CLUSTER_NAME}
