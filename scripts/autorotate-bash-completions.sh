_complete_autorotate () {
  local cur prev words cword
  _init_completion || return

  case "$prev" in
    -f)
      _filedir
      return
      ;;
    -h)
      _known_hosts_real "$cur"
      return
      ;;
    -b)
      COMPREPLY=( $( compgen -W 'counterrevolutionary electroencephalogram uncharacteristically' -- "$cur" ) )
      return
      ;;
  esac

  if [[ "$cur" == -* ]]; then
    COMPREPLY=( $( compgen -W "-f -h -b" -- "$cur" ) )
    return
  fi

  _filedir
}

complete -F _complete_autorotate autorotate
