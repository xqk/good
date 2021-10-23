package p2c_test

import (
	"context"
	"fmt"
	"net"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/xqk/good/pkg/client/grpc/balancer/p2c"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/peer"
	"google.golang.org/grpc/resolver"
	"google.golang.org/grpc/resolver/manual"
	"google.golang.org/grpc/status"
	testpb "google.golang.org/grpc/test/grpc_testing"
)

type testServer struct {
	testpb.TestServiceServer
}

func (s *testServer) EmptyCall(ctx context.Context, in *testpb.Empty) (*testpb.Empty, error) {
	// time.Sleep(1 * time.Millisecond)
	return &testpb.Empty{}, nil
}

func (s *testServer) FullDuplexCall(stream testpb.TestService_FullDuplexCallServer) error {
	return nil
}

type test struct {
	servers   []*grpc.Server
	addresses []string
}

func (t *test) cleanup() {
	for _, s := range t.servers {
		s.Stop()
	}
}

func startTestServers(count int) (_ *test, err error) {
	t := &test{}

	defer func() {
		if err != nil {
			t.cleanup()
		}
	}()
	for i := 0; i < count; i++ {
		lis, err := net.Listen("tcp", "localhost:0")
		if err != nil {
			return nil, fmt.Errorf("failed to listen %v", err)
		}

		s := grpc.NewServer()
		testpb.RegisterTestServiceServer(s, &testServer{})
		t.servers = append(t.servers, s)
		t.addresses = append(t.addresses, lis.Addr().String())

		go func(s *grpc.Server, l net.Listener) {
			s.Serve(l)
		}(s, lis)
	}

	return t, nil
}

func TestOneBackend(t *testing.T) {

	r, cleanup := manual.GenerateAndRegisterManualResolver()
	defer cleanup()

	test, err := startTestServers(1)
	if err != nil {
		t.Fatalf("failed to start servers: %v", err)
	}
	defer test.cleanup()

	cc, err := grpc.Dial(r.Scheme()+":///test.server", grpc.WithInsecure(), grpc.WithBalancerName(p2c.Name))
	if err != nil {
		t.Fatalf("failed to dial: %v", err)
	}
	defer cc.Close()
	testc := testpb.NewTestServiceClient(cc)
	// The first RPC should fail because there's no address.
	ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond)
	defer cancel()
	if _, err := testc.EmptyCall(ctx, &testpb.Empty{}); err == nil || status.Code(err) != codes.DeadlineExceeded {
		t.Fatalf("EmptyCall() = _, %v, want _, DeadlineExceeded", err)
	}

	r.UpdateState(resolver.State{Addresses: []resolver.Address{{Addr: test.addresses[0]}}})
	// The second RPC should succeed.
	if _, err := testc.EmptyCall(context.Background(), &testpb.Empty{}); err != nil {
		t.Fatalf("EmptyCall() = _, %v, want _, <nil>", err)
	}
}

func TestAddressesRemoved(t *testing.T) {

	r, cleanup := manual.GenerateAndRegisterManualResolver()
	defer cleanup()

	test, err := startTestServers(2)
	if err != nil {
		t.Fatalf("failed to start servers: %v", err)
	}
	defer test.cleanup()

	cc, err := grpc.Dial(r.Scheme()+":///test.server", grpc.WithInsecure(), grpc.WithBalancerName(p2c.Name))
	if err != nil {
		t.Fatalf("failed to dial: %v", err)
	}
	defer cc.Close()
	testc := testpb.NewTestServiceClient(cc)
	// The first RPC should fail because there's no address.
	ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond)
	defer cancel()
	if _, err := testc.EmptyCall(ctx, &testpb.Empty{}); err == nil || status.Code(err) != codes.DeadlineExceeded {
		t.Fatalf("EmptyCall() = _, %v, want _, DeadlineExceeded", err)
	}

	r.UpdateState(resolver.State{Addresses: []resolver.Address{{Addr: test.addresses[0]}, {Addr: test.addresses[1]}}})
	// The second RPC should succeed.
	if _, err := testc.EmptyCall(context.Background(), &testpb.Empty{}); err != nil {
		t.Fatalf("EmptyCall() = _, %v, want _, <nil>", err)
	}

	r.UpdateState(resolver.State{Addresses: []resolver.Address{}})
	// The third RPC should fail.
	if _, err := testc.EmptyCall(context.Background(), &testpb.Empty{}); err == nil {
		t.Fatalf("EmptyCall() = _, %v, want _, <nil>", err)
	}
}

func TestOneServerDown(t *testing.T) {

	r, cleanup := manual.GenerateAndRegisterManualResolver()
	defer cleanup()

	backendCount := 3
	test, err := startTestServers(backendCount)
	if err != nil {
		t.Fatalf("failed to start servers: %v", err)
	}
	defer test.cleanup()

	cc, err := grpc.Dial(r.Scheme()+":///test.server", grpc.WithInsecure(), grpc.WithBalancerName(p2c.Name))
	if err != nil {
		t.Fatalf("failed to dial: %v", err)
	}
	defer cc.Close()
	testc := testpb.NewTestServiceClient(cc)
	// The first RPC should fail because there's no address.
	ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond)
	defer cancel()
	if _, err := testc.EmptyCall(ctx, &testpb.Empty{}); err == nil || status.Code(err) != codes.DeadlineExceeded {
		t.Fatalf("EmptyCall() = _, %v, want _, DeadlineExceeded", err)
	}

	var resolvedAddrs []resolver.Address
	for i := 0; i < backendCount; i++ {
		resolvedAddrs = append(resolvedAddrs, resolver.Address{Addr: test.addresses[i]})
	}

	r.UpdateState(resolver.State{Addresses: resolvedAddrs})
	var p peer.Peer
	// Make sure connections to all servers are up.
	for si := 0; si < backendCount; si++ {
		var connected bool
		for i := 0; i < 1000; i++ {
			if _, err := testc.EmptyCall(context.Background(), &testpb.Empty{}, grpc.Peer(&p)); err != nil {
				t.Fatalf("EmptyCall() = _, %v, want _, <nil>", err)
			}
			if p.Addr.String() == test.addresses[si] {
				connected = true
				break
			}
			time.Sleep(time.Millisecond)
		}
		if !connected {
			t.Fatalf("Connection to %v was not up after more than 1 second", test.addresses[si])
		}
	}

	for i := 0; i < 3*backendCount; i++ {
		if _, err := testc.EmptyCall(context.Background(), &testpb.Empty{}, grpc.Peer(&p)); err != nil {
			t.Fatalf("EmptyCall() = _, %v, want _, <nil>", err)
		}
		if p.Addr.String() == "" {
			t.Fatalf("Index %d: want peer %v, got peer %v", i, test.addresses[i%backendCount], p.Addr.String())
		}
	}

	// Stop one server, RPCs should p2c among the remaining servers.
	backendCount--
	test.servers[backendCount].Stop()
	// Loop until see server[backendCount-1] twice without seeing server[backendCount].
	var targetSeen int
	for i := 0; i < 1000; i++ {
		if _, err := testc.EmptyCall(context.Background(), &testpb.Empty{}, grpc.Peer(&p)); err != nil {
			targetSeen = 0
			t.Logf("EmptyCall() = _, %v, want _, <nil>", err)
			// Due to a race, this RPC could possibly get the connection that
			// was closing, and this RPC may fail. Keep trying when this
			// happens.
			continue
		}
		switch p.Addr.String() {
		case test.addresses[backendCount-1]:
			targetSeen++
		case test.addresses[backendCount]:
			// Reset targetSeen if peer is server[backendCount].
			targetSeen = 0
		}
		// Break to make sure the last picked address is server[-1], so the following for loop won't be flaky.
		if targetSeen >= 2 {
			break
		}
	}
	if targetSeen != 2 {
		t.Fatal("Failed to see server[backendCount-1] twice without seeing server[backendCount]")
	}
	for i := 0; i < 3*backendCount; i++ {
		if _, err := testc.EmptyCall(context.Background(), &testpb.Empty{}, grpc.Peer(&p)); err != nil {
			t.Fatalf("EmptyCall() = _, %v, want _, <nil>", err)
		}
		if p.Addr.String() == test.addresses[backendCount] {
			t.Errorf("Index %d: want peer %v, got peer %v", i, test.addresses[backendCount-1], p.Addr.String())
		}
	}
}
func TestBackendsRandom(t *testing.T) {

	r, cleanup := manual.GenerateAndRegisterManualResolver()
	defer cleanup()

	backendCount := 5
	test, err := startTestServers(backendCount)
	if err != nil {
		t.Fatalf("failed to start servers: %v", err)
	}
	defer test.cleanup()

	cc, err := grpc.Dial(r.Scheme()+":///test.server", grpc.WithInsecure(), grpc.WithBalancerName(p2c.Name))
	if err != nil {
		t.Fatalf("failed to dial: %v", err)
	}
	defer cc.Close()
	testc := testpb.NewTestServiceClient(cc)
	// The first RPC should fail because there's no address.
	ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond)
	defer cancel()
	if _, err := testc.EmptyCall(ctx, &testpb.Empty{}); err == nil || status.Code(err) != codes.DeadlineExceeded {
		t.Fatalf("EmptyCall() = _, %v, want _, DeadlineExceeded", err)
	}

	var resolvedAddrs []resolver.Address
	for i := 0; i < backendCount; i++ {
		resolvedAddrs = append(resolvedAddrs, resolver.Address{Addr: test.addresses[i]})
	}

	r.UpdateState(resolver.State{Addresses: resolvedAddrs})
	var p peer.Peer

	countMap := make(map[string]int)
	// Make sure connections to all servers are up.
	for si := 0; si < backendCount; si++ {
		for i := 0; i < 1000; i++ {
			if _, err := testc.EmptyCall(context.Background(), &testpb.Empty{}, grpc.Peer(&p)); err != nil {
				t.Fatalf("EmptyCall() = _, %v, want _, <nil>", err)
			}
			if p.Addr.String() != "" {
				countMap[p.Addr.String()] = countMap[p.Addr.String()] + 1
			}
		}
	}

	assert.Equal(t, backendCount, len(countMap))
	total := 0
	for addr, count := range countMap {
		total = total + count
		assert.True(t, count > 900, "req count less than 900", addr, countMap)
	}

	assert.Equal(t, total, 1000*backendCount)
}

func TestCloseWithPendingRPC(t *testing.T) {

	r, cleanup := manual.GenerateAndRegisterManualResolver()
	defer cleanup()

	test, err := startTestServers(1)
	if err != nil {
		t.Fatalf("failed to start servers: %v", err)
	}
	defer test.cleanup()

	cc, err := grpc.Dial(r.Scheme()+":///test.server", grpc.WithInsecure(), grpc.WithBalancerName(p2c.Name))
	if err != nil {
		t.Fatalf("failed to dial: %v", err)
	}
	testc := testpb.NewTestServiceClient(cc)

	var wg sync.WaitGroup
	for i := 0; i < 3; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			// This RPC blocks until cc is closed.
			ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
			if _, err := testc.EmptyCall(ctx, &testpb.Empty{}); status.Code(err) == codes.DeadlineExceeded {
				t.Errorf("RPC failed because of deadline after cc is closed; want error the client connection is closing")
			}
			cancel()
		}()
	}
	cc.Close()
	wg.Wait()
}

func TestNewAddressWhileBlocking(t *testing.T) {

	r, cleanup := manual.GenerateAndRegisterManualResolver()
	defer cleanup()

	test, err := startTestServers(1)
	if err != nil {
		t.Fatalf("failed to start servers: %v", err)
	}
	defer test.cleanup()

	cc, err := grpc.Dial(r.Scheme()+":///test.server", grpc.WithInsecure(), grpc.WithBalancerName(p2c.Name))
	if err != nil {
		t.Fatalf("failed to dial: %v", err)
	}
	defer cc.Close()
	testc := testpb.NewTestServiceClient(cc)
	// The first RPC should fail because there's no address.
	ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond)
	defer cancel()
	if _, err := testc.EmptyCall(ctx, &testpb.Empty{}); err == nil || status.Code(err) != codes.DeadlineExceeded {
		t.Fatalf("EmptyCall() = _, %v, want _, DeadlineExceeded", err)
	}

	r.UpdateState(resolver.State{Addresses: []resolver.Address{{Addr: test.addresses[0]}}})
	// The second RPC should succeed.
	ctx, cancel = context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	if _, err := testc.EmptyCall(ctx, &testpb.Empty{}); err != nil {
		t.Fatalf("EmptyCall() = _, %v, want _, nil", err)
	}

	r.UpdateState(resolver.State{Addresses: []resolver.Address{}})

	var wg sync.WaitGroup
	for i := 0; i < 3; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			// This RPC blocks until NewAddress is called.
			testc.EmptyCall(context.Background(), &testpb.Empty{})
		}()
	}
	time.Sleep(50 * time.Millisecond)
	r.UpdateState(resolver.State{Addresses: []resolver.Address{{Addr: test.addresses[0]}}})
	wg.Wait()
}

func TestAllServersDown(t *testing.T) {

	r, cleanup := manual.GenerateAndRegisterManualResolver()
	defer cleanup()

	backendCount := 3
	test, err := startTestServers(backendCount)
	if err != nil {
		t.Fatalf("failed to start servers: %v", err)
	}
	defer test.cleanup()

	cc, err := grpc.Dial(r.Scheme()+":///test.server", grpc.WithInsecure(), grpc.WithBalancerName(p2c.Name))
	if err != nil {
		t.Fatalf("failed to dial: %v", err)
	}
	defer cc.Close()
	testc := testpb.NewTestServiceClient(cc)
	// The first RPC should fail because there's no address.
	ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond)
	defer cancel()
	if _, err := testc.EmptyCall(ctx, &testpb.Empty{}); err == nil || status.Code(err) != codes.DeadlineExceeded {
		t.Fatalf("EmptyCall() = _, %v, want _, DeadlineExceeded", err)
	}

	var resolvedAddrs []resolver.Address
	for i := 0; i < backendCount; i++ {
		resolvedAddrs = append(resolvedAddrs, resolver.Address{Addr: test.addresses[i]})
	}

	r.UpdateState(resolver.State{Addresses: resolvedAddrs})
	var p peer.Peer
	// Make sure connections to all servers are up.
	for si := 0; si < backendCount; si++ {
		var connected bool
		for i := 0; i < 1000; i++ {
			if _, err := testc.EmptyCall(context.Background(), &testpb.Empty{}, grpc.Peer(&p)); err != nil {
				t.Fatalf("EmptyCall() = _, %v, want _, <nil>", err)
			}
			if p.Addr.String() == test.addresses[si] {
				connected = true
				break
			}
			time.Sleep(time.Millisecond)
		}
		if !connected {
			t.Fatalf("Connection to %v was not up after more than 1 second", test.addresses[si])
		}
	}

	for i := 0; i < 3*backendCount; i++ {
		if _, err := testc.EmptyCall(context.Background(), &testpb.Empty{}, grpc.Peer(&p)); err != nil {
			t.Fatalf("EmptyCall() = _, %v, want _, <nil>", err)
		}
		if p.Addr.String() == "" {
			t.Fatalf("Index %d: want peer %v, got peer %v", i, test.addresses[i%backendCount], p.Addr.String())
		}
	}

	// All servers are stopped, failfast RPC should fail with unavailable.
	for i := 0; i < backendCount; i++ {
		test.servers[i].Stop()
	}
	time.Sleep(100 * time.Millisecond)
	for i := 0; i < 1000; i++ {
		if _, err := testc.EmptyCall(context.Background(), &testpb.Empty{}); status.Code(err) == codes.Unavailable {
			return
		}
		time.Sleep(time.Millisecond)
	}
	t.Fatalf("Failfast RPCs didn't fail with Unavailable after all servers are stopped")
}
