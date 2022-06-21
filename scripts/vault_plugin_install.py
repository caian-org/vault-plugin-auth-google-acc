from util import vault_cmd
from util import docker_cmd


PLUGIN_FILE = 'vault-plugin-auth-google-acc'
PLUGIN_INTERNAL_PATH = 'googleacc'


def get_shasum(output: str) -> str:
    segs = output.split(' ')
    shasum = segs[0].strip()
    if len(shasum) != 64: raise Exception('Malformed file hash sum string')

    return shasum


def vault_plugin_install():
    print('')
    print(' > started')

    # ...
    plugin_file_shasum = get_shasum(docker_cmd(f'sha256sum /etc/vault/plugins/{PLUGIN_FILE}'))
    print(' > plugin file shasum calculated')

    # ...
    plugin_catalog_cmd = [
        'write',
        f'sys/plugins/catalog/auth/{PLUGIN_FILE}',
        f'sha_256="{plugin_file_shasum}"',
        f'command="{PLUGIN_FILE}"'
    ]

    vault_cmd(' '.join(plugin_catalog_cmd))
    print(' > plugin added to catalog')

    # ...
    plugin_activation_cmd = [
        'auth enable',
        f'-path="{PLUGIN_INTERNAL_PATH}"',
        f'-plugin-name="{PLUGIN_FILE}"',
        'plugin'
    ]

    vault_cmd(' '.join(plugin_activation_cmd))
    print(' > plugin activated')

    # ...
    print('')
    print(' > done')


vault_plugin_install()
