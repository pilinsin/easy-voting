package pubsub

import (
	"context"
	"crypto/rand"
	"fmt"
	"log"
	"sync"
	"time"

	libp2p "github.com/libp2p/go-libp2p"
	p2pcrypt "github.com/libp2p/go-libp2p-core/crypto"
	host "github.com/libp2p/go-libp2p-core/host"
	network "github.com/libp2p/go-libp2p-core/network"
	peer "github.com/libp2p/go-libp2p-core/peer"

	"github.com/libp2p/go-libp2p-discovery"
	dht "github.com/libp2p/go-libp2p-kad-dht"
	libpubsub "github.com/libp2p/go-libp2p-pubsub"

	"EasyVoting/util"
)

type PubSub struct {
	ps    *libpubsub.PubSub
	sub   *libpubsub.Subscription
	topic string
}

func (ps *PubSub) Topic() string {
	return ps.topic
}

func newHost(ctx context.Context) host.Host {
	priv, _, err := p2pcrypt.GenerateEd25519Key(rand.Reader)
	util.CheckError(err)

	h, err := libp2p.New(
		ctx,
		libp2p.Identity(priv),
		libp2p.ListenAddrStrings(
			"/ip4/0.0.0.0/tcp/0",
			"/ip6/::/tcp/0",
		),
	)
	util.CheckError(err)

	return h
}
func newDHT(ctx context.Context, h host.Host, verbose bool) *dht.IpfsDHT {
	d, err := dht.New(ctx, h, dht.Mode(dht.ModeServer))
	util.CheckError(err)
	err = d.Bootstrap(ctx)
	util.CheckError(err)

	bootstraps := dht.GetDefaultBootstrapPeerAddrInfos()
	var wg sync.WaitGroup
	isErr := true
	for _, pInfo := range bootstraps {
		wg.Add(1)
		go func(pAddr peer.AddrInfo) {
			defer wg.Done()
			hConnErr := h.Connect(ctx, pAddr)
			isErr = isErr && (hConnErr != nil)
			if verbose {
				log.Println("Bootstrap connection:", hConnErr, pAddr)
			}
		}(pInfo)
	}
	wg.Wait()
	if isErr {
		return nil
	} else {
		return d
	}

}
func discoverPeers(ctx context.Context, h host.Host, disc *discovery.RoutingDiscovery, topic string, verbose bool) {
	discovery.Advertise(ctx, disc, topic)

	timer := time.NewTicker(time.Second)
	defer timer.Stop()
	for {
		select {
		case <-ctx.Done():
			return
		case <-timer.C:
			peerIDs, err := discovery.FindPeers(ctx, disc, topic)
			util.CheckError(err)

			numConnectPeers := 0
			for _, peerID := range peerIDs {
				if peerID.ID == h.ID() {
					if verbose {
						log.Println("self ID", peerID.ID)
					}
					numConnectPeers++
					continue
				}
				if h.Network().Connectedness(peerID.ID) != network.Connected {
					_, err = h.Network().DialPeer(ctx, peerID.ID)
					//err = h.Connect(ctx, peerID)
					if err != nil {
						if verbose {
							log.Println("connection error", peerID.ID)
						}
						continue
					} else {
						if verbose {
							log.Println("connect", peerID.ID)
						}
						numConnectPeers++
						continue
					}
				} else {
					if verbose {
						log.Println("already connected", peerID.ID)
					}
					numConnectPeers++
				}
			}
			if len(peerIDs) > 0 && len(peerIDs) == numConnectPeers {
				return
			}
		}
	}

}

func New(topic string) *PubSub {
	//publishより前にsubscribeしておけば読み込める
	//new PubSub -> Subscribe -> (Publish)
	ps := new(topic, true)
	sub, err := ps.Subscribe(topic)
	util.CheckError(err)
	return &PubSub{ps, sub, topic}
}
func new(topic string, verbose bool) *libpubsub.PubSub {
	ctx := context.Background()
	h := newHost(ctx)
	d := newDHT(ctx, h, verbose)
	disc := discovery.NewRoutingDiscovery(d)
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		discoverPeers(ctx, h, disc, topic, verbose)
		if verbose {
			fmt.Println("discovered")
		}
	}()
	wg.Wait()

	ps, err := libpubsub.NewGossipSub(ctx, h, libpubsub.WithDiscovery(disc))
	util.CheckError(err)
	return ps
}

func (ps *PubSub) Publish(data []byte) {
	for {
		err := ps.ps.Publish(ps.topic, data)
		if err == nil {
			log.Println("pub")
			return
		} else {
			<-time.After(2 * time.Second)
		}
	}
}

func (ps *PubSub) Next() []byte {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	msg, err := ps.sub.Next(ctx)
	if err != nil {
		//fmt.Println("sub.Next error:", err)
		return nil
	}

	//fmt.Println("anlvanl")
	return msg.GetData()
}

func (ps *PubSub) Close() {
	ps.sub.Cancel()
}
