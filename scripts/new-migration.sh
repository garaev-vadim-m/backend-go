#!/usr/bin/env bash
# scripts/new-migration.sh <название>
set -e
dir="internal/database/migrations"
name="$1"
if [ -z "$name" ]; then
  echo "usage: $0 <name>" >&2
  exit 1
fi

last=$(ls "$dir" | grep -oE '^[0-9]{6}' | sort | tail -1)
next=$(printf "%06d" $((10#${last:-0} + 1)))

touch "$dir/${next}_${name}.up.sql" "$dir/${next}_${name}.down.sql"
echo "created ${next}_${name}.{up,down}.sql"
