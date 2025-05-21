set dotenv-load

root_dir := justfile_directory()
src_dir := root_dir + "/src"

run server *args:
  #!/usr/bin/env bash
  set -euox pipefail
  src_file="{{ src_dir }}/main.go"
  bin_file="{{ root_dir }}/tmp/main"
  if [[ ! -f "$src_file" ]]; then
    echo "Could not find the service at $src_file"
  else
        cd "{{ src_dir }}" && air \
          -build.include_dir="go" \
          -build.bin="${bin_file}" \
          -build.cmd="go build -o ${bin_file} ${src_file}" \
          "$@"
  fi