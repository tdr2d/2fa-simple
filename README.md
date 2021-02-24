# 2fa-simple 
Go Simple implementation of two-factor login flow using a local configuration for users.
- Verification is done using a verification code sent by email
- Act as an http server for Login to any Single Page Application
- Supports i18n
- TailwindCss is used for the Web-UI


# Dev requirements
- go 1.14+
- Environment:
    - `export SENDGRID_API_KEY='7OL816snJH6yjECp4eO_DT8'`
- node/npm for building tailwind css


# TODO:
- i18n
- docker


Tests:
- api end to end tests


Optimizations:
- error page template
- compression
- caching
- cors/xss security