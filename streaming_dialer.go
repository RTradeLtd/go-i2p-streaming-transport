package i2pStreaming

import (
	"context"
	"crypto/rand"
	i2pma "github.com/eyedeekay/sam3-multiaddr"

	crypto "github.com/libp2p/go-libp2p-crypto"
	net "github.com/libp2p/go-libp2p-net"
	peer "github.com/libp2p/go-libp2p-peer"
	peerstore "github.com/libp2p/go-libp2p-peerstore"
	tpt "github.com/libp2p/go-libp2p-transport"
	ma "github.com/multiformats/go-multiaddr"
)

// GarlicDialer implements a go-libp2p compatible Dialer
type GarlicDialer struct {
	GarlicConn

	raddr   i2pma.I2PMultiaddr
	rPubKey crypto.PubKey
}

func (d GarlicDialer) ClosePeer(id peer.ID) error {
	return nil
}

func (d GarlicDialer) Connectedness(id peer.ID) net.Connectedness {
	var n net.Connectedness
	return n
}

func (d GarlicDialer) Conns() []net.Conn {
	return nil
}

func (d GarlicDialer) Peers() []peer.ID {
	return nil
}

func (d GarlicDialer) Peerstore() peerstore.Peerstore {
	return nil
}

func (d GarlicDialer) ConnsToPeer(p peer.ID) []net.Conn {
	return nil
}

func (d GarlicDialer) DialPeer(ctx context.Context, p peer.ID) (net.Conn, error) {
	c, e := d.DialGarlic(ctx, d.raddr, p)
	return c, e
}

func (d *GarlicDialer) Dial(ctx context.Context, raddr i2pma.I2PMultiaddr, p peer.ID) (tpt.Conn, error) {
	return d.DialGarlic(ctx, raddr, p)
}

func (d *GarlicDialer) DialGarlic(ctx context.Context, raddr i2pma.I2PMultiaddr, p peer.ID) (*GarlicConn, error) {
	var err error
	d.GarlicConn, err = NewGarlicConn(
		d.GarlicConn.GarlicTransport,
		d.laddr,
		d.lPrivKey,
		d.lPubKey,
		raddr,
		d.rPubKey,
	)
	if err != nil {
		return nil, err
	}
	d.GarlicConn.conn, err = d.StreamSession.Dial("ntcp", raddr.I2PAddr.Base32())
	if err != nil {
		return nil, err
	}
	return &d.GarlicConn, nil
}

func (d GarlicDialer) LocalPeer() peer.ID {
	var p peer.ID
	return p
}

func (d GarlicDialer) Notify(net.Notifiee) {

}

func (d GarlicDialer) StopNotify(net.Notifiee) {

}

func (d *GarlicDialer) Matches(a ma.Multiaddr) bool {
	return IsValidGarlicMultiAddr(a.(i2pma.I2PMultiaddr))
}

func (d *GarlicDialer) MatchesI2P(a i2pma.I2PMultiaddr) bool {
	return IsValidGarlicMultiAddr(a)
}

func NewGarlicDialer(t *GarlicTransport, laddr ma.Multiaddr, raddr i2pma.I2PMultiaddr) (*GarlicDialer, error) {
	sk, pk, err := crypto.GenerateEd25519Key(rand.Reader)
	if err != nil {
		return nil, err
	}
	garlicConn, err := NewGarlicConn(t, &laddr, sk, pk, raddr, pk)
	if err != nil {
		return nil, err
	}
	g := &GarlicDialer{
		GarlicConn: garlicConn,
	}
	return g, nil
}
