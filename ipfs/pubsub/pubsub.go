package pubsub

import (
	"context"
	"crypto/rand"
	"errors"
	"fmt"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	libp2p "github.com/libp2p/go-libp2p"
	p2pcrypt "github.com/libp2p/go-libp2p-core/crypto"
	host "github.com/libp2p/go-libp2p-core/host"
	peer "github.com/libp2p/go-libp2p-core/peer"
	protocol "github.com/libp2p/go-libp2p-core/protocol"
	routing "github.com/libp2p/go-libp2p-core/routing"
	p2pdiscovery "github.com/libp2p/go-libp2p/p2p/discovery"

	connman "github.com/libp2p/go-libp2p-connmgr"
	"github.com/libp2p/go-libp2p-discovery"
	dht "github.com/libp2p/go-libp2p-kad-dht"
	libpubsub "github.com/libp2p/go-libp2p-pubsub"
	libp2pquic "github.com/libp2p/go-libp2p-quic-transport"

	"EasyVoting/util"
)

func protocolID(topic string) protocol.ID {
	return protocol.ID(topic)
}

type PubSub struct {
	ps    *libpubsub.PubSub
	topic string
}

func (ps *PubSub) Topic() string {
	return ps.topic
}

type discoveryNotifee struct {
	h host.Host
}

func (n *discoveryNotifee) HandlePeerFound(pi peer.AddrInfo) {
	err := n.h.Connect(context.Background(), pi)
	if err != nil {
		log.Println("error connecting", err)
	}
}
func setupMdnsDiscovery(ctx context.Context, h host.Host) {
	disc, err := p2pdiscovery.NewMdnsService(ctx, h, time.Hour, "pubsub-example_oenilevbno;b")
	util.CheckError(err)

	disc.RegisterNotifee(&discoveryNotifee{h})
}
func New(ctx context.Context, topic string) *PubSub {
	priv, _, err := p2pcrypt.GenerateEd25519Key(rand.Reader)
	util.CheckError(err)

	bootstraps := dht.GetDefaultBootstrapPeerAddrInfos()
	h, err := libp2p.New(
		ctx,
		libp2p.DefaultListenAddrs,
		libp2p.DefaultTransports,
		libp2p.DefaultMuxers,
		libp2p.DefaultSecurity,
		libp2p.DefaultPeerstore,
		libp2p.DefaultEnableRelay,
		libp2p.Identity(priv),
		libp2p.ListenAddrStrings(
			"/ip4/0.0.0.0/udp/0/quic",
			"/ip6/::/udp/0/quic",
		),
		libp2p.Transport(libp2pquic.NewTransport),
		libp2p.DefaultStaticRelays(),
		libp2p.EnableAutoRelay(),
		libp2p.EnableNATService(),
		libp2p.NATPortMap(),
		libp2p.ConnectionManager(connman.NewConnManager(20, 40, time.Minute)),
		libp2p.Routing(func(h host.Host) (routing.PeerRouting, error) {
			return dht.New(ctx, h,
				dht.Mode(dht.ModeAuto),
				//dht.ProtocolPrefix(protocol.ID(topic)),
				dht.BootstrapPeers(bootstraps...),
			)
		}),
	)
	util.CheckError(err)
	defer h.Close()

	ps, err := libpubsub.NewGossipSub(ctx, h)
	util.CheckError(err)

	setupMdnsDiscovery(ctx, h)

	//peerIDs := ps.ListPeers(topic)
	//for _, peerID := range peerIDs{
	//	_, err := h.Network().DialPeer(ctx, peerID)
	//	log.Println(err)
	//}

	return &PubSub{ps, topic}
}

func New2(ctx context.Context, topic string) *PubSub {
	priv, _, err := p2pcrypt.GenerateEd25519Key(rand.Reader)
	util.CheckError(err)

	bootstraps := dht.GetDefaultBootstrapPeerAddrInfos()
	var dhtIpfs *dht.IpfsDHT
	h, err := libp2p.New(
		ctx,
		libp2p.DefaultListenAddrs,
		libp2p.DefaultTransports,
		libp2p.DefaultMuxers,
		libp2p.DefaultSecurity,
		libp2p.DefaultPeerstore,
		libp2p.DefaultEnableRelay,
		libp2p.Identity(priv),
		libp2p.ListenAddrStrings(
			"/ip4/0.0.0.0/udp/0/quic",
			"/ip6/::/udp/0/quic",
		),
		libp2p.Transport(libp2pquic.NewTransport),
		libp2p.DefaultStaticRelays(),
		libp2p.EnableAutoRelay(),
		libp2p.EnableNATService(),
		libp2p.NATPortMap(),
		libp2p.ConnectionManager(connman.NewConnManager(20, 40, time.Minute)),
		libp2p.Routing(func(h host.Host) (routing.PeerRouting, error) {
			var err error
			dhtIpfs, err = dht.New(ctx, h,
				dht.Mode(dht.ModeServer),
				//dht.ProtocolPrefix(protocol.ID(topic)),
				dht.BootstrapPeers(bootstraps...),
			)
			return dhtIpfs, err
		}),
	)
	util.CheckError(err)

	err = dhtIpfs.Bootstrap(ctx)
	util.CheckError(err)

	var wg sync.WaitGroup
	for _, peerAddr := range bootstraps {
		wg.Add(1)
		go func(pAddr peer.AddrInfo) {
			defer wg.Done()
			err := h.Connect(ctx, pAddr)
			if err == nil {
				log.Println("Bootstrap connect", pAddr)
			} else {
				log.Println("Bootstrap connect error", err, pAddr)
			}
		}(peerAddr)
	}
	wg.Wait()
	routingDiscovery := discovery.NewRoutingDiscovery(dhtIpfs)

	wg.Add(1)
	go func() {
		defer wg.Done()

		discovery.Advertise(ctx, routingDiscovery, topic)
		///*
		//TODO: No peer is found
		pAddrs, err := discovery.FindPeers(ctx, routingDiscovery, topic)
		if err != nil {
			log.Println("Find peers error", err)
			return
		}

		if len(pAddrs) == 0 {
			util.RaiseError("No peer is found")
		}
		fmt.Println(pAddrs)
		for _, pAddr := range pAddrs {
			if (len(pAddr.Addrs) == 0 && pAddr.ID == "") || (pAddr.ID == h.ID()) {
				continue
			}
			for {
				//fmt.Println(pAddr)
				err := h.Connect(ctx, pAddr) //TODO: fix Rounting connenct error
				//_, err := h.Network().DialPeer(ctx, pAddr.ID) // the same error as the above
				if err == nil {
					log.Println("Routing connect:", pAddr)
					break
				} else {
					log.Println("Routing connect error:", err, pAddr)
				}
			}
		}
		//*/
		/*
			peers := dhtIpfs.RoutingTable().ListPeers()
			fmt.Println(len(peers))
			for _, peerID := range peers{
				fmt.Println(peerID)
				for{
					_, err := h.Network().DialPeer(ctx, peerID)
					if err == nil{
						log.Println("Routing connect:", peerID)
						break
					}else{
						log.Println("Routing connect error:", err, peerID)
					}
				}
			}
		*/
	}()
	defer wg.Wait()

	ps, err := libpubsub.NewGossipSub(ctx, h, libpubsub.WithDiscovery(routingDiscovery))
	util.CheckError(err)

	return &PubSub{ps, topic}
}

func New3(ctx context.Context, topic string) *PubSub {
	return &PubSub{nil, topic}
}

func (ps *PubSub) Ls() []string {
	return ps.ps.GetTopics()
}
func (ps *PubSub) Peers() []peer.ID {
	return ps.ps.ListPeers(ps.topic)
}

func (ps *PubSub) Publish(ctx context.Context, data []byte) {
	fmt.Println("pub")
	tpc, err := ps.ps.Join(ps.topic)
	util.CheckError(err)
	defer tpc.Close()
	err = tpc.Publish(ctx, data)
	util.CheckError(err)
}

func (ps *PubSub) Subscribe(ctx context.Context) {
	fmt.Println("sub")
	tpc, err := ps.ps.Join(ps.topic)
	util.CheckError(err)
	defer tpc.Close()
	sub, err := tpc.Subscribe()
	util.CheckError(err)
	defer sub.Cancel()

	//go func(){
	for {
		msg, err := sub.Next(ctx)
		if err != nil {
			fmt.Println(err)
			return
		}
		fmt.Println(string(msg.GetData()))
	}
	//}()
}

func Publish(tpc string, data []byte) {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	// Generate key pair for identity
	priv, _, err := p2pcrypt.GenerateKeyPair(p2pcrypt.Ed25519, -1)
	if err != nil {
		log.Println(err)
		return
	}

	bootstrapPeerAddrs := dht.GetDefaultBootstrapPeerAddrInfos()
	var dhtIpfs *dht.IpfsDHT
	h, err := libp2p.New(ctx,
		libp2p.DefaultListenAddrs,
		libp2p.DefaultTransports,
		libp2p.DefaultMuxers,
		libp2p.DefaultSecurity,
		libp2p.DefaultPeerstore,
		libp2p.DefaultEnableRelay,
		libp2p.Identity(priv),
		libp2p.ListenAddrStrings(
			"/ip4/0.0.0.0/udp/0/quic",
			"/ip6/::/udp/0/quic",
		),
		libp2p.Transport(libp2pquic.NewTransport),
		libp2p.DefaultStaticRelays(),
		libp2p.EnableAutoRelay(),
		libp2p.EnableNATService(),
		libp2p.NATPortMap(),
		libp2p.ConnectionManager(connman.NewConnManager(20, 40, time.Minute)),
		libp2p.Routing(func(h host.Host) (routing.PeerRouting, error) {
			var err error
			dhtIpfs, err = dht.New(ctx, h,
				dht.Mode(dht.ModeServer),
				dht.BootstrapPeers(bootstrapPeerAddrs...))
			return dhtIpfs, err
		}))
	if err != nil {
		log.Println(err)
		return
	}
	defer h.Close()
	log.Println("Host ID:", h.ID())
	log.Println("ListenAddresses:", h.Network().ListenAddresses())
	err = dhtIpfs.Bootstrap(ctx)
	if err != nil {
		log.Println(err)
		return
	}
	var wg sync.WaitGroup
	for _, peerAddr := range bootstrapPeerAddrs {
		wg.Add(1)
		go func(peerAddr peer.AddrInfo) {
			defer wg.Done()
			err := h.Connect(ctx, peerAddr)
			if err != nil {
				log.Println("Bootstrap connect error:", err, peerAddr)
			} else {
				log.Println("Bootstrap connect:", peerAddr)
			}
		}(peerAddr)
	}
	wg.Wait()
	routingDiscovery := discovery.NewRoutingDiscovery(dhtIpfs)
	ps, err := libpubsub.NewGossipSub(ctx, h, libpubsub.WithDiscovery(routingDiscovery))
	if err != nil {
		log.Println(err)
		return
	}

	log.Println("Advertising and finding peers...")

	wg.Add(1)
	go func() {
		defer wg.Done()

		for ctx.Err() == nil {
			_, err := routingDiscovery.Advertise(ctx, tpc)
			if err != nil {
				if !errors.Is(err, context.DeadlineExceeded) {
					log.Println(err)
					return
				} else {
					log.Println("Advertise failed retrying...")
				}
			} else {
				break
			}
		}
		// see limit option
		peerAddrsChan, err := routingDiscovery.FindPeers(ctx, tpc)
		if err != nil {
			log.Println("Find peers error", err)
			return
		}
		for {
			select {
			case peerAddr, ok := <-peerAddrsChan:
				if !ok {
					log.Println("There is no peerAddr")
					return
				}
				// most the peers are empty!
				if (len(peerAddr.Addrs) == 0 && peerAddr.ID == "") || peerAddr.ID == h.ID() {
					continue
				}
				err := h.Connect(ctx, peerAddr)
				if err != nil {
					log.Println("Routing connect error:", err, peerAddr)
				} else {
					log.Println("Routing connect:", peerAddr)
				}
			case <-ctx.Done():
				return
			}
		}
	}()
	defer wg.Wait()

	topic, err := ps.Join(tpc)
	if err != nil {
		log.Println(err)
		return
	}
	defer topic.Close()
	for ctx.Err() == nil {
		log.Println("Sending message...")
		err = topic.Publish(ctx, data)
		if err != nil {
			log.Println(err)
			return
		} //else{
		//	log.Println("Published")
		//}
		select {
		case <-time.After(2 * time.Second):
		case <-ctx.Done():
			return
		}
	}
}

func Subscribe(tpc string) {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	// Generate key pair for identity
	priv, _, err := p2pcrypt.GenerateKeyPair(p2pcrypt.Ed25519, -1)
	if err != nil {
		log.Println(err)
		return
	}

	bootstrapPeerAddrs := dht.GetDefaultBootstrapPeerAddrInfos()
	var dhtIpfs *dht.IpfsDHT
	h, err := libp2p.New(ctx,
		libp2p.DefaultListenAddrs,
		libp2p.DefaultTransports,
		libp2p.DefaultMuxers,
		libp2p.DefaultSecurity,
		libp2p.DefaultPeerstore,
		libp2p.DefaultEnableRelay,
		libp2p.Identity(priv),
		libp2p.ListenAddrStrings(
			"/ip4/0.0.0.0/udp/0/quic",
			"/ip6/::/udp/0/quic",
		),
		libp2p.Transport(libp2pquic.NewTransport),
		libp2p.DefaultStaticRelays(),
		libp2p.EnableAutoRelay(),
		libp2p.EnableNATService(),
		libp2p.NATPortMap(),
		libp2p.ConnectionManager(connman.NewConnManager(20, 40, time.Minute)),
		libp2p.Routing(func(h host.Host) (routing.PeerRouting, error) {
			var err error
			dhtIpfs, err = dht.New(ctx, h,
				dht.Mode(dht.ModeServer),
				dht.BootstrapPeers(bootstrapPeerAddrs...))
			return dhtIpfs, err
		}))
	if err != nil {
		log.Println(err)
		return
	}
	defer h.Close()
	log.Println("Host ID:", h.ID())
	log.Println("ListenAddresses:", h.Network().ListenAddresses())
	err = dhtIpfs.Bootstrap(ctx)
	if err != nil {
		log.Println(err)
		return
	}
	var wg sync.WaitGroup
	for _, peerAddr := range bootstrapPeerAddrs {
		wg.Add(1)
		go func(peerAddr peer.AddrInfo) {
			defer wg.Done()
			err := h.Connect(ctx, peerAddr)
			if err != nil {
				log.Println("Bootstrap connect error:", err, peerAddr)
			} else {
				log.Println("Bootstrap connect:", peerAddr)
			}
		}(peerAddr)
	}
	wg.Wait()
	routingDiscovery := discovery.NewRoutingDiscovery(dhtIpfs)
	ps, err := libpubsub.NewGossipSub(ctx, h, libpubsub.WithDiscovery(routingDiscovery))
	if err != nil {
		log.Println(err)
		return
	}

	log.Println("Advertising and finding peers...")

	wg.Add(1)
	go func() {
		defer wg.Done()

		for ctx.Err() == nil {
			_, err := routingDiscovery.Advertise(ctx, tpc)
			if err != nil {
				if !errors.Is(err, context.DeadlineExceeded) {
					log.Println(err)
					return
				} else {
					log.Println("Advertise failed retrying...")
				}
			} else {
				break
			}
		}
		// see limit option
		peerAddrsChan, err := routingDiscovery.FindPeers(ctx, tpc)
		if err != nil {
			log.Println(err)
			return
		}
		for {
			select {
			case peerAddr, ok := <-peerAddrsChan:
				if !ok {
					log.Println("There is no peerAddr")
					return
				}
				// most of the peers are empty!
				if (len(peerAddr.Addrs) == 0 && peerAddr.ID == "") || peerAddr.ID == h.ID() {
					continue
				}
				err := h.Connect(ctx, peerAddr)
				if err != nil {
					log.Println("Routing connect error:", err, peerAddr)
				} else {
					log.Println("Routing connect:", peerAddr)
				}
			case <-ctx.Done():
				return
			}
		}
	}()
	defer wg.Wait()

	topic, err := ps.Join(tpc)
	if err != nil {
		log.Println(err)
		return
	}
	defer topic.Close()
	sub, err := topic.Subscribe()
	if err != nil {
		log.Println(err)
		return
	}
	defer sub.Cancel()

	log.Println("Listening on messages...")

	for ctx.Err() == nil {
		mess, err := sub.Next(ctx)
		if err != nil {
			log.Println(err)
			return
		}
		//if mess.GetFrom() == h.ID() {
		//	continue
		//}
		log.Println("From:", mess.GetFrom())
		log.Println("Data:", string(mess.GetData()))
	}
}
