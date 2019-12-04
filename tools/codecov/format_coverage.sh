#!/usr/bin/env bash

######################################################################################################
# Original script from : https://github.com/Azure/dcos-engine/blob/master/scripts/ginkgo.coverage.sh #
######################################################################################################

set -eo pipefail

coverage_mode=${COVERMODE:-atomic}
coverage_dir=$(mktemp -d /tmp/coverage.XXXXXXXX.tmp)
profile="${coverage_dir}/cover.out"
coverage_file="coverage.txt"

format_cover_data() {
  echo "" > ${coverage_file}
  find . -type f -name "*.coverprofile" | while read -r file;  do cat $file >> ${coverage_file} && mv $file ${coverage_dir}; done
  echo "mode: $coverage_mode" >"$profile"
  grep -h -v "^mode:" "$coverage_dir"/*.coverprofile >>"$profile"
}

format_cover_data
go tool cover -func "${profile}"