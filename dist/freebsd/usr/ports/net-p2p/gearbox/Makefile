PORTNAME=	gearbox
PORTVERSION=	0.0.1
CATEGORIES=	net-p2p

MAINTAINER=	vchimishuk@yandex.ru
COMMENT=	non-interactive console client for transmission-daemon
WWW=		https://github.com/vchimishuk/gearbox

LICENSE=	GPLv3+

FETCH_DEPENDS=	git:devel/git
USES=		go:1.21,modules

USE_GITHUB=	yes
GH_ACCOUNT=	vchimishuk
GH_PROJECT=	gearbox
GH_TAGNAME=	21697c4

# $ make gomod-vendor
GH_TUPLE=	pborzenkov:go-transmission:v0.3.0:pborzenkov_go_transmission/vendor/github.com/pborzenkov/go-transmission \
		vchimishuk:config:f2660fcdf53b:vchimishuk_config/vendor/github.com/vchimishuk/config \
		vchimishuk:opt:d823c9050e21:vchimishuk_opt/vendor/github.com/vchimishuk/opt

PLIST_FILES=	bin/gearbox \
		share/man/man1/gearbox.1.gz

post-install:
	${INSTALL_MAN} ${WRKSRC}/${PORTNAME}.1 ${STAGEDIR}${PREFIX}/share/man/man1

.include <bsd.port.mk>
