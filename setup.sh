#!/usr/bin/env bash

# This script will download the yaml files containing regexps from the piwik device-detector repository
rm -R ./piwik/
git clone -n https://github.com/piwik/device-detector.git --depth 1 piwik
cd piwik && git checkout HEAD regexes && rm -R ./.git/