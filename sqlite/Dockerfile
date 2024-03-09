FROM alpine:latest

RUN apk --no-cache add sqlite

WORKDIR /db

COPY ./sqlite/init.sql /db

RUN sqlite3 /db/1.db < init.sql
RUN sqlite3 /db/2.db < init.sql
RUN sqlite3 /db/3.db < init.sql
RUN sqlite3 /db/4.db < init.sql
RUN sqlite3 /db/5.db < init.sql

RUN sqlite3 /db/1.db "INSERT INTO u (id, l, b) VALUES (1, 100000, 0);"
RUN sqlite3 /db/2.db "INSERT INTO u (id, l, b) VALUES (2, 80000, 0);"
RUN sqlite3 /db/3.db "INSERT INTO u (id, l, b) VALUES (3, 1000000, 0);"
RUN sqlite3 /db/4.db "INSERT INTO u (id, l, b) VALUES (4, 10000000, 0);"
RUN sqlite3 /db/5.db "INSERT INTO u (id, l, b) VALUES (5, 500000, 0);"
