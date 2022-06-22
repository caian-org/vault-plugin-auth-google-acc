# modules
from webflow.resource._shared.controller import Controller
from webflow.resource._shared.exception import VaultConnectionError
from webflow.resource._shared.exception import VaultForbiddenError


class VaultDisplayRolesHandler(Controller):
    def get(self):
        try:
            vault_roles = self.service.vault.get_vault_roles()
            return self.render_page('roles.html', roles=vault_roles)

        except VaultConnectionError:
            return self.error.connection_error

        except VaultForbiddenError:
            return self.error.forbidden

        finally:
            v.logout()
