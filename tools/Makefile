########################################
### GLIDE

GLIDE = github.com/Masterminds/glide
GLIDE_CHECK := $(shell command -v glide 2> /dev/null)

check_tools:
ifndef GLIDE_CHECK
	@echo "No glide in path.  Install with 'make get_tools'."
else
	@echo "Found glide in path."
endif

get_tools:
ifdef GLIDE_CHECK
	@echo "glide is already installed.  Run 'make update_tools' to update."
else
	@echo "$(ansi_grn)Installing glide$(ansi_end)"
	go get -v $(GLIDE)
endif

update_tools:
	@echo "$(ansi_grn)Updating glide$(ansi_end)"
	go get -u -v $(GLIDE)


########################################
# ANSI colors

ansi_red=\033[0;31m
ansi_grn=\033[0;32m
ansi_yel=\033[0;33m
ansi_end=\033[0m