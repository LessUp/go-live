Feature: WHIP Stream Publishing
  As a stream publisher
  I want to publish a WebRTC stream to a room via WHIP protocol

  Background:
    Given a running live-webrtc-go server at http://localhost:8080

  Scenario: Successfully publish stream with valid auth token
    Given room "test-room" exists
    And I have a valid auth token "secret-token"
    When I send a POST request to /api/whip/publish/test-room with SDP offer
    And header Authorization: Bearer secret-token
    Then I should receive a 200 OK response
    And response body should contain SDP answer
    And a PeerConnection should be created for the publisher

  Scenario: Publish stream fails with invalid auth token
    Given room "test-room" exists
    And I have an invalid auth token "wrong-token"
    When I send a POST request to /api/whip/publish/test-room with SDP offer
    And header Authorization: Bearer wrong-token
    Then I should receive a 401 Unauthorized response

  Scenario: Publish stream fails with missing auth token
    Given room "test-room" exists
    When I send a POST request to /api/whip/publish/test-room with SDP offer
    And no Authorization header
    Then I should receive a 401 Unauthorized response

  Scenario: Publish stream fails with invalid room name
    Given I have a valid auth token "secret-token"
    When I send a POST request to /api/whip/publish/invalid room! with SDP offer
    And header Authorization: Bearer secret-token
    Then I should receive a 400 Bad Request response

  Scenario: Publish stream with JWT auth
    Given room "test-room" exists
    And I have a valid JWT token signed with the correct secret
    And JWT contains role=admin claim
    When I send a POST request to /api/whip/publish/test-room with SDP offer
    And header Authorization: Bearer <jwt>
    Then I should receive a 200 OK response
    And response body should contain SDP answer

  Scenario Outline: Publish stream with per-room token
    Given room "<room>" exists
    And per-room token for <room> is "<token>"
    When I send a POST request to /api/whip/publish/<room> with SDP offer
    And header Authorization: Bearer <token>
    Then I should receive a 200 OK response

    Examples:
      | room        | token        |
      | room1       | tok1         |
      | my-stream   | my-secret    |
      | test_room_2 | another-tok  |

  Scenario: Publisher disconnect triggers room cleanup
    Given room "test-room" exists with active publisher
    And 2 viewers are subscribed to the room
    When the publisher disconnects
    Then the room should be cleaned up
    And all 2 viewers should be notified
    And all subscriber PeerConnections should be closed
