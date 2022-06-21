from util import vault_cmd


def get_unseal_key(line: str) -> str:
    segs = line.split(' ')
    if len(segs) == 0: raise Exception('Malformed output from "operator init"')

    key = segs[-1]
    key = key.strip()
    if len(key) != 44: raise Exception('Unseal key with 44 characters expected')

    return key


def vault_init() -> None:
    print('')
    print(' > started')

    # ...
    init_output = vault_cmd('operator init')
    init_output_lines = init_output.split('\n')
    print(' > vault initialized')

    # ...
    unseal_keys = [get_unseal_key(line) for line in init_output_lines[0:5]]
    if len(unseal_keys) != 5: raise Exception('Five (5) unseal keys expected')
    print(' > got unseal keys')

    # ...
    initial_root_token = init_output_lines[6].replace('Initial Root Token: ', '').strip()
    if len(initial_root_token) != 26: raise Exception('Malformed initial root token')
    print(' > got initial root key')

    # ...
    for unseal_key in unseal_keys[0:3]: vault_cmd('operator unseal ' + unseal_key)
    print(' > vault unsealed')

    # ...
    vault_cmd('login ' + initial_root_token)
    print(' > logged successfully')

    # ...
    print('')
    print(' ! initial root token: ' + initial_root_token)

    print(' ! unseal keys:')
    for idx, unseal_key in enumerate(unseal_keys): print(f'   {idx + 1}. {unseal_key}')

    # ...
    print('')
    print(' > done')


vault_init()
