#!/usr/bin/env bash
# Copyright 2024 Nitro Agility S.r.l.
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#     http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.
#
# SPDX-License-Identifier: Apache-2.0

# Writes THIRD_PARTY_NOTICES.md from what the build actually resolves, and checks it is current.
#
# # Why this reads the resolved dependencies rather than the manifest
#
# The manifest says what this project asks for; the resolved set says what a build receives — the
# transitive closure, at the exact versions. A notices file written from the manifest names a
# fraction of what ships, which is worse than no file: it looks like disclosure and is not.
#
# # Why the output is sorted
#
# The file is the input to a CI check that fails when it drifts. That check is only meaningful if a
# regeneration with unchanged dependencies produces byte-identical output, so entries are sorted
# and nothing timestamped or machine-specific is written into it.
#
# # Where the licence comes from
#
# Go carries no licence field in its metadata, so the licence has to be read out of each module in
# the module cache. The classifier below matches only the phrases that distinguish the handful of
# licences this dependency set actually uses, and says so rather than guessing when it does not
# recognise one — a wrong SPDX identifier in a notices file is worse than an honest gap.

set -euo pipefail

cd "$(git rev-parse --show-toplevel)"

readonly OUTPUT="THIRD_PARTY_NOTICES.md"

usage() {
    cat >&2 <<'USAGE'
usage: third-party-notices.sh [--check]

  (no argument)  regenerate THIRD_PARTY_NOTICES.md in place
  --check        fail if the file is not what a regeneration would write
USAGE
}

checking="false"
case "${1-}" in
    "") ;;
    --check) checking="true" ;;
    -h | --help)
        usage
        exit 0
        ;;
    *)
        usage
        exit 2
        ;;
esac

require() {
    if ! command -v "$1" >/dev/null 2>&1; then
        printf 'error: %s is required to write the third-party notices.\n' "$1" >&2
        printf 'install it with: %s\n' "$2" >&2
        exit 1
    fi
}

require go "https://go.dev/dl/"

# `go list -deps ./...` walks the packages this module actually builds, which is a far smaller set
# than `go list -m all`: the latter is the whole module graph, including modules nothing here
# imports. A notice covers what ships, so the reachable set is the right one.
#
# The sources have to be in the module cache for their licence files to be readable.
go mod download

readonly OWN="$(go list -m 2>/dev/null || true)"

# Names the licence in a file, by the phrases that tell these licences apart. Conservative on
# purpose: an unrecognised text is reported as unrecognised, never guessed.
classify() {
    local file="$1"

    if grep -qi "Apache License" "${file}" && grep -q "Version 2.0" "${file}"; then
        printf 'Apache-2.0'
    elif grep -q "Redistribution and use in source and binary forms" "${file}"; then
        if grep -qi "Neither the name of" "${file}"; then
            printf 'BSD-3-Clause'
        else
            printf 'BSD-2-Clause'
        fi
    elif grep -q "Permission is hereby granted, free of charge" "${file}"; then
        printf 'MIT'
    elif grep -q "Permission to use, copy, modify, and/or distribute" "${file}"; then
        printf 'ISC'
    elif grep -qi "Mozilla Public License" "${file}" && grep -q "2.0" "${file}"; then
        printf 'MPL-2.0'
    else
        printf ''
    fi
}

collected="$(
    go list -deps -f '{{if .Module}}{{printf "%s\t%s\t%s" .Module.Path .Module.Version .Module.Dir}}{{end}}' ./... 2>/dev/null \
        | sort -u \
        | while IFS="$(printf '\t')" read -r path version directory; do
            [ -z "${path}" ] && continue
            [ "${path}" = "${OWN}" ] && continue

            licence=""
            if [ -n "${directory}" ] && [ -d "${directory}" ]; then
                file="$(find "${directory}" -maxdepth 1 -type f \
                    \( -iname 'LICENSE*' -o -iname 'LICENCE*' -o -iname 'COPYING*' \) \
                    | LC_ALL=C sort | head -1)"
                if [ -n "${file}" ]; then
                    licence="$(classify "${file}")"
                fi
            fi

            printf '%s\t%s\t%s\thttps://%s\n' "${path}" "${version}" "${licence}" "${path}"
        done
)"

# One line per dependency, tab separated: name, version, licence, source. Sorted here so the
# collection step above never has to care about order.
rows="$(printf '%s' "${collected}" | LC_ALL=C sort -u)"
count="$(printf '%s' "${rows}" | grep -c . || true)"

table="$(
    printf '%s\n' "${rows}" | awk -F'\t' 'NF >= 3 {
        name = $1; version = $2; licence = $3; source = $4;
        if (licence == "" || licence == "null" || licence == "UNKNOWN") licence = "not declared";
        if (source == "" || source == "null") source = "—";
        printf "| `%s` | %s | %s | %s |\n", name, version, licence, source;
    }'
)"

undeclared="$(
    printf '%s\n' "${rows}" | awk -F'\t' '
        $3 == "" || $3 == "null" || $3 == "UNKNOWN" {
            printf "- `%s` %s — %s\n", $1, $2, ($4 == "" || $4 == "null" ? "no source declared either" : $4);
        }'
)"

if [ -z "${undeclared}" ]; then
    undeclared_section="Every package above declares a licence."
else
    undeclared_section="$(
        cat <<UNDECLARED
The packages below declare no licence in their metadata. That is usually an upstream omission
rather than an absence of licence: check the licence file in each source tree before a distribution
that relies on this list.

${undeclared}
UNDECLARED
    )"
fi

rendered="$(
    cat <<HEADER
# Third-Party Notices

The Permguard Go SDK is distributed under the Apache License, Version 2.0. It depends on the
third-party modules listed below, each under its own licence.

This file is generated from the resolved dependencies and is checked in CI. Do not edit it by hand:
run \`task notices\` instead.

Development and test dependencies are excluded where the ecosystem distinguishes them: a notice
covers what is distributed, and a test harness is not.

## Packages

${count} packages.

| Package | Version | Licence | Source |
| ------- | ------- | ------- | ------ |
${table}

## Packages without a declared licence

${undeclared_section}

## Full licence texts

The full text of the Apache License 2.0 is in [LICENSE](LICENSE). The texts of the other licences
named above are published by their respective projects at the sources listed.

For licence questions, contact <opensource@permguard.com>.
HEADER
)"

if [ "${checking}" = "true" ]; then
    if [ ! -f "${OUTPUT}" ]; then
        printf 'error: %s does not exist. Run `task notices` and commit it.\n' "${OUTPUT}" >&2
        exit 1
    fi
    if ! printf '%s\n' "${rendered}" | diff -u "${OUTPUT}" - >/dev/null; then
        printf 'error: %s is out of date with the resolved dependencies.\n\n' "${OUTPUT}" >&2
        printf '%s\n' "${rendered}" | diff -u "${OUTPUT}" - >&2 || true
        printf '\nRun `task notices` and commit the result.\n' >&2
        exit 1
    fi

    printf 'ok: %s matches the resolved dependencies (%s packages)\n' "${OUTPUT}" "${count}"
    exit 0
fi

printf '%s\n' "${rendered}" >"${OUTPUT}"
printf 'wrote %s (%s packages)\n' "${OUTPUT}" "${count}"
