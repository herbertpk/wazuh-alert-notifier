
# wazuh-alert-notifier API
A Go-based API for Wazuh that receives alert data and forwards it to the Green API, enabling real-time notifications through messaging platforms like WhatsApp. This API integrates with Wazuh as an HTTP endpoint, listening for incoming alerts and automatically relaying them to the Green API for seamless notification delivery. This setup enhances alerting capabilities, making it easy to keep teams informed of security events as they happen.

## Installation
Clone the Repository and Compile the API:

```bash
git clone https://github.com/herbertpk/wazuh-alert-notifier.git
cd wazuh-alert-notifier
go build -o wazuh-alert-api main.go
```

Move the Binary to a System-Wide Location:

```bash
sudo mv wazuh-alert-api /usr/local/bin/
```

Set Up Environment Variables: Create a `.env` file to store Green API credentials and configuration values:

```bash
sudo nano /etc/wazuh-alert-api.env
```

Add the following configuration, replacing placeholders with your actual Green API details:

```bash
GREEN_API_URL="https://7103.api.greenapi.com"
GREEN_API_INSTANCE_ID="your_idInstance"
GREEN_API_TOKEN="your_apiTokenInstance"
GREEN_API_CHAT_ID="your_chatId"
```

### Step 2: Create a Systemd Service for the API Server

Create a Systemd Service File:

```bash
sudo nano /etc/systemd/system/wazuh-alert-api.service
```

Add the following content:

```ini
[Unit]
Description=Wazuh Alert Notifier API
After=network.target

[Service]
EnvironmentFile=/etc/wazuh-alert-api.env
ExecStart=/usr/local/bin/wazuh-alert-api
Restart=always
User=wazuh  # Adjust this to the user running the Wazuh service
WorkingDirectory=/usr/local/bin
StandardOutput=syslog
StandardError=syslog
SyslogIdentifier=wazuh-alert-api

[Install]
WantedBy=multi-user.target
```

Enable and Start the Service:

```bash
sudo systemctl daemon-reload
sudo systemctl enable wazuh-alert-api
sudo systemctl start wazuh-alert-api
```

Check Service Status:

```bash
sudo systemctl status wazuh-alert-api
```

Confirm that the service is running correctly.

### Step 3: Configure Wazuh Integration with the Notifier API

Edit the Wazuh Configuration File:

Open the `ossec.conf` file to configure Wazuh to use the notifier API as an integration.

```bash
sudo nano /var/ossec/etc/ossec.conf
```

Add the Integration Configuration:

Add a `<command>` and `<integration>` entry for `wazuh-alert-notifier` in the configuration file. Replace `http://localhost:8080/alert` with the endpoint where the notifier API server is listening.

```xml
<ossec_config>
    ...
    <!-- Configure the integration trigger -->
    <integration>
        <name>wazuh-alert-notifier</name>
        <hook_url>http://localhost:8080/alert</hook_url>
        <alert_format>json</alert_format>
        <level>3</level> <!-- Set to trigger on alerts with level 3 or higher -->
    </integration>

    ...
</ossec_config>
```
Save and Close the Configuration File.

Restart Wazuh Manager:

```bash
sudo systemctl restart wazuh-manager
```

### Step 4: Testing the Integration

Generate a Test Alert:

Use the `ossec-logtest` tool to generate a test alert and trigger the integration:

```bash
sudo /var/ossec/bin/ossec-logtest
```

In the interactive prompt, type a log line that would trigger an alert (e.g., with a severity level of 3 or higher):

```bash
2024/11/06 wazuh: test alert - severity level 3
```

Check Wazuh Logs:

Monitor Wazuh’s logs to confirm that the `wazuh-alert-notifier` command was triggered:

```bash
tail -f /var/ossec/logs/ossec.log
```

Check API Logs:

Monitor the logs for the `wazuh-alert-api` service to verify that it received the alert data and processed it correctly:

```bash
sudo journalctl -u wazuh-alert-api -f
```

Verify Alert Delivery:

Check the destination (such as the Green API or a messaging service like WhatsApp) to confirm that the alert message was successfully forwarded.
