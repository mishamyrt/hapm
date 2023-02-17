VERSION = 0.0.3
DIST_PATH = ./dist
VENV_PATH = ./venv
VENV = . $(VENV_PATH)/bin/activate;

.PHONY: publish
publish: clean $(DIST_PATH)
	git tag "v$(VERSION)"
	git push --tags
	$(VENV) python3 -m twine upload --repository pypi dist/* -umishamyrt

.PHONY: clean
clean:
	rm -rf *.egg-info
	rm -rf build
	rm -rf dist

.PHONY: build
build: $(DIST_PATH)

.PHONY: install
install: $(DIST_PATH)
	pip3 install .

.PHONY: install-venv
install-venv: $(DIST_PATH)
	$(VENV) pip install .

.PHONY: lint
lint:
	$(VENV) pylint $(SRC)

configure: requirements.txt
	rm -rf $(VENV_PATH)
	make $(VENV_PATH)

$(VENV_PATH):
	python3 -m venv $(VENV_PATH)
	$(VENV) pip install -r requirements.txt

$(CONFIG_PATH): config.json
	mkdir -p $(CONFIG_DIR)
	rm -f $(CONFIG_PATH)
	cp config.json $(CONFIG_PATH)

$(DIST_PATH):
	echo $(VERSION) > .version
	$(VENV) python setup.py sdist bdist_wheel