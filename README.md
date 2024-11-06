# wazuh-alert-notifier

A Go-based script for Wazuh that sends plain-text alerts to the Green API, allowing notifications via services like WhatsApp and other messaging platforms. This notifier integrates seamlessly with Wazuh to forward alert data in real-time.

## Installation

### Step 1: Download the Repository

First, clone the repository from GitHub:

```bash
git clone https://github.com/herbertpk/wazuh-alert-notifier.git
```

Navigate to the cloned repository:

```bash
cd wazuh-alert-notifier
```

### Step 2: Compile the Go Script

Ensure Go is installed on your system. Verify the installation with:

```bash
go version
```

If not installed, follow Go's official installation instructions.

Compile the Go script into a binary:

```bash
go build -o wazuh-alert-notifier main.go
```

This will generate an executable file named `wazuh-alert-notifier` in the current directory.

### Step 3: Move the Binary to Wazuh’s Integrations Directory

Move the compiled binary to the Wazuh integrations folder so it can be accessed by Wazuh:

```bash
sudo mv wazuh-alert-notifier /var/ossec/integrations/
```

Verify the binary is in the correct directory:

```bash
ls /var/ossec/integrations/
```

You should see `wazuh-alert-notifier` listed.

### Step 4: Configure Wazuh to Use the Script with Green API Credentials

To connect to the Green API, you need to configure Wazuh’s integration settings. Green API requires an `apiURL`, `idInstance`, and `apiTokenInstance` to send messages. These credentials must be provided in the Wazuh configuration.

Open the Wazuh configuration file:

```bash
sudo nano /var/ossec/etc/ossec.conf
```

Add a new `<command>` and `<integration>` entry, replacing `https://api.green-api.com` with the Green API’s URL and providing the appropriate `idInstance`, `apiTokenInstance` and `chatId` as arguments.

```xml
<ossec_config>
    ...
    
    <command>
        <name>wazuh-alert-notifier</name>
        <executable>var/ossec/integrations/wazuh-alert-notifier</executable>
        <extra_args>https://api.green-api.com your_idInstance your_apiTokenInstance your_chatId</extra_args> 
        <expect>json</expect>
    </command>

    <integration>
        <name>wazuh-alert-notifier</name>
        <hook_url></hook_url>
        <alert_format>json</alert_format>
        <levels>3</levels>
        <rule_id></rule_id>
        <alert_fields></alert_fields>
    </integration>

    ...
</ossec_config>
```

Save and close the file. Ensure you replace `your_idInstance` and `your_apiTokenInstance` with your actual Green API credentials.

### Step 5: Restart Wazuh Manager

After saving the changes, restart the Wazuh manager to apply the configuration:

```bash
sudo systemctl restart wazuh-manager
```

### Step 6: Test the Configuration

To verify that Wazuh is correctly using the `wazuh-alert-notifier` script with the Green API:

1. Generate an alert in Wazuh to trigger the notifier.
2. Check the Wazuh logs to confirm the notifier executed as expected:

```bash
tail -f /var/ossec/logs/ossec.log
```

If everything is set up correctly, `wazuh-alert-notifier` will process the alert and send it via the Green API to the specified messaging service, such as WhatsApp. Check your messaging platform for the notification.