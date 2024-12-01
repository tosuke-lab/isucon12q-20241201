#!/bin/bash

set -xe

rev=${1:-HEAD}
hash=$(git rev-parse --short "$rev")
rule='^/api/player/competition/[\w-]+/ranking$,^/api/player/player/[\w-]+$,^/api/organizer/competition/[\w-]+$,^/api/organizer/competition/[\w-]+/score$,^/api/organizer/competition/[\w-]+/finish$,^/api/organizer/player/[\w-]+/disqualified$'

sudo alp ltsv --file /var/log/nginx/access."$hash".log -m "$rule"
