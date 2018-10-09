<!-- START doctoc generated TOC please keep comment here to allow auto update -->
<!-- DON'T EDIT THIS SECTION, INSTEAD RE-RUN doctoc TO UPDATE -->
**Table of Contents**  *generated with [DocToc](https://github.com/thlorenz/doctoc)*

- [Upgrade Guide](#upgrade-guide)
  - [0.3.0](#030)

<!-- END doctoc generated TOC please keep comment here to allow auto update -->

# Upgrade Guide

## 0.3.0

To improve how errors are forwarded to clients, two `Writer` interface methods changed from:

* `WriteError(w http.ResponseWriter, r *http.Request, err error)`
* `WriteErrorCode(w http.ResponseWriter, r *http.Request, code int, err error)`

to

* `WriteError(w http.ResponseWriter, r *http.Request, err interface{})`
* `WriteErrorCode(w http.ResponseWriter, r *http.Request, code int, err interface{})`

This has no functional implications unless you are implementing the writer interface yourself.
