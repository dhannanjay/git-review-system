FROM node:22-bookworm

# Install system dependencies (git, Wails Linux deps, gh) — not apt golang-go (too old for Wails v2)
RUN apt-get update && apt-get install -y \
  git \
  curl \
  jq \
  build-essential \
  libgtk-3-dev \
  libwebkit2gtk-4.1-dev \
  && rm -rf /var/lib/apt/lists/*

# Go 1.23+ (required by wails v2.12+)
ARG GO_VERSION=1.23.6
RUN ARCH="$(dpkg --print-architecture)" \
  && case "$ARCH" in amd64) GOARCH=amd64 ;; arm64) GOARCH=arm64 ;; *) GOARCH=amd64 ;; esac \
  && curl -fsSL "https://go.dev/dl/go${GO_VERSION}.linux-${GOARCH}.tar.gz" \
  | tar -C /usr/local -xz
ENV PATH="/usr/local/go/bin:${PATH}"

# Install GitHub CLI
RUN curl -fsSL https://cli.github.com/packages/githubcli-archive-keyring.gpg \
  | dd of=/usr/share/keyrings/githubcli-archive-keyring.gpg \
  && echo "deb [arch=$(dpkg --print-architecture) signed-by=/usr/share/keyrings/githubcli-archive-keyring.gpg] https://cli.github.com/packages stable main" \
  | tee /etc/apt/sources.list.d/github-cli.list > /dev/null \
  && apt-get update && apt-get install -y gh \
  && rm -rf /var/lib/apt/lists/*

# Build-args for UID/GID alignment: sandcastle docker build-image
ARG AGENT_UID=1000
ARG AGENT_GID=1000

# Rename the base image's "node" user to "agent" and align UID/GID.
RUN groupmod -g $AGENT_GID node && usermod -u $AGENT_UID -g $AGENT_GID -d /home/agent -m -l agent node

# CLI tools on PATH for all users (root install phase)
RUN npm install -g opencode-ai@latest \
  && GOBIN=/usr/local/bin go install github.com/wailsapp/wails/v2/cmd/wails@v2.12.0

USER ${AGENT_UID}:${AGENT_GID}

WORKDIR /home/agent

# Sandcastle bind-mounts the git worktree at ${SANDBOX_REPO_DIR} at container start.
ENTRYPOINT ["sleep", "infinity"]
