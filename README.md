# Service ChatServer

## Description of functionality
Main service features:
- Register user (via login and password) (`http://localhost:8080/account`)
- Authentication user via login and password the issuance of a JWT token (`http://localhost:8080/auth`)


### Description of the service workflow
1) New user creation
   1) При регистрации пользователя по ручке http://localhost:8080/account, создаём юзера в таблице users
2) Login search
   1) При поиске по логину другого юзера и клине по нему - создаётся чат в таблице chats, в таблице  users - добавляем этот чат
   одному и второму юзеру
   2) ВАЖНО!! Если такой чат уже есть - то ничего не создаём. Для этого: когда мы создаём чат, мы идём в таблицу users и берём чаты у обоих юзеров
   если есть общий чат - значит такой чат уже есть и мы возвращаем его. (TODO: добавить флаги для чатов как личные и групповые
   и смотреть есть ли общий чат с флагом личный)


## TODO
1) ~~Messages will be sent via websockets~~
2) Authentication in a separate service
3) How to work with the database? Whether to raise it in a separate service.
4) Вебсокеты: как работать с одной и той же сессией когда у нас нескольно нод и распределённая память. Сохранять её в редис?
5) Notification of events (new messages to offline users) via redis
6) ~~Make real time messaging (without events with redis)~~
7) ~~Link Websocket Connection to Client ID~~
8) ~~When a new websocket session created block ws client channel for writing messages until the client recieve all the updates~~
   1) Notifications via redis
   ~~2) Message history via postgres~~
9) ~~When user registered he sends device info, server generates hash from info and sends to client.
After that client sends given deviceID with websocket session.~~
10) Enable group chats (in progress)