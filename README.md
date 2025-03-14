# Service: ChatServer

## Description of Functionality
Main features of the service include:
- **User Registration:** Users can register using a login and password. Endpoint: [`http://localhost:8080/account`].
- **User Authentication:** Logging in with a login and password generates a JWT token. Endpoint: [`http://localhost:8080/auth`].

---

## Description of the Service Workflow

1) **Creating a New User**
   - When registering a new user through the `http://localhost:8080/account` endpoint, the user's details are added to the `users` table.

2) **Finding a User by Login**
   - When searching for another user by their login and clicking on their profile:
      1. A new chat entry is created in the `chats` table, and this chat is linked to both users within the `users` table.
      2. **IMPORTANT:** If a chat between these two users already exists, no new chat is created. How this works:
         - Upon attempting to create a chat, check the existing chats in the `users` table for both users.
         - If a shared chat already exists, return the existing chat instead of creating a new one.

     *(TODO: Add flags for chats to distinguish between personal chats and group chats. Additionally, check if a shared personal chat already exists by examining the personal chat flag.)*

---

## TODO

1) ~~Messages should be sent via WebSockets.~~
2) Move the **Authentication** functionality to a separate service.
3) Evaluate how to handle the database setup. Should the database be managed as a separate service?
4) **WebSockets:** Address how to deal with shared sessions when multiple nodes are used and distributed memory is required. Should sessions be stored in Redis?
5) Implement **event notifications for offline users** (e.g., new message notifications) using Redis.
6) ~~Implement real-time messaging (without relying on Redis for events).~~
7) ~~Link each WebSocket connection to a specific Client ID.~~
8) ~~When a new WebSocket session is created, block the WebSocket client channel for writing messages until the client has received all pending updates.~~
   - Notifications via Redis.
   - ~~Message history stored using PostgreSQL.~~
9) ~~During user registration, the client sends device information to the server. The server generates a hash based on this information and sends it back to the client. Afterward, the client includes the generated `deviceID` with every WebSocket session.~~
10) Enable **group chats** (currently in progress).

# Websocket usefull info
1) https://www.reddit.com/r/kubernetes/comments/sssof0/k8_with_persistent_websocket_connections/
2) https://cloud.google.com/kubernetes-engine/docs/concepts/ingress#support_for_websocket

# How to do local deploy
1) Install kuberctl : choco install kubernetes-cli, brew install kubectl
2) Install minikube : choco install minikube or brew install minikube (don't forget to add it to PATH)
3) Run 'make create-app-docker-image' it will create docker image of application
4) Run 'make setup-minikube'
5) Run 'make deploy-all'
6) Add this lines to your host file (on Mac it's /etc/hosts)
```
127.0.0.1 myapp.local  
127.0.0.1 db.yourdomain.com
127.0.0.1 cache.yourdomain.com
127.0.0.1 pubsub.yourdomain.com
127.0.0.1 redis.gui.com
127.0.0.1 postgres.gui.com
```
7) Run 'minikube tunnel'

# To view logs 
1) Run 'make kuber-gui' it will automatically open a kubernetes GUI in your browser
2) On the top choose namespace 'enrollment'
3) On the right side press on pods, you will see all the backend services, including 3 pods of application. Open each of them and on the right side press on logs. 