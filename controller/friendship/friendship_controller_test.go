package friendship

import (
	"bytes"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"friend_connection_rest_api/services/friendship"
	"friend_connection_rest_api/services/user"
	"friend_connection_rest_api/utils"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestMakeFriendController(t *testing.T) {
	// Given
	testCase := []struct {
		// scenario
		scenario          string
		input             RequestFriend
		expectedErrorBody string
	}{
		{
			scenario: "Make Friend Success",
			input: RequestFriend{
				Friends: []string{
					"gema1@gmail.com",
					"gema2@gmail.com",
				},
			},
			expectedErrorBody: "",
		},
		{
			scenario: "Make Friend Fail",
			input: RequestFriend{
				Friends: []string{
					"arel1@gmail.com",
					"arel2@gmail.com",
				},
			},
			expectedErrorBody: "Any Error",
		},
		{
			scenario: "Not enough parameters",
			input: RequestFriend{
				Friends: []string{
					"gema@gmail.com",
				},
			},
			expectedErrorBody: "Request Invalid",
		},
		{
			scenario: "Same user",
			input: RequestFriend{
				Friends: []string{
					"faurelgema@gmail.com",
					"faurelgema@gmail.com",
				},
			},
			expectedErrorBody: "Request Invalid",
		},
		{
			scenario: "Email Invalid",
			input: RequestFriend{
				Friends: []string{
					"abc",
					"xyz",
				},
			},
			expectedErrorBody: "Email Invalid Format",
		},
		{
			scenario:          "Empty request body",
			expectedErrorBody: "BindJson Error, cause body request invalid",
		},
	}

	for _, tc := range testCase {
		t.Run(tc.scenario, func(t *testing.T) {

			frienshipMock := new(friendship.FrienshipMockService)
			if tc.input.Friends != nil {
				if tc.scenario != "Make Friend Fail" {
					if len(tc.input.Friends) == 2 {
						frienshipMock.On("MakeFriend", friendship.FrienshipServiceInput{RequestEmail: tc.input.Friends[0], TargetEmail: tc.input.Friends[1]}).Return(nil)
					} else {
						frienshipMock.On("MakeFriend", friendship.FrienshipServiceInput{RequestEmail: tc.input.Friends[0]}).Return(nil)
					}
				} else {
					frienshipMock.On("MakeFriend", friendship.FrienshipServiceInput{RequestEmail: tc.input.Friends[0], TargetEmail: tc.input.Friends[1]}).Return(errors.New("Any Error"))
				}
			}

			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)

			values := map[string][]string{"friends": tc.input.Friends}
			jsonValue, _ := json.Marshal(values)
			c.Request, _ = http.NewRequest("POST", "/add-friends", bytes.NewBuffer(jsonValue))
			c.Request.Header.Set("Content-Type", "application/json")
			// When

			MakeFriendController(c, frienshipMock)

			// Then
			var actualResult map[string]interface{}
			body, _ := ioutil.ReadAll(w.Result().Body)
			json.Unmarshal(body, &actualResult)

			if val1, ok1 := actualResult["success"]; ok1 {
				assert.Equal(t, val1, true)
			} else if val2, ok2 := actualResult["error"]; ok2 {
				assert.Equal(t, http.StatusBadRequest, w.Result().StatusCode)
				assert.Equal(t, val2, tc.expectedErrorBody)
			}

		})
	}

}

func TestGetFriendsList(t *testing.T) {

	// Given

	testCase := []struct {
		scenario            string
		input               user.Users
		mockRespone         []string
		mockError           error
		expectedErrorBody   string
		expectedSuccessBody string
	}{
		{
			scenario:            "Get List Friends Success",
			input:               user.Users{Email: "abc@gmail.com"},
			mockRespone:         []string{"gema1@gmail.com", "gema2@gmail.com", "gema3@yahoo.com"},
			expectedSuccessBody: `{"success":true,"friends":["gema1@gmail.com","gema2@gmail.com","gema3@yahoo.com"],"count":3}`,
		},
		{
			scenario:          "Get List Friends Fail",
			input:             user.Users{Email: "abcxxx@gmail.com"},
			mockError:         errors.New("Any error"),
			mockRespone:       nil,
			expectedErrorBody: `{"error":"Any error"}`,
		},
		{
			scenario:          "Invalid Email",
			input:             user.Users{Email: "abc"},
			mockRespone:       nil,
			expectedErrorBody: `{"error":"Email Invalid Format"}`,
		},
		{
			scenario:          "Empty request body",
			expectedErrorBody: `{"error":"BindJson Error, cause body request invalid"}`,
		},
	}

	for _, tc := range testCase {
		t.Run(tc.scenario, func(t *testing.T) {
			mockFriendship := new(friendship.FrienshipMockService)
			mockFriendship.On("GetFriendsList", tc.input).Return(tc.mockRespone, tc.mockError)

			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)

			values := map[string]string{"Email": tc.input.Email}
			jsonValue, _ := json.Marshal(values)
			c.Request, _ = http.NewRequest("POST", "/get-list-friends", bytes.NewBuffer(jsonValue))

			// When
			GetFriendsListController(c, mockFriendship)

			// Then
			var actualResult string
			body, _ := ioutil.ReadAll(w.Result().Body)
			actualResult = string(body)

			if tc.scenario == "Get List Friends Success" {
				assert.Equal(t, 200, w.Result().StatusCode)
				assert.Equal(t, tc.expectedSuccessBody, actualResult)
			} else {
				assert.Equal(t, http.StatusBadRequest, w.Result().StatusCode)
				assert.Equal(t, tc.expectedErrorBody, actualResult)
			}
		})
	}

}

func TestGetMutualFriendsController(t *testing.T) {
	testCase := []struct {
		scenario            string
		requestInput        RequestFriend
		mockRespone         []string
		mockError           error
		expectedErrorBody   string
		expectedSuccessBody string
	}{
		{
			scenario: "Get Mutual Friends Success",
			requestInput: RequestFriend{
				Friends: []string{
					"requestor@gmail.com",
					"target@gmail.com",
				},
			},
			mockRespone: []string{
				"mutual1@gmail.com",
				"mutual2@gmail.com",
				"mutual3@gmail.com",
			},
			expectedSuccessBody: `{"success":true,"friends":["mutual1@gmail.com","mutual2@gmail.com","mutual3@gmail.com"],"count":3}`,
		},
		{
			scenario: "Get Mutual Friends Fail",
			requestInput: RequestFriend{
				Friends: []string{
					"requestor@gmail.com",
					"target@gmail.com",
				},
			},
			mockError:         errors.New("Any error"),
			expectedErrorBody: `{"error":"Any error"}`,
		},
		{
			scenario: "Invalid Email",
			requestInput: RequestFriend{
				Friends: []string{
					"requestor",
					"target@gmail.com",
				},
			},
			expectedErrorBody: `{"error":"Email Invalid Format"}`,
		},
		{
			scenario: "requestor same target",
			requestInput: RequestFriend{
				Friends: []string{
					"target@gmail.com",
					"target@gmail.com",
				},
			},
			expectedErrorBody: `{"error":"Request Invalid"}`,
		},
		{
			scenario: "Not enough parameters",
			requestInput: RequestFriend{
				Friends: []string{
					"requestor@gmail.com",
				},
			},
			expectedErrorBody: `{"error":"Request Invalid"}`,
		},
		{
			scenario:          "Empty request body",
			expectedErrorBody: `{"error":"BindJson Error, cause body request invalid"}`,
		},
	}

	for _, tc := range testCase {
		t.Run(tc.scenario, func(t *testing.T) {
			mockFriendship := new(friendship.FrienshipMockService)
			if tc.requestInput.Friends != nil {
				if tc.scenario == "Not enough parameters" {
					mockFriendship.On("GetMutualFriendsList", friendship.FrienshipServiceInput{RequestEmail: tc.requestInput.Friends[0]}).Return(tc.mockRespone, tc.mockError)
				} else {
					mockFriendship.On("GetMutualFriendsList", friendship.FrienshipServiceInput{RequestEmail: tc.requestInput.Friends[0], TargetEmail: tc.requestInput.Friends[1]}).Return(tc.mockRespone, tc.mockError)
				}
			}

			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)
			jsonValue, _ := json.Marshal(tc.requestInput)

			c.Request, _ = http.NewRequest("POST", "/get-mutual-list-friends", bytes.NewBuffer(jsonValue))

			// When
			GetMutualFriendsController(c, mockFriendship)

			//Then

			body, _ := ioutil.ReadAll(w.Result().Body)
			actualResult := string(body)

			if tc.scenario == "Get Mutual Friends Success" {
				assert.Equal(t, 200, w.Result().StatusCode)
				assert.Equal(t, tc.expectedSuccessBody, actualResult)
			} else {
				assert.Equal(t, http.StatusBadRequest, w.Result().StatusCode)
				assert.Equal(t, tc.expectedErrorBody, actualResult)
			}
		})
	}
}

func TestSubscribeController(t *testing.T) {

	// Given

	testCase := []struct {
		scenario            string
		inputRequest        RequestUpdate
		mockError           error
		expectedErrorBody   string
		expectedSuccessBody string
	}{
		{
			scenario: "Subscribe Success",
			inputRequest: RequestUpdate{
				Requestor: "requestor@gmail.com",
				Target:    "target@gmail.com",
			},
			expectedSuccessBody: `{"success":true}`,
		},
		{
			scenario: "Subscribe Fail",
			inputRequest: RequestUpdate{
				Requestor: "requestor@gmail.com",
				Target:    "target@gmail.com",
			},
			mockError:         errors.New("Any error"),
			expectedErrorBody: `{"error":"Any error"}`,
		},
		{
			scenario: "Invalid Mail",
			inputRequest: RequestUpdate{
				Requestor: "requestor",
				Target:    "target",
			},
			expectedErrorBody: `{"error":"Email Invalid Format"}`,
		},
		{
			scenario: "Not enough parameters",
			inputRequest: RequestUpdate{
				Requestor: "requestor@gmail.com",
			},
			expectedErrorBody: `{"error":"BindJson Error, cause body request invalid"}`,
		},
	}

	for _, tc := range testCase {
		t.Run(tc.scenario, func(t *testing.T) {
			mockFriendship := new(friendship.FrienshipMockService)
			mockFriendship.On("Subscribe", friendship.FrienshipServiceInput{RequestEmail: tc.inputRequest.Requestor, TargetEmail: tc.inputRequest.Target}).Return(tc.mockError)

			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)

			jsonVal, _ := json.Marshal(tc.inputRequest)
			c.Request, _ = http.NewRequest("POST", "/subscribe", bytes.NewBuffer(jsonVal))

			// When

			SubscribeController(c, mockFriendship)

			// Then

			body, _ := ioutil.ReadAll(w.Result().Body)
			actualResult := string(body)

			if tc.scenario == "Subscribe Success" {
				assert.Equal(t, 201, w.Result().StatusCode)
				assert.Equal(t, tc.expectedSuccessBody, actualResult)
			} else {
				assert.Equal(t, http.StatusBadRequest, w.Result().StatusCode)
				assert.Equal(t, tc.expectedErrorBody, actualResult)
			}
		})
	}
}

func TestBlockController(t *testing.T) {
	// Given

	testCase := []struct {
		scenario            string
		inputRequest        RequestUpdate
		mockError           error
		expectedErrorBody   string
		expectedSuccessBody string
	}{
		{
			scenario: "Block Success",
			inputRequest: RequestUpdate{
				Requestor: "requestor@gmail.com",
				Target:    "target@gmail.com",
			},
			expectedSuccessBody: `{"success":true}`,
		},
		{
			scenario: "Block Fail",
			inputRequest: RequestUpdate{
				Requestor: "requestor@gmail.com",
				Target:    "target@gmail.com",
			},
			mockError:         errors.New("Any error"),
			expectedErrorBody: `{"error":"Any error"}`,
		},
		{
			scenario: "Invalid Mail",
			inputRequest: RequestUpdate{
				Requestor: "requestor",
				Target:    "target",
			},
			expectedErrorBody: `{"error":"Email Invalid Format"}`,
		},
		{
			scenario: "Not enough parameters",
			inputRequest: RequestUpdate{
				Requestor: "requestor@gmail.com",
			},
			expectedErrorBody: `{"error":"BindJson Error, cause body request invalid"}`,
		},
	}

	for _, tc := range testCase {
		t.Run(tc.scenario, func(t *testing.T) {
			mockFriendship := new(friendship.FrienshipMockService)
			mockFriendship.On("Block", friendship.FrienshipServiceInput{RequestEmail: tc.inputRequest.Requestor, TargetEmail: tc.inputRequest.Target}).Return(tc.mockError)

			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)

			jsonVal, _ := json.Marshal(tc.inputRequest)
			c.Request, _ = http.NewRequest("POST", "/block", bytes.NewBuffer(jsonVal))

			// When

			BlockController(c, mockFriendship)

			// Then

			body, _ := ioutil.ReadAll(w.Result().Body)
			actualResult := string(body)

			if tc.scenario == "Block Success" {
				assert.Equal(t, 201, w.Result().StatusCode)
				assert.Equal(t, tc.expectedSuccessBody, actualResult)
			} else {
				assert.Equal(t, http.StatusBadRequest, w.Result().StatusCode)
				assert.Equal(t, tc.expectedErrorBody, actualResult)
			}
		})
	}
}

func TestGetUsersReceiveUpdateController(t *testing.T) {
	// Given
	testCase := []struct {
		scenario            string
		inputRequest        *RequestReceiveUpdate
		mockRespone         []string
		mockError           error
		expectedErrorBody   string
		expectedSuccessBody string
	}{
		{
			scenario: "Receive Success",
			inputRequest: &RequestReceiveUpdate{
				Sender: "arel@gmail.com",
				Text:   "Hello world!, hi @gema@yahoo.com",
			},
			mockRespone: []string{
				"gema@yahoo.com",
				"muhammadfaurel.augistta@gmail.com",
				"faurellorentermcastilla@gmail.com",
			},
			expectedSuccessBody: `{"success":true,"recipients":["gema@yahoo.com","muhammadfaurel.augistta@gmail.com","faurellorentermcastilla@gmail.com"]}`,
		},
		{
			scenario: "Receive Fail",
			inputRequest: &RequestReceiveUpdate{
				Sender: "muhammadfaurel.augistta@gmail.com",
				Text:   "Hello world!, hi @gema@yahoo.com",
			},
			mockError:         errors.New("Any error"),
			expectedErrorBody: `{"error":"Any error"}`,
		},
		{
			scenario: "Invalid Email",
			inputRequest: &RequestReceiveUpdate{
				Sender: "arel",
				Text:   "Hello world!, hi @muhammadfaurel@yahoo.com",
			},
			expectedErrorBody: `{"error":"Email Invalid Format"}`,
		},
		{
			scenario:          "Empty request body",
			expectedErrorBody: `{"error":"BindJson Error, cause body request invalid"}`,
		},
	}

	for _, tc := range testCase {
		t.Run(tc.scenario, func(t *testing.T) {
			mockFriendship := new(friendship.FrienshipMockService)
			if tc.inputRequest != nil {
				mentioned := utils.ExtractMentionEmail(tc.inputRequest.Text)
				mockFriendship.On("GetUsersReceiveUpdate", tc.inputRequest.Sender, mentioned).Return(tc.mockRespone, tc.mockError)
			}

			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)

			jsonVal, _ := json.Marshal(tc.inputRequest)
			c.Request, _ = http.NewRequest("POST", "/get-list-users-receive-update", bytes.NewBuffer(jsonVal))

			// When

			GetUsersReceiveUpdateController(c, mockFriendship)

			// Then

			body, _ := ioutil.ReadAll(w.Result().Body)
			actualResult := string(body)

			if tc.scenario == "Recvice Success" {
				assert.Equal(t, 200, w.Result().StatusCode)
				assert.Equal(t, tc.expectedSuccessBody, actualResult)
			} else {
				assert.Equal(t, http.StatusBadRequest, w.Result().StatusCode)
				assert.Equal(t, tc.expectedErrorBody, actualResult)
			}
		})
	}
}
