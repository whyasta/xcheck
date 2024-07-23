#!/bin/sh

set -e

echo "start the app"
# exec "$@"

exec supervisord -c /etc/supervisor/conf.d/app.conf
