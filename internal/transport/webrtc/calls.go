package calls

import (
	"encoding/json"
	"fmt"
	"log"
	"sync"
	"time"

	errors "github.com/Trecer05/Swiftly/internal/errors/chat"
	models "github.com/Trecer05/Swiftly/internal/model/chat"
	service "github.com/Trecer05/Swiftly/internal/service/chat"

	"github.com/pion/webrtc/v4"
)

func ReadWS(chatID int, rooms map[models.CallsKey]*models.Room, room *models.Room, key models.CallsKey, currentPeerState *models.PeerState, roomsMutex *sync.RWMutex) {
	log.Printf("WS started for %s", currentPeerState.SessionID)

	for {
		_, raw, err := currentPeerState.WS.ReadMessage()
		if err != nil {
			log.Printf("error reading message from %s: %v", currentPeerState.SessionID, err)
			break
		}

		log.Printf("[WS RECV %s] %s", currentPeerState.SessionID, string(raw)[:29])

		var msg models.SignalMessage
		if err := json.Unmarshal(raw, &msg); err != nil {
			log.Printf("invalid signal json from %s: %v", currentPeerState.SessionID, err)
			continue
		}

		switch msg.Type {	
		case "join":
			log.Printf("[%s] JOIN room=%d", currentPeerState.SessionID, msg.RoomID)
			roomsMutex.Lock()
			if rooms[key] == nil {
				rooms[key] = service.NewRoom()
			}
			room = rooms[key]
			roomsMutex.Unlock()

			pc, err := createPeerConnection(currentPeerState, chatID, roomsMutex, rooms, key)
			if err != nil {
				if err == errors.ErrorNoCallRoom {
					log.Println("delete pc from room: ", chatID)
					pc.Close()
					return
				}
				log.Printf("error creating PeerConnection for %s: %v", currentPeerState.SessionID, err)
				continue
			}
			currentPeerState.PeerConnection = pc

			room.Mutex.Lock()
			room.Peers[currentPeerState.SessionID] = currentPeerState
			room.Mutex.Unlock()

			room.PubMutex.RLock()
			for trackID, pub := range room.Published {
				if _, ok := pub.Locals[currentPeerState.SessionID]; ok {
					continue
				}
				localTrack, err := webrtc.NewTrackLocalStaticRTP(pub.Codec, fmt.Sprintf("%s-copy-%s", trackID, currentPeerState.SessionID), pub.Remote.StreamID())
				if err != nil {
					log.Printf("error creating local copy track %s->%s: %v", trackID, currentPeerState.SessionID, err)
					continue
				}
				sender, err := pc.AddTrack(localTrack)
				if err != nil {
					log.Printf("error adding local copy to pc for %s: %v", currentPeerState.SessionID, err)
					continue
				}
				pub.Mu.Lock()
				pub.Locals[currentPeerState.SessionID] = &models.Track{Track: localTrack, Sender: sender}
				pub.Mu.Unlock()

				currentPeerState.TracksMu.Lock()
				currentPeerState.Tracks[trackID] = &models.Track{Track: localTrack, Sender: sender}
				currentPeerState.TracksMu.Unlock()
			}
			room.PubMutex.RUnlock()

			if err := sendOffer(currentPeerState, chatID); err != nil {
				log.Printf("error sending offer to %s: %v", currentPeerState.SessionID, err)
				continue
			}
		case "answer":
			log.Printf("[%s] ANSWER received", currentPeerState.SessionID)
			var answer webrtc.SessionDescription
			if err := json.Unmarshal(msg.Payload, &answer); err != nil {
				log.Printf("invalid answer payload from %s: %v", currentPeerState.SessionID, err)
				continue
			}
			if currentPeerState.PeerConnection == nil {
				log.Printf("received answer but PeerConnection == nil for %s", currentPeerState.SessionID)
				continue
			}

			if err := currentPeerState.PeerConnection.SetRemoteDescription(answer); err != nil {
				log.Printf("error SetRemoteDescription for %s: %v", currentPeerState.SessionID, err)
				continue
			}
			for _, c := range currentPeerState.PendingICE {
				if err := currentPeerState.PeerConnection.AddICECandidate(c); err != nil {
					log.Printf("failed to add buffered ICE for %s: %v", currentPeerState.SessionID, err)
				}
			}
			currentPeerState.PendingICE = nil
			log.Printf("Answer applied for %s", currentPeerState.SessionID)
		case "ice":
			var candidate webrtc.ICECandidateInit
			if err := json.Unmarshal(msg.Payload, &candidate); err != nil {
				log.Printf("invalid ice payload from %s: %v", currentPeerState.SessionID, err)
				continue
			}
			
			if currentPeerState.PeerConnection == nil {
				currentPeerState.PendingICE = append(currentPeerState.PendingICE, candidate)
				log.Printf("ICE buffered (no PC) for %s", currentPeerState.SessionID)
				continue
			}

			if currentPeerState.PeerConnection.RemoteDescription() == nil {
				currentPeerState.PendingICE = append(currentPeerState.PendingICE, candidate)
				log.Printf("ICE buffered (no remote desc) for %s", currentPeerState.SessionID)
				continue
			}
			if err := currentPeerState.PeerConnection.AddICECandidate(candidate); err != nil {
				log.Printf("error adding ice candidate for %s: %v", currentPeerState.SessionID, err)
			} else {
				log.Printf("Applied ICE candidate from %s", currentPeerState.SessionID)
			}
		case "leave":
			log.Printf("[%s] LEAVE room=%d", currentPeerState.SessionID, msg.RoomID)
			if currentPeerState.PeerConnection != nil {
				currentPeerState.PeerConnection.Close()
			}
			roomsMutex.RLock()
			room = rooms[key]
			roomsMutex.RUnlock()
			if room != nil {
				room.RemovePeer(currentPeerState.SessionID)
			}
		default:
			log.Printf("unknown message type from %s: %s", currentPeerState.SessionID, msg.Type)
		}
	}
}

func createPeerConnection(currentPeerState *models.PeerState, chatID int, roomsMutex *sync.RWMutex, rooms map[models.CallsKey]*models.Room, key models.CallsKey) (*webrtc.PeerConnection, error) {
	var NFound bool
	log.Println("🎥 Creating peer connection")
	config := webrtc.Configuration{
		ICEServers: []webrtc.ICEServer{
			{URLs: []string{"stun:stun.l.google.com:19302"}},
		},
	}
	pc, err := webrtc.NewPeerConnection(config)
	if err != nil {
		return nil, err
	}

	if _, err := pc.AddTransceiverFromKind(
    webrtc.RTPCodecTypeAudio,
		webrtc.RTPTransceiverInit{Direction: webrtc.RTPTransceiverDirectionRecvonly},
	); err != nil {
		log.Printf("warning: AddTransceiver audio failed: %v", err)
	}
	if _, err := pc.AddTransceiverFromKind(
		webrtc.RTPCodecTypeVideo,
		webrtc.RTPTransceiverInit{Direction: webrtc.RTPTransceiverDirectionRecvonly},
	); err != nil {
		log.Printf("warning: AddTransceiver video failed: %v", err)
	}

	pc.OnConnectionStateChange(func(state webrtc.PeerConnectionState) {
		log.Println("Connection state changed:", state)
		if state.String() == "failed" || state.String() == "closed" || state.String() == "disconnected" {
			NFound = true
			return
		}
		roomsMutex.RLock()
		room := rooms[key]
		roomsMutex.RUnlock()
		if room == nil {
			log.Printf("Room %d not found for track from %s", chatID, currentPeerState.SessionID)
			NFound = true
			return
		}

		log.Printf("[%s] Connection state: %s", currentPeerState.SessionID, state.String())
		if state == webrtc.PeerConnectionStateFailed || state == webrtc.PeerConnectionStateClosed {
			room.RemovePeer(currentPeerState.SessionID)
		}
	})

	if NFound {
		return pc, errors.ErrorNoCallRoom
	}

	pc.OnICECandidate(func(c *webrtc.ICECandidate) {
		if c == nil {
			log.Println("pc.OnICECandidate: gathering complete")
			return
		}
		
		b, _ := json.Marshal(c.ToJSON())
		msg := models.SignalMessage{ Type: "ice", Payload: b, RoomID: chatID, SessionID: currentPeerState.SessionID }

		currentPeerState.WSMu.Lock()
		if err := currentPeerState.WS.WriteJSON(msg); err != nil {
			log.Printf("❌ failed to WRITE ice to ws for %s: %v", currentPeerState.SessionID, err)
		} else {
			log.Printf("📤 Sent ICE to %s", currentPeerState.SessionID)
		}
		currentPeerState.WSMu.Unlock()
	})


	pc.OnTrack(func(track *webrtc.TrackRemote, receiver *webrtc.RTPReceiver) {
		log.Println("Track received:", track.ID(), track.Codec().MimeType)

		trackKey := track.ID()

		roomsMutex.RLock()
		room := rooms[key]
		roomsMutex.RUnlock()
		if room == nil {
			log.Printf("Room %d not found for track from %s", chatID, currentPeerState.SessionID)
			return
		}

		room.PubMutex.RLock()
		_, exists := room.Published[trackKey]
		if exists {
			room.PubMutex.RUnlock()
			log.Printf("Track %s already published, skipping", trackKey)
			return
		}

		pub := &models.PublishedTrack{
			Remote: track,
			Codec:  track.Codec().RTPCodecCapability,
			Locals: make(map[string]*models.Track),
		}
		room.Published[trackKey] = pub
		room.PubMutex.RUnlock()

		log.Println("New track published:", track.ID(), track.Codec().MimeType)

		go func(pub *models.PublishedTrack) {
			log.Printf("Start forwarding RTP for pub %s", pub.Remote.ID())
			for {
				pkt, _, err := pub.Remote.ReadRTP()
				if err != nil {
					log.Printf("Remote.ReadRTP error for %s: %v", pub.Remote.ID(), err)
					break
				}

				pub.Mu.RLock()
				for sid, lt := range pub.Locals {
					if lt == nil || lt.Track == nil {
						continue
					}
					if writeErr := lt.Track.WriteRTP(pkt); writeErr != nil {
						log.Printf("WriteRTP to local %s failed: %v", sid, writeErr)
					}
				}
				pub.Mu.RUnlock()
			}
			log.Printf("Stopped forwarding RTP for pub %s", pub.Remote.ID())
		}(pub)
		room.Mutex.RLock()
		for otherSessionID, otherPeer := range room.Peers {
			if otherSessionID == currentPeerState.SessionID || otherPeer.PeerConnection == nil {
				continue
			}

			otherPeer.TracksMu.RLock()
			_, alreadyAdded := otherPeer.Tracks[trackKey]
			otherPeer.TracksMu.RUnlock()
			
			if alreadyAdded {
				continue
			}

			localTrack, err := webrtc.NewTrackLocalStaticRTP(pub.Codec, fmt.Sprintf("%s-%s", track.ID(), otherSessionID), track.StreamID())
			if err != nil {
				log.Printf("Error creating local track for %s -> %s: %v", currentPeerState.SessionID, otherSessionID, err)
				continue
			}

			sender, err := otherPeer.PeerConnection.AddTrack(localTrack)
			if err != nil {
				log.Printf("Error adding track to peer %s: %v", otherSessionID, err)
				continue
			}

			pub.Mu.Lock()
			pub.Locals[otherSessionID] = &models.Track{
				Track: localTrack,
				Sender: sender,
			}
			pub.Mu.Unlock()

			otherPeer.TracksMu.Lock()
			otherPeer.Tracks[trackKey] = &models.Track{
				Track: localTrack,
				Sender: sender,
			}
			otherPeer.TracksMu.Unlock()

			log.Printf("Added track %s to peer %s", trackKey, otherSessionID)

			go func(peer *models.PeerState) {
				time.Sleep(500 * time.Millisecond)					
				if peer.PeerConnection.SignalingState() != webrtc.SignalingStateStable {
					return
				}

				offer, err := peer.PeerConnection.CreateOffer(nil)
				if err != nil {
					log.Printf("error creating renegotiation offer for %s: %v", otherSessionID, err)
					return
				}

				if err := peer.PeerConnection.SetLocalDescription(offer); err != nil {
					log.Printf("error setting local description for renegotiation %s: %v", otherSessionID, err)
					return
				}

				<-webrtc.GatheringCompletePromise(peer.PeerConnection)
				b, _ := json.Marshal(peer.PeerConnection.LocalDescription())
				peer.WSMu.Lock()
				defer peer.WSMu.Unlock()
				peer.WS.WriteJSON(models.SignalMessage{
					Type:       "offer",
					Payload:    b,
					RoomID:     chatID,
					SessionID:  otherSessionID,
				})
				log.Printf("Sent renegotiation offer to %s", otherSessionID)
			}(otherPeer)
		}
		room.Mutex.RUnlock()
	})

	return pc, nil
}

func sendOffer(currentPeerState *models.PeerState, roomID int) error {
	offer, err := currentPeerState.PeerConnection.CreateOffer(nil)
	if err != nil {
		log.Printf("Error creating offer for %s: %v", currentPeerState.SessionID, err)
		return err
	}

	if err := currentPeerState.PeerConnection.SetLocalDescription(offer); err != nil {
		log.Printf("Error setting local description for %s: %v", currentPeerState.SessionID, err)
		return err
	}

	<-webrtc.GatheringCompletePromise(currentPeerState.PeerConnection)

	local := currentPeerState.PeerConnection.LocalDescription()
	b, _ := json.Marshal(local)

	currentPeerState.WSMu.Lock()
	err = currentPeerState.WS.WriteJSON(models.SignalMessage{
		Type:    "offer",
		Payload: b,
		RoomID:  roomID,
		SessionID: currentPeerState.SessionID,
	})
	if err != nil {
		log.Printf("failed to WRITE offer to ws for %s: %v", currentPeerState.SessionID, err)
	} else {
		log.Printf("Offer sent to %s (room=%d)", currentPeerState.SessionID, roomID)
	}
	currentPeerState.WSMu.Unlock()
	return err
}
