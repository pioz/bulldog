#!/bin/bash

ARCHS="i386 amd64"
if ! echo $ARCHS | grep -w "$ARCH" > /dev/null ; then
  echo "You need to set a valid $ARCH variable (i386|amd64)"
  exit 2
fi

DEBEMAIL="epilotto@gmx.com"
DEBFULLNAME="Enrico Pilotto"
export DEBEMAIL DEBFULLNAME

NAME="bulldog"
VERSION=`cat VERSION`
PKGDIR=$NAME-$VERSION

rm -rf $PKGDIR
rm -rf ${NAME}_$VERSION*

mkdir -p $PKGDIR/usr/bin
cp ${NAME}_$ARCH $PKGDIR/usr/bin/$NAME
cp -r etc $PKGDIR
cd $PKGDIR
dh_make -s --createorig --copyright=gpl3 -y
cd ..
rm $PKGDIR/debian/README*
cp $NAME.1 $NAME.manpages changelog copyright install postinst postrm VERSION $PKGDIR/debian
cp control_$ARCH $PKGDIR/debian/control
cd $PKGDIR
dpkg-buildpackage -us -uc
cd ..
debc ${NAME}_$VERSION-1_$ARCH.changes
lintian ${NAME}_$VERSION-1_$ARCH.deb

rm -rf ${NAME}-$VERSION
rm ${NAME}_$VERSION-1.debian.tar.xz ${NAME}_$VERSION-1.dsc ${NAME}_$VERSION-1_amd64.buildinfo ${NAME}_$VERSION-1_amd64.changes ${NAME}_$VERSION.orig.tar.xz
