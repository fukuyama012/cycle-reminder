FROM mysql:5.7

RUN apt-get update -qq

# add custom
RUN echo '\n\
[mysql]\n\
default-character-set=utf8mb4\n\
\n\
[mysqld]\n\
character-set-server=utf8mb4\n\
collation-server=utf8mb4_bin\n\
\n\
server-id=1\n\
log_bin=mysql-bin\n\
binlog-format=ROW\n\
expire_logs_days=14\n\
\n\
gtid_mode=ON\n\
enforce-gtid-consistency\n\
\n\
slow_query_log=true\n\
slow_query_log-file=/var/log/mysql/mysql-slow.sql\n\
long_query_time=3\n\
\n\
[client]\n\
default-character-set=utf8mb4\n' >> /etc/mysql/my.cnf