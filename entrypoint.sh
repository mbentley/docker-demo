#!/bin/sh

# check to see if secrets exist
if [ -d "/run/secrets" ]
then
  echo "Found secrets"

  # export secrets to environment variables
  for i in /run/secrets/*
  do
    echo "Exporting secret (${i}) to an env var..."
    export "$(basename "${i}")"="$(cat "${i}")"
  done
else
  echo "No secrets found; assuming environment variables are being used"
fi

# check and see if postgres is running yet using nc
while [ "$(nc -z -v -w1 db 5432 > /dev/null 2>&1; echo $?)" -ne "0" ]
do
  echo "Database not ready."
  sleep 2
done

echo "Database ready! Starting app..."
sleep 2

exec "${@}"
