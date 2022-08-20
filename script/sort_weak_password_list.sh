#!/bin/bash

ROOT_DIR=$(git rev-parse --show-toplevel)
WEAK_PASSWD_LIST="${ROOT_DIR}/internal/embedded/passwd/weak.txt"
TMP_FILE=$(mktemp)

sort "${WEAK_PASSWD_LIST}" | uniq >"${TMP_FILE}"
mv "${TMP_FILE}" "${WEAK_PASSWD_LIST}"
