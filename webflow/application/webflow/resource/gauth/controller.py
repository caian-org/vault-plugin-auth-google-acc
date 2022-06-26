# modules
from webflow.resource._shared.controller import Controller
from webflow.resource._shared.exception import VaultConnectionError
from webflow.resource._shared.exception import VaultInvalidRequestError


class VaultGetGoogleOAuthLinkHandler(Controller):
    def get(self):
        try:
            res = self.service.vault.get_google_oauth_url()
            url = res['data']['url']

            return self.redirect_to(url)

        except VaultConnectionError:
            return self.error.connection_error('A001')

        except VaultInvalidRequestError:
            # plugin is enabled but the configuration has not been written
            return self.error.incorrectly_configured('A002')

        except (TypeError, KeyError):
            # plugin is probably not enabled/mounted
            return self.error.incorrectly_configured('A003')
