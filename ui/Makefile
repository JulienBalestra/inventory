NPM=npm install
NODE_MODULES=node_modules

PROXY_EXISTS=$(shell [ $PROXY ] && echo 1 || echo 0 )

ifeq ($(PROXY_EXISTS), 1)
	NPM += --proxy $(PROXY)
endif

npm:
	@$(NPM)

fclean:
	rm -Rf $(NODE_MODULES)