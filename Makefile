DUMMY_OBJ = AIDummy.o

# Add here any extra .o player files you want to link to the executable
EXTRA_OBJ = AItest7.o 

# Configuration
OPTIMIZE = 0 # Optimization level    (0 to 3)
DEBUG    = 1 # Compile for debugging (0 or 1)
PROFILE  = 0 # Compile for profile   (0 or 1)

# For debugging matches against Dummy
# OPTIMIZE = 0, DEBUG = 1

# For debugging matches among your players
# OPTIMIZE = 0, DEBUG = 1 and add -D_GLIBCXX_DEBUG at the end of DEBUG_FLAGS

# Flags
ifeq ($(strip $(PROFILE)),1)
	PROFILEFLAGS=-pg
endif

ifeq ($(strip $(DEBUG)),1)
	DEBUGFLAGS=-g -O0 -fno-inline #-D_GLIBCXX_DEBUG 
endif

CXXFLAGS = -std=c++11 -Wall -Wno-unused-variable -fPIC $(PROFILEFLAGS) $(DEBUGFLAGS) -O$(strip $(OPTIMIZE)) 
LDFLAGS  = -std=c++11                            $(PROFILEFLAGS) $(DEBUGFLAGS) -O$(strip $(OPTIMIZE))


# The following two lines will detect all your players (files matching "AI*.cc")

PLAYERS_SRC = $(wildcard AI*.cc)
PLAYERS_OBJ = $(patsubst %.cc, %.o, $(PLAYERS_SRC)) $(EXTRA_OBJ) $(DUMMY_OBJ)

# Rules

OBJ = Structs.o Settings.o State.o Info.o Random.o Board.o Action.o Player.o Registry.o Utils.o 

all: Game

clean:
	rm -rf Game  *.o *.exe Makefile.deps

Game:  $(OBJ) Game.o Main.o $(PLAYERS_OBJ) 
	$(CXX) $^ -o $@ $(LDFLAGS)

SecGame: $(OBJ) SecGame.o SecMain.o
	$(CXX) $^ -o $@ $(LDFLAGS) -lrt

%.exe: %.o $(OBJ) SecGame.o SecMain.o 
	$(CXX) $^ -o $@ $(LDFLAGS) -lrt

Makefile.deps: *.cc
	$(CXX) $(CXXFLAGS) -MM *.cc > Makefile.deps

-include Makefile.deps

# ─── Web app targets ───────────────────────────────────
.PHONY: web-install web-dev web-build docker-build docker-up docker-down web-help

web-install: ## Install web frontend dependencies
	cd web/frontend && npm install

web-dev: ## Start Go backend + frontend dev servers
	cd web/backend-go && GAME_BIN=../../Game CONFIG_FILE=../backend/data/default.cnf DATA_DIR=../backend/data STATIC_DIR=../frontend/dist go run . &
	cd web/frontend && npm run dev

web-build: ## Build web frontend for production
	cd web/frontend && npm run build

docker-build: ## Build Docker image
	docker compose build

docker-up: ## Start Docker container
	docker compose up -d

docker-down: ## Stop Docker container
	docker compose down

web-help: ## Show web targets
	@grep -E '^(web-|docker-)[a-zA-Z_-]+:.*##' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*##"}; {printf "  \033[36m%-16s\033[0m %s\n", $$1, $$2}'
