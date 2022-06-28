# modules
from webflow.resource._shared.controller import Controller
from webflow.resource._shared.exception import VaultConnectionError
from webflow.resource._shared.exception import VaultForbiddenError
from webflow.resource._shared.exception import VaultInvalidRequestError


class VaultLoginWithGoogleAccHandler(Controller):
    def post(self):
        try:
            google_oauth_code = self.fetch_form_val('code')
            vault_role = self.fetch_form_val('role')

            vault_token = self.service.vault.get_vault_token(google_oauth_code, vault_role)
            return self.render_page('token.html', token=vault_token)

        except VaultConnectionError as e:
            return self.error.connection_error('L001', e)

        except VaultForbiddenError as e:
            return self.error.forbidden('L002', e)

        except VaultInvalidRequestError as e:
            return self.error.invalid_request('L003', e)
