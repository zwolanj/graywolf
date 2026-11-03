//go:build android

package platformsvc

import pb "github.com/chrissnell/graywolf/pkg/platformproto"

// MessageNotification is the payload for Client.NotifyMessage.
type MessageNotification = pb.MessageNotification

// NotifyMessage asks the Kotlin side to post a system notification. It is
// fire-and-forget: there is no ack, and a send error just means it was dropped.
func (c *clientImpl) NotifyMessage(n *MessageNotification) error {
	return c.send(&pb.PlatformMessage{Body: &pb.PlatformMessage_MessageNotification{
		MessageNotification: n,
	}})
}
