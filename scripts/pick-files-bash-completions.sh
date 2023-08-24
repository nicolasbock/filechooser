#!/bin/bash

_complete_pick_files () {
  local cur prev
  local known_options=(
    --block-selection
    --config
    --debug
    --destination
    --destination-option
    --dump-configuration
    --dry-run
    --folder
    --journald
    --print-database
    --print-database-format
    --print-database-statistics
    --reset-database
    --suffix
    --verbose
    --version
    -h --help
    -N --number
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

  _filedir
}

complete -F _complete_pick_files pick-files
