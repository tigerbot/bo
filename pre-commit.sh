#!/bin/sh
#
# A git pre-commit hook to format the go files added
# and make sure no files with invalid names are committed.
#
# To enable this hook, add a symbolic link to this file in your
# repository's .git/hooks directory with the name "pre-commit"

exec 1>&2
set -e
set -u

# The below check to make sure file names don't include unusual characters was
# modified from the example pre-commit hook found in any recently created repo
NEW_NAMES=$(git diff --cached --name-only --diff-filter=ACR -z)
if echo "${NEW_NAMES}" | grep -q [^a-zA-Z0-9\.\_\\\/\-]; then
  echo "Error: Attempt to add an invalid file name"
  echo "${NEW_NAMES}" | grep [^a-zA-Z0-9\.\_\\\/\-]
  echo
  echo "File names can only include letters, numbers,"
  echo "hyphens, underscores and periods"
  echo
  echo "Please rename your file before committing it"
  exit 1
fi

CHANGED_GO=$(git diff --cached --name-only --diff-filter=ACRM | grep '\.go$') || true
if ! [ -n "${CHANGED_GO}" ]; then
  # Our unit tests are currently only for go and we can only programatically format go code,
  # so if nothing in the go source changed we really have nothing to check for this commit.
  exit 0
fi

# We want to test what is being committed not what's on the file system, so we stash everything.
# This does mean that we can't automatically fix any formatting problems since that would prevent
# us from cleanly restoring the stash, but this is a relatively small price to pay to test
# exactly what's being committed.
old_stash=$(git rev-parse -q --verify refs/stash)
git stash save --include-untracked --keep-index --quiet "pre-commit"
new_stash=$(git rev-parse -q --verify refs/stash)

# http://stackoverflow.com/questions/20479794/how-do-i-properly-git-stash-pop-in-pre-commit-hooks-to-get-a-clean-working-tree
# see above link for explanation of the check on stash states and the reset
if [ "${old_stash}" != "${new_stash}" ]; then
  cleanup() {
    set +e
    git reset --hard -q
    git stash pop --quiet --index
  }
  trap cleanup EXIT
fi

# Run the unit tests with the code that will be committed
./test.sh
echo

# find any of the staged files that differ from what gofmt produces
for file in "${CHANGED_GO}"; do
  # I am primarily developing on Windows, so I need to make sure gofmt doesn't report the
  # formatting as different because of the line endings.
  if [ -n "$(sed 's|\r\n|\n|g' ${file} | gofmt -l)" ]; then
    UNFORMATTED="${UNFORMATTED:-}\n${file}"
  fi
done
if [ -n "${UNFORMATTED:-}" ]; then
  gofmt -w ${CHANGED_GO}
  echo -e "improperly formatted GO files:${UNFORMATTED}"
  exit 1
fi

exit 0
