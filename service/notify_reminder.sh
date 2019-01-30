#!/usr/bin/env bash

# cron exec
# example) 0 12-18/2 * * * bash ~/cycle-reminder/service/notify_reminder.sh 2>> ~/cron_err.log
/usr/local/bin/docker-compose -f ~/cycle-reminder/service/docker-compose.yml exec web go run app/cron/notifyreminder.go