# 2fa-simple 
Go Simple implementation of two-factor login flow using a local configuration for users. <br>
Verification is done using a verification code sent by email. <br>
Act as an http server for Login to any Single Page Application <br>
TailwindCss is used for the Web-UI<br>

# Requirement
- go 1.14+
- Environment:
    - `export SENDGRID_API_KEY='7OL816snJH6yjECp4eO_DT8'`
- node/npm


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