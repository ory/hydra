**This file is no longer being updated and kept for historical reasons. Please check
the [GitHub releases](https://github.com/ory/fosite/releases) instead!**

<!-- START doctoc generated TOC please keep comment here to allow auto update -->
<!-- DON'T EDIT THIS SECTION, INSTEAD RE-RUN doctoc TO UPDATE -->

**Table of Contents**

- [0.0.0 (2022-09-22)](#000-2022-09-22)
  - [Breaking Changes](#breaking-changes)
    - [Bug Fixes](#bug-fixes)
    - [Code Refactoring](#code-refactoring)
    - [Features](#features)
    - [Tests](#tests)
    - [Unclassified](#unclassified)
- [0.42.2 (2022-04-17)](#0422-2022-04-17)
  - [Bug Fixes](#bug-fixes-1)
  - [Code Generation](#code-generation)
  - [Documentation](#documentation)
  - [Features](#features-1)
- [0.42.1 (2022-02-03)](#0421-2022-02-03)
  - [Code Generation](#code-generation-1)
  - [Features](#features-2)
- [0.42.0 (2022-01-06)](#0420-2022-01-06)
  - [Code Generation](#code-generation-2)
  - [Features](#features-3)
- [0.41.0 (2021-11-13)](#0410-2021-11-13)
  - [Bug Fixes](#bug-fixes-2)
  - [Code Generation](#code-generation-3)
  - [Code Refactoring](#code-refactoring-1)
  - [Documentation](#documentation-1)
  - [Features](#features-4)
- [0.40.2 (2021-05-28)](#0402-2021-05-28)
  - [Features](#features-5)
- [0.40.1 (2021-05-23)](#0401-2021-05-23)
  - [Bug Fixes](#bug-fixes-3)
- [0.40.0 (2021-05-21)](#0400-2021-05-21)
  - [Bug Fixes](#bug-fixes-4)
  - [Code Refactoring](#code-refactoring-2)
  - [Documentation](#documentation-2)
  - [Features](#features-6)
  - [Tests](#tests-1)
- [0.39.0 (2021-03-08)](#0390-2021-03-08)
  - [Features](#features-7)
- [0.38.0 (2021-02-23)](#0380-2021-02-23)
  - [Breaking Changes](#breaking-changes-1)
    - [Bug Fixes](#bug-fixes-5)
    - [Features](#features-8)
- [0.37.0 (2021-02-05)](#0370-2021-02-05)
  - [Bug Fixes](#bug-fixes-6)
  - [Features](#features-9)
- [0.36.1 (2021-01-11)](#0361-2021-01-11)
  - [Bug Fixes](#bug-fixes-7)
  - [Chores](#chores)
  - [Code Refactoring](#code-refactoring-3)
- [0.36.0 (2020-11-16)](#0360-2020-11-16)
  - [Breaking Changes](#breaking-changes-2)
    - [Bug Fixes](#bug-fixes-8)
    - [Code Refactoring](#code-refactoring-4)
    - [Documentation](#documentation-3)
    - [Features](#features-10)
- [0.35.1 (2020-10-11)](#0351-2020-10-11)
  - [Bug Fixes](#bug-fixes-9)
  - [Code Generation](#code-generation-4)
  - [Documentation](#documentation-4)
  - [Features](#features-11)
- [0.35.0 (2020-10-06)](#0350-2020-10-06)
  - [Breaking Changes](#breaking-changes-3)
    - [Bug Fixes](#bug-fixes-10)
    - [Code Generation](#code-generation-5)
- [0.34.1 (2020-10-02)](#0341-2020-10-02)
  - [Bug Fixes](#bug-fixes-11)
  - [Documentation](#documentation-5)
- [0.34.0 (2020-09-24)](#0340-2020-09-24)
  - [Breaking Changes](#breaking-changes-4)
    - [Bug Fixes](#bug-fixes-12)
    - [Chores](#chores-1)
    - [Features](#features-12)
    - [Unclassified](#unclassified-1)
- [0.33.0 (2020-09-16)](#0330-2020-09-16)
  - [Breaking Changes](#breaking-changes-5)
    - [Features](#features-13)
- [0.32.4 (2020-09-15)](#0324-2020-09-15)
  - [Code Generation](#code-generation-6)
  - [Code Refactoring](#code-refactoring-5)
  - [Documentation](#documentation-6)
- [0.32.3 (2020-09-12)](#0323-2020-09-12)
  - [Bug Fixes](#bug-fixes-13)
  - [Code Refactoring](#code-refactoring-6)
  - [Documentation](#documentation-7)
  - [Features](#features-14)
- [0.32.2 (2020-06-22)](#0322-2020-06-22)
  - [Features](#features-15)
- [0.32.1 (2020-06-05)](#0321-2020-06-05)
  - [Bug Fixes](#bug-fixes-14)
  - [Features](#features-16)
- [0.32.0 (2020-05-28)](#0320-2020-05-28)
  - [Bug Fixes](#bug-fixes-15)
  - [Documentation](#documentation-8)
  - [Features](#features-17)
- [0.31.3 (2020-05-09)](#0313-2020-05-09)
  - [Bug Fixes](#bug-fixes-16)
  - [Features](#features-18)
- [0.31.2 (2020-04-16)](#0312-2020-04-16)
  - [Bug Fixes](#bug-fixes-17)
- [0.31.1 (2020-04-16)](#0311-2020-04-16)
  - [Bug Fixes](#bug-fixes-18)
  - [Documentation](#documentation-9)
- [0.31.0 (2020-03-29)](#0310-2020-03-29)
  - [Unclassified](#unclassified-2)
- [0.30.6 (2020-03-26)](#0306-2020-03-26)
  - [Bug Fixes](#bug-fixes-19)
  - [Documentation](#documentation-10)
- [0.30.5 (2020-03-25)](#0305-2020-03-25)
  - [Bug Fixes](#bug-fixes-20)
- [0.30.4 (2020-03-17)](#0304-2020-03-17)
  - [Bug Fixes](#bug-fixes-21)
- [0.30.3 (2020-03-04)](#0303-2020-03-04)
  - [Bug Fixes](#bug-fixes-22)
  - [Documentation](#documentation-11)
  - [Features](#features-19)
- [0.30.2 (2019-11-21)](#0302-2019-11-21)
  - [Unclassified](#unclassified-3)
- [0.30.1 (2019-09-23)](#0301-2019-09-23)
  - [Unclassified](#unclassified-4)
- [0.30.0 (2019-09-16)](#0300-2019-09-16)
  - [Unclassified](#unclassified-5)
- [0.29.8 (2019-08-29)](#0298-2019-08-29)
  - [Documentation](#documentation-12)
  - [Unclassified](#unclassified-6)
- [0.29.7 (2019-08-06)](#0297-2019-08-06)
  - [Documentation](#documentation-13)
  - [Unclassified](#unclassified-7)
- [0.29.6 (2019-04-26)](#0296-2019-04-26)
  - [Unclassified](#unclassified-8)
- [0.29.5 (2019-04-25)](#0295-2019-04-25)
  - [Unclassified](#unclassified-9)
- [0.29.3 (2019-04-17)](#0293-2019-04-17)
  - [Unclassified](#unclassified-10)
- [0.29.2 (2019-04-11)](#0292-2019-04-11)
  - [Unclassified](#unclassified-11)
- [0.29.1 (2019-03-27)](#0291-2019-03-27)
  - [Unclassified](#unclassified-12)
- [0.29.0 (2018-12-23)](#0290-2018-12-23)
  - [Unclassified](#unclassified-13)
- [0.28.1 (2018-12-04)](#0281-2018-12-04)
  - [Unclassified](#unclassified-14)
- [0.28.0 (2018-11-16)](#0280-2018-11-16)
  - [Unclassified](#unclassified-15)
- [0.27.4 (2018-11-12)](#0274-2018-11-12)
  - [Documentation](#documentation-14)
  - [Unclassified](#unclassified-16)
- [0.27.3 (2018-11-08)](#0273-2018-11-08)
  - [Unclassified](#unclassified-17)
- [0.27.2 (2018-11-07)](#0272-2018-11-07)
  - [Unclassified](#unclassified-18)
- [0.27.1 (2018-11-03)](#0271-2018-11-03)
  - [Unclassified](#unclassified-19)
- [0.27.0 (2018-10-31)](#0270-2018-10-31)
  - [Unclassified](#unclassified-20)
- [0.26.1 (2018-10-25)](#0261-2018-10-25)
  - [Unclassified](#unclassified-21)
- [0.26.0 (2018-10-24)](#0260-2018-10-24)
  - [Unclassified](#unclassified-22)
- [0.25.1 (2018-10-23)](#0251-2018-10-23)
  - [Documentation](#documentation-15)
  - [Unclassified](#unclassified-23)
- [0.25.0 (2018-10-08)](#0250-2018-10-08)
  - [Unclassified](#unclassified-24)
- [0.24.0 (2018-09-27)](#0240-2018-09-27)
  - [Unclassified](#unclassified-25)
- [0.23.0 (2018-09-22)](#0230-2018-09-22)
  - [Unclassified](#unclassified-26)
- [0.22.0 (2018-09-19)](#0220-2018-09-19)
  - [Unclassified](#unclassified-27)
- [0.21.5 (2018-08-31)](#0215-2018-08-31)
  - [Unclassified](#unclassified-28)
- [0.21.4 (2018-08-26)](#0214-2018-08-26)
  - [Unclassified](#unclassified-29)
- [0.21.3 (2018-08-22)](#0213-2018-08-22)
  - [Unclassified](#unclassified-30)
- [0.21.2 (2018-08-07)](#0212-2018-08-07)
  - [Unclassified](#unclassified-31)
- [0.21.1 (2018-07-22)](#0211-2018-07-22)
  - [Unclassified](#unclassified-32)
- [0.21.0 (2018-06-23)](#0210-2018-06-23)
  - [Documentation](#documentation-16)
  - [Unclassified](#unclassified-33)
- [0.20.3 (2018-06-07)](#0203-2018-06-07)
  - [Unclassified](#unclassified-34)
- [0.20.2 (2018-05-29)](#0202-2018-05-29)
  - [Unclassified](#unclassified-35)
- [0.20.1 (2018-05-29)](#0201-2018-05-29)
  - [Unclassified](#unclassified-36)
- [0.20.0 (2018-05-28)](#0200-2018-05-28)
  - [Unclassified](#unclassified-37)
- [0.19.8 (2018-05-24)](#0198-2018-05-24)
  - [Unclassified](#unclassified-38)
- [0.19.7 (2018-05-24)](#0197-2018-05-24)
  - [Unclassified](#unclassified-39)
- [0.19.6 (2018-05-24)](#0196-2018-05-24)
  - [Unclassified](#unclassified-40)
- [0.19.5 (2018-05-23)](#0195-2018-05-23)
  - [Unclassified](#unclassified-41)
- [0.19.4 (2018-05-20)](#0194-2018-05-20)
  - [Unclassified](#unclassified-42)
- [0.19.3 (2018-05-20)](#0193-2018-05-20)
  - [Unclassified](#unclassified-43)
- [0.19.2 (2018-05-19)](#0192-2018-05-19)
  - [Unclassified](#unclassified-44)
- [0.19.1 (2018-05-19)](#0191-2018-05-19)
  - [Unclassified](#unclassified-45)
- [0.19.0 (2018-05-17)](#0190-2018-05-17)
  - [Unclassified](#unclassified-46)
- [0.18.1 (2018-05-01)](#0181-2018-05-01)
  - [Unclassified](#unclassified-47)
- [0.18.0 (2018-04-30)](#0180-2018-04-30)
  - [Unclassified](#unclassified-48)
- [0.17.2 (2018-04-26)](#0172-2018-04-26)
  - [Unclassified](#unclassified-49)
- [0.17.1 (2018-04-22)](#0171-2018-04-22)
  - [Unclassified](#unclassified-50)
- [0.17.0 (2018-04-08)](#0170-2018-04-08)
  - [Documentation](#documentation-17)
  - [Unclassified](#unclassified-51)
- [0.16.5 (2018-03-17)](#0165-2018-03-17)
  - [Documentation](#documentation-18)
  - [Unclassified](#unclassified-52)
- [0.16.4 (2018-02-07)](#0164-2018-02-07)
  - [Unclassified](#unclassified-53)
- [0.16.3 (2018-02-07)](#0163-2018-02-07)
  - [Unclassified](#unclassified-54)
- [0.16.2 (2018-01-25)](#0162-2018-01-25)
  - [Unclassified](#unclassified-55)
- [0.16.1 (2017-12-23)](#0161-2017-12-23)
  - [Unclassified](#unclassified-56)
- [0.16.0 (2017-12-23)](#0160-2017-12-23)
  - [Unclassified](#unclassified-57)
- [0.15.6 (2017-12-21)](#0156-2017-12-21)
  - [Unclassified](#unclassified-58)
- [0.15.5 (2017-12-17)](#0155-2017-12-17)
  - [Unclassified](#unclassified-59)
- [0.15.4 (2017-12-17)](#0154-2017-12-17)
  - [Unclassified](#unclassified-60)
- [0.15.3 (2017-12-17)](#0153-2017-12-17)
  - [Unclassified](#unclassified-61)
- [0.15.2 (2017-12-10)](#0152-2017-12-10)
  - [Unclassified](#unclassified-62)
- [0.15.1 (2017-12-10)](#0151-2017-12-10)
  - [Unclassified](#unclassified-63)
- [0.15.0 (2017-12-09)](#0150-2017-12-09)
  - [Documentation](#documentation-19)
  - [Unclassified](#unclassified-64)
- [0.14.2 (2017-12-06)](#0142-2017-12-06)
  - [Unclassified](#unclassified-65)
- [0.14.1 (2017-12-06)](#0141-2017-12-06)
  - [Unclassified](#unclassified-66)
- [0.14.0 (2017-12-06)](#0140-2017-12-06)
  - [Unclassified](#unclassified-67)
- [0.13.1 (2017-12-04)](#0131-2017-12-04)
  - [Unclassified](#unclassified-68)
- [0.13.0 (2017-10-25)](#0130-2017-10-25)
  - [Unclassified](#unclassified-69)
- [0.12.0 (2017-10-25)](#0120-2017-10-25)
  - [Unclassified](#unclassified-70)
- [0.11.4 (2017-10-10)](#0114-2017-10-10)
  - [Documentation](#documentation-20)
  - [Unclassified](#unclassified-71)
- [0.11.3 (2017-08-21)](#0113-2017-08-21)
  - [Documentation](#documentation-21)
  - [Unclassified](#unclassified-72)
- [0.11.2 (2017-07-09)](#0112-2017-07-09)
  - [Unclassified](#unclassified-73)
- [0.11.1 (2017-07-09)](#0111-2017-07-09)
  - [Unclassified](#unclassified-74)
- [0.11.0 (2017-07-09)](#0110-2017-07-09)
  - [Unclassified](#unclassified-75)
- [0.10.0 (2017-07-06)](#0100-2017-07-06)
  - [Unclassified](#unclassified-76)
- [0.9.7 (2017-06-28)](#097-2017-06-28)
  - [Unclassified](#unclassified-77)
- [0.9.6 (2017-06-21)](#096-2017-06-21)
  - [Documentation](#documentation-22)
  - [Unclassified](#unclassified-78)
- [0.9.5 (2017-06-08)](#095-2017-06-08)
  - [Unclassified](#unclassified-79)
- [0.9.4 (2017-06-05)](#094-2017-06-05)
  - [Unclassified](#unclassified-80)
- [0.9.3 (2017-06-05)](#093-2017-06-05)
  - [Unclassified](#unclassified-81)
- [0.9.2 (2017-06-05)](#092-2017-06-05)
  - [Unclassified](#unclassified-82)
- [0.9.1 (2017-06-04)](#091-2017-06-04)
  - [Unclassified](#unclassified-83)
- [0.9.0 (2017-06-03)](#090-2017-06-03)
  - [Documentation](#documentation-23)
  - [Unclassified](#unclassified-84)
- [0.8.0 (2017-05-18)](#080-2017-05-18)
  - [Documentation](#documentation-24)
  - [Unclassified](#unclassified-85)
- [0.7.0 (2017-05-03)](#070-2017-05-03)
  - [Documentation](#documentation-25)
  - [Unclassified](#unclassified-86)
- [0.6.19 (2017-05-03)](#0619-2017-05-03)
  - [Unclassified](#unclassified-87)
- [0.6.18 (2017-04-14)](#0618-2017-04-14)
  - [Unclassified](#unclassified-88)
- [0.6.17 (2017-02-24)](#0617-2017-02-24)
  - [Unclassified](#unclassified-89)
- [0.6.15 (2017-02-11)](#0615-2017-02-11)
  - [Unclassified](#unclassified-90)
- [0.6.14 (2017-01-08)](#0614-2017-01-08)
  - [Unclassified](#unclassified-91)
- [0.6.13 (2017-01-08)](#0613-2017-01-08)
  - [Unclassified](#unclassified-92)
- [0.6.12 (2017-01-02)](#0612-2017-01-02)
  - [Unclassified](#unclassified-93)
- [0.6.11 (2017-01-02)](#0611-2017-01-02)
  - [Unclassified](#unclassified-94)
- [0.6.10 (2016-12-29)](#0610-2016-12-29)
  - [Unclassified](#unclassified-95)
- [0.6.9 (2016-12-29)](#069-2016-12-29)
  - [Documentation](#documentation-26)
  - [Unclassified](#unclassified-96)
- [0.6.8 (2016-12-20)](#068-2016-12-20)
  - [Unclassified](#unclassified-97)
- [0.6.7 (2016-12-06)](#067-2016-12-06)
  - [Unclassified](#unclassified-98)
- [0.6.6 (2016-12-06)](#066-2016-12-06)
  - [Unclassified](#unclassified-99)
- [0.6.5 (2016-12-04)](#065-2016-12-04)
  - [Unclassified](#unclassified-100)
- [0.6.4 (2016-11-29)](#064-2016-11-29)
  - [Unclassified](#unclassified-101)
- [0.6.2 (2016-11-25)](#062-2016-11-25)
  - [Unclassified](#unclassified-102)
- [0.6.1 (2016-11-17)](#061-2016-11-17)
  - [Unclassified](#unclassified-103)
- [0.6.0 (2016-11-17)](#060-2016-11-17)
  - [Unclassified](#unclassified-104)
- [0.5.1 (2016-10-22)](#051-2016-10-22)
  - [Unclassified](#unclassified-105)
- [0.5.0 (2016-10-17)](#050-2016-10-17)
  - [Unclassified](#unclassified-106)
- [0.4.0 (2016-10-16)](#040-2016-10-16)
  - [Documentation](#documentation-27)
  - [Unclassified](#unclassified-107)
- [0.3.6 (2016-10-07)](#036-2016-10-07)
  - [Unclassified](#unclassified-108)
- [0.3.5 (2016-10-06)](#035-2016-10-06)
  - [Unclassified](#unclassified-109)
- [0.3.4 (2016-10-04)](#034-2016-10-04)
  - [Unclassified](#unclassified-110)
- [0.3.3 (2016-10-03)](#033-2016-10-03)
  - [Documentation](#documentation-28)
  - [Unclassified](#unclassified-111)
- [0.3.2 (2016-09-22)](#032-2016-09-22)
  - [Unclassified](#unclassified-112)
- [0.3.1 (2016-09-22)](#031-2016-09-22)
  - [Unclassified](#unclassified-113)
- [0.3.0 (2016-08-22)](#030-2016-08-22)
  - [Unclassified](#unclassified-114)
- [0.2.4 (2016-08-09)](#024-2016-08-09)
  - [Unclassified](#unclassified-115)
- [0.2.3 (2016-08-08)](#023-2016-08-08)
  - [Unclassified](#unclassified-116)
- [0.2.2 (2016-08-08)](#022-2016-08-08)
  - [Unclassified](#unclassified-117)
- [0.2.1 (2016-08-08)](#021-2016-08-08)
  - [Unclassified](#unclassified-118)
- [0.2.0 (2016-08-06)](#020-2016-08-06)
  - [Unclassified](#unclassified-119)
- [0.1.0 (2016-08-01)](#010-2016-08-01)
  - [Code Refactoring](#code-refactoring-7)
  - [Documentation](#documentation-29)
  - [Unclassified](#unclassified-120)

<!-- END doctoc generated TOC please keep comment here to allow auto update -->

# [0.0.0](https://github.com/ory/fosite/compare/v0.42.2...v0.0.0) (2022-09-22)

## Breaking Changes

Please be aware that several internal APIs have changed, as well as public methods. Most notably, we added the context to all `Write*` metods.

```patch
 type OAuth2Provider interface {
-    WriteAuthorizeError(rw http.ResponseWriter, requester AuthorizeRequester, err error)
+    WriteAuthorizeError(ctx context.Context, rw http.ResponseWriter, requester AuthorizeRequester, err error)

-    WriteAuthorizeResponse(rw http.ResponseWriter, requester AuthorizeRequester, responder AuthorizeResponder)
+    WriteAuthorizeResponse(ctx context.Context, rw http.ResponseWriter, requester AuthorizeRequester, responder AuthorizeResponder)

-    WriteAccessError(rw http.ResponseWriter, requester AccessRequester, err error)
+    WriteAccessError(ctx context.Context, rw http.ResponseWriter, requester AccessRequester, err error)

-    WriteAccessResponse(rw http.ResponseWriter, requester AccessRequester, responder AccessResponder)
+    WriteAccessResponse(ctx context.Context, rw http.ResponseWriter, requester AccessRequester, responder AccessResponder)

-    WriteRevocationResponse(rw http.ResponseWriter, err error)
+    WriteRevocationResponse(ctx context.Context, rw http.ResponseWriter, err error)

-    WriteIntrospectionError(rw http.ResponseWriter, err error)
+    WriteIntrospectionError(ctx context.Context, rw http.ResponseWriter, err error)

-    WriteIntrospectionResponse(rw http.ResponseWriter, r IntrospectionResponder)
+    WriteIntrospectionResponse(ctx context.Context, rw http.ResponseWriter, r IntrospectionResponder)
 }
```

The default config struct has moved from package `github.com/ory/fosite/compose.Config` to `github.com/ory/fosite.Config`. Struct `github.com/ory/fosite.Fosite` no longer has any configuration parameters
itself.

Please note that the HMAC / global secret has to be set no longer in the compose call, but in the config initialization:

```patch
-compose.ComposeAllEnabled(&compose.Config{}, store, secret, privateKey)
+compose.ComposeAllEnabled(&fosite.Config{GlobalSecret: secret}, store, privateKey)
```

Many internal interfaces have been changed, usually adding `ctx context.Context` as the first parameter.

### Bug Fixes

- Bump dependencies ([5dab818](https://github.com/ory/fosite/commit/5dab818f9707e364dcfe56bc6fd2245049417cc1))
- Cves in deps ([f5782c3](https://github.com/ory/fosite/commit/f5782c33814ec738ea188b0ffac50ef45e7f3eb8))
- Include `at_hash` claim in authcode flow's ID token ([#679](https://github.com/ory/fosite/issues/679)) ([c3b7bab](https://github.com/ory/fosite/commit/c3b7bab41db24b000f8e1416e1475e0aae4c310c))
- Linting ([222ca97](https://github.com/ory/fosite/commit/222ca97805edfb52a655969841c2ac2958cc6d36))
- **rfc7523:** Comment mentioned incorrect granttype ([#668](https://github.com/ory/fosite/issues/668)) ([b41f187](https://github.com/ory/fosite/commit/b41f187703bc1c8dc43ac0ec1ea23569779974bb))
- State check for hybrid flow ([#670](https://github.com/ory/fosite/issues/670)) ([37f8a0a](https://github.com/ory/fosite/commit/37f8a0ac12e47893459528cabb38b9879600286d))

### Code Refactoring

- **config:** Support hot reloading ([1661401](https://github.com/ory/fosite/commit/16614014a42b3905d065188c8f1f45433c4353f9)), closes [#666](https://github.com/ory/fosite/issues/666):

  This patch updates the config system to be replacable and uses functions instead of struct fields. This allows implementing hot reloading mechanisms easily.

- Move to go 1.17 ([d9d0fed](https://github.com/ory/fosite/commit/d9d0fedaad87044e4d38ba82e01d9e430d09514c))

### Features

- Add `ory_at|pt|ac` prefixes to HMAC tokens ([b652335](https://github.com/ory/fosite/commit/b652335c965d5cc523faebad9c9792c4135cfb75)):

  See https://github.com/ory/hydra/issues/2845

- Add json mappings to default session and its contents ([#688](https://github.com/ory/fosite/issues/688)) ([d8ecac4](https://github.com/ory/fosite/commit/d8ecac4077c446b71842372169abc37a02f9e1b7))
- Add json mappings to generic session to match openid session ([#690](https://github.com/ory/fosite/issues/690)) ([2386b25](https://github.com/ory/fosite/commit/2386b259837ab89983f6d0ee37b147b36b171f5b))
- Implement client token lifespan customization ([#684](https://github.com/ory/fosite/issues/684)) ([cfffe8c](https://github.com/ory/fosite/commit/cfffe8cec67a986e2abc736b940f9f0bab9ad7d9)):

  This change introduces the ability to control the lifespan of tokens for each valid combination of Client, GrantType, and TokenType.

- Introduce cache strategy for JWKS fetcher ([452f377](https://github.com/ory/fosite/commit/452f37728890c68524b9aa190e1cdb279414f802))
- Make http source contextualized ([9fc89e9](https://github.com/ory/fosite/commit/9fc89e9007c71354f7fe2d036ea6e175a2e5860b))
- PAR implementation ([#660](https://github.com/ory/fosite/issues/660)) ([3de78db](https://github.com/ory/fosite/commit/3de78db805fe1c69b0fc5b853bfabeb19433feba)), closes [#628](https://github.com/ory/fosite/issues/628):

  Implements [RFC9126 - Pushed Authorization Request](https://www.rfc-editor.org/rfc/rfc9126.html).

- Support variety of JWT formats when `jose.JSONWebKey` is used ([2590eb8](https://github.com/ory/fosite/commit/2590eb83d1e66df998053bc2fb7381b9043c232e))

### Tests

- Fix assertions ([#683](https://github.com/ory/fosite/issues/683)) ([551b8b8](https://github.com/ory/fosite/commit/551b8b827cf0b7033aac80818a516ee3c5b8523e))
- Fix panic ([fe60766](https://github.com/ory/fosite/commit/fe60766cdb1f0d22df7d9c4543b06cfd6dc7aea1))

### Unclassified

- Revert "chore: delete .circleci folder (#699)" (#705) ([ef753d5](https://github.com/ory/fosite/commit/ef753d550d59b077f6ae349d6c795e2c142ec676)), closes [#699](https://github.com/ory/fosite/issues/699) [#705](https://github.com/ory/fosite/issues/705):

  This reverts commit 2eea63bddcbdf50771adf670391e495e339f619f since CircleCI is still used here.

# [0.42.2](https://github.com/ory/fosite/compare/v0.42.1...v0.42.2) (2022-04-17)

autogen(docs): regenerate and update changelog

### Bug Fixes

- Always rollback ([#638](https://github.com/ory/fosite/issues/638)) ([7edf673](https://github.com/ory/fosite/commit/7edf673f20aece260f9ba677a07086c48835fba8)), closes [#637](https://github.com/ory/fosite/issues/637)
- Empty client secret via basic auth header means "none" authn ([#655](https://github.com/ory/fosite/issues/655)) ([7a2d972](https://github.com/ory/fosite/commit/7a2d9721f4b6da0e3b2b829ec4312de1e3d66b6f)), closes [/github.com/golang/oauth2/blob/ee480838109b20d468babcb00b7027c82f962065/internal/token.go#L174-L176](https://github.com//github.com/golang/oauth2/blob/ee480838109b20d468babcb00b7027c82f962065/internal/token.go/issues/L174-L176):

  The existing client authentication code treats an empty client_secret
  query parameter to be equivalent to "none" authentication instead of
  "client_secret_post."

  This change updates the basic auth check to be consistent with this.
  That is, an empty secret via the basic auth header is considered to
  mean "none" instead of "client_secret_basic."

  The "golang.org/x/oauth2" library probes for both methods of
  authentication, starting with the basic auth header approach first.

  As required, both client ID and secret are encoded in one header:

- Handle invalid_token error for refresh_token is expired ([#664](https://github.com/ory/fosite/issues/664)) ([76bb274](https://github.com/ory/fosite/commit/76bb274e95585d4552789abbd1c1f123463ff47e))
- Handle token_inactive error for multiple concurrent refresh requests ([#652](https://github.com/ory/fosite/issues/652)) ([7c8f4ae](https://github.com/ory/fosite/commit/7c8f4ae49550c61ff43d1a86adace4ed08c71e3e)):

  See https://github.com/ory/hydra/issues/3004

- Url-encode the fragment in the redirect URL of the authorize response ([#649](https://github.com/ory/fosite/issues/649)) ([beec138](https://github.com/ory/fosite/commit/beec13889c431ff06348c032dd260d00db253dd2)), closes [#648](https://github.com/ory/fosite/issues/648):

  This patch reverts the encoding logic for the fragment of the redirect URL returned as part of the authorize response to what was the one before version `0.36.0`. In that version, the code was refactored and the keys and values of the fragment ceased to be url-encoded. This in turn reflected on all Ory Hydra versions starting from `1.9.0` and provoked a breaking change that made the parsing of the fragment impossible if any of the params contain a character like `&` or `=` because they get treated as separators instead of as text

- Use the correct algorithm for at_hash and c_hash ([#659](https://github.com/ory/fosite/issues/659)) ([8cb4b4b](https://github.com/ory/fosite/commit/8cb4b4b0c57be8944e403a0f3ec588b19f49f6f7)), closes [#630](https://github.com/ory/fosite/issues/630)

### Code Generation

- **docs:** Regenerate and update changelog ([5dbfa9a](https://github.com/ory/fosite/commit/5dbfa9a56d36061d5bf80149e1801c36a371bafd))

### Documentation

- Add deprecation to communicate ropc discouragement ([#665](https://github.com/ory/fosite/issues/665)) ([df491be](https://github.com/ory/fosite/commit/df491beb5e82ca66bf5c5825c91ded0ca9d67b57)):

  This adds godoc deprecations to the compose.OAuth2ResourceOwnerPasswordCredentialsFactory and oauth2.ResourceOwnerPasswordCredentialsGrantHandler in order to clearly communicate the discouragement of the ROPC grant type to users implementing this library.

### Features

- Use custom hash.Hash in hmac.HMACStrategy ([#663](https://github.com/ory/fosite/issues/663)) ([d09a8c3](https://github.com/ory/fosite/commit/d09a8c39284fecce47933ff3b53d90d35b646b0c)), closes [#654](https://github.com/ory/fosite/issues/654)

# [0.42.1](https://github.com/ory/fosite/compare/v0.42.0...v0.42.1) (2022-02-03)

autogen(docs): regenerate and update changelog

### Code Generation

- **docs:** Regenerate and update changelog ([dcc6550](https://github.com/ory/fosite/commit/dcc6550b807980faca740b261790b3be339632c7))

### Features

- Support FormPostHTMLTemplate config for fosite ([#647](https://github.com/ory/fosite/issues/647)) ([570ce3f](https://github.com/ory/fosite/commit/570ce3f6e3bf4e54781a6bfffc2ce777f0ac5194)), closes [#646](https://github.com/ory/fosite/issues/646)

# [0.42.0](https://github.com/ory/fosite/compare/v0.41.0...v0.42.0) (2022-01-06)

autogen(docs): regenerate and update changelog

### Code Generation

- **docs:** Regenerate and update changelog ([cf2c545](https://github.com/ory/fosite/commit/cf2c545540c12bfa5cfbf752bc84c03a8a515ecc))

### Features

- Add new function to TokenRevocationStorage to support refresh token grace-period ([#635](https://github.com/ory/fosite/issues/635)) ([9b40d03](https://github.com/ory/fosite/commit/9b40d036e6494dfe9942b513b8bc4a50c7c9f730))

# [0.41.0](https://github.com/ory/fosite/compare/v0.40.2...v0.41.0) (2021-11-13)

autogen(docs): regenerate and update changelog

### Bug Fixes

- Force HTTP GET for redirect responses ([#636](https://github.com/ory/fosite/issues/636)) ([f6c6523](https://github.com/ory/fosite/commit/f6c6523a09e7733d5ca263bccb7fd4fdb80172b2))
- Include `typ` in jwt header ([#607](https://github.com/ory/fosite/issues/607)) ([7644a74](https://github.com/ory/fosite/commit/7644a74bd48accb46d8578f6846b3e509dfd4b03)), closes [#606](https://github.com/ory/fosite/issues/606)
- Make `amr` claim an array to match the OIDC spec ([#625](https://github.com/ory/fosite/issues/625)) ([8a6f66a](https://github.com/ory/fosite/commit/8a6f66ab5d9f74140f4ce94210f09ccb0e27f56d))
- Resolve nancy warning ([b6cf0a6](https://github.com/ory/fosite/commit/b6cf0a641d1169595ceb3110f76be0788e778521))

### Code Generation

- **docs:** Regenerate and update changelog ([1777ad5](https://github.com/ory/fosite/commit/1777ad52e68b20ce57ed7f2f7d085895c3c157c6))

### Code Refactoring

- Upgrade go-jose to decode JSON numbers into int64 ([#603](https://github.com/ory/fosite/issues/603)) ([c02d327](https://github.com/ory/fosite/commit/c02d3273e30ca9b29285d1641b252e6c29598ea5)), closes [#602](https://github.com/ory/fosite/issues/602)

### Documentation

- Add missing word ([#626](https://github.com/ory/fosite/issues/626)) ([c7a553b](https://github.com/ory/fosite/commit/c7a553bb4945013be17d2bbd2ec126ae93113a72))
- Document that DeleteOpenIDConnectSession is deprecated ([#634](https://github.com/ory/fosite/issues/634)) ([4e2c03d](https://github.com/ory/fosite/commit/4e2c03d3f6dcb3a3b50e7ea245128edde7ebf959))

### Features

- Add client secret rotation support ([#608](https://github.com/ory/fosite/issues/608)) ([a4ce354](https://github.com/ory/fosite/commit/a4ce3544c2996a99b65350d4b200967df9fc0d45)), closes [#590](https://github.com/ory/fosite/issues/590)
- Add prettier and format ([d682bdf](https://github.com/ory/fosite/commit/d682bdf51c22c211ee1aceb06fb7c4a7e43db326))
- Add ResponseModeHandler to support custom response modes ([#592](https://github.com/ory/fosite/issues/592)) ([10ec003](https://github.com/ory/fosite/commit/10ec003fb414fd3fcbd3e2e6d250cb2da51a0304)), closes [#591](https://github.com/ory/fosite/issues/591)
- I18n support added ([#627](https://github.com/ory/fosite/issues/627)) ([cf02af9](https://github.com/ory/fosite/commit/cf02af977681fd667b33f8e131891f6746d0b9da)), closes [#615](https://github.com/ory/fosite/issues/615)
- Support jose.opaquesigner for JWTs ([#611](https://github.com/ory/fosite/issues/611)) ([1121a0a](https://github.com/ory/fosite/commit/1121a0aa4155e9216abb989ab008df8cff67830d))
- Use bitwise comparison for jwt validation errors ([#633](https://github.com/ory/fosite/issues/633)) ([52ee93f](https://github.com/ory/fosite/commit/52ee93fe976152457482870b4ebb487560ca93e0))

# [0.40.2](https://github.com/ory/fosite/compare/v0.40.1...v0.40.2) (2021-05-28)

feat: use int64 type for claims with timestamps (#600)

Co-authored-by: Nestor <nesterran@gmail.com>

### Features

- Use int64 type for claims with timestamps ([#600](https://github.com/ory/fosite/issues/600)) ([c370994](https://github.com/ory/fosite/commit/c370994c007be101a388f825f1a4d6b38393756e))

# [0.40.1](https://github.com/ory/fosite/compare/v0.40.0...v0.40.1) (2021-05-23)

fix: revert float64 auth_time claim (#599)

Closes #598

### Bug Fixes

- Revert float64 auth_time claim ([#599](https://github.com/ory/fosite/issues/599)) ([e609d91](https://github.com/ory/fosite/commit/e609d9196070050adf39b9bdb3cbfbba2edda0d5)), closes [#598](https://github.com/ory/fosite/issues/598)

# [0.40.0](https://github.com/ory/fosite/compare/v0.39.0...v0.40.0) (2021-05-21)

feat: transit from jwt-go to go-jose (#593)

Closes #514

Co-authored-by: hackerman <3372410+aeneasr@users.noreply.github.com>

### Bug Fixes

- 582memory store authentication error code ([#583](https://github.com/ory/fosite/issues/583)) ([51b4424](https://github.com/ory/fosite/commit/51b44248275128ca83e1899522f2cd412e5c466e))
- Do not include nonce in ID tokens when not used ([#570](https://github.com/ory/fosite/issues/570)) ([795dee2](https://github.com/ory/fosite/commit/795dee246f26c1fef16dcd52da37e3df75e73772))
- Sha alg name in error message and go doc ([#571](https://github.com/ory/fosite/issues/571)) ([0f2e289](https://github.com/ory/fosite/commit/0f2e289973ad22d14c5d5bedd4fc9bb886134354))
- Upgrade gogo protubuf ([#573](https://github.com/ory/fosite/issues/573)) ([9a9467a](https://github.com/ory/fosite/commit/9a9467a20391059534df859b2b295711918bfd08))

### Code Refactoring

- Generate claims in the same way ([#595](https://github.com/ory/fosite/issues/595)) ([4c7b13f](https://github.com/ory/fosite/commit/4c7b13f2f1234128c53e8fc3e6cc3981e10d3069))

### Documentation

- Add client credentials grant how-to ([#589](https://github.com/ory/fosite/issues/589)) ([893aae4](https://github.com/ory/fosite/commit/893aae4348cfef78cb3d7f9aa70568e2137b4b3f)), closes [#566](https://github.com/ory/fosite/issues/566)

### Features

- Allow extra fields in introspect response ([#579](https://github.com/ory/fosite/issues/579)) ([294a0bf](https://github.com/ory/fosite/commit/294a0bf7f4cb01739a560480364403118d1408bf)), closes [#441](https://github.com/ory/fosite/issues/441)
- Allow omitting scope in authorization redirect uri ([#588](https://github.com/ory/fosite/issues/588)) ([6ad9264](https://github.com/ory/fosite/commit/6ad92642f0f01ff4d3662f3680a825db22594366))
- Pass requests through context ([#596](https://github.com/ory/fosite/issues/596)) ([2f96bb8](https://github.com/ory/fosite/commit/2f96bb8a2623fe7b4abb31db870582b555df6db8)), closes [#537](https://github.com/ory/fosite/issues/537)
- Transit from jwt-go to go-jose ([#593](https://github.com/ory/fosite/issues/593)) ([d022bbc](https://github.com/ory/fosite/commit/d022bbc2b45fd603cb12575e28bbe884170bf788)), closes [#514](https://github.com/ory/fosite/issues/514)

### Tests

- Change sha algorithm name acc to standard naming ([#572](https://github.com/ory/fosite/issues/572)) ([a3594a3](https://github.com/ory/fosite/commit/a3594a3cb0eb70e912a7268d2d396d19a45116c6))

# [0.39.0](https://github.com/ory/fosite/compare/v0.38.0...v0.39.0) (2021-03-08)

feat: token reuse detection (#567)

See https://github.com/ory/hydra/issues/2022

### Features

- Token reuse detection ([#567](https://github.com/ory/fosite/issues/567)) ([db7f981](https://github.com/ory/fosite/commit/db7f9817ee19878c4bf650e97b49be7e3b268ee0)):

  See https://github.com/ory/hydra/issues/2022

# [0.38.0](https://github.com/ory/fosite/compare/v0.37.0...v0.38.0) (2021-02-23)

feat: add ClientAuthenticationStrategy extension point (#565)

Closes #564

## Breaking Changes

Replaces `token_expired` error ID with `invalid_token` which is the correct value according to https://tools.ietf.org/html/rfc6750#section-3.1

### Bug Fixes

- Use correct error code for expired token ([#562](https://github.com/ory/fosite/issues/562)) ([56a71e5](https://github.com/ory/fosite/commit/56a71e5f9797abe35a9566c86f9ce9c1f485c11a))

### Features

- Add ClientAuthenticationStrategy extension point ([#565](https://github.com/ory/fosite/issues/565)) ([ec0bec2](https://github.com/ory/fosite/commit/ec0bec2d8462bae2dc545defbd21190dfe832024)), closes [#564](https://github.com/ory/fosite/issues/564)

# [0.37.0](https://github.com/ory/fosite/compare/v0.36.1...v0.37.0) (2021-02-05)

feat: add support for urn:ietf:params:oauth:grant-type:jwt-bearer grant type RFC 7523 (#560)

Closes #546
Closes #305

Co-authored-by: Vladimir Kalugin <v.p.kalugin@tinkoff.ru>
Co-authored-by: i.seliverstov <i.seliverstov@tinkoff.ru>

### Bug Fixes

- Resolve regression ([#561](https://github.com/ory/fosite/issues/561)) ([173d60e](https://github.com/ory/fosite/commit/173d60e5324c19c2323d2b8a731e201bf26845ce))

### Features

- Add support for urn:ietf:params:oauth:grant-type:jwt-bearer grant type RFC 7523 ([#560](https://github.com/ory/fosite/issues/560)) ([9720241](https://github.com/ory/fosite/commit/9720241c57e2154ed9fdb44fcf25e8c6b50410ee)), closes [#546](https://github.com/ory/fosite/issues/546) [#305](https://github.com/ory/fosite/issues/305)

# [0.36.1](https://github.com/ory/fosite/compare/v0.36.0...v0.36.1) (2021-01-11)

chore: bump deps

### Bug Fixes

- Broken dependency to reflection package ([#555](https://github.com/ory/fosite/issues/555)) ([a103222](https://github.com/ory/fosite/commit/a1032221363726bdcdc2f9b1c1898f99c62e8932))

### Chores

- Bump deps ([c2375de](https://github.com/ory/fosite/commit/c2375de6ff3229493b6a6ad628bf4e4961c8d989))

### Code Refactoring

- Use constructor ([#535](https://github.com/ory/fosite/issues/535)) ([2da54e3](https://github.com/ory/fosite/commit/2da54e3620a467e20d67ae05d0d3885a2383e4d4))
- Use provided context ([#536](https://github.com/ory/fosite/issues/536)) ([35d4f13](https://github.com/ory/fosite/commit/35d4f133faa87076c7eb1c5e8384f3653643de9e))

# [0.36.0](https://github.com/ory/fosite/compare/v0.35.1...v0.36.0) (2020-11-16)

fix: be more permissive in time checks

Time equality should not cause failures in OpenID Connect validation.

## Breaking Changes

This patch removes fields `error_hint`, `error_debug` from error responses. To use the legacy error format where these fields are included, set `UseLegacyErrorFormat` to true in your compose config or directly on the `Fosite` struct. If `UseLegacyErrorFormat` is set, the `error_description` no longer merges `error_hint` nor `error_debug` messages which reverts a change introduced in `v0.33.0`. Instead, `error_hint` and `error_debug` are included and the merged message can be constructed from those fields.

As part of this change, the error interface and its fields have changed:

- `RFC6749Error.Name` was renamed to `RFC6749Error.ErrorField`.
- `RFC6749Error.Description` was renamed to `RFC6749Error.DescriptionField`.
- `RFC6749Error.Hint` was renamed to `RFC6749Error.HintField`.
- `RFC6749Error.Code` was renamed to `RFC6749Error.CodeField`.
- `RFC6749Error.Hint` was renamed to `RFC6749Error.HintField`.
- `RFC6749Error.WithCause()` was renamed to `RFC6749Error.WithWrap() *RFC6749Error` and alternatively to `RFC6749Error.Wrap()` (without return value) to standardize naming conventions around the new Go 1.14+ error interfaces.

As part of this change, methods `GetResponseMode`, `SetDefaultResponseMode`, `GetDefaultResponseMode ` where added to interface `AuthorizeRequester`. Also, methods `GetQuery`, `AddQuery`, and `GetFragment` were merged into one function `GetParameters` and `AddParameter` on the `AuthorizeResponder` interface. Methods on `AuthorizeRequest` and `AuthorizeResponse` changed accordingly and will need to be updated in your codebase. Additionally, the field `Debug` was renamed to `DebugField` and a new method `Debug() string` was added to `RFC6749Error`.

Co-authored-by: hackerman <3372410+aeneasr@users.noreply.github.com>

### Bug Fixes

- Allow all request object algs when client value is unset ([1d14636](https://github.com/ory/fosite/commit/1d14636e61b2047e5eee6d1d740249b819fc0794)):

  Allows all request object signing algorithms when the client has not explicitly allowed a certain algorithm. This follows the spec:

  > \*request_object_signing_alg - OPTIONAL. JWS [JWS] alg algorithm [JWA] that MUST be used for signing Request Objects sent to the OP. All Request Objects from this Client MUST be rejected, if not signed with this algorithm. Request Objects are described in Section 6.1 of OpenID Connect Core 1.0 [OpenID.Core]. This algorithm MUST be used both when the Request Object is passed by value (using the request parameter) and when it is passed by reference (using the request_uri parameter). Servers SHOULD support RS256. The value none MAY be used. The default, if omitted, is that any algorithm supported by the OP and the RP MAY be used.

- Always return non-error response for inactive tokens ([#517](https://github.com/ory/fosite/issues/517)) ([5f2cae3](https://github.com/ory/fosite/commit/5f2cae3eabb83da898e1b5515176e65dda4da862))
- Be more permissive in time checks ([839d000](https://github.com/ory/fosite/commit/839d00093a2ed8c590d910f113186cd96fad9185)):

  Time equality should not cause failures in OpenID Connect validation.

- Do not accidentally leak jwks fetching errors ([6d2092d](https://github.com/ory/fosite/commit/6d2092da1e8699e43fd6dccb4c3a33b885cec7f8)), closes [/github.com/ory/fosite/pull/526#discussion_r517491738](https://github.com//github.com/ory/fosite/pull/526/issues/discussion_r517491738)
- Do not require nonce for hybrid flows ([de5c8f9](https://github.com/ory/fosite/commit/de5c8f90e8ccae0849fa6426d53563ef7520880d)):

  This patch resolves an issue where nonce was required for hybrid flows, which does not comply with the OpenID Connect conformity test suite, specifically the `oidcc-ensure-request-without-nonce-succeeds-for-code-flow` test.

- Guess default response mode in `NewAuthorizeRequest` ([a2952d7](https://github.com/ory/fosite/commit/a2952d7ad09fbd83a354b22dbcc0cef8a15f50f7))
- Improve claims handling for jwts ([a72ca9a](https://github.com/ory/fosite/commit/a72ca9a978e60d7c4b000c41357719f0e2b61f8e))
- Improve error stack wrapping ([620d4c1](https://github.com/ory/fosite/commit/620d4c148307f7be7b2674fe420141b33aef6075))
- Kid header is not required for key lookup ([27cc5c0](https://github.com/ory/fosite/commit/27cc5c0e935ecb8bca23dd8c2670c8a93f7b829d))
- Modernized JWT stateless introspection ([#519](https://github.com/ory/fosite/issues/519)) ([a6bfb92](https://github.com/ory/fosite/commit/a6bfb921ebc746ba7a1215e32fb42a2c0530a2bf))
- Only use allowed characters in error_description ([431f9a5](https://github.com/ory/fosite/commit/431f9a56ed03648ea4ef637fe6c2b6d74e765dad)), closes [#525](https://github.com/ory/fosite/issues/525):

  Replace LF and quotes with `.` and `'` to match allowed and recommended character set defined in various RFCs.

- Prevent debug details from leaking during key lookup ([c0598fb](https://github.com/ory/fosite/commit/c0598fb8d8ce75b7f0ad645420caea641e64a4d2)), closes [/github.com/ory/fosite/pull/526#discussion_r517490461](https://github.com//github.com/ory/fosite/pull/526/issues/discussion_r517490461)
- Reset jti and hash ID token claims on refresh ([#523](https://github.com/ory/fosite/issues/523)) ([ce2de73](https://github.com/ory/fosite/commit/ce2de73ff979b02be32d850c1c695067a35576c7))
- Use state from request object ([8cac1a0](https://github.com/ory/fosite/commit/8cac1a00a6f87523b88fea6962ab1194049cbacd)):

  Resolves failing OIDC conformity test "oidcc-request-uri-unsigned".

### Code Refactoring

- Use rfc compliant error formating ([edbbda3](https://github.com/ory/fosite/commit/edbbda3c4cf70a77cdcd1383c55762c73613f87e))

### Documentation

- Document Session interface methods ([#512](https://github.com/ory/fosite/issues/512)) ([11a95ba](https://github.com/ory/fosite/commit/11a95ba00f562b3864fc0d6878c9d93943cc4273))
- Updates banner in readme.md ([#529](https://github.com/ory/fosite/issues/529)) ([9718eb6](https://github.com/ory/fosite/commit/9718eb6ce63983ade0689908b5cce3e27c8838bc))

### Features

- Add support for response_mode=form_post ([#509](https://github.com/ory/fosite/issues/509)) ([3e3290f](https://github.com/ory/fosite/commit/3e3290f811f849881f1c6bafabc1c765d9a42ac7)):

  This patch introduces support for `response_mode=form_post` as well as `response_mode` of `none` and `query` and `fragment`.

  To support this new feature your OAuth2 Client must implement the `fosite.ResponseModeClient` interface. We suggest to always return all response modes there unless you want to explicitly disable one of the response modes:

  ```go
  func (c *Client) GetResponseModes() []fosite.ResponseModeType {
  	return []fosite.ResponseModeType{
  		fosite.ResponseModeDefault,
  		fosite.ResponseModeFormPost,
  		fosite.ResponseModeQuery,
  		fosite.ResponseModeFragment,
  	}
  }
  ```

- Improve error messages ([#513](https://github.com/ory/fosite/issues/513)) ([fcac5a6](https://github.com/ory/fosite/commit/fcac5a6457c92d1eb1a389192cd0c7fb590ab8b3))
- Introduce WithExposeDebug to error interface ([625a521](https://github.com/ory/fosite/commit/625a5214c4a002b4d0f86e49555edf8755703968))
- Support passing repeated audience parameter in URL query ([#518](https://github.com/ory/fosite/issues/518)) ([47f2a31](https://github.com/ory/fosite/commit/47f2a31fbed137b58e4866f78ec8b9f591134f98)), closes [#504](https://github.com/ory/fosite/issues/504):

  Added `GetAudiences` helper function which tries to have current behavior and also support multiple/repeated audience parameters. If there are parameter is repeated, then it is not split by space. If there is only one then it is split by space. I think this is the best balance between standard/backwards behavior and allowing repeated parameter and allowing also URIs/audiences with spaces in them (which we probably all agree is probably not something anyone should be doing).

  Also added `ExactAudienceMatchingStrategy` which is slightly more suitable to use for audiences which are not URIs. In [OIDC spec](https://openid.net/specs/openid-connect-core-1_0.html) audience is described as:

  > Audience(s) that this ID Token is intended for. It MUST contain the OAuth 2.0 client_id of the Relying Party as an audience value. It MAY also contain identifiers for other audiences. In the general case, the aud value is an array of case sensitive strings. In the common special case when there is one audience, the aud value MAY be a single case sensitive string.

  `client_id` is generally not an URI, but some UUID or some other random string.

# [0.35.1](https://github.com/ory/fosite/compare/v0.35.0...v0.35.1) (2020-10-11)

autogen(docs): regenerate and update changelog

### Bug Fixes

- Uniform audience parsing ([#505](https://github.com/ory/fosite/issues/505)) ([e3f331d](https://github.com/ory/fosite/commit/e3f331d0d8e4470eef3dd7ecb46e66eeebfbe4c7))

### Code Generation

- **docs:** Regenerate and update changelog ([c598cc7](https://github.com/ory/fosite/commit/c598cc7fae17e70db2bad555cff94e97b2ca185b))

### Documentation

- Improved test descriptions ([#507](https://github.com/ory/fosite/issues/507)) ([29e9336](https://github.com/ory/fosite/commit/29e9336be5673530ae00e735c3dc7d191f4b03a6))

### Features

- Allow configuring redirect secure checker everywhere ([#489](https://github.com/ory/fosite/issues/489)) ([e87d091](https://github.com/ory/fosite/commit/e87d0910f3ee960dbc7b1bc0fef124c9b928a55c))
- Scope can now be space delimited in access tokens ([#482](https://github.com/ory/fosite/issues/482)) ([8225935](https://github.com/ory/fosite/commit/8225935276d40a24da400d46ee7e7b63976488a1)), closes [#362](https://github.com/ory/fosite/issues/362)

# [0.35.0](https://github.com/ory/fosite/compare/v0.34.1...v0.35.0) (2020-10-06)

autogen(docs): regenerate and update changelog

## Breaking Changes

Type `fosite.TokenType` has been renamed to `fosite.TokenUse`.

### Bug Fixes

- Redirct_url with query escape character outside of query is failing ([#480](https://github.com/ory/fosite/issues/480)) ([6e49c57](https://github.com/ory/fosite/commit/6e49c57c8f7a46a78eda4d3091765d631f427845)):

  See https://github.com/ory/hydra/issues/2055

  Co-authored-by: ajanthan <ca52ca6fe18c44787827017e14ca2d0c3c5bdb58>

- Rename TokenType to TokenUse in introspection ([#486](https://github.com/ory/fosite/issues/486)) ([4b81316](https://github.com/ory/fosite/commit/4b81316a1dbb0c5246bac39ecbaff749b00e4efa)), closes [ory/hydra#1762](https://github.com/ory/hydra/issues/1762)
- Return allowed redirect url with preference ([f0badc4](https://github.com/ory/fosite/commit/f0badc4919e00fa179dd54edcbd7385fac14fa19))

### Code Generation

- **docs:** Regenerate and update changelog ([3f0bc87](https://github.com/ory/fosite/commit/3f0bc875af230342d161de8516b7c0050f89d648))

# [0.34.1](https://github.com/ory/fosite/compare/v0.34.0...v0.34.1) (2020-10-02)

fix: make redirect URL checking more strict

The OAuth 2.0 Client's Redirect URL and the Redirect URL used in the OAuth 2.0 flow do not check if the query string is equal:

1. Registering a client with allowed redirect URL `https://example.com/callback`
2. Performing OAuth2 flow and requesting redirect URL `https://example.com/callback?bar=foo`
3. Instead of an error, the browser is redirected to `https://example.com/callback?bar=foo` with a potentially successful OAuth2 response.

Additionally, matching Redirect URLs used `strings.ToLower` normalization:

1. Registering a client with allowed redirect URL `https://example.com/callback`
2. Performing OAuth2 flow and requesting redirect URL `https://example.com/CALLBACK`
3. Instead of an error, the browser is redirected to `https://example.com/CALLBACK ` with a potentially successful OAuth2 response.

This patch addresses all of these issues and adds regression tests to keep the implementation secure in future releases.

### Bug Fixes

- Make redirect URL checking more strict ([cdee51e](https://github.com/ory/fosite/commit/cdee51ebe721bfc8acca0fd0b86b030ca70867bf)):

  The OAuth 2.0 Client's Redirect URL and the Redirect URL used in the OAuth 2.0 flow do not check if the query string is equal:

  1. Registering a client with allowed redirect URL `https://example.com/callback`
  2. Performing OAuth2 flow and requesting redirect URL `https://example.com/callback?bar=foo`
  3. Instead of an error, the browser is redirected to `https://example.com/callback?bar=foo` with a potentially successful OAuth2 response.

  Additionally, matching Redirect URLs used `strings.ToLower` normalization:

  1. Registering a client with allowed redirect URL `https://example.com/callback`
  2. Performing OAuth2 flow and requesting redirect URL `https://example.com/CALLBACK`
  3. Instead of an error, the browser is redirected to `https://example.com/CALLBACK ` with a potentially successful OAuth2 response.

  This patch addresses all of these issues and adds regression tests to keep the implementation secure in future releases.

### Documentation

- Added missing dot ([#487](https://github.com/ory/fosite/issues/487)) ([a822244](https://github.com/ory/fosite/commit/a82224430292b2f209d011f107998273d568912b))

# [0.34.0](https://github.com/ory/fosite/compare/v0.33.0...v0.34.0) (2020-09-24)

chore: fix unused const linter error (#484)

## Breaking Changes

`fosite.ErrRevocationClientMismatch` was removed because it is not part of [RFC 6749](https://tools.ietf.org/html/rfc6749#section-5.2). Instead, `fosite.ErrUnauthorizedClient` will be returned when calling `RevokeToken` with an OAuth2 Client which does not match the Access or Refresh Token to be revoked.

### Bug Fixes

- Full JSON escaping ([#481](https://github.com/ory/fosite/issues/481)) ([0943a10](https://github.com/ory/fosite/commit/0943a1095a209fdfb2f8a29524b676ee9c9650a1))
- Ignore x/net false positives ([#483](https://github.com/ory/fosite/issues/483)) ([aead149](https://github.com/ory/fosite/commit/aead1499deb8b08f48bcc196a88e5715702b5431))

### Chores

- Fix unused const linter error ([#484](https://github.com/ory/fosite/issues/484)) ([3540462](https://github.com/ory/fosite/commit/354046265cd4ffcbff8465e4b7a7ea7b6741c5e4))

### Features

- Errors now wrap underlying errors ([#479](https://github.com/ory/fosite/issues/479)) ([b53f8f5](https://github.com/ory/fosite/commit/b53f8f58f0b9889d044cf9a8e2604316f0559ff6)), closes [#458](https://github.com/ory/fosite/issues/458)

### Unclassified

- Merge pull request from GHSA-7mqr-2v3q-v2wm ([03dd558](https://github.com/ory/fosite/commit/03dd55813f5521985f7dd64277b7ba0cf1441319))

# [0.33.0](https://github.com/ory/fosite/compare/v0.32.4...v0.33.0) (2020-09-16)

feat: error_hint and error_debug are now exposed through error_description (#460)

BREAKING CHANGE: Merges the error description with error hint and error debug, making it easier to consume error messages in standardized OAuth2 clients.

## Breaking Changes

Merges the error description with error hint and error debug, making it easier to consume error messages in standardized OAuth2 clients.

### Features

- Error_hint and error_debug are now exposed through error_description ([#460](https://github.com/ory/fosite/issues/460)) ([8daab21](https://github.com/ory/fosite/commit/8daab21f97c513101d224a7ad7a44b871440be57))

# [0.32.4](https://github.com/ory/fosite/compare/v0.32.3...v0.32.4) (2020-09-15)

autogen(docs): regenerate and update changelog

### Code Generation

- **docs:** Regenerate and update changelog ([1f16df0](https://github.com/ory/fosite/commit/1f16df0862bbcdfba98644d1c8fce8a9f92bbbec))

### Code Refactoring

- Fix inconsistent spelling of revocation ([#477](https://github.com/ory/fosite/issues/477)) ([7a55edb](https://github.com/ory/fosite/commit/7a55edbb67738a721c5f1a8f58d2db67f6738f65))

### Documentation

- Fix minor typos ([#475](https://github.com/ory/fosite/issues/475)) ([23cc9c1](https://github.com/ory/fosite/commit/23cc9c1d29f35a73acbf05fe6b505b692f6fe49c))

# [0.32.3](https://github.com/ory/fosite/compare/v0.32.2...v0.32.3) (2020-09-12)

fix: add missing OAuth2TokenRevocationFactory to ComposeAllEnabled (#472)

### Bug Fixes

- Add missing OAuth2TokenRevocationFactory to ComposeAllEnabled ([#472](https://github.com/ory/fosite/issues/472)) ([88587fd](https://github.com/ory/fosite/commit/88587fde8fc92137660383c401250e492716c396))
- Align error returned when a grant_type was requested that's not allowed for a client ([#467](https://github.com/ory/fosite/issues/467)) ([3c30c0d](https://github.com/ory/fosite/commit/3c30c0d9f1e62b237acc845d5b3a42d1ea9a80c0)), closes [/tools.ietf.org/html/rfc6749#section-5](https://github.com//tools.ietf.org/html/rfc6749/issues/section-5):

  Returned error was 'invalid_grant'.

- All responses now contain headers to not cache them ([#465](https://github.com/ory/fosite/issues/465)) ([2012cb7](https://github.com/ory/fosite/commit/2012cb7ec6feb504d1faa6e393fce8d25edafebb))
- No cache headers followup ([#466](https://github.com/ory/fosite/issues/466)) ([1627c6a](https://github.com/ory/fosite/commit/1627c6ab31cb151f01671cd3403bc3c7de6fcfbd))

### Code Refactoring

- Copy all values when sanitizing ([#455](https://github.com/ory/fosite/issues/455)) ([c80d0d4](https://github.com/ory/fosite/commit/c80d0d42a34f8cf664d44c687d7cfea576a0b232))

### Documentation

- Add empty session example explanation ([#450](https://github.com/ory/fosite/issues/450)) ([36d65cb](https://github.com/ory/fosite/commit/36d65cbc061ff6cae38e90b0a6954646c8daf5d7))
- Better section reference for GetRedirectURIFromRequestValues ([#463](https://github.com/ory/fosite/issues/463)) ([48a3daf](https://github.com/ory/fosite/commit/48a3daf45bd1885c4412eeb9b2bc3117b6075de9))
- Deprecate history.md ([b0d5fea](https://github.com/ory/fosite/commit/b0d5feacfcbeedf609563fa8567bd0e031b179b5)), closes [/github.com/ory/fosite/issues/414#issuecomment-662538622](https://github.com//github.com/ory/fosite/issues/414/issues/issuecomment-662538622)

### Features

- Add locking to memory storage ([#471](https://github.com/ory/fosite/issues/471)) ([4687147](https://github.com/ory/fosite/commit/46871476b1f47cefc09888615f70dd9fdd5af8b3))
- Make MinParameterEntropy configurable ([#461](https://github.com/ory/fosite/issues/461)) ([2c793e6](https://github.com/ory/fosite/commit/2c793e6c010ac6cbc552200197ae1262d91c2bda)), closes [#267](https://github.com/ory/fosite/issues/267)
- New compose strategies for ES256 ([#446](https://github.com/ory/fosite/issues/446)) ([39053ee](https://github.com/ory/fosite/commit/39053eedaa687fe1d8dbe8b928fb98cd5ce8c021))

# [0.32.2](https://github.com/ory/fosite/compare/v0.32.1...v0.32.2) (2020-06-22)

feat: new factory with default issuer for JWT tokens (#444)

### Features

- New factory with default issuer for JWT tokens ([#444](https://github.com/ory/fosite/issues/444)) ([901e206](https://github.com/ory/fosite/commit/901e206d03b615c189e12f94607d92c10d6909fa))

# [0.32.1](https://github.com/ory/fosite/compare/v0.32.0...v0.32.1) (2020-06-05)

feat: makeRemoveEmpty public (#443)

### Bug Fixes

- Improved error messages in client authentication ([#440](https://github.com/ory/fosite/issues/440)) ([c06e560](https://github.com/ory/fosite/commit/c06e5608c7ae6a0243428252e6ec80bc37ae33ca)), closes [#436](https://github.com/ory/fosite/issues/436)

### Features

- MakeRemoveEmpty public ([#443](https://github.com/ory/fosite/issues/443)) ([17b0756](https://github.com/ory/fosite/commit/17b075688f9a012b09e650e90d765de6d4d538cf))

# [0.32.0](https://github.com/ory/fosite/compare/v0.31.3...v0.32.0) (2020-05-28)

feat: added support for ES256 token strategy and client authentication (#439)

I added to `DefaultOpenIDConnectClient` a field `TokenEndpointAuthSigningAlgorithm` to be able to configure what `GetTokenEndpointAuthSigningAlgorithm` returns. I also cleaned some other places where there were assumptions about only RSA keys.

Closes #429

### Bug Fixes

- **arguments:** Fixes a logic bug in MatchesExact and adds documentation ([#433](https://github.com/ory/fosite/issues/433)) ([10fd67b](https://github.com/ory/fosite/commit/10fd67bf84118affc9269ca0c0dbc8da4b0bf2cd)):

- Double-decoding of client credentials in request body ([#434](https://github.com/ory/fosite/issues/434)) ([48c9b41](https://github.com/ory/fosite/commit/48c9b41ea2dc89ec2bf58ba918c45c8430bb0ccd)):

  I noticed that client credentials are URL-decoded after being extracted from the POST body form, which was already URL-decoded by Go. The accompanying error message suggests this was copied and pasted from the HTTP basic authorization header handling, which is the only place where the extra URL-decoding was needed (as per the OAuth 2.0 spec). The result is that client credentials containing %-prefixed sequences, whether valid sequences or not, are going to fail validation.

  Remove the extra URL decoding. Add tests that ensure client credentials work with special characters in both the HTTP basic auth header and in the request body.

### Documentation

- Update github templates ([#432](https://github.com/ory/fosite/issues/432)) ([b393832](https://github.com/ory/fosite/commit/b393832765e0c97661bb5495e3a3d51a8019afd7))
- Update repository templates ([a840a62](https://github.com/ory/fosite/commit/a840a62e401b4111f8304fa8b963006a866a20f8))

### Features

- Added support for ES256 token strategy and client authentication ([#439](https://github.com/ory/fosite/issues/439)) ([36eb661](https://github.com/ory/fosite/commit/36eb661cc8b609877d8e81c849c34631bbab245a)), closes [#429](https://github.com/ory/fosite/issues/429):

  I added to `DefaultOpenIDConnectClient` a field `TokenEndpointAuthSigningAlgorithm` to be able to configure what `GetTokenEndpointAuthSigningAlgorithm` returns. I also cleaned some other places where there were assumptions about only RSA keys.

# [0.31.3](https://github.com/ory/fosite/compare/v0.31.2...v0.31.3) (2020-05-09)

feat(pkce): add EnforcePKCEForPublicClients config flag (#431)

Alternative proposal for the issue discussed in #389 and #391, where enforcement of PKCE is wanted only for certain clients.

Add a new flag EnforcePKCEForPublicClients which enforces PKCE only for public clients. The error hint is slightly different, as it mentions PKCE is enforced for "this client" rather than "clients". (It intentionally does not mention why it's enforced, as I think basing it on public clients is an implementation detail that servers may want to change without adding to the error hints).

Closes #389
Closes #391

### Bug Fixes

- Do not issue refresh tokens to clients who cannot use it ([#430](https://github.com/ory/fosite/issues/430)) ([792670d](https://github.com/ory/fosite/commit/792670d0e81ff83f2b345502ea7adadf99bcaa9b)), closes [#370](https://github.com/ory/fosite/issues/370)

### Features

- **pkce:** Add EnforcePKCEForPublicClients config flag ([#431](https://github.com/ory/fosite/issues/431)) ([9f53c84](https://github.com/ory/fosite/commit/9f53c843e4a72d0ff34acb084e5a920d7114278f)), closes [#389](https://github.com/ory/fosite/issues/389) [#391](https://github.com/ory/fosite/issues/391) [#389](https://github.com/ory/fosite/issues/389) [#391](https://github.com/ory/fosite/issues/391)

# [0.31.2](https://github.com/ory/fosite/compare/v0.31.1...v0.31.2) (2020-04-16)

fix: introduce better linting pipeline and resolve Go issues (#428)

### Bug Fixes

- Introduce better linting pipeline and resolve Go issues ([#428](https://github.com/ory/fosite/issues/428)) ([e02f731](https://github.com/ory/fosite/commit/e02f731a41fb82ac8d6b62ea3f6fd8a915526090))

# [0.31.1](https://github.com/ory/fosite/compare/v0.31.0...v0.31.1) (2020-04-16)

fix: return invalid_grant instead of invalid_request in refresh flow (#427)

Return invalid_grant instead of invalid_request when in authorization code flow when the user is not the owner of the authorization code or if the redirect uri doesn't match from the authorization request.

Co-authored-by: Damien Bravin <damienbr@users.noreply.github.com>

### Bug Fixes

- List all response types in example memory store ([#413](https://github.com/ory/fosite/issues/413)) ([427d40d](https://github.com/ory/fosite/commit/427d40dcaadab6933a4e571def7d9729fd442581)), closes [#304](https://github.com/ory/fosite/issues/304)
- Return invalid_grant instead of invalid_request in refresh flow ([#427](https://github.com/ory/fosite/issues/427)) ([f5a0e96](https://github.com/ory/fosite/commit/f5a0e9696750e3f1d67bd919a6588b175e7cc2bb)):

  Return invalid_grant instead of invalid_request when in authorization code flow when the user is not the owner of the authorization code or if the redirect uri doesn't match from the authorization request.

- **storage:** Remove unused field ([#422](https://github.com/ory/fosite/issues/422)) ([d2eb3b9](https://github.com/ory/fosite/commit/d2eb3b9ff5f52810067ac59969a3c4272772bdb3)), closes [#417](https://github.com/ory/fosite/issues/417)
- **storage:** Remove unused methods ([#417](https://github.com/ory/fosite/issues/417)) ([023bdcf](https://github.com/ory/fosite/commit/023bdcf1217b8f86de250f53391ad3b1e356949d))

### Documentation

- Fix various typos ([#415](https://github.com/ory/fosite/issues/415)) ([719aaa0](https://github.com/ory/fosite/commit/719aaa0b695f02556167f02fc94133a380ccfa16))
- Replace Discord with Slack ([#412](https://github.com/ory/fosite/issues/412)) ([d8591bb](https://github.com/ory/fosite/commit/d8591bba33d16b61e6c611b7042d695166bd94e5))
- Update github templates ([#424](https://github.com/ory/fosite/issues/424)) ([d37fc4b](https://github.com/ory/fosite/commit/d37fc4babe43b52c92eb081b9ea78c0fa9f51865))
- Update github templates ([#425](https://github.com/ory/fosite/issues/425)) ([0399871](https://github.com/ory/fosite/commit/039987119ea78d69fe991bbb0edb6735b88b16cc))
- Update SetSession comment ([#423](https://github.com/ory/fosite/issues/423)) ([32951ab](https://github.com/ory/fosite/commit/32951ab56fb3400ff6980519c2e6e20802292f2f))
- Updates issue and pull request templates ([#419](https://github.com/ory/fosite/issues/419)) ([d804da1](https://github.com/ory/fosite/commit/d804da1e3dfda46872d358d2987bd19462c03e98))

# [0.31.0](https://github.com/ory/fosite/compare/v0.30.6...v0.31.0) (2020-03-29)

Merge pull request from GHSA-v3q9-2p3m-7g43

- u

- u

### Unclassified

- Merge pull request from GHSA-v3q9-2p3m-7g43 ([0c9e0f6](https://github.com/ory/fosite/commit/0c9e0f6d654913ad57c507dd9a36631e1858a3e9)):

  - u

  - u

# [0.30.6](https://github.com/ory/fosite/compare/v0.30.5...v0.30.6) (2020-03-26)

fix: handle serialization errors that can be thrown by call to 'Commit' (#403)

### Bug Fixes

- Handle serialization errors that can be thrown by call to 'Commit' ([#403](https://github.com/ory/fosite/issues/403)) ([35a1558](https://github.com/ory/fosite/commit/35a1558d8d845ac15bc6ec99fb4be062716b231a))

### Documentation

- Update forum and chat links ([b1ba04e](https://github.com/ory/fosite/commit/b1ba04e447d6dfdaf9f0c84336d3bacab41b2c8d))

# [0.30.5](https://github.com/ory/fosite/compare/v0.30.4...v0.30.5) (2020-03-25)

fix: handle concurrent transactional errors in the refresh token grant handler (#402)

This commit provides the functionality required to address https://github.com/ory/hydra/issues/1719 & https://github.com/ory/hydra/issues/1735 by adding error checking to the RefreshTokenGrantHandler's PopulateTokenEndpointResponse method so it can deal with errors due to concurrent access. This will allow the authorization server to render a better error to the user-agent.

No longer returns fosite.ErrServerError in the event the storage. Instead a wrapped fosite.ErrNotFound is returned when fetching the refresh token fails due to it no longer being present. This scenario is caused when the user sends two or more request to refresh using the same token and one request gets into the handler just after the prior request finished and successfully committed its transaction.

Adds unit test coverage for transaction error handling logic added to the RefreshTokenGrantHandler's PopulateTokenEndpointResponse method

### Bug Fixes

- Handle concurrent transactional errors in the refresh token grant handler ([#402](https://github.com/ory/fosite/issues/402)) ([b17190b](https://github.com/ory/fosite/commit/b17190b4964e911d6f94379873139cdfc3def5bd)):

  This commit provides the functionality required to address https://github.com/ory/hydra/issues/1719 & https://github.com/ory/hydra/issues/1735 by adding error checking to the RefreshTokenGrantHandler's PopulateTokenEndpointResponse method so it can deal with errors due to concurrent access. This will allow the authorization server to render a better error to the user-agent.

  No longer returns fosite.ErrServerError in the event the storage. Instead a wrapped fosite.ErrNotFound is returned when fetching the refresh token fails due to it no longer being present. This scenario is caused when the user sends two or more request to refresh using the same token and one request gets into the handler just after the prior request finished and successfully committed its transaction.

  Adds unit test coverage for transaction error handling logic added to the RefreshTokenGrantHandler's PopulateTokenEndpointResponse method

# [0.30.4](https://github.com/ory/fosite/compare/v0.30.3...v0.30.4) (2020-03-17)

fix: add ability to specify amr values natively in id_token payload (#401)

See ory/hydra#1756

### Bug Fixes

- Add ability to specify amr values natively in id_token payload ([#401](https://github.com/ory/fosite/issues/401)) ([f99bb80](https://github.com/ory/fosite/commit/f99bb8012a583b25fd591718a51308c208cb9a55)), closes [ory/hydra#1756](https://github.com/ory/hydra/issues/1756)

# [0.30.3](https://github.com/ory/fosite/compare/v0.30.2...v0.30.3) (2020-03-04)

fix: Support RFC8252#section-7.3 Loopback Interface Redirection (#400)

Closes #284

### Bug Fixes

- Merge request ID as well ([#398](https://github.com/ory/fosite/issues/398)) ([67c081c](https://github.com/ory/fosite/commit/67c081cb5cb650e7095d7343a618484103cf8bb5)), closes [#386](https://github.com/ory/fosite/issues/386)
- Support RFC8252#section-7.3 Loopback Interface Redirection ([#400](https://github.com/ory/fosite/issues/400)) ([4104135](https://github.com/ory/fosite/commit/41041350c06853d490e94849b25d0fee87a95a32)), closes [RFC8252#section-7](https://github.com/RFC8252/issues/section-7) [#284](https://github.com/ory/fosite/issues/284)

### Documentation

- Add undocumented ExactScopeStrategy ([#395](https://github.com/ory/fosite/issues/395)) ([387cade](https://github.com/ory/fosite/commit/387cade4c6e96e0b83df274da5835691e54d07af))
- Updates issue and pull request templates ([#393](https://github.com/ory/fosite/issues/393)) ([cdefb3e](https://github.com/ory/fosite/commit/cdefb3e99e73b69e62a449de489b0e806d5158af))
- Updates issue and pull request templates ([#394](https://github.com/ory/fosite/issues/394)) ([119e6ab](https://github.com/ory/fosite/commit/119e6ab6d83ab8dee3fd31085153f64ca008582a))

### Features

- Add ExactOne and MatchesExact to Arguments ([#399](https://github.com/ory/fosite/issues/399)) ([cf23400](https://github.com/ory/fosite/commit/cf23400930e63a6d5244262d284ddc79943775e6)):

  Previously Arguments.Exact had vague semantic where
  it coudln't distinguish between value with a space and multiple
  values. Split it into 2 functions with clear semantic.

  Old .Exact() remains for compatibility and marked as deprecated

# [0.30.2](https://github.com/ory/fosite/compare/v0.30.1...v0.30.2) (2019-11-21)

Return state parameter in authorization error conditions (#388)

Related to ory/hydra#1642

### Unclassified

- Return state parameter in authorization error conditions (#388) ([3ece795](https://github.com/ory/fosite/commit/3ece795f3080db5de3529cea9bfa670e70704686)), closes [#388](https://github.com/ory/fosite/issues/388) [ory/hydra#1642](https://github.com/ory/hydra/issues/1642)
- Revert incorrect license changes ([40a49f7](https://github.com/ory/fosite/commit/40a49f743dff60d07b6314667933a47dbf2635aa))

# [0.30.1](https://github.com/ory/fosite/compare/v0.30.0...v0.30.1) (2019-09-23)

pkce: Enforce verifier formatting (#383)

### Unclassified

- Enforce verifier formatting ([#383](https://github.com/ory/fosite/issues/383)) ([024667a](https://github.com/ory/fosite/commit/024667ac1905a4d0274294ab552f3566e2eb3b6a))

# [0.30.0](https://github.com/ory/fosite/compare/v0.29.8...v0.30.0) (2019-09-16)

handler/pkce: Enable PKCE for private clients (#382)

### Unclassified

- handler/pkce: Enable PKCE for private clients (#382) ([e21830e](https://github.com/ory/fosite/commit/e21830ec0c0c37ca6ca5544b1362c85abe38b80f)), closes [#382](https://github.com/ory/fosite/issues/382)
- Add RefreshTokenScopes Config (#371) ([bcc7859](https://github.com/ory/fosite/commit/bcc78599eadbff38dc0efc9370e5ef64eadfefa9)), closes [#371](https://github.com/ory/fosite/issues/371):

  When set to true, this will return refresh tokens even if the user did
  not ask for the offline or offline_access Oauth Scope.

# [0.29.8](https://github.com/ory/fosite/compare/v0.29.7...v0.29.8) (2019-08-29)

handler/revoke: respecting ErrInvalidRequest code (#380)

This commit modifies the case for ErrInvalidRequest in
WriteRevocationResponse to respect the 400 error code
and not fallthrough to ErrInvalidClient.

Author: DefinitelyNotAGoat <baldrich@protonmail.com>

### Documentation

- Updates issue and pull request templates ([#376](https://github.com/ory/fosite/issues/376)) ([165e93e](https://github.com/ory/fosite/commit/165e93eeff7d187af682f7f958b39e2393d15821))
- Updates issue and pull request templates ([#377](https://github.com/ory/fosite/issues/377)) ([40590cb](https://github.com/ory/fosite/commit/40590cbaa45167dff2085483ccf5b4bddb37e422))
- Updates issue and pull request templates ([#378](https://github.com/ory/fosite/issues/378)) ([54426bb](https://github.com/ory/fosite/commit/54426bbf3d3bb125753aaf7fda5a7ded5effdf4c))

### Unclassified

- handler/revoke: respecting ErrInvalidRequest code (#380) ([cc34bfb](https://github.com/ory/fosite/commit/cc34bfb4f970d25f59948dcdcbc0eb587ae78d6d)), closes [#380](https://github.com/ory/fosite/issues/380):

  This commit modifies the case for ErrInvalidRequest in
  WriteRevocationResponse to respect the 400 error code
  and not fallthrough to ErrInvalidClient.

  Author: DefinitelyNotAGoat <baldrich@protonmail.com>

# [0.29.7](https://github.com/ory/fosite/compare/v0.29.6...v0.29.7) (2019-08-06)

pkce: Return error when PKCE is used with private clients (#375)

### Documentation

- Fix method/struct documents ([#360](https://github.com/ory/fosite/issues/360)) ([ad06f22](https://github.com/ory/fosite/commit/ad06f2266b28b3d1844f36e97c1118822fd2a46c))
- Updates issue and pull request templates ([#361](https://github.com/ory/fosite/issues/361)) ([35157e2](https://github.com/ory/fosite/commit/35157e2a5174f1a8ee9074452b77953e35c4161c))
- Updates issue and pull request templates ([#365](https://github.com/ory/fosite/issues/365)) ([90a3c50](https://github.com/ory/fosite/commit/90a3c509e718445b799821fac400aad28d9de928))
- Updates issue and pull request templates ([#366](https://github.com/ory/fosite/issues/366)) ([27c64ec](https://github.com/ory/fosite/commit/27c64ec1b7d12ee1b1e1e0d35dc6b24f7ade92e0))
- Updates issue and pull request templates ([#367](https://github.com/ory/fosite/issues/367)) ([01cd955](https://github.com/ory/fosite/commit/01cd955efe9a00c014a5ef7488774c3913e7218d))
- Updates issue and pull request templates ([#373](https://github.com/ory/fosite/issues/373)) ([5962474](https://github.com/ory/fosite/commit/5962474c904f80517d1a9c2731e703ffda972d6a))
- Updates issue and pull request templates ([#374](https://github.com/ory/fosite/issues/374)) ([9f7cf40](https://github.com/ory/fosite/commit/9f7cf409a643b72cfa25dd2f1340f1aa1c17c443))

### Unclassified

- Create FUNDING.yml ([1b7b479](https://github.com/ory/fosite/commit/1b7b479ca040f95f3ea4cff642c7f678df5cb0ab))
- Return error when PKCE is used with private clients ([#375](https://github.com/ory/fosite/issues/375)) ([7219387](https://github.com/ory/fosite/commit/72193870c9914dc97c1117a566c68bede0bf5290))

# [0.29.6](https://github.com/ory/fosite/compare/v0.29.5...v0.29.6) (2019-04-26)

openid: Allow promp=none for https/localhost (#359)

Signed-off-by: aeneasr <aeneas@ory.sh>

### Unclassified

- Allow promp=none for https/localhost ([#359](https://github.com/ory/fosite/issues/359)) ([27bbe00](https://github.com/ory/fosite/commit/27bbe0033273157ea449310c064675127e2550e6))

# [0.29.5](https://github.com/ory/fosite/compare/v0.29.4...v0.29.5) (2019-04-25)

core: Add debug log to invalid_client error(#358)

Signed-off-by: nerocrux <nerocrux@gmail.com>

### Unclassified

- Add debug log to invalid_client error([#358](https://github.com/ory/fosite/issues/358)) ([dce3111](https://github.com/ory/fosite/commit/dce3111ad0dac62911c19d9b6ea4cb776f087c4d))

# [0.29.3](https://github.com/ory/fosite/compare/v0.29.2...v0.29.3) (2019-04-17)

Export IsLocalhost

Signed-off-by: aeneasr <aeneas@ory.sh>

### Unclassified

- Export IsLocalhost ([a95ea09](https://github.com/ory/fosite/commit/a95ea092ef682cd5fe3449c23245d211444f28cb))
- Improve IsRedirectURISecure check ([d6f8962](https://github.com/ory/fosite/commit/d6f8962de5336ce17128b1fd238cba13862c85a7))

# [0.29.2](https://github.com/ory/fosite/compare/v0.29.1...v0.29.2) (2019-04-11)

Allow providing a custom redirect URI checker (#355)

Signed-off-by: aeneasr <aeneas@ory.sh>

### Unclassified

- Allow providing a custom redirect URI checker (#355) ([3d16e39](https://github.com/ory/fosite/commit/3d16e39a3b25cb5d77b8b10cb568c9bc2a835356)), closes [#355](https://github.com/ory/fosite/issues/355)

# [0.29.1](https://github.com/ory/fosite/compare/v0.29.0...v0.29.1) (2019-03-27)

token: Improve rotated secret error reporting in HMAC strategy (#354)

Signed-off-by: aeneasr <aeneas@ory.sh>

### Unclassified

- Improve rotated secret error reporting in HMAC strategy ([#354](https://github.com/ory/fosite/issues/354)) ([f21d930](https://github.com/ory/fosite/commit/f21d930291ada9e609ea5018693d6e4745815f03))
- Propagate session data properly ([#353](https://github.com/ory/fosite/issues/353)) ([5ba0f04](https://github.com/ory/fosite/commit/5ba0f0465039e7072593205b1252e630d340d6ab)):

  This example is slightly inaccurate; the session data will need to come from the returned AccessRequester, not the pre-created session. The session passed to IntrospectToken isn't mutated.

- Remove useless details fn receiver ([#349](https://github.com/ory/fosite/issues/349)) ([af403c6](https://github.com/ory/fosite/commit/af403c6fac913736a05ca0c44765b10baaf89295))
- Update HISTORY.md, README.md, CONTRIBUTING.md ([#347](https://github.com/ory/fosite/issues/347)) ([de5e61e](https://github.com/ory/fosite/commit/de5e61e0eb445af57e692964057ea8e661f98618)):

  - README: Breaks out `0.26.0` as was stuck inside a code block.
  - README: Ensures the later versions formats code blocks as Go code.
  - Runs doctoc to ensure TOCs are up to date.

# [0.29.0](https://github.com/ory/fosite/compare/v0.28.1...v0.29.0) (2018-12-23)

oauth2: add test coverage to exercise the transactional support in the AuthorizeExplicitGrantHandler's PopulateTokenEndpointResponse method.

Signed-off-by: Amir Aslaminejad <aslaminejad@gmail.com>

### Unclassified

- Add mock for storage.Transactional + update generate-mocks.sh ([03f7bc8](https://github.com/ory/fosite/commit/03f7bc8e59f15d7b9c0df47c8c77c106f3fd4a0c))
- Add test coverage to exercise the transactional support in the AuthorizeExplicitGrantHandler's PopulateTokenEndpointResponse method. ([2f58f9e](https://github.com/ory/fosite/commit/2f58f9e0ea1a197c8b7eb62dc545d9467ed2ff10))
- Add test coverage to exercise the transactional support in the RefreshTokenGrantHandler's PopulateTokenEndpointResponse method. ([b38d7c8](https://github.com/ory/fosite/commit/b38d7c89b9a45b7576af379b2dc479ddb880195c))
- Adds new interface `Transactional` which is to be implemented by storage providers that can support transactions. ([c364b33](https://github.com/ory/fosite/commit/c364b33eefe813da4da02fc78d9e72e1d5301234))
- Don't double encode URL fragments ([#346](https://github.com/ory/fosite/issues/346)) ([1f41934](https://github.com/ory/fosite/commit/1f419341886c8e37a10c68d7a5c8d576176e666a)), closes [#345](https://github.com/ory/fosite/issues/345)
- Use transactions in the auth code token flow (if the storage implementation implements the `Transactional` interface) to address [#309](https://github.com/ory/fosite/issues/309) ([e00c567](https://github.com/ory/fosite/commit/e00c5675182eb5d90644160c0f3f1b10f0f287f4))
- Use transactions in the refresh token flow (if the storage implementation implements the `Transactional` interface) to address [#309](https://github.com/ory/fosite/issues/309) ([07d1a39](https://github.com/ory/fosite/commit/07d1a3974ff6d53c239c4050703b09928f484e01))

# [0.28.1](https://github.com/ory/fosite/compare/v0.28.0...v0.28.1) (2018-12-04)

compose: Expose token entropy setting (#342)

Signed-off-by: nerocrux <nerocrux@gmail.com>

### Unclassified

- Remove cryptopasta dependency (#339) ([b156e6b](https://github.com/ory/fosite/commit/b156e6b48383926974a560bb416a9ac7507347ec)), closes [#339](https://github.com/ory/fosite/issues/339)
- Expose token entropy setting ([#342](https://github.com/ory/fosite/issues/342)) ([0761fca](https://github.com/ory/fosite/commit/0761fcae7e6ecd0f7d16c51a3c7fa3891d85d85b))

# [0.28.0](https://github.com/ory/fosite/compare/v0.27.4...v0.28.0) (2018-11-16)

oauth2: Add ability to specify refresh token lifespan (#337)

Set it to `-1` to disable this feature. Defaults to 30 days.

Closes #319

Signed-off-by: arekkas <aeneas@ory.am>

### Unclassified

- Add ability to specify refresh token lifespan ([#337](https://github.com/ory/fosite/issues/337)) ([fa65408](https://github.com/ory/fosite/commit/fa654089e09900d842117827ec2f6258323ec436)), closes [#319](https://github.com/ory/fosite/issues/319):

  Set it to `-1` to disable this feature. Defaults to 30 days.

# [0.27.4](https://github.com/ory/fosite/compare/v0.27.3...v0.27.4) (2018-11-12)

docs: Fix quickstart (#335)

- replace NewMemoryStore with NewExampleStore
- fix length of signing key
- fix config type

Signed-off-by: Peter Schultz <peter.schultz@classmarkets.com>

### Documentation

- Fix quickstart ([#335](https://github.com/ory/fosite/issues/335)) ([25cc6c4](https://github.com/ory/fosite/commit/25cc6c42e2befe3b200d79c9d8edac47cc6d3f86)):

  - replace NewMemoryStore with NewExampleStore
  - fix length of signing key
  - fix config type

### Unclassified

- Omit exp if ExpiresAt is zero value ([#334](https://github.com/ory/fosite/issues/334)) ([6d50176](https://github.com/ory/fosite/commit/6d501761a17bc3a720e2a0b72ff5f218fa72660c))

# [0.27.3](https://github.com/ory/fosite/compare/v0.27.2...v0.27.3) (2018-11-08)

oauth2: Set exp for authorize code issued by hybrid flow (#333)

Signed-off-by: nerocrux <nerocrux@gmail.com>

### Unclassified

- Set exp for authorize code issued by hybrid flow ([#333](https://github.com/ory/fosite/issues/333)) ([d275e84](https://github.com/ory/fosite/commit/d275e84dc6f4bf4e71393672e0e16d54b401bc3c))

# [0.27.2](https://github.com/ory/fosite/compare/v0.27.1...v0.27.2) (2018-11-07)

pkce: Allow hybrid flows (#328)

Signed-off-by: Adam Shannon <adamkshannon@gmail.com>
Signed-off-by: Wenhao Ni <niwenhao@gmail.com>

### Unclassified

- Allow hybrid flows ([#328](https://github.com/ory/fosite/issues/328)) ([cdfddc8](https://github.com/ory/fosite/commit/cdfddc8b06d861708ebe3494a35d65da2d2fcef8)):

  Signed-off-by: Wenhao Ni <niwenhao@gmail.com>

# [0.27.1](https://github.com/ory/fosite/compare/v0.27.0...v0.27.1) (2018-11-03)

oauth2: Improve refresh security and reliability (#332)

This patch resolves several issues regarding the refresh flow. First,
an issue has been resolved which caused the audience to not be
set in the refreshed access tokens.

Second, scope and audience are validated against the client's
whitelisted values and if the values are no longer allowed,
the grant is canceled.

Closes #331
Closes #325
Closes #324

### Unclassified

- Improve refresh security and reliability ([#332](https://github.com/ory/fosite/issues/332)) ([4e4121b](https://github.com/ory/fosite/commit/4e4121bac5cda8efa7d3eb6aaf7720f3ff59c329)), closes [#331](https://github.com/ory/fosite/issues/331) [#325](https://github.com/ory/fosite/issues/325) [#324](https://github.com/ory/fosite/issues/324):

  This patch resolves several issues regarding the refresh flow. First,
  an issue has been resolved which caused the audience to not be
  set in the refreshed access tokens.

  Second, scope and audience are validated against the client's
  whitelisted values and if the values are no longer allowed,
  the grant is canceled.

# [0.27.0](https://github.com/ory/fosite/compare/v0.26.1...v0.27.0) (2018-10-31)

oauth2: Update jwt access token interface (#330)

The interface needed to change in order to natively handle the audience claim.

Signed-off-by: arekkas <aeneas@ory.am>

### Unclassified

- Introduce audience capabilities ([#327](https://github.com/ory/fosite/issues/327)) ([e2441d2](https://github.com/ory/fosite/commit/e2441d231a19cd1133b3316d3477b84d7b649522)), closes [#326](https://github.com/ory/fosite/issues/326):

  This patch allows clients to whitelist audiences and request that audiences are set for oauth2 access and refresh tokens

- Update jwt access token interface ([#330](https://github.com/ory/fosite/issues/330)) ([2da9764](https://github.com/ory/fosite/commit/2da976477fcd41493103ea478541d68ca04083ae)):

  The interface needed to change in order to natively handle the audience claim.

# [0.26.1](https://github.com/ory/fosite/compare/v0.26.0...v0.26.1) (2018-10-25)

hash: Raise bcrypt cost factor lower bound (#321)

Users of this library can easily create the following:

hasher := fosite.BCrypt{}
hasher.Hash(..)

This is a problem because WorkFactor will default to 0 and x/crypto/bcrypt will default that to 4 (See https://godoc.org/golang.org/x/crypto/bcrypt).

Instead this should be some higher cost factor. Callers who need a lower WorkFactor can still lower the cost, if needed.

Signed-off-by: Adam Shannon <adamkshannon@gmail.com>

### Unclassified

- Fix Config.GetHashCost godoc comment ([#320](https://github.com/ory/fosite/issues/320)) ([4d2b119](https://github.com/ory/fosite/commit/4d2b119b7a302bf7e6a4d9b600697e08cf089b02))
- Fix doc typo ([#322](https://github.com/ory/fosite/issues/322)) ([239b1ed](https://github.com/ory/fosite/commit/239b1ed4b9b406287fa49e01f8316e5fc4eb7923))
- Raise bcrypt cost factor lower bound ([#321](https://github.com/ory/fosite/issues/321)) ([799fc70](https://github.com/ory/fosite/commit/799fc70a48b68b3403eb150084c28d4e78c035e4)):

  Users of this library can easily create the following:

  hasher := fosite.BCrypt{}
  hasher.Hash(..)

  This is a problem because WorkFactor will default to 0 and x/crypto/bcrypt will default that to 4 (See https://godoc.org/golang.org/x/crypto/bcrypt).

  Instead this should be some higher cost factor. Callers who need a lower WorkFactor can still lower the cost, if needed.

# [0.26.0](https://github.com/ory/fosite/compare/v0.25.1...v0.26.0) (2018-10-24)

all: Rearrange commits with goreturns

Signed-off-by: aeneasr <aeneas@ory.sh>

### Unclassified

- Allow customization of JWT claims ([f97e451](https://github.com/ory/fosite/commit/f97e45118fbf7a87129ee40c8a56e97efc30c8b9))
- Rearrange commits with goreturns ([211b43b](https://github.com/ory/fosite/commit/211b43b4c04c732adc5fbfa7cab339f44fbea7d7))

# [0.25.1](https://github.com/ory/fosite/compare/v0.25.0...v0.25.1) (2018-10-23)

handler/openid: Populate at_hash in explicit/refresh flows (#315)

Signed-off-by: Wenhao Ni <niwenhao@gmail.com>

### Documentation

- Updates issue and pull request templates ([#313](https://github.com/ory/fosite/issues/313)) ([53c7b55](https://github.com/ory/fosite/commit/53c7b55dba903cdb8071417f39ebc01e00921cd4))
- Updates issue and pull request templates ([#314](https://github.com/ory/fosite/issues/314)) ([73ae623](https://github.com/ory/fosite/commit/73ae6238fc6db4737d5b529ceeb08b26dbab88ea))
- Updates issue and pull request templates ([#316](https://github.com/ory/fosite/issues/316)) ([64299bb](https://github.com/ory/fosite/commit/64299bb72fe0f9f7886bdd061519cc7e9c9081da))

### Unclassified

- handler/openid: Populate at_hash in explicit/refresh flows (#315) ([189589c](https://github.com/ory/fosite/commit/189589c400467460029424226398da709eb9ec48)), closes [#315](https://github.com/ory/fosite/issues/315)
- Fix typo in README.md (#312) ([dcb83ae](https://github.com/ory/fosite/commit/dcb83ae59f984edeb1dfda19d0c0851e2e1574ae)), closes [#312](https://github.com/ory/fosite/issues/312)

# [0.25.0](https://github.com/ory/fosite/compare/v0.24.0...v0.25.0) (2018-10-08)

Fix broken go modules tests (#311)

Signed-off-by: arekkas <aeneas@ory.am>

### Unclassified

- Fix broken go modules tests (#311) ([02ea4b1](https://github.com/ory/fosite/commit/02ea4b186a6384bfe2a36741842f49f7370e0991)), closes [#311](https://github.com/ory/fosite/issues/311)
- Switch from dep to go modules (#310) ([ac46a67](https://github.com/ory/fosite/commit/ac46a67863cb0842d48c83413789a9d6bf595f8a)), closes [#310](https://github.com/ory/fosite/issues/310)

# [0.24.0](https://github.com/ory/fosite/compare/v0.23.0...v0.24.0) (2018-09-27)

Propagate context in jwt strategies (#308)

Closes #307

Signed-off-by: Prateek Malhotra <someone1@gmail.com>

### Unclassified

- Propagate context in jwt strategies (#308) ([e1e18d6](https://github.com/ory/fosite/commit/e1e18d6b22697abeceff6e22d4741c3bf04174f8)), closes [#308](https://github.com/ory/fosite/issues/308) [#307](https://github.com/ory/fosite/issues/307)
- Use test tables for Hasher unit tests (#306) ([499af11](https://github.com/ory/fosite/commit/499af11c14eb4f09f630ce84e971389ab668e94a)), closes [#306](https://github.com/ory/fosite/issues/306)

# [0.23.0](https://github.com/ory/fosite/compare/v0.22.0...v0.23.0) (2018-09-22)

Add breaking change to the Hasher interface to the change log

Signed-off-by: Amir Aslaminejad <aslaminejad@gmail.com>

### Unclassified

- Add breaking change to the Hasher interface to the change log ([805e0e9](https://github.com/ory/fosite/commit/805e0e9a36aa254b18e853b8a9c7881738deb010))
- Update BCrypt to adhere to new Hasher interface ([938e50a](https://github.com/ory/fosite/commit/938e50a32024693670d1a8180b33c5c4a0df470b))
- Update Hasher to take in context ([02f19fa](https://github.com/ory/fosite/commit/02f19fa3a9db72c54c2be6a904f8a2d35792974e))

# [0.22.0](https://github.com/ory/fosite/compare/v0.21.5...v0.22.0) (2018-09-19)

jwt: update JWTStrategy to take in context (#302)

Signed-off-by: Amir Aslaminejad <aslaminejad@gmail.com>

### Unclassified

- Update PR template ([3920be2](https://github.com/ory/fosite/commit/3920be20e78ed304ee3752ffcb997ade12862734))
- Add github issue and PR templates ([b630f54](https://github.com/ory/fosite/commit/b630f54bbd5f01891b2f3cce462819e13136d94c))
- Update JWTStrategy to take in context ([#302](https://github.com/ory/fosite/issues/302)) ([514fdbd](https://github.com/ory/fosite/commit/514fdbd20393c2175c66f3a69eb7bb849b3d5dfa))

# [0.21.5](https://github.com/ory/fosite/compare/v0.21.4...v0.21.5) (2018-08-31)

openid: Allow JWT from id_token_hint to be expired (#299)

Signed-off-by: arekkas <aeneas@ory.am>

### Unclassified

- Allow JWT from id_token_hint to be expired ([#299](https://github.com/ory/fosite/issues/299)) ([1ad9cd3](https://github.com/ory/fosite/commit/1ad9cd36069f61b2ace0fec097fe4bdc92e9f6c6))

# [0.21.4](https://github.com/ory/fosite/compare/v0.21.3...v0.21.4) (2018-08-26)

token/hmac: Add ability to rotate HMAC keys (#298)

Signed-off-by: arekkas <aeneas@ory.am>

### Unclassified

- token/hmac: Add ability to rotate HMAC keys (#298) ([2134650](https://github.com/ory/fosite/commit/213465099b72b6e5afd0e69a7916a95f65e17481)), closes [#298](https://github.com/ory/fosite/issues/298)

# [0.21.3](https://github.com/ory/fosite/compare/v0.21.2...v0.21.3) (2018-08-22)

compose: Pass ID Token configuration to strategy (#297)

Resolves an issue where expiry and issuer where not properly configurable in the strategy.

See https://github.com/ory/hydra/issues/985

Signed-off-by: arekkas <aeneas@ory.am>

### Unclassified

- Pass ID Token configuration to strategy ([#297](https://github.com/ory/fosite/issues/297)) ([a07ce27](https://github.com/ory/fosite/commit/a07ce27c814538c7d0e6228ae814482be2e96e7e)):

  Resolves an issue where expiry and issuer where not properly configurable in the strategy.

  See https://github.com/ory/hydra/issues/985

# [0.21.2](https://github.com/ory/fosite/compare/v0.21.1...v0.21.2) (2018-08-07)

openid: Validate id_token_hint only via ID claims (#296)

Signed-off-by: arekkas <aeneas@ory.am>

### Unclassified

- Validate id_token_hint only via ID claims ([#296](https://github.com/ory/fosite/issues/296)) ([0fcbfea](https://github.com/ory/fosite/commit/0fcbfea741d0f0bb2a96d5fa08a2797a109a4a33))

# [0.21.1](https://github.com/ory/fosite/compare/v0.21.0...v0.21.1) (2018-07-22)

Improve token_endpoint_auth_method error message (#294)

Signed-off-by: arekkas <aeneas@ory.am>

### Unclassified

- Improve token_endpoint_auth_method error message (#294) ([7820fb2](https://github.com/ory/fosite/commit/7820fb2e380ca737277095876c7f91b5ebee1467)), closes [#294](https://github.com/ory/fosite/issues/294)
- Gofmt ([#290](https://github.com/ory/fosite/issues/290)) ([f02884b](https://github.com/ory/fosite/commit/f02884ba0b236d81e338fd3bcd3e8ebc6d65538f)):

  Run standard gofmt command on project root.

  - go version go1.10.3 darwin/amd64

# [0.21.0](https://github.com/ory/fosite/compare/v0.20.3...v0.21.0) (2018-06-23)

Makes error messages easier to debug for end-users

### Documentation

- Fixes header image in README ([4907d60](https://github.com/ory/fosite/commit/4907d60537202e3aa04e81d87efe2c5e17c2e492))

### Unclassified

- Makes error messages easier to debug for end-users ([5688a1c](https://github.com/ory/fosite/commit/5688a1c5acbafad5eabe649ce56e06e922c36a60))
- Adds errors for request and registration parameters ([920ed71](https://github.com/ory/fosite/commit/920ed71a538f7fa5e7531660d76e076b655bf48b))
- Adds OIDC request/request_uri support ([c7abcca](https://github.com/ory/fosite/commit/c7abcca923175f85833473508684c209b1151f5a))
- Adds private_key_jwt authentication method ([baa4cf1](https://github.com/ory/fosite/commit/baa4cf15e1f30da0a52c9314730279302a15a7a4))
- Adds proper error responses to request object ([f483262](https://github.com/ory/fosite/commit/f4832621071290773fca25e8992fc283d76f390b))
- Disallow empty response_type in request ([cf2eb85](https://github.com/ory/fosite/commit/cf2eb85ed17c8d51d1c2e90c3349d4f51662a8f0))
- Do not require id_token response type for auth_code ([#288](https://github.com/ory/fosite/issues/288)) ([edc4910](https://github.com/ory/fosite/commit/edc491045155abbdbc54409889d7ccc7c3999019)):

  Before this patch, the `id_token` response type was required whenever an ID Token was requested. This patch changes that.

- Implements oidc compliant response_type validation ([f950b9e](https://github.com/ory/fosite/commit/f950b9ea63f10b7ecfe0fa47ec3716b543450dc5))
- Return unsupported_response_type in validator ([a24708e](https://github.com/ory/fosite/commit/a24708e8044268b324b1aec443a09940ae998c2f))
- Uses JWTStrategy in oauth2.DefaultStrategy ([e2d2e75](https://github.com/ory/fosite/commit/e2d2e7511931d17fd92e627c65eaabd9598b185d))
- Uses JWTStrategy interface in openid.DefaultStrategy ([517fdc5](https://github.com/ory/fosite/commit/517fdc5002ccef00a5a105b1a19bcba4c5e6839f)), closes [#252](https://github.com/ory/fosite/issues/252)

# [0.20.3](https://github.com/ory/fosite/compare/v0.20.2...v0.20.3) (2018-06-07)

Allows multipart content type as alternative to x-www-form-urlencoded (#285)

### Unclassified

- Allows multipart content type as alternative to x-www-form-urlencoded (#285) ([2edf8f8](https://github.com/ory/fosite/commit/2edf8f828b99cbabefa7f00066b49e081fab4920)), closes [#285](https://github.com/ory/fosite/issues/285)

# [0.20.2](https://github.com/ory/fosite/compare/v0.20.1...v0.20.2) (2018-05-29)

openid: Merge duplicate aud claim values (#283)

### Unclassified

- Merge duplicate aud claim values ([#283](https://github.com/ory/fosite/issues/283)) ([93618d6](https://github.com/ory/fosite/commit/93618d66a99d2756e0a4c638727b728afc62520f))

# [0.20.1](https://github.com/ory/fosite/compare/v0.20.0...v0.20.1) (2018-05-29)

Uses query instead of fragment when handling unsupported response type (#282)

### Unclassified

- Uses query instead of fragment when handling unsupported response type (#282) ([57b1471](https://github.com/ory/fosite/commit/57b14710c9aa845f2fa87322e0a3f3fa1e3e09c9)), closes [#282](https://github.com/ory/fosite/issues/282)
- Updates upgrade guide ([a958ab8](https://github.com/ory/fosite/commit/a958ab8218d13c4b0533eb38d07203f2da7ac114))

# [0.20.0](https://github.com/ory/fosite/compare/v0.19.8...v0.20.0) (2018-05-28)

oauth2: Resolves several issues related to revokation (#281)

This patch resolves several issues related to token revokation as well as duplicate authorize code usage:

- oauth2: Revoking access or refresh tokens should revoke past and future tokens too
- oauth2: Revoke access and refresh tokens when authorize code is used twice

Additionally, this patch resolves an issue where refreshing a token would not revoke previous tokens.

Closes #278
Closes #280

### Unclassified

- Resolves several issues related to revokation ([#281](https://github.com/ory/fosite/issues/281)) ([72bff7f](https://github.com/ory/fosite/commit/72bff7f33ee8c3a4a8806cc266ca7299ff1785d4)), closes [#278](https://github.com/ory/fosite/issues/278) [#280](https://github.com/ory/fosite/issues/280):

  This patch resolves several issues related to token revokation as well as duplicate authorize code usage:

  - oauth2: Revoking access or refresh tokens should revoke past and future tokens too
  - oauth2: Revoke access and refresh tokens when authorize code is used twice

  Additionally, this patch resolves an issue where refreshing a token would not revoke previous tokens.

- Sets audience to a string array ([#279](https://github.com/ory/fosite/issues/279)) ([2d58a58](https://github.com/ory/fosite/commit/2d58a585d6b53831b17bcd3ed31e67d5b2637d4a)), closes [#215](https://github.com/ory/fosite/issues/215)

# [0.19.8](https://github.com/ory/fosite/compare/v0.19.7...v0.19.8) (2018-05-24)

authorize: Fixes implicit detection in error writer (#277)

### Unclassified

- Fixes implicit detection in error writer ([#277](https://github.com/ory/fosite/issues/277)) ([608bf5f](https://github.com/ory/fosite/commit/608bf5fff7f5f7fc0dde0b3aecd03534974ba982))

# [0.19.7](https://github.com/ory/fosite/compare/v0.19.6...v0.19.7) (2018-05-24)

openid: Use claims.RequestedAt for a reference of "now" (#276)

Previously, time.Now() was used to get a reference of "now". However, this caused short max_age values to fail if, for example, the consent screen took a long time. This patch now uses the "requested_at" claim value to determine a sense of "now" which should resolve the mentioned issue.

### Unclassified

- Use claims.RequestedAt for a reference of "now" ([#276](https://github.com/ory/fosite/issues/276)) ([91e7a4c](https://github.com/ory/fosite/commit/91e7a4c236caccbea211c7790ad8194b7bd5f8a2)):

  Previously, time.Now() was used to get a reference of "now". However, this caused short max_age values to fail if, for example, the consent screen took a long time. This patch now uses the "requested_at" claim value to determine a sense of "now" which should resolve the mentioned issue.

# [0.19.6](https://github.com/ory/fosite/compare/v0.19.5...v0.19.6) (2018-05-24)

openid: Issue ID Token on implicit code flow as well

### Unclassified

- Issue ID Token on implicit code flow as well ([180c749](https://github.com/ory/fosite/commit/180c74965cb128059d63e894ba2dd04184458a33))

# [0.19.5](https://github.com/ory/fosite/compare/v0.19.4...v0.19.5) (2018-05-23)

jwt: Add JTI to counter missing nonce

### Unclassified

- Add JTI to counter missing nonce ([28822d7](https://github.com/ory/fosite/commit/28822d7b686c3a48ca9afec5291699b758c5f6cf))
- Enforce nonce on implicit/hybrid flows ([3b44eb3](https://github.com/ory/fosite/commit/3b44eb3538d4faff5fc05a74c8b9fa88ddb48202))

# [0.19.4](https://github.com/ory/fosite/compare/v0.19.3...v0.19.4) (2018-05-20)

core: Checks scopes before dispatching handlers (#272)

### Unclassified

- Checks scopes before dispatching handlers ([#272](https://github.com/ory/fosite/issues/272)) ([0f18305](https://github.com/ory/fosite/commit/0f18305e742c17db1eee6784ce3451837b5fd09a))

# [0.19.3](https://github.com/ory/fosite/compare/v0.19.2...v0.19.3) (2018-05-20)

openid: Resolves timing issues in JWT strategy (#271)

### Unclassified

- Resolves timing issues in JWT strategy ([#271](https://github.com/ory/fosite/issues/271)) ([aaec994](https://github.com/ory/fosite/commit/aaec9940e2c3fc5a696b3d174d517a6ff1490a6f))

# [0.19.2](https://github.com/ory/fosite/compare/v0.19.1...v0.19.2) (2018-05-19)

openid: Resolves timing issues by setting now to the future (#270)

### Unclassified

- Resolves timing issues by setting now to the future ([#270](https://github.com/ory/fosite/issues/270)) ([e9339d7](https://github.com/ory/fosite/commit/e9339d73eb39b15ffdb4b9a62ddc1ff1ba512530))

# [0.19.1](https://github.com/ory/fosite/compare/v0.19.0...v0.19.1) (2018-05-19)

openid: Improves validation errors and uses UTC everywhere (#269)

### Unclassified

- Improves validation errors and uses UTC everywhere ([#269](https://github.com/ory/fosite/issues/269)) ([eee3dad](https://github.com/ory/fosite/commit/eee3dad91e571a5b09217cc00caf485165f5a7d7))

# [0.19.0](https://github.com/ory/fosite/compare/v0.18.1...v0.19.0) (2018-05-17)

openid: Improves prompt, max_age and id_token_hint validation (#268)

This patch improves the OIDC prompt, max_age, and id_token_hint
validation.

### Unclassified

- Improves prompt, max_age and id_token_hint validation ([#268](https://github.com/ory/fosite/issues/268)) ([7ccad77](https://github.com/ory/fosite/commit/7ccad77095dbf8d094b2f3151634f074b0903dbc)):

  This patch improves the OIDC prompt, max_age, and id_token_hint
  validation.

# [0.18.1](https://github.com/ory/fosite/compare/v0.18.0...v0.18.1) (2018-05-01)

openid: Adds a validator used to validate OIDC parameters (#266)

The validator, for now, validates the prompt parameter of OIDC requests.

### Unclassified

- Adds a validator used to validate OIDC parameters ([#266](https://github.com/ory/fosite/issues/266)) ([91c9d19](https://github.com/ory/fosite/commit/91c9d194a88e6b395668211df60cb512eab08541)):

  The validator, for now, validates the prompt parameter of OIDC requests.

# [0.18.0](https://github.com/ory/fosite/compare/v0.17.2...v0.18.0) (2018-04-30)

oauth2: Introspection should return token type (#265)

Closes #264

This patch allows the introspection handler to return the token type (e.g. `access_token`, `refresh_token`) of the
introspected token. To achieve that, some breaking API changes have been introduced:

- `OAuth2.IntrospectToken(ctx context.Context, token string, tokenType TokenType, session Session, scope ...string) (AccessRequester, error)` is now `OAuth2.IntrospectToken(ctx context.Context, token string, tokenType TokenType, session Session, scope ...string) (TokenType, AccessRequester, error)`.
- `TokenIntrospector.IntrospectToken(ctx context.Context, token string, tokenType TokenType, accessRequest AccessRequester, scopes []string) (error)` is now `TokenIntrospector.IntrospectToken(ctx context.Context, token string, tokenType TokenType, accessRequest AccessRequester, scopes []string) (TokenType, error)`.

This patch also resolves a misconfigured json key in the `IntrospectionResponse` struct. `AccessRequester AccessRequester json:",extra"` is now properly declared as `AccessRequester AccessRequester json:"extra"`.

### Unclassified

- Introspection should return token type ([#265](https://github.com/ory/fosite/issues/265)) ([2bf9b6c](https://github.com/ory/fosite/commit/2bf9b6c4177be3050ff9ba3b82c6474e4c324c39)), closes [#264](https://github.com/ory/fosite/issues/264)

# [0.17.2](https://github.com/ory/fosite/compare/v0.17.1...v0.17.2) (2018-04-26)

core: Regression fix for request ID in refresh token flow (#262)

Signed-off-by: Beorn Facchini <beorn@lade.io>

### Unclassified

- handler/oauth2: Returns request unauthorized error on invalid password credentials (#261) ([cca6af4](https://github.com/ory/fosite/commit/cca6af4161818682edb98936cae9249db814db27)), closes [#261](https://github.com/ory/fosite/issues/261)
- Regression fix for request ID in refresh token flow ([#262](https://github.com/ory/fosite/issues/262)) ([99029e0](https://github.com/ory/fosite/commit/99029e0e1bc4b1d6dfa1ca8b85a46d79cffad6e8))

# [0.17.1](https://github.com/ory/fosite/compare/v0.17.0...v0.17.1) (2018-04-22)

core: Adds ExactScopeStrategy (#260)

The ExactScopeStrategy performs a simple string match (case sensitive)
of scopes.

### Unclassified

- Adds ExactScopeStrategy ([#260](https://github.com/ory/fosite/issues/260)) ([0fcdf33](https://github.com/ory/fosite/commit/0fcdf33fb52551e02798b4e6733110024b7d24d9)):

  The ExactScopeStrategy performs a simple string match (case sensitive)
  of scopes.

# [0.17.0](https://github.com/ory/fosite/compare/v0.16.5...v0.17.0) (2018-04-08)

core: Sanitizes request body before sending it to the storage adapter (#258)

This release resolves a security issue (reported by [platform.sh](https://www.platform.sh)) related to potential storage implementations. This library used to pass
all of the request body from both authorize and token endpoints to the storage adapters. As some of these values
are needed in consecutive requests, some storage adapters chose to drop the full body to the database. This in turn caused,
with the addition of enabling POST-body based client authentication, the client secret to be leaked.

The issue has been resolved by sanitizing the request body and only including those values truly required by their
respective handlers. This lead to two breaking changes in the API:

1. The `fosite.Requester` interface has a new method `Sanitize(allowedParameters []string) Requester` which returns
   a sanitized clone of the method receiver. If you do not use your own `fosite.Requester` implementation, this won't affect you.
2. If you use the PKCE handler, you will have to add three new methods to your storage implementation. The methods
   to be added work exactly like, for example `CreateAuthorizeCodeSession`. The method signatures are as follows:

```go
type PKCERequestStorage interface {
	GetPKCERequestSession(ctx context.Context, signature string, session fosite.Session) (fosite.Requester, error)
	CreatePKCERequestSession(ctx context.Context, signature string, requester fosite.Requester) error
	DeletePKCERequestSession(ctx context.Context, signature string) error
}
```

We encourage you to upgrade to this release and check your storage implementations and potentially remove old data.

We would like to thank [platform.sh](https://www.platform.sh) for sponsoring the development of a patch that resolves this
issue.

### Documentation

- Fixes eaxmple errors in README ([#257](https://github.com/ory/fosite/issues/257)) ([b138f59](https://github.com/ory/fosite/commit/b138f5997d535151b3541a15b8c4f7a304cea4eb))
- Updates banner in readme ([#253](https://github.com/ory/fosite/issues/253)) ([07ac5b8](https://github.com/ory/fosite/commit/07ac5b89878e07fd54edf267f23ebc7059c8bb48))

### Unclassified

- Sanitizes request body before sending it to the storage adapter ([#258](https://github.com/ory/fosite/issues/258)) ([018b5c1](https://github.com/ory/fosite/commit/018b5c12b71b0da443255f4a5cf0ac9543bbf9f7)):

  This release resolves a security issue (reported by [platform.sh](https://www.platform.sh)) related to potential storage implementations. This library used to pass
  all of the request body from both authorize and token endpoints to the storage adapters. As some of these values
  are needed in consecutive requests, some storage adapters chose to drop the full body to the database. This in turn caused,
  with the addition of enabling POST-body based client authentication, the client secret to be leaked.

  The issue has been resolved by sanitizing the request body and only including those values truly required by their
  respective handlers. This lead to two breaking changes in the API:

  1. The `fosite.Requester` interface has a new method `Sanitize(allowedParameters []string) Requester` which returns
     a sanitized clone of the method receiver. If you do not use your own `fosite.Requester` implementation, this won't affect you.
  2. If you use the PKCE handler, you will have to add three new methods to your storage implementation. The methods
     to be added work exactly like, for example `CreateAuthorizeCodeSession`. The method signatures are as follows:

  ```go
  type PKCERequestStorage interface {
  	GetPKCERequestSession(ctx context.Context, signature string, session fosite.Session) (fosite.Requester, error)
  	CreatePKCERequestSession(ctx context.Context, signature string, requester fosite.Requester) error
  	DeletePKCERequestSession(ctx context.Context, signature string) error
  }
  ```

  We encourage you to upgrade to this release and check your storage implementations and potentially remove old data.

  We would like to thank [platform.sh](https://www.platform.sh) for sponsoring the development of a patch that resolves this
  issue.

# [0.16.5](https://github.com/ory/fosite/compare/v0.16.4...v0.16.5) (2018-03-17)

introspection: Improves debug messages (#254)

### Documentation

- Resolves minor code documentation misspellings ([#248](https://github.com/ory/fosite/issues/248)) ([c580d79](https://github.com/ory/fosite/commit/c580d79aaa54f2aec179df400a3365ca711ead66))
- Resolves minor spelling mistakes ([#250](https://github.com/ory/fosite/issues/250)) ([7fbd246](https://github.com/ory/fosite/commit/7fbd2468dfb83cf7288643958db9890af5ffd3d1))
- Updates chat badge to discord ([b6380be](https://github.com/ory/fosite/commit/b6380be3365fc9703135f6ef3ee747d60d835915))

### Unclassified

- docs : Fixes typo in README (#249) ([d05fadf](https://github.com/ory/fosite/commit/d05fadfa7c4fa88ec58175fef146c7cc9c6c120c)), closes [#249](https://github.com/ory/fosite/issues/249)
- Adds email to license notice ([77fa262](https://github.com/ory/fosite/commit/77fa262093d783bc3f0e302ebddd1a2da3f2581d))
- Improves debug messages ([#254](https://github.com/ory/fosite/issues/254)) ([338399b](https://github.com/ory/fosite/commit/338399becb5114f84e6dc7166a95f6d036a6b748))
- Updates license header ([85bdbcb](https://github.com/ory/fosite/commit/85bdbcb4c34c646c7eae56c0a1dc41dc1f75b470))
- Updates license notice ([917401c](https://github.com/ory/fosite/commit/917401cdf0b891afa9a3aa65edb2539ff0f0a5ba))
- Updates years in license headers ([77df218](https://github.com/ory/fosite/commit/77df218b30566ab7cd513b723a7e44f9f6afbe7e))
- Updates years in license headers ([d8458ab](https://github.com/ory/fosite/commit/d8458abe997f70c743a7e2fa3cc27c2cb1d38c9e))

# [0.16.4](https://github.com/ory/fosite/compare/v0.16.3...v0.16.4) (2018-02-07)

handler: Adds PKCE implementation for none and S256 (#246)

This patch adds support for PKCE (https://tools.ietf.org/html/rfc7636) which is used by native apps (mobile) and prevents eavesdropping attacks against authorization codes.

PKCE is enabled by default but not enforced. Challenge method plain is disabled by default. Both settings can be changed using `compose.Config.EnforcePKCE` and `compose.config.EnablePKCEPlainChallengeMethod`.

Closes #213

### Unclassified

- Adds PKCE implementation for none and S256 ([#246](https://github.com/ory/fosite/issues/246)) ([4512853](https://github.com/ory/fosite/commit/45128532dc4bbb40a56bf6250a58f9c5d57a9c7a)), closes [#213](https://github.com/ory/fosite/issues/213):

  This patch adds support for PKCE (https://tools.ietf.org/html/rfc7636) which is used by native apps (mobile) and prevents eavesdropping attacks against authorization codes.

  PKCE is enabled by default but not enforced. Challenge method plain is disabled by default. Both settings can be changed using `compose.Config.EnforcePKCE` and `compose.config.EnablePKCEPlainChallengeMethod`.

# [0.16.3](https://github.com/ory/fosite/compare/v0.16.2...v0.16.3) (2018-02-07)

introspection: Adds missing http header to response writer (#247)

The introspection response writer was missing `application/json`
in header `Content-Type`. This patch fixes that.

Closes #209

### Unclassified

- Adds missing http header to response writer ([#247](https://github.com/ory/fosite/issues/247)) ([f345ec1](https://github.com/ory/fosite/commit/f345ec1413aa0fc2ba4588a482e469fa19cc08aa)), closes [#209](https://github.com/ory/fosite/issues/209):

  The introspection response writer was missing `application/json`
  in header `Content-Type`. This patch fixes that.

# [0.16.2](https://github.com/ory/fosite/compare/v0.16.1...v0.16.2) (2018-01-25)

introspection: Decodes of Basic Authorization username/password (#245)

Signed-off-by: Dmitry Dolbik <dolbik@gmail.com>

### Unclassified

- Decodes of Basic Authorization username/password ([#245](https://github.com/ory/fosite/issues/245)) ([b94312e](https://github.com/ory/fosite/commit/b94312e25f011b54894da69256416271c23b5d14))

# [0.16.1](https://github.com/ory/fosite/compare/v0.16.0...v0.16.1) (2017-12-23)

compose: Makes SendDebugMessages first class citizen (#243)

### Unclassified

- Makes SendDebugMessages first class citizen ([#243](https://github.com/ory/fosite/issues/243)) ([1ef3041](https://github.com/ory/fosite/commit/1ef3041c4da40d27ea25d56710e59d5f9352df5f))

# [0.16.0](https://github.com/ory/fosite/compare/v0.15.6...v0.16.0) (2017-12-23)

Adds ability to forward hints and debug messages to clients (#242)

### Unclassified

- Adds ability to forward hints and debug messages to clients (#242) ([7216c4f](https://github.com/ory/fosite/commit/7216c4f2711c79cf3d8a2c75ad7da4f54103988f)), closes [#242](https://github.com/ory/fosite/issues/242)

# [0.15.6](https://github.com/ory/fosite/compare/v0.15.5...v0.15.6) (2017-12-21)

handler/oauth2: Adds offline_access alias for refresh flow

### Unclassified

- handler/oauth2: Adds offline_access alias for refresh flow ([2aa8e70](https://github.com/ory/fosite/commit/2aa8e70bb88aa6bafde8d4ea949c5d514c6f568e))

# [0.15.5](https://github.com/ory/fosite/compare/v0.15.4...v0.15.5) (2017-12-17)

Returns the correct error on duplicate auth code use

### Unclassified

- Returns the correct error on duplicate auth code use ([95d5f58](https://github.com/ory/fosite/commit/95d5f580c939eea0e6e93cdb4bae4cdbf5082869))

# [0.15.4](https://github.com/ory/fosite/compare/v0.15.3...v0.15.4) (2017-12-17)

Improves http error codes

### Unclassified

- Improves http error codes ([6831f75](https://github.com/ory/fosite/commit/6831f7543000b3704879e52d8c9a4555653b4bd5))

# [0.15.3](https://github.com/ory/fosite/compare/v0.15.2...v0.15.3) (2017-12-17)

Resolves overriding auth_time with wrong value

### Unclassified

- Resolves overriding auth_time with wrong value ([c85b32d](https://github.com/ory/fosite/commit/c85b32d355a183dac3e46e50aac8b2c344cbd2d7))

# [0.15.2](https://github.com/ory/fosite/compare/v0.15.1...v0.15.2) (2017-12-10)

Adds ability to catch non-conform OIDC authorizations

Fosite is now capable of detecting authorization flows that
are not conformant with the OpenID Connect spec.

### Unclassified

- Adds ability to catch non-conform OIDC authorizations ([97fbeb3](https://github.com/ory/fosite/commit/97fbeb333e353d5d7d7d2368f51899262338ce62)):

  Fosite is now capable of detecting authorization flows that
  are not conformant with the OpenID Connect spec.

- Forces use of UTC time zone everywhere ([4c7e4e5](https://github.com/ory/fosite/commit/4c7e4e5512061e9add22cc246882c78d2b06599c))

# [0.15.1](https://github.com/ory/fosite/compare/v0.15.0...v0.15.1) (2017-12-10)

token/jwt: Adds ability to specify acr value natively in id token payload

### Unclassified

- token/jwt: Adds ability to specify acr value natively in id token payload ([b87ca49](https://github.com/ory/fosite/commit/b87ca49b9418b99f492077f8ba78bf00e6c29180))

# [0.15.0](https://github.com/ory/fosite/compare/v0.14.2...v0.15.0) (2017-12-09)

Upgrades history.md

### Documentation

- Updates history.md ([9fc25a8](https://github.com/ory/fosite/commit/9fc25a86c4d8609aafa382e5eab32d3d087ec9d8))

### Unclassified

- Upgrades history.md ([87c37c3](https://github.com/ory/fosite/commit/87c37c3d6929b1edd2ab52a28d51ed1890628f51))
- Improves test coverage report by removing internal package from it ([831f56a](https://github.com/ory/fosite/commit/831f56a9e6774b1e80c13cd301583edea6378245))
- Resolves test issues and reverts auth code revokation patch ([59fc47b](https://github.com/ory/fosite/commit/59fc47bbeb8093ab3652149ef6789a4e1564e1d8))
- Improves error debug messages across the project ([7ec8d19](https://github.com/ory/fosite/commit/7ec8d19815d10913ef8cfd8ced9b9794f578dbf4))
- handler/oauth2: Adds token revokation on authorize code reuse ([2341dec](https://github.com/ory/fosite/commit/2341dec8febeda9da535dc898c7d19aa3ecc8c00))
- handler/oauth2: Improves authorization code error handling ([d6e0fbd](https://github.com/ory/fosite/commit/d6e0fbd9bdde624fa2e9feada3dec5b4266c4b9e))
- Allows client credentials in POST body and solves public client auth ([392c191](https://github.com/ory/fosite/commit/392c191bc1859ca57e3d0cf4d2b996d5ab382530)), closes [#231](https://github.com/ory/fosite/issues/231) [#217](https://github.com/ory/fosite/issues/217)
- Updates mocks and mock generation ([1f9d07d](https://github.com/ory/fosite/commit/1f9d07d15e8f70986ed12cfb3ac9fac4a6e7e278))

# [0.14.2](https://github.com/ory/fosite/compare/v0.14.1...v0.14.2) (2017-12-06)

Makes use of rfcerr in access error endpoint writer explicit

### Unclassified

- Makes use of rfcerr in access error endpoint writer explicit ([701d850](https://github.com/ory/fosite/commit/701d85072d1ea5c35c7d05acf19bccdef626ba3c))

# [0.14.1](https://github.com/ory/fosite/compare/v0.14.0...v0.14.1) (2017-12-06)

Exports ErrorToRFC6749Error again (#228)

### Unclassified

- Exports ErrorToRFC6749Error again (#228) ([8d35b66](https://github.com/ory/fosite/commit/8d35b668079db8642ede3b1d345d74692926515f)), closes [#228](https://github.com/ory/fosite/issues/228)

# [0.14.0](https://github.com/ory/fosite/compare/v0.13.1...v0.14.0) (2017-12-06)

Simplifies error contexts (#227)

Simplifies how errors are instantiated. Errors now contain all necessary information without relying on `fosite.ErrorToRFC6749Error` any more. `fosite.ErrorToRFC6749Error` is now an internal method and was renamed to `fosite.errorToRFC6749Error`.

### Unclassified

- Simplifies error contexts (#227) ([8961d86](https://github.com/ory/fosite/commit/8961d861814862f9432f0608bcd14dfbcd4ec979)), closes [#227](https://github.com/ory/fosite/issues/227):

  Simplifies how errors are instantiated. Errors now contain all necessary information without relying on `fosite.ErrorToRFC6749Error` any more. `fosite.ErrorToRFC6749Error` is now an internal method and was renamed to `fosite.errorToRFC6749Error`.

# [0.13.1](https://github.com/ory/fosite/compare/v0.13.0...v0.13.1) (2017-12-04)

handler/oauth2: Client IDs in revokation requests must match now (#226)

Closes #225

### Unclassified

- handler/oauth2: Client IDs in revokation requests must match now (#226) ([83136a3](https://github.com/ory/fosite/commit/83136a3ed5ed99b3a525f0ad87d693eadf273e8a)), closes [#226](https://github.com/ory/fosite/issues/226) [#225](https://github.com/ory/fosite/issues/225)
- Add license header to all source files (#222) ([dd9398e](https://github.com/ory/fosite/commit/dd9398ea0553b07d63022af50ee2090d1616c5a9)), closes [#222](https://github.com/ory/fosite/issues/222) [#221](https://github.com/ory/fosite/issues/221)
- Update go version ([#220](https://github.com/ory/fosite/issues/220)) ([ff751ee](https://github.com/ory/fosite/commit/ff751ee3691f79886ccfc6afa3936c2c3b506a9e))

# [0.13.0](https://github.com/ory/fosite/compare/v0.12.0...v0.13.0) (2017-10-25)

vendor: replace glide with dep

### Unclassified

- Replace glide with dep ([ec43e3a](https://github.com/ory/fosite/commit/ec43e3a05da49d45ebe8a98b28b14f8817c507f4))

# [0.12.0](https://github.com/ory/fosite/compare/v0.11.4...v0.12.0) (2017-10-25)

scripts: fix goimports import path

### Unclassified

- token/hmac: replace custom logic with copypasta ([b4b9be5](https://github.com/ory/fosite/commit/b4b9be5640c9d814b35f54b2c8621137364209ca))
- Add 0.12.0 to TOC ([a2e3a47](https://github.com/ory/fosite/commit/a2e3a474b2439e4ad68a641152639f7921e610a6))
- Add format helper scripts ([92c73ae](https://github.com/ory/fosite/commit/92c73aee93b5d1fe2acf3395b495caf912453368))
- Add goimports to install section ([4f5df70](https://github.com/ory/fosite/commit/4f5df700e3c220f3aa5f7eb79a4b4f19d2f4576e))
- Fix goimports import path ([65743b4](https://github.com/ory/fosite/commit/65743b40c69ccc76f07fd3eb4c45837d3b4a1505))
- Format files with goimports ([c87defe](https://github.com/ory/fosite/commit/c87defe18676b36d880fa834c10e2cbd5464e061))
- Replace nil checks with Error/NoError ([7fe1f94](https://github.com/ory/fosite/commit/7fe1f946af7b4921da008f245da84b85ea3f26d0))
- Update to go 1.9 ([c17222c](https://github.com/ory/fosite/commit/c17222c854198a7a388a2656a710bf13a5c3c3b9))
- Use go-acc and test format ([47fd477](https://github.com/ory/fosite/commit/47fd477814c7826a9e9e89a02c248cfbad6b5a7a))

# [0.11.4](https://github.com/ory/fosite/compare/v0.11.3...v0.11.4) (2017-10-10)

handler/oauth2: set expiration time before the access token is generated (#216)

Signed-off-by: Nikita Vorobey <nikita@vorobey.by>

### Documentation

- Update banner ([d6cf027](https://github.com/ory/fosite/commit/d6cf027401e828c8e608b042615f982acdf6d915))

### Unclassified

- handler/oauth2: set expiration time before the access token is generated (#216) ([0911eb0](https://github.com/ory/fosite/commit/0911eb0d643d77105e0126bf2303bdfd7190ccd3)), closes [#216](https://github.com/ory/fosite/issues/216)

# [0.11.3](https://github.com/ory/fosite/compare/v0.11.2...v0.11.3) (2017-08-21)

oauth2/ropc: Set expires at for password credentials flow (#210)

Signed-off-by: Beorn Facchini <beornf@gmail.com>

### Documentation

- Fixes documentation oauth2 variable and updates old method ([#205](https://github.com/ory/fosite/issues/205)) ([fa50c80](https://github.com/ory/fosite/commit/fa50c80d36bbc8dda2633b59617689d8ef21042c)):

  It seems that the documentation was declaring as OAuth2Provider the variable `oauth2Provider` whereas it used a non-declared variable `oauth2`. I renamed `oauth2` into the variable declared `oauth2Provider`.

  Furthermore, on line 333, the IntrospectToken method was called without the TokenType argument. I added the fosite.AccessToken type.

- Update docs on scope strategy ([68119ca](https://github.com/ory/fosite/commit/68119ca5e282c356284a6dc7a2edb2b632d57a47))

### Unclassified

- oauth2/ropc: Set expires at for password credentials flow (#210) ([461b38f](https://github.com/ory/fosite/commit/461b38fd07e47dad709667f024e98a71bfd3792b)), closes [#210](https://github.com/ory/fosite/issues/210)
- oauth2/introspection: configure core validator with access only option (#208) ([80cae74](https://github.com/ory/fosite/commit/80cae74590bfdf7d3f9439073a4a5aac21d7fd45)), closes [#208](https://github.com/ory/fosite/issues/208)
- Add more test cases ([c45a37d](https://github.com/ory/fosite/commit/c45a37d3bb9e3e79d16323f42d76ef96b624dbd0))

# [0.11.2](https://github.com/ory/fosite/compare/v0.11.1...v0.11.2) (2017-07-09)

scope: resolve haystack needle mixup - closes #201

### Unclassified

- Resolve haystack needle mixup - closes [#201](https://github.com/ory/fosite/issues/201) ([2c7cdff](https://github.com/ory/fosite/commit/2c7cdff9d2e677f5f892d6107a3c0b8b9ce61632))

# [0.11.1](https://github.com/ory/fosite/compare/v0.11.0...v0.11.1) (2017-07-09)

token/jwt: add claims tests

### Unclassified

- token/jwt: add claims tests ([c55d679](https://github.com/ory/fosite/commit/c55d67903fdc5b2f4b200b663d4f1a0cb1d21dca))
- handler/openid: only refresh id token with id_token response type ([dd2463a](https://github.com/ory/fosite/commit/dd2463a1a262600096f040867dcabe2a28e1a56c)), closes [#199](https://github.com/ory/fosite/issues/199)
- Add tests for nil sessions ([d67d52d](https://github.com/ory/fosite/commit/d67d52df200dfc72c9eb79e38ae6e91a1fb701f4))

# [0.11.0](https://github.com/ory/fosite/compare/v0.10.0...v0.11.0) (2017-07-09)

handler/oauth2: update docs

### Unclassified

- handler/oauth2: update docs ([63f329b](https://github.com/ory/fosite/commit/63f329b104c36dcbe2ee2f2a5562c6422f36224b))
- handler/oauth2: remove code validity check from test ([664d1a6](https://github.com/ory/fosite/commit/664d1a6c0177abfb4d8f780f28ecd69cb2d44d87))
- handler/oauth2: first retrieve, then validate ([ab72cba](https://github.com/ory/fosite/commit/ab72cba1799accc7b50990908139fa762eb2efc1))
- handler/oauth2: set requested at date in auth code test ([edd4084](https://github.com/ory/fosite/commit/edd4084b43ed88135fb60a4581283d8abaf92384))
- handler/oauth2: resolve travis time mismatch ([ec6534c](https://github.com/ory/fosite/commit/ec6534cfebf24d716aba28dee43e6ec268c0918b))
- handler/oauth2: simplify storage interface ([361b368](https://github.com/ory/fosite/commit/361b3683552bcadf62d1d1c42baf6d5cc1ca1409)), closes [#194](https://github.com/ory/fosite/issues/194)
- handler/oauth2: use hmac strategy for jwt refresh tokens (#190) ([56c88c0](https://github.com/ory/fosite/commit/56c88c04d4819aec08cb068a5fb7697dbaeb3288)), closes [#190](https://github.com/ory/fosite/issues/190) [#180](https://github.com/ory/fosite/issues/180)
- handler/openid: refresh token handler for oidc (#193) ([04888c5](https://github.com/ory/fosite/commit/04888c5448382612a55fb0c57ccf2c0e3d841c2c)), closes [#193](https://github.com/ory/fosite/issues/193) [#181](https://github.com/ory/fosite/issues/181)
- Gofmt ([7a998fe](https://github.com/ory/fosite/commit/7a998fece7ea2fd63ad7943266e67954ab81aaf6))
- Implement new wildcard strategy - closes [#188](https://github.com/ory/fosite/issues/188) ([e03e99e](https://github.com/ory/fosite/commit/e03e99e653454ab7cc997aacd162374bdbf38c75))
- Revoke access tokens when refreshing ([bb74955](https://github.com/ory/fosite/commit/bb74955ead77dbadf2f7b99ec3bff9b27f2a4388)), closes [#167](https://github.com/ory/fosite/issues/167)
- Run goimports ([35941c2](https://github.com/ory/fosite/commit/35941c2f3ed0436019429d9657d9dab59cae93e1))
- Use deepcopy not gob encoding - closes [#191](https://github.com/ory/fosite/issues/191) ([823db5b](https://github.com/ory/fosite/commit/823db5b65cd7c0c356b211c920ca06ec10cfa8b6))

# [0.10.0](https://github.com/ory/fosite/compare/v0.9.7...v0.10.0) (2017-07-06)

oauth2/introspector: remove auth code, refresh scopes (#187)

Removes authorize code introspection in the HMAC-based strategy and now checks scopes of refresh tokens as well.

### Unclassified

- oauth2/introspector: remove auth code, refresh scopes (#187) ([ef8f175](https://github.com/ory/fosite/commit/ef8f1757f0c26317fd7dbb46f66fde7516a3b4bb)), closes [#187](https://github.com/ory/fosite/issues/187):

  Removes authorize code introspection in the HMAC-based strategy and now checks scopes of refresh tokens as well.

- Separate test dependencies ([#186](https://github.com/ory/fosite/issues/186)) ([71451f0](https://github.com/ory/fosite/commit/71451f05fa2b572c4467a9bca26ec3d018a74cd3)):

  - vendor: Move testify to testImport
  - test: Move Assert/Require helpers to \_test pkg

# [0.9.7](https://github.com/ory/fosite/compare/v0.9.6...v0.9.7) (2017-06-28)

handler/openid: remove forced nonce (#185)

Signed-off-by: Wyatt Anderson <wanderson@gmail.com>

### Unclassified

- handler/openid: remove forced nonce (#185) ([6c91a21](https://github.com/ory/fosite/commit/6c91a21b540c534c9a2330922e357e24c7d5fda9)), closes [#185](https://github.com/ory/fosite/issues/185)

# [0.9.6](https://github.com/ory/fosite/compare/v0.9.5...v0.9.6) (2017-06-21)

oauth2: basic auth should decode client id and secret

closes #182

### Documentation

- Update test command in README and CONTRIBUTING ([#183](https://github.com/ory/fosite/issues/183)) ([c1ab029](https://github.com/ory/fosite/commit/c1ab029745520914fae525f150e91dfe7ae76142))

### Unclassified

- Basic auth should decode client id and secret ([92b75d9](https://github.com/ory/fosite/commit/92b75d93070fdb96f0ec9975dc24b69243d8f894)), closes [#182](https://github.com/ory/fosite/issues/182)

# [0.9.5](https://github.com/ory/fosite/compare/v0.9.4...v0.9.5) (2017-06-08)

handler/oauth2: grant scopes before the access token is generated (#177)

Signed-off-by: Nikita Vorobey <nikita@vorobey.by>

### Unclassified

- handler/oauth2: grant scopes before the access token is generated (#177) ([3497260](https://github.com/ory/fosite/commit/349726028d42f3c60aeefc67aef06f9f907ccf94)), closes [#177](https://github.com/ory/fosite/issues/177)

# [0.9.4](https://github.com/ory/fosite/compare/v0.9.3...v0.9.4) (2017-06-05)

introspection: return with active set false on token error (#176)

### Unclassified

- Return with active set false on token error ([#176](https://github.com/ory/fosite/issues/176)) ([82944aa](https://github.com/ory/fosite/commit/82944aaa42ddc9c718ee072d5a11635ec982394d))

# [0.9.3](https://github.com/ory/fosite/compare/v0.9.2...v0.9.3) (2017-06-05)

vendor: remove unnecessary go-jose import (#175)

### Unclassified

- Remove unnecessary go-jose import ([#175](https://github.com/ory/fosite/issues/175)) ([d26aa4a](https://github.com/ory/fosite/commit/d26aa4a76fda898677f333c38242a9049e448e1a))

# [0.9.2](https://github.com/ory/fosite/compare/v0.9.1...v0.9.2) (2017-06-05)

Resolve issues with error handling (#174)

- errors: do not convert errors compliant with rfcerrors

- handler/oauth2: improve redirect message for insecure http

### Unclassified

- Resolve issues with error handling (#174) ([9abdfd0](https://github.com/ory/fosite/commit/9abdfd04261f472f34c9d6a545ccaa2d491c4f06)), closes [#174](https://github.com/ory/fosite/issues/174):

  - errors: do not convert errors compliant with rfcerrors

  - handler/oauth2: improve redirect message for insecure http

# [0.9.1](https://github.com/ory/fosite/compare/v0.9.0...v0.9.1) (2017-06-04)

vendor: clean up dependencies (#173)

- vendor: remove stray github.com/Sirupsen/logrus
- vendor: remove common lib

### Unclassified

- Clean up dependencies ([#173](https://github.com/ory/fosite/issues/173)) ([524d3b6](https://github.com/ory/fosite/commit/524d3b6fb51e81330608f727c63dbf41980de7ae)):

  - vendor: remove stray github.com/Sirupsen/logrus
  - vendor: remove common lib

# [0.9.0](https://github.com/ory/fosite/compare/v0.8.0...v0.9.0) (2017-06-03)

docs: add 0.9.0 release note

### Documentation

- Add 0.9.0 release note ([852cf82](https://github.com/ory/fosite/commit/852cf82344c2d78863508eaa0fca32f468cd7fab))

### Unclassified

- Enable fosite composing with custom hashers. ([#170](https://github.com/ory/fosite/issues/170)) ([d70d882](https://github.com/ory/fosite/commit/d70d882d0b125e386e52cd1aee3712d48538fd66))
- Removed implicit storage as its never used - closes [#165](https://github.com/ory/fosite/issues/165) ([#171](https://github.com/ory/fosite/issues/171)) ([fe74027](https://github.com/ory/fosite/commit/fe74027ee70292a72fe453095603cca060ff6290))

# [0.8.0](https://github.com/ory/fosite/compare/v0.7.0...v0.8.0) (2017-05-18)

docs: add notes for breaking changes that come with 0.8.0

### Documentation

- Add notes for breaking changes that come with 0.8.0 ([d5fafb8](https://github.com/ory/fosite/commit/d5fafb87b04ddf2ced6b58a063eac71892bcd5c9))

### Unclassified

- Added context to GetClient storage interface ([#162](https://github.com/ory/fosite/issues/162)) ([974585d](https://github.com/ory/fosite/commit/974585d4f809f96c8bf9ee3f0f1540bf9478b8a9)), closes [#161](https://github.com/ory/fosite/issues/161)
- Removed \*http.Request from interfaces that access request objects ([786b971](https://github.com/ory/fosite/commit/786b971ca1d36a8f0bd0a5c0bfa798802d5c0c26)):

  - removed the requirement to \*http.Request for endpoints and response object, they are resolvable trough the request.GetRequestForm

  - updated readme to reflect changes to implementation

  - run goimports on internal dir
    added goimports command to generate-mocks.sh to force first run after generating the mock files

- Set authorize code expire time before persist ([#166](https://github.com/ory/fosite/issues/166)) ([305a74f](https://github.com/ory/fosite/commit/305a74fe20649bde7150509ec072a43b958e0ee9))
- Set expiry date on implicit access tokens ([#164](https://github.com/ory/fosite/issues/164)) ([0785b07](https://github.com/ory/fosite/commit/0785b072dba9a9cf65bc8b7304af4e7691f96a96))

# [0.7.0](https://github.com/ory/fosite/compare/v0.6.19...v0.7.0) (2017-05-03)

vendor: glide update

### Documentation

- Add breaking changes note ([7d726e1](https://github.com/ory/fosite/commit/7d726e13800667a32372bb7f97a7f652c7eb9f3e))

### Unclassified

- Glide update ([575dd79](https://github.com/ory/fosite/commit/575dd791f9f11cd8e5471178b1ec3a7638653cae))
- Goimports ([1cb7e26](https://github.com/ory/fosite/commit/1cb7e26e164c1f11b7cb6ab64191d680d19e7ca0))
- Move to new org ([bd13085](https://github.com/ory/fosite/commit/bd1308540c519a09d4228048d3d9a028d363a7bd))
- Replace golang.org/x/net/context with context ([6b1d931](https://github.com/ory/fosite/commit/6b1d93124be24d4b2949060a4c3428c220667738))

# [0.6.19](https://github.com/ory/fosite/compare/v0.6.18...v0.6.19) (2017-05-03)

access: revert regression issue introduced by #150

### Unclassified

- Revert regression issue introduced by [#150](https://github.com/ory/fosite/issues/150) ([6f13d58](https://github.com/ory/fosite/commit/6f13d58533573ec847dca6e5cfa1d4338aef95b1))
- Revert regression issue introduced by [#150](https://github.com/ory/fosite/issues/150) ([6bb4135](https://github.com/ory/fosite/commit/6bb4135523c4e2fcf7b3a0630e233ccb7a806fc8))

# [0.6.18](https://github.com/ory/fosite/compare/v0.6.17...v0.6.18) (2017-04-14)

oauth2: basic auth should www-url-decode client id and secret - closes #150

### Unclassified

- handler/oauth2: removes RevokeHandler from JWT introspector (#155) ([344dbef](https://github.com/ory/fosite/commit/344dbeff15cfce9990c0ccfd687a0c44f6a81569)), closes [#155](https://github.com/ory/fosite/issues/155):

  - Removes RevokeHandler from JWT Introspector

  RevokeHandler has been removed because it conflicts with Stateless JWT
  accesstokens and revocable hmac refresh tokens. The readme has been
  updated to warn users about possible misconfiguration.

  - Moves text back to correct section

- Allow localhost subdomains such as blog.localhost:1234 ([5e1c890](https://github.com/ory/fosite/commit/5e1c890fd144ce1ec12ee26d7ebfe02862af067e))
- Basic auth should www-url-decode client id and secret - closes [#150](https://github.com/ory/fosite/issues/150) ([ad395bf](https://github.com/ory/fosite/commit/ad395bf323137e30ce12d40646a9229a42695863))
- Get the token from the access_token query parameter ([#156](https://github.com/ory/fosite/issues/156)) ([9edac04](https://github.com/ory/fosite/commit/9edac0441f4f9c8400e0cbd9cd637e9d2bfcae05))

# [0.6.17](https://github.com/ory/fosite/compare/v0.6.15...v0.6.17) (2017-02-24)

readme: update badges to ory

### Unclassified

- revert unintentional change ([14a18a7](https://github.com/ory/fosite/commit/14a18a714c419b31d4bf1341e1017159bc17540f))
- make stateless validator return an error on revocation ([f8f7978](https://github.com/ory/fosite/commit/f8f797869eaa1895791ed1bba3b0f3c3a06a03ca))
- dont client id for aud ([a39200b](https://github.com/ory/fosite/commit/a39200b3eb08b77d0181586454e5d7348d519aa5))
- handler/oauth2: allow stateless introspection of jwt access tokens ([c2d2ac2](https://github.com/ory/fosite/commit/c2d2ac258ecb1378493c0d60add2967e510fbc6b))
- Redirect uris should ignore cases during matching - closes [#144](https://github.com/ory/fosite/issues/144) ([4b88774](https://github.com/ory/fosite/commit/4b887746fde977a0f5cf8fbbe06c90577f416fca))
- Update badges to ory ([9b33931](https://github.com/ory/fosite/commit/9b33931ee14ae0768ea46a423d569330a85b482e))

# [0.6.15](https://github.com/ory/fosite/compare/v0.6.14...v0.6.15) (2017-02-11)

errors: fixed typo in acccess_error

### Unclassified

- Fixed typo in acccess_error ([08b2242](https://github.com/ory/fosite/commit/08b2242b66a8d430084c6aada57018f8c2dabea6))

# [0.6.14](https://github.com/ory/fosite/compare/v0.6.13...v0.6.14) (2017-01-08)

allow public clients to revoke tokens with just an ID

This functionality is described in the OAuth2 spec here: https://tools.ietf.org/html/rfc7009#section-5

### Unclassified

- allow public clients to revoke tokens with just an ID ([7b94f47](https://github.com/ory/fosite/commit/7b94f470bede7cf5e94d11e05aa3364d0db75fe2)), closes [/tools.ietf.org/html/rfc7009#section-5](https://github.com//tools.ietf.org/html/rfc7009/issues/section-5)
- Conform to RFC 6749 ([c404554](https://github.com/ory/fosite/commit/c4045541ae19c88634d79818a0060d71c9ef07ec)), closes [/tools.ietf.org/html/rfc6749#section-5](https://github.com//tools.ietf.org/html/rfc6749/issues/section-5):

  Section 5.2 specifies the parameters for access error responses;
  the "error" and "error_description" parameters are misnamed.

# [0.6.13](https://github.com/ory/fosite/compare/v0.6.12...v0.6.13) (2017-01-08)

request: fix SetRequestedScopes (#139)

Signed-off-by: Peter Schultz <peter.schultz@classmarkets.com>

### Unclassified

- Fix SetRequestedScopes ([#139](https://github.com/ory/fosite/issues/139)) ([d02c427](https://github.com/ory/fosite/commit/d02c427a76d5d8ef2f099bae79b7af69be3f643a))

# [0.6.12](https://github.com/ory/fosite/compare/v0.6.11...v0.6.12) (2017-01-02)

authorize: allow custom redirect url schemas

### Unclassified

- Allow custom redirect url schemas ([c740b70](https://github.com/ory/fosite/commit/c740b703399e7a1479dac9f261baec4b341f6cff))
- Properly wrap errors ([e054b6e](https://github.com/ory/fosite/commit/e054b6e04a9253e3d1d333064998045b3ab649fe))

# [0.6.11](https://github.com/ory/fosite/compare/v0.6.10...v0.6.11) (2017-01-02)

openid: c_hash / at_hash should use url-safe base64 encoding

### Unclassified

- C_hash / at_hash should use url-safe base64 encoding ([33d4414](https://github.com/ory/fosite/commit/33d44146ef17f9c176a2a74e7ee77eaae98ee5c1))

# [0.6.10](https://github.com/ory/fosite/compare/v0.6.9...v0.6.10) (2016-12-29)

openid: c_hash / at_hash should be string not byte slice

### Unclassified

- C_hash / at_hash should be string not byte slice ([b489cc9](https://github.com/ory/fosite/commit/b489cc95b87d74785c5e9b8ea5eb48e975559f63))

# [0.6.9](https://github.com/ory/fosite/compare/v0.6.8...v0.6.9) (2016-12-29)

oauth2/implicit: fix redirect url on error
Signed-off-by: Nikita Vorobey <nikita@vorobey.by>

### Documentation

- Fix missing protocol in link in readme ([#132](https://github.com/ory/fosite/issues/132)) ([37ef374](https://github.com/ory/fosite/commit/37ef374aec940d6b9fdcc33800c09ba08b830f39))

### Unclassified

- oauth2/implicit: fix redirect url on error ([435288c](https://github.com/ory/fosite/commit/435288ccdee2aed2447a5a0babf885dbfeae6b55))

# [0.6.8](https://github.com/ory/fosite/compare/v0.6.7...v0.6.8) (2016-12-20)

lint: gofmt -w -s .

### Unclassified

- Add id_token + code flow ([3f347e3](https://github.com/ory/fosite/commit/3f347e35b603fdde805a8b7a4fdaeff6bcddaa02))
- Fix typos ([#130](https://github.com/ory/fosite/issues/130)) ([e6b410d](https://github.com/ory/fosite/commit/e6b410d519a0944cd52ffde656f7b21c4682b5a6))
- Gofmt -w -s . ([95caa96](https://github.com/ory/fosite/commit/95caa96835a1254ba3f8f4a21e635fe6da34f0fe))

# [0.6.7](https://github.com/ory/fosite/compare/v0.6.6...v0.6.7) (2016-12-06)

access: response expires in should be int, not string

### Unclassified

- Response expires in should be int, not string ([a2080a3](https://github.com/ory/fosite/commit/a2080a30c04abf6a9b3f7dee63026cb5816f8bbd))

# [0.6.6](https://github.com/ory/fosite/compare/v0.6.5...v0.6.6) (2016-12-06)

errors: add inactive token error

### Unclassified

- Add content type to error response ([75aad53](https://github.com/ory/fosite/commit/75aad53be3dfda8a02a47bd8f574dc23914b4b65))
- Add inactive token error ([0151f1e](https://github.com/ory/fosite/commit/0151f1e17dda1c81185d00b388c83b25b7c5f72c))
- Resolve broken test ([51ab7bb](https://github.com/ory/fosite/commit/51ab7bb960640bcd8722e2731af72c6c26e3bacd))

# [0.6.5](https://github.com/ory/fosite/compare/v0.6.4...v0.6.5) (2016-12-04)

introspection: always return the error

### Unclassified

- Always return the error ([366b4c1](https://github.com/ory/fosite/commit/366b4c1a06369b2cecaf6f71b720273e686d520d))

# [0.6.4](https://github.com/ory/fosite/compare/v0.6.3...v0.6.4) (2016-11-29)

token/jwt: Allow single element string arrays to be treated as strings

This commit allows `aud` to be passed in as a single element array
during consent validation on Hydra. This fixes
https://github.com/ory-am/hydra/issues/314.

Signed-off-by: Son Dinh <son.dinh@blacksquaremedia.com>

### Unclassified

- token/jwt: Allow single element string arrays to be treated as strings ([5388e10](https://github.com/ory/fosite/commit/5388e107ac994650eb1623efb6c88d14d045e325)):

  This commit allows `aud` to be passed in as a single element array
  during consent validation on Hydra. This fixes
  https://github.com/ory-am/hydra/issues/314.

# [0.6.2](https://github.com/ory/fosite/compare/v0.6.1...v0.6.2) (2016-11-25)

oauth2/introspection: endpoint responds to invalid requests appropriately (#126)

### Unclassified

- oauth2/introspection: endpoint responds to invalid requests appropriately (#126) ([9360f64](https://github.com/ory/fosite/commit/9360f6473249324e2c2c2f6e94b3f123bdb929fa)), closes [#126](https://github.com/ory/fosite/issues/126)

# [0.6.1](https://github.com/ory/fosite/compare/v0.6.0...v0.6.1) (2016-11-17)

core: resolve issues with token introspection and sessions

### Unclassified

- Resolve issues with token introspection and sessions ([895d169](https://github.com/ory/fosite/commit/895d16935bd97831eecff66b1d775af9b91a2506))

# [0.6.0](https://github.com/ory/fosite/compare/v0.5.1...v0.6.0) (2016-11-17)

core: resolve session referencing issue (#125)

### Unclassified

- Comply with Go license terms - closes [#123](https://github.com/ory/fosite/issues/123) ([4c4507f](https://github.com/ory/fosite/commit/4c4507f865e0968e0a06c961aef9176bd8e7b7e3))
- Resolve session referencing issue ([#125](https://github.com/ory/fosite/issues/125)) ([81a3229](https://github.com/ory/fosite/commit/81a3229706c38e29c7745acf930272f4711547f4))

# [0.5.1](https://github.com/ory/fosite/compare/v0.5.0...v0.5.1) (2016-10-22)

handler/oauth2: set JWT ExpiresAt claim per TokenType from the session (#121)

Signed-off-by: Cristian Graziano <cristian.graziano@gmail.com>

### Unclassified

- handler/oauth2: set JWT ExpiresAt claim per TokenType from the session (#121) ([66170ae](https://github.com/ory/fosite/commit/66170ae25a3ac26abcd2ab27d687434d4e2a60a7)), closes [#121](https://github.com/ory/fosite/issues/121)
- oauth2/introspection: do not include the session in the response ([daad271](https://github.com/ory/fosite/commit/daad27179358c71aeb89dc8d7d6fdd2c04a15871))

# [0.5.0](https://github.com/ory/fosite/compare/v0.4.0...v0.5.0) (2016-10-17)

0.5.0 (#119)

- all: resolve regression issues introduced by 0.4.0 - closes #118
- oauth2: introspection handler excess calls - closes #117
- oauth2: inaccurate expires_in time - closes #72

### Unclassified

- 0.5.0 (#119) ([eb9077f](https://github.com/ory/fosite/commit/eb9077f6608d776ae50eb2ad4205705bad6ee0eb)), closes [#119](https://github.com/ory/fosite/issues/119) [#118](https://github.com/ory/fosite/issues/118) [#117](https://github.com/ory/fosite/issues/117) [#72](https://github.com/ory/fosite/issues/72)

# [0.4.0](https://github.com/ory/fosite/compare/v0.3.6...v0.4.0) (2016-10-16)

all: clean up, resolve broken tests

### Documentation

- Add danilobuerger and jrossiter to hall of fame ([f864e26](https://github.com/ory/fosite/commit/f864e26f6b22726ad592742e8654b099729a4b46))
- Add offline note to readme ([60a7672](https://github.com/ory/fosite/commit/60a767221625d0f6541f203e41a7ef20a1782eb0))
- Document reasoning for interface{} in compose package - closes [#94](https://github.com/ory/fosite/issues/94) ([f193012](https://github.com/ory/fosite/commit/f1930124e072153f9d5ec8dc4f14733f9bdc20a1))

### Unclassified

- Allow public clients to access token endpoint - closes [#78](https://github.com/ory/fosite/issues/78) ([cbe433e](https://github.com/ory/fosite/commit/cbe433e1985d782217cb973261a3b1677af1f664))
- Clean up, resolve broken tests ([1041e67](https://github.com/ory/fosite/commit/1041e67f395480fd334446bd8b13f09dfbeeb658))
- Flatten package hierarchy and merge files - closes [#93](https://github.com/ory/fosite/issues/93) ([9b7ba80](https://github.com/ory/fosite/commit/9b7ba808064d33a5251cb6cd3d30d2d4b8f3ff25))
- Reduce third party dependencies - closes [#116](https://github.com/ory/fosite/issues/116) ([5ec5cff](https://github.com/ory/fosite/commit/5ec5cff534008820671e56f6b062dc2aa1e364e6))
- Split library and example - closes [#92](https://github.com/ory/fosite/issues/92) ([6d76d35](https://github.com/ory/fosite/commit/6d76d35018159d830a9b050f99c15b099a6975e2))

# [0.3.6](https://github.com/ory/fosite/compare/v0.3.5...v0.3.6) (2016-10-07)

oauth2: added refresh token generation for password grant type (#107)

- oauth2: added refresh token generation for password grant type when offline scope is requested

Signed-off-by: Jason Rossiter <jrossiter403@gmail.com>

### Unclassified

- Added refresh token generation for password grant type ([#107](https://github.com/ory/fosite/issues/107)) ([81c3cbd](https://github.com/ory/fosite/commit/81c3cbdb6b00399219b57c9e1aa1b4cbebf888d8)):

  - oauth2: added refresh token generation for password grant type when offline scope is requested

# [0.3.5](https://github.com/ory/fosite/compare/v0.3.4...v0.3.5) (2016-10-06)

handler/oauth2: resolve issues with refresh token flow (#110)

- handler/oauth2/refresh: requestedAt time is not reset - closes #109
- handler/oauth2/refresh: session is not transported to new access token - closes #108

### Unclassified

- handler/oauth2: resolve issues with refresh token flow (#110) ([bef6197](https://github.com/ory/fosite/commit/bef61973fdee1a18aedba4e42a1d8977c3f8cc1c)), closes [#110](https://github.com/ory/fosite/issues/110) [#109](https://github.com/ory/fosite/issues/109) [#108](https://github.com/ory/fosite/issues/108)
- Add tests to request state ([8c7c77e](https://github.com/ory/fosite/commit/8c7c77e1f2116c38ed1765cc846c4b7c0bdc94b8))

# [0.3.4](https://github.com/ory/fosite/compare/v0.3.3...v0.3.4) (2016-10-04)

handler/oauth2: refresh token does not migrate original access data - closes #103 (#104)

### Unclassified

- handler/oauth2: refresh token does not migrate original access data - closes #103 (#104) ([8ffa0bc](https://github.com/ory/fosite/commit/8ffa0bc825179bbffbd3a548219062846f9b0250)), closes [#103](https://github.com/ory/fosite/issues/103) [#104](https://github.com/ory/fosite/issues/104)

# [0.3.3](https://github.com/ory/fosite/compare/v0.3.2...v0.3.3) (2016-10-03)

authorize: scopes should be separated by %20 and not +, to ensure javascript compatibility - closes #101 (#102)

### Documentation

- Fix reference to store example in readme ([#87](https://github.com/ory/fosite/issues/87)) ([b1e2cda](https://github.com/ory/fosite/commit/b1e2cda5bb64ffdcce40aed52af5c9be0852c8ef))

### Unclassified

- Scopes should be separated by %20 and not +, to ensure javascript compatibility - closes [#101](https://github.com/ory/fosite/issues/101) ([#102](https://github.com/ory/fosite/issues/102)) ([e61a25f](https://github.com/ory/fosite/commit/e61a25f3e3d3f067141c3f6464ab4213f4e14d45))

# [0.3.2](https://github.com/ory/fosite/compare/v0.3.1...v0.3.2) (2016-09-22)

openid: resolves an issue with the explicit token flow

### Unclassified

- Resolves an issue with the explicit token flow ([aa1b854](https://github.com/ory/fosite/commit/aa1b8548678e5807399d35b5bcad4f62a83cf6e4))

# [0.3.1](https://github.com/ory/fosite/compare/v0.3.0...v0.3.1) (2016-09-22)

0.3.1 (#98)

- all: better error handling - closes #100
- oauth2/implicit: bad HTML encoding of the scope parameter - closes #95
- oauth2: state parameter is missing when response_type=id_token - closes #96
- oauth2: id token hashes are not base64 url encoded - closes #97
- openid: hybrid flow using `token+code+id_token` returns multiple tokens of the same type - closes #99

### Unclassified

- 0.3.1 (#98) ([b16e3fc](https://github.com/ory/fosite/commit/b16e3fcfdf8f3f47802cd87b2388235186b9f108)), closes [#98](https://github.com/ory/fosite/issues/98) [#100](https://github.com/ory/fosite/issues/100) [#95](https://github.com/ory/fosite/issues/95) [#96](https://github.com/ory/fosite/issues/96) [#97](https://github.com/ory/fosite/issues/97) [#99](https://github.com/ory/fosite/issues/99)
- Add additional tests to HierarchicScopeStrategy ([#81](https://github.com/ory/fosite/issues/81)) ([64e869c](https://github.com/ory/fosite/commit/64e869cb9b69a4b027bfc0284bfeb33b2836ea41))
- Corrected grant type in comment ([#82](https://github.com/ory/fosite/issues/82)) ([27ddd19](https://github.com/ory/fosite/commit/27ddd19e9b07101b712b4b7d82443b3c9d53fa69))
- Removed unnecessary logging ([#86](https://github.com/ory/fosite/issues/86)) ([cb328ca](https://github.com/ory/fosite/commit/cb328caca6287c7995ee5285c6446bffd4ef496b))
- Simplify scope comparison logic ([7fb850e](https://github.com/ory/fosite/commit/7fb850ef530b3445adb07406f8bc773e6ad38884))

# [0.3.0](https://github.com/ory/fosite/compare/v0.2.4...v0.3.0) (2016-08-22)

vendor: jwt-go is now v3.0.0 (#77)

Signed-off-by: Alexander Widerberg <alexander.widerberg@cybercom.com>

### Unclassified

- HierarchicScopeStrategy worngly accepts missing scopes ([7faee6b](https://github.com/ory/fosite/commit/7faee6bbd53ee762ddfe194fb2ea5e7d0205e46d))
- Jwt-go is now v3.0.0 ([#77](https://github.com/ory/fosite/issues/77)) ([76ef7ea](https://github.com/ory/fosite/commit/76ef7ea8f51735d63476cd91e1f9a9f367d544cb))

# [0.2.4](https://github.com/ory/fosite/compare/v0.2.3...v0.2.4) (2016-08-09)

all: resolve race condition and package fosite with glide

### Unclassified

- Resolve race condition and package fosite with glide ([66b53a9](https://github.com/ory/fosite/commit/66b53a903c03950ac5180dc30c3f69e477344205))

# [0.2.3](https://github.com/ory/fosite/compare/v0.2.2...v0.2.3) (2016-08-08)

vendor: commit missing lock file

### Unclassified

- Commit missing lock file ([be30574](https://github.com/ory/fosite/commit/be30574ee5f5f51cb22faf0a187231141f1c2f63))

# [0.2.2](https://github.com/ory/fosite/compare/v0.2.1...v0.2.2) (2016-08-08)

vendor: updated go-jwt to use semver instead of gopkg

### Unclassified

- Updated go-jwt to use semver instead of gopkg ([3b66309](https://github.com/ory/fosite/commit/3b663092771e796816c1f9ac2169139f27b70c4b))

# [0.2.1](https://github.com/ory/fosite/compare/v0.2.0...v0.2.1) (2016-08-08)

core: remove unused fields and methods from client

### Unclassified

- Remove unused fields and methods from client ([5f1851b](https://github.com/ory/fosite/commit/5f1851b088e9f087a7bd3e7beca4c3112418fcfc))
- Resolved package naming issue ([4d1caeb](https://github.com/ory/fosite/commit/4d1caeb18275f2a4a5f40a7cdd06a74cfc1c3e73))

# [0.2.0](https://github.com/ory/fosite/compare/v0.1.0...v0.2.0) (2016-08-06)

all: composable factories, better token validation, better scope handling and simplify structure

- readme: add gitter chat badge closes #67
- handler: flatten packages closes #70
- openid: don't autogrant openid scope - closes #68
- all: clean up scopes / arguments - closes #66
- all: composable factories - closes #64
- all: refactor token validation - closes #63
- all: remove mandatory scope - closes #62

### Unclassified

- Composable factories, better token validation, better scope handling and simplify structure ([a92c755](https://github.com/ory/fosite/commit/a92c75531cf5bb89524cd719c9bc2c98fe709c62)), closes [#67](https://github.com/ory/fosite/issues/67) [#70](https://github.com/ory/fosite/issues/70) [#68](https://github.com/ory/fosite/issues/68) [#66](https://github.com/ory/fosite/issues/66) [#64](https://github.com/ory/fosite/issues/64) [#63](https://github.com/ory/fosite/issues/63) [#62](https://github.com/ory/fosite/issues/62)

# [0.1.0](https://github.com/ory/fosite/compare/7adad58c327cf52530d8c1e08059564ca0b51538...v0.1.0) (2016-08-01)

oauth2: implicit handlers do not require tls over https (#61)

closes #60

### Code Refactoring

- New api signatures ([8a830d3](https://github.com/ory/fosite/commit/8a830d34405f3b3d50734f5258151426dc61a94b))

### Documentation

- Add -d option to go get ([0e63038](https://github.com/ory/fosite/commit/0e630382425e6d1a7e9177828eeb59f6748e856f))
- Define implicitHandler ([745a4df](https://github.com/ory/fosite/commit/745a4df7758caa8c3338d006a60f4948120f00bf)):

  Someone forgot to rename the variable name when copy-pasting in the example.

- Document new token generation and validation ([ddef55b](https://github.com/ory/fosite/commit/ddef55ba96b6c533681b7a1953da5c33ed30587a))
- Drafted workflows ([4ad1d14](https://github.com/ory/fosite/commit/4ad1d146d67c0e17c545d1c3959dc697777b9828))
- Explain what handlers are ([48ca03b](https://github.com/ory/fosite/commit/48ca03b9026843f1047e510c3b66ccb6a54def2c))
- Fix typos in readme ([b9ed7ac](https://github.com/ory/fosite/commit/b9ed7acf8b00f05fcc99578f7a49d55275041515))
- Readme ([a5aa697](https://github.com/ory/fosite/commit/a5aa69736505502303bc99ee180539033d5ba886))
- Readme ([f77fd41](https://github.com/ory/fosite/commit/f77fd412ea7f2be15b0f0c5ac801ac177e7d3dc4))
- Readme ([e143d8c](https://github.com/ory/fosite/commit/e143d8ca506f7cf2f70c92710b2fc123e003a12d))
- Readme ([d483568](https://github.com/ory/fosite/commit/d483568c06d9542bbf383771dee3ea44b60dff0e))
- Updated authorize section ([9c21afb](https://github.com/ory/fosite/commit/9c21afbc38fbd35f951c127beb2623ae4d2590e7))
- Updated readme docs ([336a2cd](https://github.com/ory/fosite/commit/336a2cd10ac08ca6867952555802c225c475c17a))

### Unclassified

- updated gif ([39c239f](https://github.com/ory/fosite/commit/39c239faca97882da9d5293306dfdcbabf8ee0cc))
- gofmt ([f813288](https://github.com/ory/fosite/commit/f813288911ba653b197589edc4206b52d6c11545))
- updated example gif ([29b39ea](https://github.com/ory/fosite/commit/29b39ea32fee62b1013ee383ce56c653a7ef33d9))
- added open id connect to example ([6f0ce68](https://github.com/ory/fosite/commit/6f0ce681147428b51c3673a4c46ab018cf46cf81))
- added integration tests ([8d47f80](https://github.com/ory/fosite/commit/8d47f80420c288a25ba846927c532e156d27a23b))
- added doc to fix travis ([a0db129](https://github.com/ory/fosite/commit/a0db129b0a063fe9438560b1a339f973736327f7))
- Add go report card ([204c5d6](https://github.com/ory/fosite/commit/204c5d60b6f42b0e8f918bdd96214698ad3717da))
- Clean-up fosite-example/main.go link in README.md ([497ff80](https://github.com/ory/fosite/commit/497ff807a10a9fb41b697c5f91ed9eeb26375b24)):

  The README url to the suggested example was broken.

- Added jti as parameter to claims helper to privide better interface to developers ([bde3822](https://github.com/ory/fosite/commit/bde38221ed4d32c2f175a60540ac529b306a2ced))
- Added missing jti claim ([26f41a0](https://github.com/ory/fosite/commit/26f41a06689bd12f7165044a2de7d9332fea3759))
- Added NOTE ([64516f8](https://github.com/ory/fosite/commit/64516f8e2e0154f46358723d710447380f6d5dc2))
- Removed unnecessary print. Added bugfix from Arekkas. ([96458b6](https://github.com/ory/fosite/commit/96458b6cf8ee46edbef35598b6d3d877fb63ff87))
- Example updated ([5022339](https://github.com/ory/fosite/commit/50223396d01d742b1a0a3f0be1252e339cf22985))
- Added working example of jwt token ([9410fca](https://github.com/ory/fosite/commit/9410fca73dfb00f1dc1e3aa6ec580554ec3daaba))
- Added tests. Still need to verify implemtation with test ([1ebdd88](https://github.com/ory/fosite/commit/1ebdd88746c875bff1a6d074437c5742c812a200))
- WIP ([caaa43a](https://github.com/ory/fosite/commit/caaa43a184a66b78972fa3725d3636837da1cd68))
- readme ([c97d844](https://github.com/ory/fosite/commit/c97d84471bc3941e479a79ef2eed4b1ddc07f21c))
- readme ([fe24f26](https://github.com/ory/fosite/commit/fe24f261de60711d91c016c435ce83938d367609))
- readme ([be8cd23](https://github.com/ory/fosite/commit/be8cd2333d3eaaf266b56c30951741d7f88edc5e))
- refactor done (unstaged) ([625f168](https://github.com/ory/fosite/commit/625f1683a0449384877823c2dae1464718c0b264))
- unstaged ([6c616b1](https://github.com/ory/fosite/commit/6c616b12198419ed33035dabd9e33d1e2afffff2))
- unstaged ([17ad70b](https://github.com/ory/fosite/commit/17ad70b88ff6ba2add1136762428340d21b86126))
- Include user session data in all calls to storage handlers. ([2be3fc1](https://github.com/ory/fosite/commit/2be3fc18f5a35646f7cd001eb6b4b92cbb07ef16))
- unstaged ([fde7c80](https://github.com/ory/fosite/commit/fde7c803798b1f7fa2056bb434dd74d9a4ebeea7))
- unstaged ([e775aad](https://github.com/ory/fosite/commit/e775aadbc33ec8f15adc7f3b78de5eca53b349f5))
- unstaged ([ae2fc16](https://github.com/ory/fosite/commit/ae2fc169e663486248f6518a3497b0245754892e))
- handler/core: fixed tests ([7f5938a](https://github.com/ory/fosite/commit/7f5938adc4f79380239292cd3b6f6e0064df39ef))
- core handlers: added tests ([e9affb7](https://github.com/ory/fosite/commit/e9affb77442c46fb4647c9a22c1a5eb60945d21d))
- authorize/explicit ✓ ([d61635b](https://github.com/ory/fosite/commit/d61635b26e3cd34822d4f3ffc0fe25bd4774bd45))
- authorize/explicit: minor name refactoring and tests for authorize endpoint ([4736e28](https://github.com/ory/fosite/commit/4736e284b327f0941e58073bf860caca4117c545))
- plugin/token: fix import path ([fdba2f7](https://github.com/ory/fosite/commit/fdba2f7b5bdec0e77faa804066abe1b8895b909e))
- unstaged ([f939597](https://github.com/ory/fosite/commit/f939597f3f3e6ad4eb582a56b643589271cbf646))
- Initial commit ([7adad58](https://github.com/ory/fosite/commit/7adad58c327cf52530d8c1e08059564ca0b51538))
- Access code request workflow finalized ([0232918](https://github.com/ory/fosite/commit/0232918e250eeee93bdab98502a5a30273510c49))
- Access request api draft ([9f482ef](https://github.com/ory/fosite/commit/9f482ef50711b608dbfb72022ef998f947f0487a))
- Add api stability section ([3ca6ec9](https://github.com/ory/fosite/commit/3ca6ec936d6b3a8dab0add136b3a2fbfefa4b4df))
- Add go-rethink tags ([49c82bc](https://github.com/ory/fosite/commit/49c82bc9fe0c4edbb90579e1746e0dad1ae01c5c))
- Add ValidateToken to CoreValidator ([4c2b9d8](https://github.com/ory/fosite/commit/4c2b9d8f0c84f19ae11f59cb07927ceb59598adc))
- Added authorize code grant example ([269c5fa](https://github.com/ory/fosite/commit/269c5fab1109bb4cd2e624940dac1b9467663507))
- Added client grant and did some renaming ([75c8179](https://github.com/ory/fosite/commit/75c8179ef537e6ea87b16cdd87016fca6d389490))
- Added cristiangraz to the hall of fame ([1b6e2b4](https://github.com/ory/fosite/commit/1b6e2b470f8f477fdfb2ec1f914e64293bdc7b1b))
- Added danielchatfield to the hall of fame ([2b988a8](https://github.com/ory/fosite/commit/2b988a8b2abd3dea619e31e174b306e45a62fcc1))
- Added go 1.6 ([ae41a0a](https://github.com/ory/fosite/commit/ae41a0ace8f74480fec08c83fb1c7bda35830f35))
- Added go1.4 to allowed failures ([49aa920](https://github.com/ory/fosite/commit/49aa920401a3cf62f16541d8fa4f9fb488270cf3))
- Added grant and response type validation ([f524fc2](https://github.com/ory/fosite/commit/f524fc2b026621192407ce22e71f2b062635b134))
- Added json and gorethink tags ([99c836c](https://github.com/ory/fosite/commit/99c836cd526c276419e31db25b695dd0097f0656))
- Added JWT generator and validator. ([58acd68](https://github.com/ory/fosite/commit/58acd688530666f4720eeacb598da72a475282d5)), closes [#16](https://github.com/ory/fosite/issues/16)
- Added missing file ([8fc1615](https://github.com/ory/fosite/commit/8fc1615bf40777c2c456e1ec4515a269e348e3b4))
- Added owner method ([78012ed](https://github.com/ory/fosite/commit/78012ed85819caaf154fe9dc4afd212f068fc0a1))
- Added tests fragment capabilities to writeresponse ([6df0eca](https://github.com/ory/fosite/commit/6df0eca1d74d79e807a77910776ff2249340f103))
- Api cleanup, gofmt ([3d6e8b6](https://github.com/ory/fosite/commit/3d6e8b6281c6d170a77103b89cfabdd3086a03f0))
- Api refactor ([d936c91](https://github.com/ory/fosite/commit/d936c914253c58297dcc462a14fb6ddb87bfcac4))
- Basic draft ([480af91](https://github.com/ory/fosite/commit/480af9165fef8a5e8bcc4896ed680cbf5afbe23c))
- Defined OAuth2.HandleResponseTypes ([30b6e74](https://github.com/ory/fosite/commit/30b6e74b13f567237ea770bf6a4e99dd95085dcc)):

  Incorporated feedback from GitHub, did refactoring and renaming, added tests

- Enforce https for all redirect endpoints except localhost ([d65b45a](https://github.com/ory/fosite/commit/d65b45a192cd3a2073f8e6118c005ac93f0bb974))
- Enforce use of scopes ([12d76dd](https://github.com/ory/fosite/commit/12d76dd7c86408e52f85a3099f6063c462e0b97b)), closes [#14](https://github.com/ory/fosite/issues/14)
- Finalized auth endpoint, added tests, added integration tests ([c6dcb90](https://github.com/ory/fosite/commit/c6dcb90ccbd1d7a179a601e0e6d46cc1004cde92))
- Finalized token endpoint api ([8de3f10](https://github.com/ory/fosite/commit/8de3f10d89b47ad0d23cf13b425442393f51e104))
- Finished up integration tests ([a6d027e](https://github.com/ory/fosite/commit/a6d027e3a4f817bb72706fbf0d7e3245f8823b27))
- Fix broken test ([653e324](https://github.com/ory/fosite/commit/653e3248c0a1aae3bb2c33f64f21854155304e1a))
- Fix config ([82e9332](https://github.com/ory/fosite/commit/82e9332815579e538089dff61281a7a446f0f6cd))
- Fix deps ([bcc6a07](https://github.com/ory/fosite/commit/bcc6a07fef6f4036643e79eaf3cdd1f485a682fb))
- Fix jwt strategy interface mismatch ([#58](https://github.com/ory/fosite/issues/58)) ([4d0a545](https://github.com/ory/fosite/commit/4d0a5450dd3b44e44f5169f90b3591566a6eef1d))
- Fix unique scope tests ([3ac3a79](https://github.com/ory/fosite/commit/3ac3a798cd1ad5fcd0a53abb45fbb93c7321d154))
- Fixed granted scope match ([13b7efa](https://github.com/ory/fosite/commit/13b7efae68b4f68171422b876e8df197b3453e42))
- Fixed racy tests ([f0b691d](https://github.com/ory/fosite/commit/f0b691dac03f455ae429116cf121a1ae9054c3e3))
- Fixed tests ([8bf73e3](https://github.com/ory/fosite/commit/8bf73e3bb4b12e098f63b1007d4ce9a25e0221b7))
- Fixed tests refactor broke ([5da857b](https://github.com/ory/fosite/commit/5da857b4bcf76b3cc87aa5c9c1f8ee2c0c814992))
- Fixed typos ([a5391de](https://github.com/ory/fosite/commit/a5391deaa543441f1e3838b0c5093692be247015)), closes [#10](https://github.com/ory/fosite/issues/10)
- Fixed urls ([58908b8](https://github.com/ory/fosite/commit/58908b8cd323434dce944119c5a300f1196634f2))
- Fixed wrongfully set constant ErrTemporaryUnvailableName ([71a9105](https://github.com/ory/fosite/commit/71a9105a1e4afde3eed0a3ef80239140f6674d15)), closes [#9](https://github.com/ory/fosite/issues/9)
- Generic claims and headers ([1f2e97f](https://github.com/ory/fosite/commit/1f2e97ff847921939fe1f93f6dfdfcbb7bfb0792))
- Glide ([#43](https://github.com/ory/fosite/issues/43)) ([de85e2a](https://github.com/ory/fosite/commit/de85e2a7ebce57a804ae0beef42b1f1b9017914c))
- Godep save ([c457104](https://github.com/ory/fosite/commit/c45710465f990e74e8cddf5190f2e309da592297))
- Goimports ([8b9816c](https://github.com/ory/fosite/commit/8b9816cb1ecbc7befef924b6a923bd52530141f3))
- Goimports ([96be194](https://github.com/ory/fosite/commit/96be194cae6562fe35696c6ee6c7c547ce20388d))
- Implemented all core grant types ([ce0a849](https://github.com/ory/fosite/commit/ce0a8496942259d6fe518104bab0dfd3dfea9856))
- Implemented and documented examples ([8c625c9](https://github.com/ory/fosite/commit/8c625c9cd1e9854eddecafc36e4502577c113ef0))
- Implemented new token generator based on hmac-sha256 ([01f9ede](https://github.com/ory/fosite/commit/01f9ede7e69588caf12940979a1fc0586d5aac3c)), closes [#11](https://github.com/ory/fosite/issues/11)
- Implemented validator for access tokens ([4140422](https://github.com/ory/fosite/commit/414042259d6f7b1aefe4244bc3f8eb80a83a2d2c))
- Implicit handlers do not require tls over https ([#61](https://github.com/ory/fosite/issues/61)) ([6c40c08](https://github.com/ory/fosite/commit/6c40c086a1f082d466bac21721571558c32de97c)), closes [#60](https://github.com/ory/fosite/issues/60)
- Improve handling of expiry and include a protected api example ([dfb047d](https://github.com/ory/fosite/commit/dfb047d52b75b5d8a28bcd8d70a3e139da289da1))
- Improve strategy API ([21f5e8c](https://github.com/ory/fosite/commit/21f5e8ce68097959ef97b1b8dca268f2a9a5d276))
- Increased coverage ([83194b6](https://github.com/ory/fosite/commit/83194b6b2849292da041385e2274d42a06b36120))
- Issue refresh token only when 'offline' scope is set ([34068b9](https://github.com/ory/fosite/commit/34068b951d8deea523c40f792608b75d2b4c656f)), closes [#47](https://github.com/ory/fosite/issues/47)
- Jwt signing and client changes ([#44](https://github.com/ory/fosite/issues/44)) ([fae3c96](https://github.com/ory/fosite/commit/fae3c96e89cd364f21bee00f8d5384cd053ab9c1))
- Made hybrid flow optional ([08ddbae](https://github.com/ory/fosite/commit/08ddbae46bca5ef18e4a8c7560a46d6238d6a3e9))
- Major refactor, use enigma, finalized authorize skeleton ([38bacd3](https://github.com/ory/fosite/commit/38bacd340eed991d69dc95f8a7bf6c0f328d8f47)), closes [#8](https://github.com/ory/fosite/issues/8) [#11](https://github.com/ory/fosite/issues/11)
- More test cases ([1188750](https://github.com/ory/fosite/commit/1188750e06c6ba30ebc783a8297aab75a0f95247))
- More tests ([164506a](https://github.com/ory/fosite/commit/164506a23a3105a37b60b1154052589d1be6c7b1))
- Moved to root package, updated docs ([1871702](https://github.com/ory/fosite/commit/18717023c4d6b5c02691f94fe80714f2e5e9862d))
- Moved to root package, updated docs ([5b9b20c](https://github.com/ory/fosite/commit/5b9b20cd6b91a5cf72d054dc9afa2afc9d6dfd15))
- No "session" secret required ([d1f45ad](https://github.com/ory/fosite/commit/d1f45ad9dcbb0b2866f7c8fa0fe99bc77fb93506))
- Preview ([ba84987](https://github.com/ory/fosite/commit/ba849870e24070ea44fec9cbcf99cc04a281ffef))
- Refactor ([eb9153c](https://github.com/ory/fosite/commit/eb9153c389b1c7ca14af78b091705d84e5bba68c))
- Refactor, fixed tests, incorporated feedback ([9e59df2](https://github.com/ory/fosite/commit/9e59df23353964644bfcc0d148745f8dca691b39))
- Refactoring, more tests ([df79a81](https://github.com/ory/fosite/commit/df79a81577ec8a9b7517af794ea6f04da71abf91))
- Refactoring, renaming, docs ([e5476d1](https://github.com/ory/fosite/commit/e5476d15413c7bf96b5a1c282f9d079f538dcc83))
- Refactoring, renaming, more tests ([9467ca8](https://github.com/ory/fosite/commit/9467ca8ac7b7b7785c96f049a422ed1d16e639b4))
- Remove duplicate field ([e134351](https://github.com/ory/fosite/commit/e13435109928d11ae9eeb13f1e347043e8be0d53))
- Remove store mock ([80c14f7](https://github.com/ory/fosite/commit/80c14f786b4a1ed4f1379a5fd6deaf036ece4b47))
- Rename fields name to client_name and secret to client_secret ([99ce066](https://github.com/ory/fosite/commit/99ce0662f10c82ce034c9c21c8041aa29c460883))
- Renaming and refactoring ([d3697bd](https://github.com/ory/fosite/commit/d3697bd15cc05bbc8bf3a6833911c3cc5dd1b2f8))
- Replace internal import ([#52](https://github.com/ory/fosite/issues/52)) ([1290282](https://github.com/ory/fosite/commit/1290282d421ee999ff8e5c2d5d6d0f762dba599c))
- Replace pkg.ErrNotFound with fosite.ErrNotFound ([4390c49](https://github.com/ory/fosite/commit/4390c495a1794fc7cf26cbeb47969f92d19f0ecc))
- Request should return unique scopes ([af66918](https://github.com/ory/fosite/commit/af66918f0c91a451659fa2bf01d2c804e14799eb))
- Resolve an issue where query params could be used instead of post body ([7eb85c6](https://github.com/ory/fosite/commit/7eb85c6e4ae2bb4a67c2e6f6166824351cc17f1d))
- Resolve danger of not reading enough bytes ([c68a3e9](https://github.com/ory/fosite/commit/c68a3e9bea4bb5a6550e55b2ce2beb59eb48782a))
- Resolve id token issues with empty claims ([89c60c9](https://github.com/ory/fosite/commit/89c60c9f2898345fd3d75044c8e41eacbf0d4fd5))
- Resolve scope issues ([#55](https://github.com/ory/fosite/issues/55)) ([9d54b98](https://github.com/ory/fosite/commit/9d54b989c8d04c4d586e7810cce2e6d4f03d7c48)):

  handler: resolve scope issues

- Sanitized tests and apis ([12c70bb](https://github.com/ory/fosite/commit/12c70bb4f167afe8d39e85d3ef0e0f13b5761070))
- Tests for client credentials flow ([c13298c](https://github.com/ory/fosite/commit/c13298cbf165c873f9463a6bbad91b962762f3b0))
- Tests for resource owner password credentials grant ([f503615](https://github.com/ory/fosite/commit/f5036150f90d7d73e85088400cda9f7de2722a20))
- Update ([88e84de](https://github.com/ory/fosite/commit/88e84de2676281bb5a7a1e6b5051faa1feb14c2e))
- Update installation instructions ([201c6aa](https://github.com/ory/fosite/commit/201c6aa6c15d35da14022f7ec43d0e9b87b2bc68)), closes [#33](https://github.com/ory/fosite/issues/33)
- Updated example and added implicit grant ([d12fa5c](https://github.com/ory/fosite/commit/d12fa5ca89cfebb351e023d53b0c57420725195b))
- Use jwt-go.v2 and fix bc break ([f731d88](https://github.com/ory/fosite/commit/f731d8892ca50501fdc054023f0b7b77d9ecb6ef))
