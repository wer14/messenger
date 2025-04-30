SERVICES := auth profile messaging gateway contact session

.PHONY: all build apply restart logs

all: build apply

build:
	for svc in $(SERVICES); do \
		$(MAKE) -C services/$$svc build; \
	done

apply:
	for svc in $(SERVICES); do \
		$(MAKE) -C services/$$svc apply; \
	done

restart:
	for svc in $(SERVICES); do \
		$(MAKE) -C services/$$svc restart; \
	done

logs:
	for svc in $(SERVICES); do \
		$(MAKE) -C services/$$svc logs; \
	done
