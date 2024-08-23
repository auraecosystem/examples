#!/bin/bash

set -eux

echo '{
  "content": "tell me about fuel savings"
}' | tr -d "\n" | curl \
    -X GET \
    -H 'Content-Type: application/json' \
    -d @- \
    http://localhost:9020/query/ | sed 's/\\n/\n/g'
