# standard
import os

# 3rd-party
from hvac import Client as VaultClient


class Vault:
    def __init__(self):
        auth_path = os.environ.get('VAULT_AUTH_PATH')
        endpoint = os.environ.get('VAULT_SERVER_ENDPOINT')
        token = os.environ.get('VAULT_AUTH_TOKEN')

        if not (auth_path and endpoint and token):
            raise Exception('Missing Vault configuration data')

        self._auth_path = auth_path
        self._client = VaultClient(url=endpoint, token=token)

    def get_google_oauth_url(self):
        return self._client.read(f'auth/{self._auth_path}/code_url')

    def get_vault_roles(self):
        return self._client.auth.approle.list_roles(mount_point=self._auth_path)

    def get_vault_token(self, code, role):
        res = self._client.write(f'auth/{self._auth_path}/login', code=code, role=role)
        return res['auth']['client_token']
