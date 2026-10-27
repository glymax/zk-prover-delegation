#!/usr/bin/env bash
set -euo pipefail

# Run on an already allocated experiment node, not on the management node.
# No PoS operations, reboot, disk formatting, or project dependency installation.
readonly GO_VERSION=1.27.2
readonly GO_DIR=/usr/local/go
readonly PROFILE_FILE=/etc/profile.d/bachelor-go.sh
readonly -a PACKAGES=(
    git
    ca-certificates
    curl
)

die() { printf 'Error: %s\n' "$*" >&2; exit 1; }

usage() {
    cat <<EOF
Usage: bash scripts/setup-node.sh [--dry-run | --apply]

Run this on the reserved and allocated experiment node as root for --apply.
The default is --dry-run: no network or writes.

Installs git, ca-certificates, curl and Go ${GO_VERSION}.
Assumes an Ubuntu 22.04 amd64 node; platform is not checked.
Does not check calendar/PoS ownership; verify those on coinbase first.
Does not reset/free nodes, touch NVMe, install gnark, or run experiments.
EOF
}

profile_content() {
    printf '# Managed by bachelor implementation/scripts/setup-node.sh\n'
    printf 'export PATH="%s/bin:$PATH"\n' "$GO_DIR"
    # Keep Go from silently downloading a different toolchain for a module.
    printf 'export GOTOOLCHAIN=local\n'
}

check_go() {
    [[ -x "$1/bin/go" ]] &&
        [[ $(GOTOOLCHAIN=local "$1/bin/go" version) == "go version go${GO_VERSION} linux/amd64" ]]
}

main() {
    (( $# <= 1 )) || die 'Use only one option: --dry-run or --apply.'
    local mode=${1:---dry-run}
    case $mode in
        --dry-run|--apply) ;;
        -h|--help) usage; return 0 ;;
        *) die "Unknown argument: $mode" ;;
    esac

    if [[ $mode == --dry-run ]]; then
        printf 'Ensure packages: %s\n' "${PACKAGES[*]}"
        printf 'Ensure Go %s in %s\n' "$GO_VERSION" "$GO_DIR"
        printf 'Write shell environment: %s\n' "$PROFILE_FILE"
        printf 'No changes made; use --apply to install.\n'
        return 0
    fi
    (( EUID == 0 )) || die '--apply requires root; use sudo or the node root shell.'

    apt-get update
    DEBIAN_FRONTEND=noninteractive apt-get install -y --no-install-recommends \
        --no-upgrade "${PACKAGES[@]}"

    # Reuse the requested Go version on repeated runs.
    if ! check_go "$GO_DIR"; then
        local archive="/usr/local/go${GO_VERSION}.linux-amd64.tar.gz"
        curl --fail --show-error --location --proto '=https' --proto-redir '=https' \
            --connect-timeout 30 --max-time 600 --retry 3 \
            "https://go.dev/dl/go${GO_VERSION}.linux-amd64.tar.gz" --output "$archive"
        # The archive contains go/: extract into its parent after removing old files.
        rm -rf -- "$GO_DIR"
        tar -xzf "$archive" -C /usr/local --no-same-owner
        check_go "$GO_DIR" || die 'Installed Go failed the version check.'
        rm -f -- "$archive"
    fi

    profile_content > "$PROFILE_FILE"
    chmod 0644 "$PROFILE_FILE"

    printf 'Setup complete.\n'
    printf 'For this shell run: source %s\n' "$PROFILE_FILE"
    printf 'For non-login SSH commands use: %s/bin/go\n' "$GO_DIR"
}

if [[ ${BASH_SOURCE[0]} == "$0" ]]; then
    main "$@"
fi
