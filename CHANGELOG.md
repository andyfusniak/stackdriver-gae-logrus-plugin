# CHANGELOG

## v0.2.1 (Tue October 29, 2024)
+ golang.org/x/sys/unix has Incorrect privilege reporting in syscall
  + Go before 1.17.10 and 1.18.x before 1.18.2 has Incorrect Privilege Reporting in syscall. When called with a non-zero flags parameter, the Faccessat function could incorrectly report that a file is accessible.
  + Specific Go Packages Affected
golang.org/x/sys/unix

## v0.2.0 (Tue October 29, 2024)
+ Support for v2 `X-Cloud-Trace-Context` header that uses 16-digit hex strings for its Span ID. See the README.md for more information.

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
