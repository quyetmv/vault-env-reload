# Project Name

Service Interval Pull from Hashicorp Vault to Environment Variables

## Description

This project consists of a systemd service that runs at an interval of 10 seconds. It checks for changes in the version of secrets stored in HashiCorp Vault and pulls the corresponding JSON from the Vault. The JSON data is then parsed and written to environment variables, which are used by the system. The configuration for the service is read from a JSON file.

## Installation

### Ubuntu

1. Download the `vault-env-reload.deb` package.
2. Install the package using the following command: sudo dpkg -i vault-env-reload.deb
3. Update the configuration file `/etc/vault-env-reload/config/vault.json` with your desired settings.
4. Restart the service using the following command: sudo systemctl restart vault-env-reload
### Other Linux Distributions

1. Download the `vault-env-reload` binary file.
2. Place the binary file in `/etc/vault-env-reload/bin/vault-env-reload`.
3. Create a configuration file `vault.json` using the provided `vault.json.sample` file as a reference. Place it in `/etc/vault-env-reload/config/vault.json`.
4. Create a systemd service for the `vault-env-reload` binary and start it.

## Support

If you find this project helpful or valuable, please consider supporting its development by buying me a coffee!

<a href="https://www.buymeacoffee.com/quyetmv" target="_blank"><img src="https://cdn.buymeacoffee.com/buttons/v2/default-yellow.png" alt="Buy Me A Coffee" width="180" height="45" ></a>

## Contact

For any questions or inquiries, please contact me:

- Email: [<img src="https://upload.wikimedia.org/wikipedia/commons/7/7e/Gmail_icon_%282020%29.svg" alt="Email" height="15" width="15"> quyetmv@gmail.com](mailto:quyetmv@gmail.com)
- Telegram: [<img src="https://upload.wikimedia.org/wikipedia/commons/8/82/Telegram_logo.svg" alt="Telegram" height="15" width="15"> quyetmv](https://t.me/quyetmv)
