#!/usr/bin/env zsh

# Ensure that the version was supplied in the argument.
if [[ $# -lt 1 ]]; then
  echo "Release version must be supplied as the first argument to this script"
  return 1
fi

# Ensure that the git branch is at `main` and that there is no pending change to be committed and pushed
if ! command -v git >/dev/null; then
  echo "git is not installed. Cannot proceed"
  return 1
fi

if ! command -v docker >/dev/null; then
  echo "docker is not installed. Cannot proceed"
  return 1
fi

GIT_BRANCH=$(git branch  | grep '^*' | cut -d ' ' -f 2 | xargs)

if [[ $GIT_BRANCH != "main" ]]; then
  echo "Branch should be set to 'main' but it is set to ${GIT_BRANCH}"
  return 1
else
  echo "Branch is at main (as it should be)"
fi

if [ -z "$(git status --porcelain)" ]; then
  echo "Working directory clean (as it should be)"
else
  echo "Working directory is not clean. Cannot proceed"
  return 1
fi

# Check that there is a tmp directory or not. There must be.
if [ ! -d "tmp" ]; then
  echo "'tmp' directory does not exist. Are you running the script from the project root?"
  echo "Cannot continue"
  return 1
fi

# Ask the user to enter the version again so that a previous version does not start getting built and pushed by mistake
RELEASEVERSION=$1

echo "We are about to release version: '$RELEASEVERSION'"
printf "If you want to proceed, please input the version again: "

read inputversion

if [[ $RELEASEVERSION != $inputversion ]]; then
  echo "Versions passed via argument and typed do not match. Cannot continue"
  return 1
fi

# Create a new Directory named `bark_release_<version>` in the `tmp` directory. e.g. `barkk_release_v1.1.0`
mkdir -p "tmp/bark_release_$RELEASEVERSION"

if [ $? -ne 0 ]; then
  echo "Making release directory failed. Cannot continue"
  return 1
fi

# Now compile the codebase for all the specified targets. Output the resulting binaries into the newly created folder
echo "Building for GOOS=linux GOARCH=arm64"
GOOS=linux GOARCH=arm64 go build -o tmp/bark_release_$RELEASEVERSION/bark_${RELEASEVERSION}_linux_arm64 ./cmd/server
if [ $? -ne 0 ]; then
  echo "Could not compile for linux/arm64"
  return 1
fi

echo "Building for GOOS=linux GOARCH=amd64"
GOOS=linux GOARCH=amd64 go build -o tmp/bark_release_$RELEASEVERSION/bark_${RELEASEVERSION}_linux_amd64 ./cmd/server
if [ $? -ne 0 ]; then
  echo "Could not compile for linux/amd64"
  return 1
fi

echo "Building for GOOS=darwin GOARCH=arm64"
GOOS=darwin GOARCH=arm64 go build -o tmp/bark_release_$RELEASEVERSION/bark_${RELEASEVERSION}_macos_arm64 ./cmd/server
if [ $? -ne 0 ]; then
  echo "Could not compile for darwin/arm64"
  return 1
fi

echo "Building for GOOS=darwin GOARCH=amd64"
GOOS=darwin GOARCH=amd64 go build -o tmp/bark_release_$RELEASEVERSION/bark_${RELEASEVERSION}_macos_amd64 ./cmd/server
if [ $? -ne 0 ]; then
  echo "Could not compile for darwin/amd64"
  return 1
fi

# Push the tag to GitHub.
echo "Tagging"
git tag -a $RELEASEVERSION -m "Release $RELEASEVERSION"
echo "Pushing"
git push -f origin refs/tags/$RELEASEVERSION

# Run a Docker multi-platform build along with the push
echo "Building the docker image for the release"
docker buildx build --platform linux/amd64,linux/arm64 -t techrail/bark:$RELEASEVERSION --push .

printf "Do you build the images with the 'latest' tag? (yes/no) "
read choice

if [[ $choice != "yes" ]]; then
  echo "You chose not to tag this release as the latest one"
  return 0
else
  docker buildx build --platform linux/amd64,linux/arm64 -t techrail/bark:latest --push .
fi


