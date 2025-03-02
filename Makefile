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
.PHONY: flush-cache
flush-cache:
	docker exec -it $(REDIS_CONTAINER_CACHE) redis-cli -a $(REDIS_PASSWORD_CACHE) FLUSHALL

.PHONY: flush-pubsub
flush-pubsub:
	docker exec -it $(REDIS_CONTAINER_PUBSUB) redis-cli -a $(REDIS_PASSWORD_PUBSUB) FLUSHALL

.PHONY: flush-postgres
flush-postgres:
	docker exec -i $(CONTAINER_NAME) env PGPASSWORD=$(DATABASE_PASSWORD) psql -U $(USERNAME) -d $(DATABASE_NAME) -c $(SQL_QUERY)

.PHONY: flush-all
flush-all:
	docker exec -it $(REDIS_CONTAINER) redis-cli -a $(REDIS_PASSWORD) FLUSHALL