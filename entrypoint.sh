#!/bin/sh

# check to see if secrets exist
if [ -d "/run/secrets" ]
then
  echo "Found secrets"

  # export secrets to environment variables
  for i in $(ls /run/secrets)
  do
    echo "Exporting secret (${i}) to an env var..."
    export ${i}="$(cat /run/secrets/${i})"
  done
else
  echo "No secrets found; assuming environment variables are being used"
fi

exec "${@}"
