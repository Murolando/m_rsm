package entitites

import (
	"log"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func setupAcceptors() []*Acceptor {
	result := make([]*Acceptor, 3)
	result[0] = NewAcceptor(1)
	result[1] = NewAcceptor(2)
	result[2] = NewAcceptor(3)
	return result
}
func randomLogs(ballotNum BallotNumber) []Log {
	logs := setupLogs()
	logs[0].LastAcceptedBallotNumber = &ballotNum
	logs[1].LastAcceptedBallotNumber = &ballotNum
	logs[2].LastAcceptedBallotNumber = &ballotNum
	return logs
}
func TestPrepare(t *testing.T) {
	t.Run("Quorum", func(t *testing.T) {
		acceptors := setupAcceptors()
		proposer := NewProposer(1, acceptors, time.Duration(3))

		testPrepareMessage := &PrepareMessage{
			ProposerBallotNumber: &BallotNumber{
				Number:     11,
				ProposerID: 1,
			},
		}
		promises, promise, err := proposer.Prepare(testPrepareMessage)

		require.NoError(t, err)
		assert.True(t, promise)
		assert.GreaterOrEqual(t, len(promises), 2)
	})
	t.Run("Waiter expired", func(t *testing.T) {
		acceptors := setupAcceptors()
		proposer := NewProposer(1, acceptors[2:], time.Duration(0))

		testPrepareMessage := &PrepareMessage{
			ProposerBallotNumber: &BallotNumber{
				Number:     11,
				ProposerID: 1,
			},
		}
		_, promise, err := proposer.Prepare(testPrepareMessage)

		assert.Error(t, err)
		assert.Equal(t, ProposerWaiterExpiredError, err.Error())
		assert.False(t, promise)
	})
	t.Run("False", func(t *testing.T) {
		acceptors := setupAcceptors()

		acceptors[0].MaxPromissedBallotNumber.Number = 22
		acceptors[1].MaxPromissedBallotNumber.Number = 22
		acceptors[2].MaxPromissedBallotNumber.Number = 22

		proposer := NewProposer(1, acceptors, time.Duration(3))
		testPrepareMessage := &PrepareMessage{
			ProposerBallotNumber: &BallotNumber{
				Number:     1,
				ProposerID: 1,
			},
		}
		_, promise, err := proposer.Prepare(testPrepareMessage)

		assert.NoError(t, err)
		assert.False(t, promise)
	})
	t.Run("Parallel prepare", func(t *testing.T) {
		acceptors := setupAcceptors()

		proposer1 := NewProposer(1, acceptors, time.Duration(3))
		proposer2 := NewProposer(2, acceptors, time.Duration(3))

		acceptors[0].MaxPromissedBallotNumber.Number = 22
		acceptors[1].MaxPromissedBallotNumber.Number = 22
		acceptors[2].MaxPromissedBallotNumber.Number = 22

		testPrepareMessage1 := &PrepareMessage{
			ProposerBallotNumber: &BallotNumber{
				Number:     11,
				ProposerID: 1,
			},
		}

		testPrepareMessage2 := &PrepareMessage{
			ProposerBallotNumber: &BallotNumber{
				Number:     33,
				ProposerID: 2,
			},
		}

		wg := sync.WaitGroup{}
		wg.Add(2)
		go func() {
			promises, promise, err := proposer2.Prepare(testPrepareMessage2)
			require.NoError(t, err)
			assert.True(t, promise)
			assert.GreaterOrEqual(t, len(promises), 2)
			wg.Done()
		}()
		go func() {
			_, promise, err := proposer1.Prepare(testPrepareMessage1)
			assert.NoError(t, err)
			assert.False(t, promise)
			wg.Done()
		}()
		wg.Wait()
		assert.Equal(t, int64(33), proposer1.acceptors[0].MaxPromissedBallotNumber.Number)
	})
}

func TestPropose(t *testing.T) {
	t.Run("Quorum", func(t *testing.T) {
		acceptors := setupAcceptors()
		proposer := NewProposer(1, acceptors, time.Duration(1))

		testProposeMessage := &ProposeMessage{
			NewAcceptedBallotNumber: &BallotNumber{
				Number:     11,
				ProposerID: 1,
			},
		}
		promises, accepted, err := proposer.Propose(testProposeMessage)

		require.NoError(t, err)
		assert.True(t, accepted)
		assert.GreaterOrEqual(t, len(promises), 2)
	})
	t.Run("Waiter expired", func(t *testing.T) {
		acceptors := setupAcceptors()
		proposer := NewProposer(1, acceptors, time.Duration(0))

		testProposeMessage := &ProposeMessage{
			NewAcceptedBallotNumber: &BallotNumber{
				Number:     11,
				ProposerID: 1,
			},
		}
		_, accepted, err := proposer.Propose(testProposeMessage)

		require.Error(t, err)
		require.Equal(t, err.Error(), ProposerWaiterExpiredError)
		assert.False(t, accepted)
	})
	t.Run("Random Quorum", func(t *testing.T) {
		acceptors := setupAcceptors()

		proposer := NewProposer(1, acceptors, time.Duration(1))

		testProposeMessage := &ProposeMessage{
			NewAcceptedBallotNumber: &BallotNumber{
				Number:     11,
				ProposerID: 1,
			},
		}
		messages, accepted, err := proposer.Propose(testProposeMessage)

		require.NoError(t, err)
		assert.True(t, accepted)
		assert.GreaterOrEqual(t, len(messages), 2)
	})
	t.Run("False", func(t *testing.T) {
		acceptors := setupAcceptors()

		acceptors[0].MaxPromissedBallotNumber.Number = 22
		acceptors[1].MaxPromissedBallotNumber.Number = 22
		acceptors[2].MaxPromissedBallotNumber.Number = 22

		proposer := NewProposer(1, acceptors, time.Duration(1))
		testProposeMessage := &ProposeMessage{
			NewAcceptedBallotNumber: &BallotNumber{
				Number:     1,
				ProposerID: 1,
			},
		}
		_, accepted, err := proposer.Propose(testProposeMessage)

		assert.NoError(t, err)
		assert.False(t, accepted)
	})
	t.Run("Parallel propose", func(t *testing.T) {
		acceptors := setupAcceptors()

		proposer1 := NewProposer(1, acceptors, time.Duration(1))
		proposer2 := NewProposer(2, acceptors, time.Duration(1))

		acceptors[0].MaxPromissedBallotNumber.Number = 22
		acceptors[1].MaxPromissedBallotNumber.Number = 22
		acceptors[2].MaxPromissedBallotNumber.Number = 22

		testProposeMessage1 := &ProposeMessage{
			NewAcceptedBallotNumber: &BallotNumber{
				Number:     11,
				ProposerID: 1,
			},
		}

		testProposeMessage2 := &ProposeMessage{
			NewAcceptedBallotNumber: &BallotNumber{
				Number:     22,
				ProposerID: 2,
			},
		}
		wg := sync.WaitGroup{}
		wg.Add(2)
		go func() {
			messages, accepted, err := proposer2.Propose(testProposeMessage2)
			require.NoError(t, err)
			assert.True(t, accepted)
			assert.GreaterOrEqual(t, len(messages), 2)
			log.Println(messages)
			wg.Done()
		}()
		go func() {
			_, accepted, err := proposer1.Propose(testProposeMessage1)
			assert.NoError(t, err)
			assert.False(t, accepted)
			wg.Done()
		}()
		wg.Wait()
		assert.Equal(t, int64(22), proposer1.acceptors[0].LastAcceptedBallotNumber.Number)
	})
}

func TestProposeValue(t *testing.T) {
	t.Run("Test Log builder logic", func(t *testing.T) {
		acceptors := setupAcceptors()

		acceptors[0].LastAcceptedBallotNumber.Number = 21
		acceptors[1].LastAcceptedBallotNumber.Number = 15
		acceptors[2].LastAcceptedBallotNumber.Number = 7

		proposer := NewProposer(1, acceptors, time.Duration(1))
		proposer.LastBallotNumber = &BallotNumber{
			Number:     22,
			ProposerID: 1,
		}
		acceptors[0].Logs = randomLogs(BallotNumber{Number: 21, ProposerID: 1})

		acceptors[1].Logs = randomLogs(BallotNumber{Number: 15, ProposerID: 1})[:2]

		acceptors[2].Logs = randomLogs(BallotNumber{Number: 7, ProposerID: 1})[:1]

		logs := randomLogs(BallotNumber{Number: 23, ProposerID: 1})
		logs = append(logs, Log{
			LastAcceptedBallotNumber: &BallotNumber{
				Number:     23,
				ProposerID: 1,
			},
			LastAcceptedValue: "murolando",
		})
		_, err := proposer.ProposeValue("murolando")
		require.NoError(t, err)
		assert.Equal(t, logs, acceptors[0].Logs)
		assert.Equal(t, int64(23), acceptors[1].LastAcceptedBallotNumber.Number)
	})
	t.Run("Test Retries to propose", func(t *testing.T) {
		acceptors := setupAcceptors()

		acceptors[0].LastAcceptedBallotNumber.Number = 21
		acceptors[0].MaxPromissedBallotNumber.Number = 21

		acceptors[1].LastAcceptedBallotNumber.Number = 21
		acceptors[1].MaxPromissedBallotNumber.Number = 21

		acceptors[2].LastAcceptedBallotNumber.Number = 7
		acceptors[2].MaxPromissedBallotNumber.Number = 7

		proposer := NewProposer(1, acceptors, time.Duration(1))
		proposer.LastBallotNumber = &BallotNumber{
			Number:     19,
			ProposerID: 1,
		}
		acceptors[0].Logs = randomLogs(BallotNumber{Number: 21, ProposerID: 1})

		acceptors[1].Logs = randomLogs(BallotNumber{Number: 21, ProposerID: 1})[:2]

		acceptors[2].Logs = randomLogs(BallotNumber{Number: 7, ProposerID: 1})[:1]

		logs := randomLogs(BallotNumber{Number: 22, ProposerID: 1})
		logs = append(logs, Log{
			LastAcceptedBallotNumber: &BallotNumber{
				Number:     22,
				ProposerID: 1,
			},
			LastAcceptedValue: "murolando",
		})
		for i := range 5 {
			log.Println(i, proposer.LastBallotNumber.Number)
			_, err := proposer.ProposeValue("murolando")
			if err == nil {
				break
			}
			assert.Equal(t, err.Error(), FailedToProposeNewValue)
		}

		assert.Equal(t, logs, acceptors[1].Logs)
		assert.Equal(t, int64(22), acceptors[1].LastAcceptedBallotNumber.Number)
	})

	t.Run("Test Parallel propose", func(t *testing.T) {
		acceptors := setupAcceptors()

		acceptors[0].LastAcceptedBallotNumber.Number = 21
		acceptors[0].MaxPromissedBallotNumber.Number = 21

		acceptors[1].LastAcceptedBallotNumber.Number = 21
		acceptors[1].MaxPromissedBallotNumber.Number = 21

		acceptors[2].LastAcceptedBallotNumber.Number = 7
		acceptors[2].MaxPromissedBallotNumber.Number = 7

		proposer1 := NewProposer(1, acceptors, time.Duration(1))
		proposer1.LastBallotNumber = &BallotNumber{
			Number:     19,
			ProposerID: 1,
		}

		proposer2 := NewProposer(1, acceptors, time.Duration(1))
		proposer2.LastBallotNumber = &BallotNumber{
			Number:     227,
			ProposerID: 1,
		}

		acceptors[0].Logs = randomLogs(BallotNumber{Number: 21, ProposerID: 1})

		acceptors[1].Logs = randomLogs(BallotNumber{Number: 21, ProposerID: 1})[:2]

		acceptors[2].Logs = randomLogs(BallotNumber{Number: 7, ProposerID: 1})[:1]

		logs := randomLogs(BallotNumber{Number: 228, ProposerID: 1})
		logs = append(logs, Log{
			LastAcceptedBallotNumber: &BallotNumber{
				Number:     228,
				ProposerID: 1,
			},
			LastAcceptedValue: "murolando",
		})

		wg := sync.WaitGroup{}
		wg.Add(2)
		go func() {
			_, err := proposer2.ProposeValue("murolando")
			require.NoError(t, err)
			wg.Done()
		}()
		go func() {
			_, err := proposer1.ProposeValue("murolando")
			assert.Error(t, err)
			wg.Done()
		}()
		wg.Wait()

		assert.Equal(t, logs, acceptors[1].Logs)
		assert.Equal(t, int64(228), acceptors[1].LastAcceptedBallotNumber.Number)
	})
}
