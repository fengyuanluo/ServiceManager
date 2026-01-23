#!/usr/bin/env bash
set -euo pipefail

APP_NAME="${SM_APP_NAME:-service-manager}"
REPO="${SM_REPO:-fengyuanluo/ServiceManager}"
INSTALL_DIR="${SM_INSTALL_DIR:-/usr/local/bin}"
VERSION="latest"
AUTO_YES="false"

usage() {
  cat <<EOF
usage: install.sh [--version <tag>] [--install-dir <path>] [--yes]

options:
  --version       Git tag to install (default: latest)
  --install-dir   Install directory (default: /usr/local/bin)
  --yes           Skip confirmation prompt
  -h, --help      Show this help

env:
  SM_APP_NAME     Override app name (default: service-manager)
  SM_REPO         Override repo (default: fengyuanluo/ServiceManager)
  SM_INSTALL_DIR  Override install dir (default: /usr/local/bin)
EOF
}

die() {
  echo "error: $1" >&2
  exit 1
}

need_cmd() {
  command -v "$1" >/dev/null 2>&1 || die "missing required command: $1"
}

while [[ $# -gt 0 ]]; do
  case "$1" in
    -h|--help)
      usage
      exit 0
      ;;
    --version)
      VERSION="${2:-}"
      shift 2
      ;;
    --install-dir)
      INSTALL_DIR="${2:-}"
      shift 2
      ;;
    --yes)
      AUTO_YES="true"
      shift
      ;;
    *)
      die "unknown option: $1"
      ;;
  esac
done

if [[ -z "${VERSION}" ]]; then
  die "--version requires a value"
fi

need_cmd uname
need_cmd curl
need_cmd tar
need_cmd install

OS="$(uname -s | tr '[:upper:]' '[:lower:]')"
case "${OS}" in
  linux|darwin) ;;
  *) die "unsupported OS: ${OS}" ;;
esac

ARCH_RAW="$(uname -m)"
case "${ARCH_RAW}" in
  x86_64|amd64) ARCH="amd64" ;;
  aarch64|arm64) ARCH="arm64" ;;
  armv7l|armv7) ARCH="armv7" ;;
  *) die "unsupported arch: ${ARCH_RAW}" ;;
esac

if [[ "${VERSION}" == "latest" ]]; then
  TAG="$(curl -fsSL "https://api.github.com/repos/${REPO}/releases/latest" \
    | grep -m1 '"tag_name"' \
    | sed -E 's/.*"tag_name": "([^"]+)".*/\1/')"
  [[ -n "${TAG}" ]] || die "failed to resolve latest version"
  VERSION="${TAG}"
fi

ARCHIVE="${APP_NAME}_${VERSION}_${OS}_${ARCH}.tar.gz"
URL="https://github.com/${REPO}/releases/download/${VERSION}/${ARCHIVE}"

if [[ "${AUTO_YES}" != "true" ]]; then
  echo "install ${APP_NAME} ${VERSION} to ${INSTALL_DIR}? [y/N]"
  read -r confirm
  case "${confirm}" in
    y|Y|yes|YES) ;;
    *) echo "cancelled"; exit 0 ;;
  esac
fi

TMP_DIR="$(mktemp -d)"
trap 'rm -rf "${TMP_DIR}"' EXIT

curl -fSL "${URL}" -o "${TMP_DIR}/${ARCHIVE}" || die "download failed: ${URL}"
tar -xzf "${TMP_DIR}/${ARCHIVE}" -C "${TMP_DIR}"

BIN_PATH="${TMP_DIR}/${APP_NAME}"
[[ -f "${BIN_PATH}" ]] || die "binary not found in archive"

SUDO=""
if [[ ! -w "${INSTALL_DIR}" ]]; then
  if command -v sudo >/dev/null 2>&1; then
    SUDO="sudo"
  else
    die "install dir not writable and sudo not available"
  fi
fi

${SUDO} install -m 0755 "${BIN_PATH}" "${INSTALL_DIR}/${APP_NAME}"
echo "installed ${APP_NAME} to ${INSTALL_DIR}/${APP_NAME}"
