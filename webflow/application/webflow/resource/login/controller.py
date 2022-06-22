# modules
from webflow.resource._shared.controller import Controller
from webflow.resource._shared.exception import VaultConnectionError
from webflow.resource._shared.exception import VaultForbiddenError
from webflow.resource._shared.exception import VaultInvalidRequestError


class VaultLoginWithGoogleAccHandler(Controller):
    def post(self):
        try:
            google_oauth_code = request.form.get('code')
            vault_role = request.form.get('role')

            vault_token = self.service.vault.get_vault_token(google_oauth_code, vault_role)
            return self.render_page('token.html', token=vault_token)

        except VaultConnectionError:
            return self.error.connection_error

        except VaultForbiddenError:
            return self.error.forbidden

        except VaultInvalidRequestError:
            return self.error.invalid_request

        finally:
            v.logout()
