# 2fa-simple 
Go Simple implementation of two-factor login flow using a local configuration for users. <br>
Verification is done using a verification code sent by email. <br>
Act as an http server for Login to any Single Page Application <br>

# Requirement
- go 1.14+
- Environment:
    - `export SENDGRID_API_KEY='7OL816snJH6yjECp4eO_DT8'`


# TODO:
- forgot password
- reset config url
- UI design
- frontend form validation
- error template
- graceful stop/restart
- i18n


Tests:
- api end to end tests


Optimizations:
- compression
- caching
- cors/xss security