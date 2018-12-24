#!/usr/bin/env bash

echo "GRANT ALL ON *.* TO '"$MYSQL_USER"'@'%' ;" | "${mysql[@]}"
echo 'FLUSH PRIVILEGES ;' | "${mysql[@]}"