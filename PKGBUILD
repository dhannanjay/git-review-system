# Maintainer: Dhannanjay Raje Vaid <dhannanjay19@gmail.com>
# Contributor: review-diff developers
#
# To bump the version:
#   1. Update pkgver and pkgrel below
#   2. Run: updpkgsums PKGBUILD
#   3. Test: makepkg -fi
#   4. Commit the updated PKGBUILD
#
# Depends on: git (>=2.30), glibc
# This PKGBUILD downloads the pre-built Linux tarball from the official
# release — it does not build from source.  See also:
#   https://github.com/dhannanjay/git-review-system

pkgname=review-diff
pkgver=0.1.0
pkgrel=1
pkgdesc="A native desktop diff viewer for local Git branches"
arch=('x86_64')
url="https://github.com/dhannanjay/git-review-system"
license=('MIT')
depends=('git>=2.30' 'glibc')
source=("$url/releases/download/v$pkgver/review-diff-$pkgver-linux-amd64.tar.gz")
sha256sums=('SKIP')
# ^ Replace SKIP with the actual SHA-256 checksum after publishing the release.
#   Run: curl -sL https://github.com/dhannanjay/git-review-system/releases/download/v$pkgver/review-diff-$pkgver-linux-amd64.tar.gz | sha256sum

package() {
  install -Dm755 "$srcdir/review-diff" "$pkgdir/usr/bin/review-diff"
}
