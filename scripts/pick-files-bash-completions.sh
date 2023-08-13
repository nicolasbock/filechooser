#!/bin/bash

_complete_pick_files () {
  local cur prev
  local known_options=(
    -N --number
    --block-selection
    --debug
    --destination
    --destination-option
    --dry-run
    --folder
    -h --help
    --print-database
    --print-database-format
    --reset-database
    --suffix
    --verbose
    --version
  )

  _init_completion || return

  case "$prev" in
    --folder|--destination)
      _filedir
      return
      ;;
    --print-database-format)
      readarray -t COMPREPLY < <(compgen -W 'CSV JSON YAML' -- "${cur}")
      return
      ;;
    --destination-option)
      readarray -t COMPREPLY < <(compgen -W 'panic delete append' -- "${cur}")
  esac

  if [[ "$cur" == -* ]]; then
    readarray -t COMPREPLY < <(compgen -W "${known_options[*]}" -- "${cur}" )
    return
  fi

  _filedir
}

complete -F _complete_pick_files pick-files
