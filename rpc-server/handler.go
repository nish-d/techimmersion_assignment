package main

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/TikTokTechImmersion/assignment_demo_2023/rpc-server/kitex_gen/rpc"
)

// IMServiceImpl implements the last service interface defined in the IDL.
type IMServiceImpl struct{}

func getGroupID(chat string) (string, error) {
	var groupID string

	lowercase := strings.ToLower(chat)
	senders := strings.Split(lowercase, "-")
	sender1, sender2 := senders[0], senders[1]
	// Compare the sender and receiver alphabetically, and sort them asc to form the room ID
	if comp := strings.Compare(sender1, sender2); comp == 1 {
		groupID = fmt.Sprintf("%s-%s", sender2, sender1)
	} else {
		groupID = fmt.Sprintf("%s-%s", sender1, sender2)
	}

	return groupID, nil
}

func (s *IMServiceImpl) Send(ctx context.Context, req *rpc.SendRequest) (*rpc.SendResponse, error) {
	resp := rpc.NewSendResponse()
	sendTime := time.Now()

	message := req.GetMessage()
	message.SetSendTime(sendTime.Unix())

	groupId, err := getGroupID(req.Message.Chat)
	redisClient.SaveMessageToRedis(ctx, groupId, message)
	resp.Code, resp.Msg = areYouLucky(err)
	return resp, err
}

func (s *IMServiceImpl) Pull(ctx context.Context, req *rpc.PullRequest) (*rpc.PullResponse, error) {
	resp := rpc.NewPullResponse()

	groupId, err := getGroupID(req.Chat)
	start := req.GetCursor()
	end := start + int64(req.GetLimit()) // did not minus 1 on purpose for hasMore check later on

	messages, err := redisClient.GetMessagesByGroupID(ctx, groupId, start, end, req.GetReverse())
	if err != nil {
		return nil, err
	}

	respMessages := make([]*rpc.Message, 0)
	var counter int32 = 0
	var nextCursor int64 = 0
	hasMore := false
	for _, msg := range messages {
		if counter+1 > req.GetLimit() {
			// having extra value here means it has more data
			hasMore = true
			nextCursor = end
			break // do not return the last message
		}
		temp := &rpc.Message{
			Chat:     req.GetChat(),
			Text:     msg.Text,
			Sender:   msg.Sender,
			SendTime: msg.SendTime,
		}
		respMessages = append(respMessages, temp)
		counter += 1
	}

	resp.Messages = respMessages
	resp.Code, resp.Msg = areYouLucky(err)
	resp.HasMore = &hasMore
	resp.NextCursor = &nextCursor

	return resp, nil
}

func areYouLucky(err error) (int32, string) {
	if err != nil {
		return 500, err.Error()
	} else {
		return 0, "success"
	}
}
