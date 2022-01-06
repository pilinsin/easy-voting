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

func newHost() (host.Host, error) {
	priv, _, _ := p2pcrypt.GenerateEd25519Key(rand.Reader)
	return libp2p.New(
		libp2p.Identity(priv),
		libp2p.ListenAddrStrings(
			"/ip4/0.0.0.0/tcp/0",
			"/ip6/::/tcp/0",
		),
	)
}
func newDHT(ctx context.Context, h host.Host, verbose bool) (*dht.IpfsDHT, error) {
	d, err := dht.New(ctx, h, dht.Mode(dht.ModeServer))
	if err != nil {
		return nil, err
	}
	err = d.Bootstrap(ctx)
	if err != nil {
		return nil, err
	}

	bootstraps := dht.GetDefaultBootstrapPeerAddrInfos()
	ctx2, cancel := context.WithTimeout(ctx, 2*time.Minute)
	defer cancel()
	var wg sync.WaitGroup
	for _, pInfo := range bootstraps {
		wg.Add(1)
		go func(pAddr peer.AddrInfo) {
			defer wg.Done()
			for {
				hConnErr := h.Connect(ctx2, pAddr)
				if verbose {
					log.Println("Bootstrap connection:", hConnErr, pAddr)
				}
				if hConnErr == nil {
					break
				}
				time.After(10 * time.Second)
			}
		}(pInfo)
	}
	wg.Wait()

	err = ctx2.Err()
	if err != nil {
		return nil, err
	}
	return d, nil
}
func discoverPeers(ctx context.Context, h host.Host, disc *discovery.RoutingDiscovery, topic string, parent, verbose bool) error {
	discovery.Advertise(ctx, disc, topic)

	timer := time.NewTicker(time.Second)
	defer timer.Stop()
	for {
		select {
		case <-ctx.Done():
			return nil
		case <-timer.C:
			peerIDs, err := discovery.FindPeers(ctx, disc, topic)
			if err != nil {
				return err
			}

			numConnectOthers := 0
			for _, peerID := range peerIDs {
				if peerID.ID == h.ID() {
					if verbose {
						log.Println("self ID", peerID.ID)
					}
					continue
				}
				if h.Network().Connectedness(peerID.ID) != network.Connected {
					//_, err = h.Network().DialPeer(ctx, peerID.ID)
					if err = h.Connect(ctx, peerID); err != nil {
						if verbose {
							log.Println("connection error", peerID.ID)
						}
						continue
					} else {
						if verbose {
							log.Println("connect", peerID.ID)
						}
						numConnectOthers++
						continue
					}
				} else {
					if verbose {
						log.Println("already connected", peerID.ID)
					}
					numConnectOthers++
				}
			}
			if len(peerIDs) <= 0 {
				return util.NewError("no peers are found")
			}
			if parent || numConnectOthers > 0 {
				return nil
			} else {
				return util.NewError("cannot connect to any other peers.")
			}
		}
	}

}

func New(topic string, parent bool) (*PubSub, error) {
	//publishより前にsubscribeしておけば読み込める
	//new PubSub -> Subscribe -> (Publish)
	ps, err := newPubSub(topic, parent, true)
	if err != nil {
		return nil, err
	}
	sub, err := ps.Subscribe(topic)
	if err != nil {
		return nil, err
	}
	return &PubSub{ps, sub, topic}, nil
}
func newPubSub(topic string, parent, verbose bool) (*libpubsub.PubSub, error) {
	ctx := context.Background()
	h, err := newHost()
	if err != nil {
		return nil, err
	}
	d, err := newDHT(ctx, h, verbose)
	if err != nil {
		return nil, err
	}
	disc := discovery.NewRoutingDiscovery(d)
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		err = discoverPeers(ctx, h, disc, topic, parent, verbose)
		if err != nil {
			return
		}
		if verbose {
			fmt.Println("discovered")
		}
	}()
	wg.Wait()
	if err != nil {
		return nil, err
	}

	ps, err := libpubsub.NewGossipSub(ctx, h, libpubsub.WithDiscovery(disc))
	if err != nil {
		return nil, err
	}
	return ps, nil
}
func (ps *PubSub) Close() {
	ps.sub.Cancel()
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

	return msg.GetData()
}

func (ps *PubSub) SubTest() {
	var dataset []string
	for {
		data := ps.Next()
		if data == nil {
			fmt.Println(dataset)
			return
		}

		dataset = append(dataset, string(data))
	}

}
func (ps *PubSub) Subscribe() [][]byte {
	var dataset [][]byte
	for {
		data := ps.Next()
		if data == nil {
			if len(dataset) > 0 {
				return dataset
			} else {
				return nil
			}
		}
		dataset = append(dataset, data)
	}
}
