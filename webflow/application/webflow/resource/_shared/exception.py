# 3rd-party
import hvac.exceptions
import requests.exceptions


VaultConnectionError = requests.exceptions.ConnectionError

VaultForbiddenError = hvac.exceptions.Forbidden

VaultInvalidRequestError = hvac.exceptions.InvalidRequest
