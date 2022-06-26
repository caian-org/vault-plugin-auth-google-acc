# modules
from webflow.resource._shared.controller import Controller
from webflow.resource._shared.exception import VaultConnectionError
from webflow.resource._shared.exception import VaultForbiddenError
from webflow.resource._shared.exception import VaultInvalidPathError


class VaultDisplayRolesHandler(Controller):
    def get(self):
        try:
            res = self.service.vault.get_vault_roles()
            roles = res['data']['keys']

            return self.render_page('roles.html', roles=roles)

        except VaultConnectionError:
            return self.error.connection_error('R001')

        except VaultInvalidPathError:
            # there is no configured role for the google auth plugin
            return self.error.incorrectly_configured('R002')

        except VaultForbiddenError:
            # the token used by this service don't have enough permissions
            return self.error.incorrectly_configured('R003')
