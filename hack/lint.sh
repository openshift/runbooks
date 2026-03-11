#!/bin/bash
# shellcheck disable=SC2015,SC1091
set -e

usage(){
  echo "
  usage: $0
  "
}

failed(){
  echo "
    failed: ${1:-lint}
  "
  return 1
}

py_setup_venv(){
  python3 -m venv venv
  source venv/bin/activate
  pip install -q -U pip

  py_check_venv || usage
}

py_check_venv(){
  # activate python venv
  [ -d venv ] && source venv/bin/activate || py_setup_venv
  [ -e $(dirname "$0")/requirements.txt ] && pip install -q -r $(dirname "$0")/requirements.txt
}

py_bin_checks(){
  which python || exit 0
  which pip || exit 0
}

lint_spelling(){
  which aspell || return
  which pyspelling || return
  [ -e .pyspelling.yml ] || return
  [ -e .wordlist-md ] || return

  pyspelling -c .pyspelling.yml
}

lint_markdown(){
  which pymarkdown || fail "no pymarkdown"
  pymarkdown scan --recurse .
}

lint_init(){
  mkdir -p scratch
}

fix_markdown(){
  which pymarkdown || failed "no pymarkdown"
  pymarkdown fix --recurse .
}

lint(){
  lint_spelling
  lint_markdown
}

py_check_venv
py_bin_checks

lint_init

FUNCTION=$(basename -s .sh ${0})
"${FUNCTION}" 0 || failed "${FUNCTION}"
