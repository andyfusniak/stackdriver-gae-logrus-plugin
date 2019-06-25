# CHANGELOG
## v0.1.3 (Tue June 25, 2019)
+ Fix off by one error when parsing `X-Cloud-Trace-Context`.

## v0.1.2
+ Fix panic occuring with `X-Cloud-Trace-Context` values with missing `;o=1` part.

## v0.1.1
+ Adds `app.yaml` file and fixes Makefile for gcloud app deploy.
+ Prevents crash that occured from outside a non HTTP request.
+ Fix internal typo.

## v0.1.0
+ Initial release
