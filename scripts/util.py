import sys
from subprocess import getstatusoutput as run


def docker_cmd(args: str) -> str:
    cmd = f'docker exec {sys.argv[1]} {args}'
    exit_code, output = run(cmd)
    if exit_code != 0:
        raise Exception(f'Command "{cmd}" failed with code {exit_code}; got: \n\n{output}')

    return output


def vault_cmd(args: str) -> str:
    return docker_cmd(f'vault {args}')
