#!/bin/sh
#
# Script for first-time setup and deployment of this application to Heroku.
#

[ -e "`which heroku`" ] || (
	echo "Could not find Heroku toolbelt."
	echo "See https://toolbelt.heroku.com/ for installation instructions."
); exit 1

set -x

heroku create --buildpack git://github.com/kr/heroku-buildpack-go.git
heroku addons:add mongolab:starter
heroku addons:add memcachier:25

time git push heroku master

heroku open
