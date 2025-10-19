package chat

import (
	"context"
	"encoding/json"
	"sync"

	"github.com/gorilla/websocket"
	"github.com/pion/webrtc/v4"
)

type PeerState struct {
	WS             *websocket.Conn
	WSMu           sync.Mutex
	PeerConnection *webrtc.PeerConnection
	SessionID      string
	Tracks         map[string]*Track
	TracksMu       sync.RWMutex
	PendingICE     []webrtc.ICECandidateInit
}

type PublishedTrack struct {
	Remote *webrtc.TrackRemote
	Codec  webrtc.RTPCodecCapability
	Locals map[string]*Track
	Mu     sync.RWMutex
	ctx    context.Context
	cancel context.CancelFunc
}

type Track struct {
	Track  *webrtc.TrackLocalStaticRTP
	Sender *webrtc.RTPSender
}

type Room struct {
	Peers     map[string]*PeerState
	Mutex     sync.RWMutex
	Published map[string]*PublishedTrack
	PubMutex  sync.RWMutex
}

func (r *Room) RemovePeer(sessionID string) {
	r.Mutex.Lock()
	defer r.Mutex.Unlock()
	delete(r.Peers, sessionID)
}

type SignalMessage struct {
	Type      string          `json:"type"`
	Payload   json.RawMessage `json:"payload"`
	RoomID    int          `json:"roomId"`
	SessionID string          `json:"sessionId,omitempty"`
}
