from webflow.resource._shared.controller import Controller
from webflow.resource._shared.exception import VaultConnectionError


class VaultGetGoogleOAuthLinkHandler(Controller):
    def get(self):
        try:
            oauth_url = self.service.vault.get_google_oauth_url()
            return self.redirect_to(oauth_url)

        except VaultConnectionError:
            return self.error.connection_error

        finally:
            v.logout()
