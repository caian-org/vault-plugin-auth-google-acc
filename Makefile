include plugin/.env

.DEFAULT_GOAL := start


start:
	docker compose \
		--env-file ./plugin/.env \
		up --build \
		--exit-code-from vault-server

vault-init:
	python3 scripts/vault_init.py "$(VAULT_PLUGIN_NAME)"

vault-plugin-install:
	python3 scripts/vault_plugin_install.py "$(VAULT_PLUGIN_NAME)"
