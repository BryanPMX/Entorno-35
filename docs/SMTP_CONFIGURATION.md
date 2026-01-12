# SMTP Email Configuration Guide

This guide explains how to configure SMTP settings for sending assessment invitation emails in Entorno35.

## Required Environment Variables

The following environment variables must be configured for email functionality:

```bash
SMTP_HOST=          # SMTP server hostname (e.g., smtp.gmail.com)
SMTP_PORT=          # SMTP server port (e.g., 587 for TLS, 465 for SSL)
SMTP_USER=          # SMTP username/email address
SMTP_PASSWORD=      # SMTP password or app-specific password
SMTP_FROM_ADDRESS=  # Sender email address (usually same as SMTP_USER)
SMTP_FROM_NAME=     # Sender display name (e.g., "Entorno35 - NOM-035")
APP_BASE_URL=       # Base URL of your application (e.g., https://app.entorno35.com)
```

## Quick Setup Guide

### Option 1: Gmail (Recommended for Development/Testing)

Gmail is free and easy to set up, perfect for development and small deployments.

#### Step 1: Enable 2-Factor Authentication
1. Go to [Google Account Security](https://myaccount.google.com/security)
2. Enable 2-Step Verification if not already enabled

#### Step 2: Generate an App Password
1. Go to [Google App Passwords](https://myaccount.google.com/apppasswords)
2. Select "Mail" as the app and "Other (Custom name)" as the device
3. Enter "Entorno35" as the name
4. Click "Generate"
5. Copy the 16-character password (you'll use this as `SMTP_PASSWORD`)

#### Step 3: Configure Environment Variables

```bash
SMTP_HOST=smtp.gmail.com
SMTP_PORT=587
SMTP_USER=your-email@gmail.com
SMTP_PASSWORD=xxxx xxxx xxxx xxxx  # The 16-character app password from Step 2
SMTP_FROM_ADDRESS=your-email@gmail.com
SMTP_FROM_NAME=Entorno35 - NOM-035
APP_BASE_URL=http://localhost:3000  # For development
```

**Note**: Gmail has sending limits (500 emails/day for free accounts, 2000/day for Google Workspace).

---

### Option 2: SendGrid (Recommended for Production)

SendGrid offers a free tier (100 emails/day) and is designed for transactional emails.

#### Step 1: Create a SendGrid Account
1. Sign up at [SendGrid](https://sendgrid.com/)
2. Verify your email address
3. Complete the account setup

#### Step 2: Create an API Key
1. Go to **Settings** → **API Keys**
2. Click **Create API Key**
3. Name it "Entorno35 Production" (or similar)
4. Select **Full Access** or **Restricted Access** (Mail Send permissions)
5. Click **Create & View**
6. **Copy the API key immediately** (you won't be able to see it again)

#### Step 3: Verify Sender Identity
1. Go to **Settings** → **Sender Authentication**
2. Choose **Single Sender Verification**
3. Enter your sender details (name, email, company)
4. Verify the email address by clicking the verification link sent to your email

#### Step 4: Configure Environment Variables

```bash
SMTP_HOST=smtp.sendgrid.net
SMTP_PORT=587
SMTP_USER=apikey  # Literally the word "apikey"
SMTP_PASSWORD=SG.xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx  # Your API key from Step 2
SMTP_FROM_ADDRESS=noreply@yourdomain.com  # Must be verified in Step 3
SMTP_FROM_NAME=Entorno35 - NOM-035
APP_BASE_URL=https://app.entorno35.com  # Your production URL
```

**Note**: For SendGrid, `SMTP_USER` must be exactly `apikey` (not your SendGrid username).

---

### Option 3: AWS SES (Amazon Simple Email Service)

AWS SES is cost-effective for high-volume sending (very low cost per email).

#### Step 1: Set Up AWS Account
1. Sign up for [AWS](https://aws.amazon.com/)
2. Navigate to **SES** (Simple Email Service) in the AWS Console

#### Step 2: Verify Email Address or Domain
1. Go to **Verified identities** → **Create identity**
2. Choose **Email address** or **Domain** (domain recommended for production)
3. Follow verification steps (verify email or configure DNS records)

#### Step 3: Create SMTP Credentials
1. Go to **SMTP settings** in SES console
2. Click **Create SMTP credentials**
3. Enter a name (e.g., "entorno35-smtp")
4. Click **Create**
5. **Download the credentials** (save both username and password)

#### Step 4: Move Out of SES Sandbox (Optional but Recommended)
- By default, SES is in "sandbox mode" (can only send to verified emails)
- To send to any email, submit a request to move out of sandbox mode
- Go to **Account dashboard** → **Request production access**

#### Step 5: Configure Environment Variables

```bash
SMTP_HOST=email-smtp.us-east-1.amazonaws.com  # Replace with your AWS region
SMTP_PORT=587
SMTP_USER=AKIAIOSFODNN7EXAMPLE  # Your SES SMTP username
SMTP_PASSWORD=wJalrXUtnFEMI/K7MDENG/bPxRfiCYEXAMPLEKEY  # Your SES SMTP password
SMTP_FROM_ADDRESS=noreply@yourdomain.com  # Must be verified in Step 2
SMTP_FROM_NAME=Entorno35 - NOM-035
APP_BASE_URL=https://app.entorno35.com
```

**AWS SES SMTP Endpoints by Region:**
- US East (N. Virginia): `email-smtp.us-east-1.amazonaws.com`
- US West (Oregon): `email-smtp.us-west-2.amazonaws.com`
- EU (Ireland): `email-smtp.eu-west-1.amazonaws.com`
- EU (Frankfurt): `email-smtp.eu-central-1.amazonaws.com`

---

### Option 4: Mailgun

Mailgun offers a free tier (100 emails/day for first 3 months, then 5,000 emails/month).

#### Step 1: Create a Mailgun Account
1. Sign up at [Mailgun](https://www.mailgun.com/)
2. Verify your email address
3. Add and verify your domain (or use sandbox domain for testing)

#### Step 2: Get SMTP Credentials
1. Go to **Sending** → **Domain Settings**
2. Click on your domain
3. Navigate to **SMTP credentials** section
4. Your SMTP credentials will be displayed:
   - **SMTP hostname**: `smtp.mailgun.org`
   - **Username**: Your Mailgun account email
   - **Password**: Your Mailgun account password (or API key)

#### Step 3: Configure Environment Variables

```bash
SMTP_HOST=smtp.mailgun.org
SMTP_PORT=587
SMTP_USER=postmaster@mg.yourdomain.com  # Or your Mailgun account email
SMTP_PASSWORD=your-mailgun-password  # Or API key
SMTP_FROM_ADDRESS=noreply@yourdomain.com  # Must be verified domain
SMTP_FROM_NAME=Entorno35 - NOM-035
APP_BASE_URL=https://app.entorno35.com
```

---

## Setting Environment Variables

### Development (Local)

#### Option A: Using `.env` file (Recommended)

1. Create a `.env` file in the project root (if it doesn't exist):

```bash
# Copy the example if available
cp .env.example .env
```

2. Add your SMTP configuration to `.env`:

```bash
# SMTP Configuration
SMTP_HOST=smtp.gmail.com
SMTP_PORT=587
SMTP_USER=your-email@gmail.com
SMTP_PASSWORD=your-app-password
SMTP_FROM_ADDRESS=your-email@gmail.com
SMTP_FROM_NAME=Entorno35 - NOM-035
APP_BASE_URL=http://localhost:3000
```

3. Load the `.env` file when running the application (many frameworks do this automatically)

#### Option B: Export in Shell

```bash
export SMTP_HOST=smtp.gmail.com
export SMTP_PORT=587
export SMTP_USER=your-email@gmail.com
export SMTP_PASSWORD=your-app-password
export SMTP_FROM_ADDRESS=your-email@gmail.com
export SMTP_FROM_NAME="Entorno35 - NOM-035"
export APP_BASE_URL=http://localhost:3000

# Then run your application
go run cmd/api/main.go
```

### Production

#### Option A: Environment Variables in Deployment Platform

**Docker/Docker Compose:**
```yaml
services:
  api:
    environment:
      SMTP_HOST: ${SMTP_HOST}
      SMTP_PORT: ${SMTP_PORT}
      SMTP_USER: ${SMTP_USER}
      SMTP_PASSWORD: ${SMTP_PASSWORD}
      SMTP_FROM_ADDRESS: ${SMTP_FROM_ADDRESS}
      SMTP_FROM_NAME: ${SMTP_FROM_NAME}
      APP_BASE_URL: ${APP_BASE_URL}
```

**Kubernetes:**
```yaml
env:
  - name: SMTP_HOST
    valueFrom:
      secretKeyRef:
        name: smtp-credentials
        key: host
  - name: SMTP_PASSWORD
    valueFrom:
      secretKeyRef:
        name: smtp-credentials
        key: password
```

**Heroku:**
```bash
heroku config:set SMTP_HOST=smtp.sendgrid.net
heroku config:set SMTP_PORT=587
heroku config:set SMTP_USER=apikey
heroku config:set SMTP_PASSWORD=your-api-key
```

**AWS ECS/Fargate:**
- Use AWS Systems Manager Parameter Store or Secrets Manager
- Reference in task definition

#### Option B: Using a Secrets Manager

For production, use a secrets manager (AWS Secrets Manager, HashiCorp Vault, etc.) to securely store credentials.

---

## Testing Email Configuration

After configuring SMTP settings, test that emails are being sent:

1. **Check Application Logs**: The email service logs whether emails are sent or logged (if SMTP is not configured):
   ```
   [EMAIL] Email sent successfully to: user@example.com
   ```
   Or (if not configured):
   ```
   [EMAIL] Would send email to: user@example.com
   Subject: Evaluacion NOM-035
   Body length: 1234 bytes
   ```

2. **Send a Test Email**: Use the assessment email feature in the application to send a test email to a known address.

3. **Check Email Provider Dashboard**: Most email providers (SendGrid, Mailgun, AWS SES) have dashboards showing sent emails, bounces, and errors.

---

## Troubleshooting

### Emails Not Being Received

1. **Check SMTP Configuration**:
   - Verify all environment variables are set correctly
   - Check that `SMTP_PASSWORD` is correct (especially for Gmail app passwords)
   - Ensure `SMTP_PORT` matches your provider's requirements (587 for TLS, 465 for SSL)

2. **Check Provider Limits**:
   - Gmail: 500 emails/day (free), 2000/day (Workspace)
   - SendGrid: 100/day (free tier)
   - AWS SES: Sandbox mode limits to verified emails only

3. **Check Spam Folder**: Emails may be filtered as spam, especially if using a new domain or sending from a free email service.

4. **Check Email Provider Logs**:
   - SendGrid: **Activity** dashboard shows sent, bounces, and errors
   - AWS SES: **Sending statistics** and **Suppression list**
   - Mailgun: **Logs** section shows delivery status

5. **Verify Sender Identity**: Most providers require sender verification. Ensure your `SMTP_FROM_ADDRESS` is verified.

### Common Error Messages

**"Error al enviar correo: failed to send email: 535 Authentication failed"**
- Incorrect `SMTP_USER` or `SMTP_PASSWORD`
- For Gmail: Make sure you're using an App Password, not your regular password
- For SendGrid: Make sure `SMTP_USER` is exactly `apikey`

**"Error al enviar correo: failed to send email: connection refused"**
- Check `SMTP_HOST` and `SMTP_PORT` are correct
- Ensure your network/firewall allows SMTP connections
- Try different ports (587 vs 465)

**"Error al enviar correo: failed to send email: 550 Relaying denied"**
- Your `SMTP_FROM_ADDRESS` may not be verified
- Check sender authentication in your email provider dashboard

---

## Security Best Practices

1. **Never Commit Credentials**: Never commit `.env` files or credentials to version control
2. **Use App Passwords**: For Gmail, always use App Passwords, never your main account password
3. **Use Secrets Managers**: In production, use secrets managers (AWS Secrets Manager, HashiCorp Vault) instead of environment variables
4. **Rotate Credentials**: Regularly rotate SMTP passwords and API keys
5. **Limit Permissions**: When possible, use restricted API keys (SendGrid) instead of full access
6. **Monitor Usage**: Set up alerts for unusual sending patterns or quota limits

---

## Default Values

If SMTP configuration is not provided, the email service will use these defaults:

- `SMTP_HOST`: `smtp.gmail.com`
- `SMTP_PORT`: `587`
- `SMTP_USER`: (empty, emails will be logged but not sent)
- `SMTP_PASSWORD`: (empty, emails will be logged but not sent)
- `SMTP_FROM_ADDRESS`: `noreply@entorno35.com`
- `SMTP_FROM_NAME`: `Entorno35 - NOM-035`
- `APP_BASE_URL`: `http://localhost:3000`

**Important**: Without valid `SMTP_USER` and `SMTP_PASSWORD`, emails will only be logged to the console, not actually sent.

---

## Additional Resources

- [Gmail App Passwords](https://support.google.com/accounts/answer/185833)
- [SendGrid SMTP Documentation](https://docs.sendgrid.com/for-developers/sending-email/getting-started-smtp)
- [AWS SES SMTP Documentation](https://docs.aws.amazon.com/ses/latest/dg/send-email-smtp.html)
- [Mailgun SMTP Documentation](https://documentation.mailgun.com/en/latest/user_manual.html#sending-via-smtp)

---

**Last Updated**: 2026-01-12
