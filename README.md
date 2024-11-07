# Wazuh Alert Notifier

## Project Description

The Wazuh Alert Notifier is a Go-based integration tool designed to forward Wazuh alerts to WhatsApp using the Green API. This integration allows real-time notifications of security alerts by automatically sending messages to a specified WhatsApp chat or group whenever a significant event is detected by Wazuh.

## Key Features

- **Automated Alerting**: Receives Wazuh alert data and sends it as WhatsApp messages to specified chat IDs.
- **Configurable Integration**: Filters alerts based on severity level and specific Wazuh rule groups.
- **Seamless Green API Integration**: Uses the Green API to communicate with WhatsApp, ensuring timely delivery of critical alerts.

## Requirements

To use this integration:

- **Green API Account**: You must have an account with Green API and a configured instance to use this integration.
- **Instance Setup**: Obtain your instance ID, API token, and chat ID from Green API. These will be used to authenticate and specify the message recipient for alerts.

For more details on setting up your Green API instance and obtaining necessary credentials, please refer to the Green API documentation.

With Wazuh Alert Notifier, you'll stay informed of potential security issues directly through WhatsApp, providing immediate visibility into critical alerts.

## Installation Guide

### Step 1: Ensure Go is Installed

#### Check if Go is Installed:

Run the following command to verify if Go is installed:

```bash
go version
```

If Go is installed, you should see the version information (e.g., `go version go1.17.2 linux/amd64`). If it’s not installed, download Go from the [official Go website](https://golang.org/dl/) and follow the installation instructions for your operating system.

### Step 2: Clone the Repository and Build the Binary

#### Clone the Repository:

```bash
git clone https://github.com/herbertpk/wazuh-alert-notifier.git
cd wazuh-alert-notifier
```

#### Build the Binary:

```bash
go build -o custom-whatsapp-notifier whatsapp-alert-notifier.go
```

#### Move the Binary to the Wazuh Integrations Directory:

```bash
sudo mv custom-whatsapp-notifier /var/ossec/integrations/
```

#### Set Permissions:

```bash
sudo chmod 750 /var/ossec/integrations/custom-whatsapp-notifier
sudo chown root:wazuh /var/ossec/integrations/custom-whatsapp-notifier
```

### Step 3: Configure Wazuh to Use the Notifier

#### Edit Wazuh Configuration:

Open the `ossec.conf` file on the Wazuh manager:

```bash
sudo nano /var/ossec/etc/ossec.conf
```

#### Add the Integration Block:

```xml
<integration>
    <name>custom-whatsapp-notifier</name>
    <hook_url>https://{your_apiUrl}.com/waInstance{your_InstanceId}/sendMessage/{your_apiToken}</hook_url> <!-- Replace with your API endpoint -->
    <api_key>your_chat_id</api_key> <!-- Replace with your chat ID -->
    <level>3</level> <!-- Filter for alerts of level 3 or higher -->
    <alert_format>json</alert_format>
</integration>
```

In this configuration, the `<api_key>` field is being used to pass the chat ID due to a requirement by the WhatsApp API (or similar messaging API) that the chat ID be sent separately in the body of the request, rather than as part of the URL or headers.

#### Integration Configuration Breakdown

- `<hook_url>`: Specifies the URL of the external API endpoint where alerts will be sent.
- `<api_key>`: Holds the chat ID needed by the API, which is included in the body of the request to specify the message recipient.
- `<level>`: Defines the minimum alert level required to trigger the notifier, filtering out lower-severity alerts.
- `<alert_format>`: Set to `json` for structured alerts, making them easier to parse and process in the integration script.

### Step 4: Restart Wazuh Manager

Apply the configuration changes by restarting the Wazuh manager:

```bash
sudo systemctl restart wazuh-manager
```

### Step 5: Test the Integration

#### Generate a Test Alert:

Use `ossec-logtest` to simulate an alert:

```bash
sudo /var/ossec/bin/ossec-logtest
```

#### Confirm Alert Delivery:

Check the `ossec.log` file for confirmation that the alert was processed and sent by the notifier:

```bash
tail -f /var/ossec/logs/ossec.log
```

This setup will allow Wazuh to forward specific alerts to an external API using `custom-whatsapp-notifier`. Adjust configuration as needed based on your alerting requirements.