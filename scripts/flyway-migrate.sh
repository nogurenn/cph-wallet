#!/bin/bash -e

PORT=5432
if ! [[ -z $DB_PORT ]]; then
  PORT=$DB_PORT
fi

DB_URL=jdbc:postgresql://$DB_HOST:$PORT/$DB_NAME

echo "*** migrating schema: $SCHEMA"
bash flyway -url=$DB_URL -user=$DB_USER -password=$DB_PASSWORD -schemas=$SCHEMA -cleanDisabled=true migrate