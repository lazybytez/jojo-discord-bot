#!/usr/bin/env sh
set -e

PUID=${PUID:=0}
PGID=${PGID:=0}

chown -Rc "$PUID":"$PGID" /app/data
chown -Rc "$PUID":"$PGID" /app/log

exec setpriv --reuid "$PUID" --regid "$PGID" --clear-groups "$@"
