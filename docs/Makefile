default: install

help:
	@egrep '(^\S)|^$$' Makefile

install:
	bundle config set --local path vendor/bundle
	bundle install

upgrade:
	bundle update

b build:
	bundle exec jekyll build --source . --destination build/

s serve:
	bundle exec jekyll serve --source . --destination build/ --livereload --trace
