#!/bin/bash

# Path to the config file
CONFIG_FILE="./mmhealth/files/config.yaml"

if [[ -z "$GITHUB_TOKEN" ]]; then
  echo "Error: GITHUB_TOKEN is not set or is empty"
  exit 1
fi

echo "Current working directory: $(pwd)"
echo "Current working directory: $GITHUB_TOKEN"
# Loop over each plugin
for plugin in $(yq e '.plugins | keys' $CONFIG_FILE); do

  # Skip if the plugin name is '-'
  if [ "$plugin" == "-" ]; then
    continue
  fi

  # Remove quotes from the plugin name
  plugin=${plugin//\"/}

  # Get the repo URL for the plugin
  repo=$(yq e ".plugins.[\"$plugin\"].repo" $CONFIG_FILE)

  # Remove quotes from the repo URL
  repo=${repo//\"/}

  # Extract the owner and repo name from the repo URL
  owner=$(echo $repo | cut -d'/' -f4)
  repo_name=$(echo $repo | cut -d'/' -f5)

  # Call the GitHub API to get the latest release
  response=$(curl --silent -H "Authorization: token $GITHUB_TOKEN" "https://api.github.com/repos/$owner/$repo_name/releases/latest" )
  # Check if the "tag_name" field is available
  if echo "$response" | jq -e .tag_name > /dev/null; then
    release=$(echo "$response")
  else
    # Fall back to the /releases endpoint if the latest release isn't available
    response=$(curl --silent -H "Authorization: token $GITHUB_TOKEN" "https://api.github.com/repos/$owner/$repo_name/releases")
    release=$(echo "$response" | jq -r '.[0]')
  fi
  echo "https://api.github.com/repos/$owner/$repo_name/releases/latest"

  latest_release=$(echo "$release" | jq -r .tag_name)
  # Remove 'v' from the latest release
  latest_release=$(echo $latest_release | sed 's/v//')
  latest_release=${latest_release//\"/}

  release_date=$(echo "$release" | jq -r .published_at)

  echo "Latest release: $latest_release for $plugin date: $release_date"

  # Update the 'latest' key for the plugin
  yq -i ".plugins[\"$plugin\"].latest = \"$latest_release\""  $CONFIG_FILE
  yq -i ".plugins[\"$plugin\"].release_date = \"$release_date\""  $CONFIG_FILE

done