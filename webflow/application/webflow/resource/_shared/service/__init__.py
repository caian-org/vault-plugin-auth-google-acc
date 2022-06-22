# modules
from .vault import Vault


class Service:
    def __init__(self):
        self.vault = Vault()
