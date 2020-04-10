## Overview

This documentation describes how to enable the email-related capabilities of Remark.

- email authentication for users:

    enabling this will let the user log in using their emails:

    ![Email authentication](/docs/images/email_auth.png?raw=true)

- email notifications for any users except anonymous:

    GitHub or Google or Twitter or any other kind of user gets the ability to get email notifications about new replies to their comments:

    ![Email notifications subscription](/docs/images/email_notifications.png?raw=true)

## Setup email provider

To enable any of email functionality you need to set up an email provider.

all providers are **WEB API** based except for smtp.

the recommended way is to use the **WEB API** provider for many good reasons

current supported providers
```ini
mailgun
sendgrid
smtp
```
common variables:

```
EMAIL_PROVIDER
```

### `mailgun` provider setup

The `mailgun` provider requires these variables:

```ini
EMAIL_PROVIDER
EMAIL_MG_DOMAIN
EMAIL_MG_API_KEY
EMAIL_MG_TIMEOUT
```

### `sendgrid` provider setup

The `sendgrid` provider requires these variables:

```ini
EMAIL_PROVIDER
EMAIL_SG_API_KEY
EMAIL_SG_TIMEOUT
```

### `smtp` provider setup

The `smtp` provider requires these variables:

```
EMAIL_PROVIDER
EMAIL_SMTP_HOST
EMAIL_SMTP_PORT
EMAIL_SMTP_TLS
EMAIL_SMTP_USERNAME
EMAIL_SMTP_PASSWORD
EMAIL_SMTP_TIMEOUT
```

### Mailgun SMTP

This is an example of a configuration using [Mailgun](https://www.mailgun.com/) email service:

```
      - EMAIL_PROVIDER=smtp
      - EMAIL_SMTP_HOST=smtp.eu.mailgun.org
      - EMAIL_SMTP_PORT=465
      - EMAIL_SMTP_TLS=true
      - EMAIL_SMTP_USERNAME=postmaster@mg.example.com
      - EMAIL_SMTP_PASSWORD=secretpassword
      - AUTH_EMAIL_FROM=notify@example.com
```

### Gmail SMTP

Configuration example for Gmail:

```
      - EMAIL_PROVIDER=smtp
      - EMAIL_SMTP_HOST=smtp.gmail.com
      - EMAIL_SMTP_PORT=465
      - EMAIL_SMTP_TLS=true
      - EMAIL_SMTP_USERNAME=example.user@gmail.com
      - EMAIL_SMTP_PASSWORD=secretpassword
      - AUTH_EMAIL_FROM=example.user@gmail.com
```

### Amazon SES SMTP

Configuration example for [Amazon SES](https://aws.amazon.com/ses/) (us-east-1 region):
```
      - EMAIL_PROVIDER=smtp
      - EMAIL_SMTP_HOST=email-smtp.us-east-1.amazonaws.com
      - EMAIL_SMTP_PORT=465
      - EMAIL_SMTP_TLS=true
      - EMAIL_SMTP_USERNAME=access_key_id
      - EMAIL_SMTP_PASSWORD=secret_access_key
      - AUTH_EMAIL_FROM=notify@example.com

```

A domain or an email that will be used in `AUTH_EMAIL_FROM` or `NOTIFY_EMAIL_FROM` must first be [verified](https://docs.aws.amazon.com/ses/latest/DeveloperGuide/verify-domain-procedure.html).

[SMTP Credentials](https://docs.aws.amazon.com/ses/latest/DeveloperGuide/smtp-credentials.html) must first be obtained from [Amazon SES Console](https://console.aws.amazon.com/ses/home?region=us-east-1#smtp-settings:):

## Setup email authentication

Here is the list of variables which affect email authentication:

```
AUTH_EMAIL_ENABLE
AUTH_EMAIL_FROM
AUTH_EMAIL_SUBJ
AUTH_EMAIL_CONTENT_TYPE
AUTH_EMAIL_TEMPLATE
```

After `SMTP_` variables are set, you can allow email authentication by setting these two variables:

```
      - AUTH_EMAIL_ENABLE=true
      - AUTH_EMAIL_FROM=notify@example.com
```


Usually, you don't need to change/set anything else. In case if you want to use a different email template set `AUTH_EMAIL_TEMPLATE`, for instance
`- AUTH_EMAIL_TEMPLATE="Confirmation email, token: {{.Token}}"`. See [verified-authentication](https://github.com/go-pkgz/auth#verified-authentication) for more details.

## Setup email notifications

Here is the list of variables which affect email notifications:

```
NOTIFY_TYPE
NOTIFY_EMAIL_FROM
NOTIFY_EMAIL_VERIFICATION_SUBJ
# for administrator notifications for new comments on their site
ADMIN_SHARED_EMAIL
NOTIFY_EMAIL_ADMIN
```

After `SMTP_` variables are set, you can allow email notifications by setting these two variables:

```
      - NOTIFY_TYPE=email
      # - NOTIFY_TYPE=email,telegram # this is in case you want to have both email and telegram notifications enabled
      - NOTIFY_EMAIL_FROM=notify@example.com
```
