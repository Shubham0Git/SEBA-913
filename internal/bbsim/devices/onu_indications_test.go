/*
 * Copyright 2018-present Open Networking Foundation

 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at

 * http://www.apache.org/licenses/LICENSE-2.0

 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package devices

import (
	"context"
	"errors"
	"github.com/opencord/voltha-protos/v2/go/openolt"
	"google.golang.org/grpc"
	"gotest.tools/assert"
	"testing"
	"time"
)

type mockStream struct {
	grpc.ServerStream
	CallCount int
	Calls     map[int]*openolt.OnuDiscIndication
	channel   chan int
	fail      bool
}

func (s *mockStream) Send(ind *openolt.Indication) error {
	s.CallCount++
	if s.fail {
		return errors.New("fake-error")
	}
	s.Calls[s.CallCount] = ind.GetOnuDiscInd()
	s.channel <- s.CallCount
	return nil
}

// test that we're sending a Discovery indication to VOLTHA
func Test_Onu_DiscoverIndication_send_on_discovery(t *testing.T) {
	onu := createTestOnu()
	stream := &mockStream{
		CallCount: 0,
		Calls:     make(map[int]*openolt.OnuDiscIndication),
		fail:      false,
		channel:   make(chan int, 10),
	}
	ctx, cancel := context.WithCancel(context.TODO())
	go onu.ProcessOnuMessages(ctx, stream, nil)
	onu.InternalState.SetState("initialized")
	onu.InternalState.Event("discover")

	select {
	case <-time.After(90 * time.Millisecond):
		assert.Equal(t, stream.CallCount, 1)
		assert.Equal(t, stream.Calls[1].IntfId, onu.PonPortID)
		assert.Equal(t, stream.Calls[1].SerialNumber, onu.SerialNumber)
	}
	cancel()
}

// test that if the discovery indication is not acknowledge we'll keep sending new ones
func Test_Onu_DiscoverIndication_retry_on_discovery(t *testing.T) {
	onu := createTestOnu()
	stream := &mockStream{
		CallCount: 0,
		Calls:     make(map[int]*openolt.OnuDiscIndication),
		fail:      false,
		channel:   make(chan int, 10),
	}
	ctx, cancel := context.WithCancel(context.TODO())
	go onu.ProcessOnuMessages(ctx, stream, nil)
	onu.InternalState.SetState("initialized")
	onu.InternalState.Event("discover")

	select {
	case <-time.After(400 * time.Millisecond):
		assert.Equal(t, stream.CallCount, 4)
	}
	cancel()
}

// test that if the discovery indication is not acknowledge we'll send a new one
func Test_Onu_DiscoverIndication_retry_on_discovery_stops(t *testing.T) {
	onu := createTestOnu()
	onu.DiscoveryRetryDelay = 500 * time.Millisecond
	stream := &mockStream{
		CallCount: 0,
		Calls:     make(map[int]*openolt.OnuDiscIndication),
		fail:      false,
		channel:   make(chan int, 10),
	}
	ctx, cancel := context.WithCancel(context.TODO())
	go onu.ProcessOnuMessages(ctx, stream, nil)
	onu.InternalState.SetState("initialized")
	onu.InternalState.Event("discover")

	go func() {
		for calls := range stream.channel {
			if calls == 2 {
				onu.InternalState.SetState("enabled")
			}
		}
	}()

	select {
	case <-time.After(1 * time.Second):

		assert.Equal(t, stream.CallCount, 2)
	}
	cancel()
}
