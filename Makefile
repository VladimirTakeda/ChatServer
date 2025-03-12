# Variables
REDIS_CONTAINER_CACHE=chatserver-cache-1
REDIS_PASSWORD_CACHE=eYVX7EwVmmxKPCDmwMtyKVge8oLd2t81

REDIS_CONTAINER_PUBSUB=chatserver-redis_pubsub-1
REDIS_PASSWORD_PUBSUB=eYVX7EwVmmxKPCDmwMtyKVge8oLd2t81


CONTAINER_NAME="chatserver-db-1"
USERNAME="postgres"
DATABASE_NAME="postgres"
DATABASE_PASSWORD="password"

# SQL Script
SQL_QUERY="TRUNCATE TABLE\
               message,\
               user_last_seen,\
               last_message,\
               chats,\
               users,\
               devices\
           RESTART IDENTITY CASCADE;"


# Clean Cache (because the data storage is persistent even if container has been deleted)
.PHONY: flush-cache-docker
flush-cache-docker:
	docker exec -it $(REDIS_CONTAINER_CACHE) redis-cli -a $(REDIS_PASSWORD_CACHE) FLUSHALL

.PHONY: flush-pubsub-docker
flush-pubsub-docker:
	docker exec -it $(REDIS_CONTAINER_PUBSUB) redis-cli -a $(REDIS_PASSWORD_PUBSUB) FLUSHALL

.PHONY: flush-postgres-docker
flush-postgres-docker:
	docker exec -i $(CONTAINER_NAME) env PGPASSWORD=$(DATABASE_PASSWORD) psql -U $(USERNAME) -d $(DATABASE_NAME) -c $(SQL_QUERY)

.PHONY: flush-all-docker
flush-all:
	flush-cache-docker
	flush-pubsub-docker
	flush-postgres-docker


# Variables
REDIS_CACHE_POD=$(shell kubectl get pod -l app=cache -n enrollment -o jsonpath='{.items[0].metadata.name}')
REDIS_PUBSUB_POD=$(shell kubectl get pod -l app=pubsub -n enrollment -o jsonpath='{.items[0].metadata.name}')
POSTGRES_POD=$(shell kubectl get pod -l app=db -n enrollment -o jsonpath='{.items[0].metadata.name}')


REDIS_PASSWORD_CACHE=eYVX7EwVmmxKPCDmwMtyKVge8oLd2t81
REDIS_PASSWORD_PUBSUB=eYVX7EwVmmxKPCDmwMtyKVge8oLd2t81

USERNAME="postgres"
DATABASE_NAME="postgres"
DATABASE_PASSWORD="password"

# Clean Redis Cache
.PHONY: flush-cache-kuber
flush-cache-kuber:
	kubectl config set-context --current --namespace=enrollment
	kubectl exec -it $(REDIS_CACHE_POD) -- redis-cli -a $(REDIS_PASSWORD_CACHE) FLUSHALL

# Clean Redis PubSub
.PHONY: flush-pubsub-kuber
flush-pubsub-kuber:
	kubectl exec -it $(REDIS_PUBSUB_POD) -- redis-cli -a $(REDIS_PASSWORD_PUBSUB) FLUSHALL

# Clean Postgres
.PHONY: flush-postgres-kuber
flush-postgres-kuber:
	kubectl exec -i $(POSTGRES_POD) -- env PGPASSWORD=$(DATABASE_PASSWORD) psql -U $(USERNAME) -d $(DATABASE_NAME) -c $(SQL_QUERY)

# Clean All
.PHONY: flush-all-kuber
flush-all-kuber: flush-cache-kuber flush-pubsub-kuber flush-postgres-kuber


 #(this builds image direcly in minikube) eval $$(minikube -p minikube docker-env)
.PHONY: create-app-docker-image
create-app-docker-image:
	eval $$(minikube -p minikube docker-env) && \
	docker build -t myapp:latest -f ./docker/backend_docker/Dockerfile .

.PHONY: verify-image
verify-image:
	minikube ssh "docker images | grep myapp"

.PHONY: deploy-all
deploy-all:
	kubectl	create namespace enrollment
	kubectl	apply -f k8s/ -n enrollment
	kubectl get all -n enrollment

.PHONY: deploy-update
deploy-update:
	kubectl	apply -f k8s/ -n enrollment
	kubectl get all -n enrollment

.PHONY: redeploy
redeploy:
	kubectl rollout restart deployment application -n enrollment

.PHONY: check
check:
	kubectl get pods -n enrollment
	kubectl get svc -n enrollment

.PHONY: logs
logs:
	kubectl logs -f deployment/application -n enrollment

.PHONY: clean
clean:
	kubectl delete namespace enrollment

.PHONY: check-kuber-status
check-kuber-status:
	kubectl get pods -n enrollment -w

.PHONY: setup-minikube
setup-minikube:
	minikube stop
	minikube start --cpus=6 --memory=12288 --driver=docker
	minikube status
	minikube addons enable ingress
	minikube addons list
	minikube status

.PHONY: check-ports
check-ports:
	minikube service application -n enrollment

.PHONY: run-tunnel
run-tunnel:
	minikube tunnel

.PHONY: kuber-gui
kuber-gui:
	minikube dashboard