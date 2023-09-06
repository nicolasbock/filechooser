#!/bin/bash

_complete_pick_files () {
  local cur prev
  local known_options=(
    -N --number
    --append
    --block-selection
    --config
    --debug
    --delete-existing
    --destination
    --destination-option
    --dry-run
    --dump-configuration
    --folder
    -h --help
    --journald
    --print-database
    --print-database-format
    --print-database-statistics
    --reset-database
    --suffix
    --verbose
    --version
  )

  _init_completion || return

  case "$prev" in
    --config)
      _filedir
      return
      ;;
    --destination-option)
      readarray -t COMPREPLY < <(compgen -W 'panic delete append' -- "${cur}")
      return
      ;;
    --folder|--destination)
      _filedir
      return
      ;;
    --print-database-format)
      readarray -t COMPREPLY < <(compgen -W 'CSV JSON YAML' -- "${cur}")
      return
      ;;
  esac

  if [[ "$cur" == -* ]]; then
    readarray -t COMPREPLY < <(compgen -W "${known_options[*]}" -- "${cur}" )
    return
  fi
}

complete -F _complete_pick_files pick-files
