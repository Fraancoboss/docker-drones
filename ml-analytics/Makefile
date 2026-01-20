PYTHON ?= python3
VENV := .venv
BIN := $(VENV)/bin
PIP := $(BIN)/pip
PY := $(BIN)/python

.PHONY: check-python venv deps run clean

check-python:
	@command -v $(PYTHON) >/dev/null 2>&1 || (echo "PYTHON not found: $(PYTHON)"; exit 1)
	@$(PYTHON) -c "import sys; m=sys.version_info[:2]; \
		(0 if (m[0]==3 and 11<=m[1]<14) else (_ for _ in ()).throw(SystemExit('Python 3.11 or 3.12 required (pandas wheels).')))"

venv: check-python
	@$(PYTHON) -m venv $(VENV)
	@$(PIP) install -U pip

deps: venv
	@$(PIP) install -r requirements.txt

run: deps
	@$(PY) main.py

clean:
	rm -rf $(VENV)
