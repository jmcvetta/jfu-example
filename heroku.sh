#!/bin/sh
#
# Script for first-time setup and deployment of this application to Heroku.
#

set -x

heroku create --buildpack git://github.com/kr/heroku-buildpack-go.git
heroku addons:add mongolab:starter
heroku addons:add memcachier:25

time git push heroku master

heroku open
