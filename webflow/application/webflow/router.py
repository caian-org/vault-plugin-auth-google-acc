# modules
from webflow.resource.gauth.controller import VaultGetGoogleOAuthLinkHandler
from webflow.resource.login.controller import VaultLoginWithGoogleAccHandler
from webflow.resource.roles.controller import VaultDisplayRolesHandler


def init_router(api) -> None:
    api.add_resource(VaultGetGoogleOAuthLinkHandler, '/')
    api.add_resource(VaultDisplayRolesHandler, '/roles')
    api.add_resource(VaultLoginWithGoogleAccHandler, '/login')
