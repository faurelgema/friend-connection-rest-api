package friendship

import (
	"net/http"

	httpRes "friend_connection_rest_api/controller/common_respone"
	"friend_connection_rest_api/services/friendship"
	"friend_connection_rest_api/services/user"
	"friend_connection_rest_api/utils"

	"github.com/gin-gonic/gin"
)

func MakeFriendController(c *gin.Context, service friendship.FrienshipServices) {
	var reqFriend RequestFriend

	if err := c.BindJSON(&reqFriend); err != nil {
		c.JSON(http.StatusBadRequest, httpRes.HTTPError{Message: "BindJson Error, cause body request invalid"})
		return
	}

	if len(reqFriend.Friends) != 2 || reqFriend.Friends[0] == reqFriend.Friends[1] {
		c.JSON(http.StatusBadRequest, httpRes.HTTPError{Message: "Request Invalid"})
		return
	}

	firstUser := reqFriend.Friends[0]
	secondUser := reqFriend.Friends[1]

	if utils.ValidateEmail(firstUser) == false || utils.ValidateEmail(secondUser) == false {
		c.JSON(http.StatusBadRequest, httpRes.HTTPError{Message: "Email Invalid Format"})
		return
	}

	rs := service.MakeFriend(friendship.FrienshipServiceInput{RequestEmail: firstUser, TargetEmail: secondUser})

	if rs == nil {
		c.JSON(201, httpRes.HTTPSuccess{Success: true})
		return
	}

	c.JSON(400, httpRes.HTTPError{Message: rs.Error()})
}

func GetFriendsListController(c *gin.Context, service friendship.FrienshipServices) {
	email := RequestListFriends{}

	if err := c.BindJSON(&email); err != nil {
		c.JSON(http.StatusBadRequest, httpRes.HTTPError{Message: "BindJson Error, cause body request invalid"})
		return
	}

	if utils.ValidateEmail(email.Mail) == false {
		c.JSON(http.StatusBadRequest, httpRes.HTTPError{Message: "Email Invalid Format"})
		return
	}

	rs, err := service.GetFriendsList(user.Users{Email: email.Mail})

	if err != nil {
		c.JSON(400, httpRes.HTTPError{Message: err.Error()})
		return
	}

	c.JSON(200, toListFriendsStruct(rs))
}

func GetMutualFriendsController(c *gin.Context, service friendship.FrienshipServices) {
	reqFriend := RequestFriend{}

	if err := c.BindJSON(&reqFriend); err != nil {
		c.JSON(http.StatusBadRequest, httpRes.HTTPError{Message: "BindJson Error, cause body request invalid"})
		return
	}

	if len(reqFriend.Friends) != 2 || reqFriend.Friends[0] == reqFriend.Friends[1] {
		c.JSON(http.StatusBadRequest, httpRes.HTTPError{Message: "Request Invalid"})
		return
	}

	firstUser := reqFriend.Friends[0]
	secondUser := reqFriend.Friends[1]

	if utils.ValidateEmail(firstUser) == false || utils.ValidateEmail(secondUser) == false {
		c.JSON(http.StatusBadRequest, httpRes.HTTPError{Message: "Email Invalid Format"})
		return
	}

	rs, err := service.GetMutualFriendsList(friendship.FrienshipServiceInput{RequestEmail: firstUser, TargetEmail: secondUser})

	if err != nil {
		c.JSON(400, httpRes.HTTPError{Message: err.Error()})
		return
	}

	c.JSON(200, toListFriendsStruct(rs))
}

func SubscribeController(c *gin.Context, service friendship.FrienshipServices) {
	reqSubscribe := RequestUpdate{}

	if err := c.BindJSON(&reqSubscribe); err != nil {
		c.JSON(http.StatusBadRequest, httpRes.HTTPError{Message: "BindJson Error, cause body request invalid"})
		return
	}

	if reqSubscribe.Requestor == reqSubscribe.Target {
		c.JSON(http.StatusBadRequest, httpRes.HTTPError{Message: "Request Invalid"})
		return
	}

	firstUser := reqSubscribe.Requestor
	secondUser := reqSubscribe.Target

	if utils.ValidateEmail(firstUser) == false || utils.ValidateEmail(secondUser) == false {
		c.JSON(http.StatusBadRequest, httpRes.HTTPError{Message: "Email Invalid Format"})
		return
	}

	rs := service.Subscribe(friendship.FrienshipServiceInput{RequestEmail: firstUser, TargetEmail: secondUser})

	if rs != nil {
		c.JSON(400, httpRes.HTTPError{Message: rs.Error()})
		return
	}

	c.JSON(201, httpRes.HTTPSuccess{Success: true})
}

func BlockController(c *gin.Context, service friendship.FrienshipServices) {
	reqSubscribe := RequestUpdate{}

	if err := c.BindJSON(&reqSubscribe); err != nil {
		c.JSON(http.StatusBadRequest, httpRes.HTTPError{Message: "BindJson Error, cause body request invalid"})
		return
	}

	if reqSubscribe.Requestor == reqSubscribe.Target {
		c.JSON(http.StatusBadRequest, httpRes.HTTPError{Message: "Request Invalid"})
		return
	}

	firstUser := reqSubscribe.Requestor
	secondUser := reqSubscribe.Target

	if utils.ValidateEmail(firstUser) == false || utils.ValidateEmail(secondUser) == false {
		c.JSON(http.StatusBadRequest, httpRes.HTTPError{Message: "Email Invalid Format"})
		return
	}

	rs := service.Block(friendship.FrienshipServiceInput{RequestEmail: firstUser, TargetEmail: secondUser})

	if rs != nil {
		c.JSON(400, httpRes.HTTPError{Message: rs.Error()})
		return
	}

	c.JSON(201, httpRes.HTTPSuccess{Success: true})
}

func GetUsersReceiveUpdateController(c *gin.Context, service friendship.FrienshipServices) {
	reqRecvUpdate := RequestReceiveUpdate{}

	if err := c.BindJSON(&reqRecvUpdate); err != nil {
		c.JSON(http.StatusBadRequest, httpRes.HTTPError{Message: "BindJson Error, cause body request invalid"})
		return
	}

	if utils.ValidateEmail(reqRecvUpdate.Sender) == false {
		c.JSON(http.StatusBadRequest, httpRes.HTTPError{Message: "Email Invalid Format"})
		return
	}

	// rename
	mentionedUsers := utils.ExtractMentionEmail(reqRecvUpdate.Text)

	rs, err := service.GetUsersReceiveUpdate(reqRecvUpdate.Sender, mentionedUsers)

	if err != nil {
		c.JSON(400, httpRes.HTTPError{Message: err.Error()})
		return
	}

	c.JSON(200, toUsersCanReceiveUpdate(removeDuplicates(rs)))
}

func toListFriendsStruct(list []string) ResponeListFriends {
	listFriendsRespone := ResponeListFriends{}
	listFriendsRespone.Count = uint(len(list))
	listFriendsRespone.Success = true
	listFriendsRespone.Friends = append(listFriendsRespone.Friends, list...)
	return listFriendsRespone
}

func toUsersCanReceiveUpdate(list []string) ResponeReceiveUpdate {
	listUsersRecvUpdate := ResponeReceiveUpdate{}
	listUsersRecvUpdate.Success = true
	listUsersRecvUpdate.Recipients = append(listUsersRecvUpdate.Recipients, list...)
	return listUsersRecvUpdate
}

func removeDuplicates(elements []string) []string {
	// Use map to record duplicates as we find them.
	encountered := map[string]bool{}
	result := []string{}

	for v := range elements {
		if encountered[elements[v]] == true {
			// Do not add duplicate.
		} else {
			// Record this element as an encountered element.
			encountered[elements[v]] = true
			// Append to result slice.
			result = append(result, elements[v])
		}
	}
	// Return the new slice.
	return result
}
