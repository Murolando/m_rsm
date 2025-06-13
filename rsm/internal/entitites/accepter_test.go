package entitites

import (
	"log"
	"testing"

	"github.com/stretchr/testify/assert"
)

func setupLogs() []Log {
	return []Log{
		{
			LastAcceptedBallotNumber: &BallotNumber{
				Number:     2,
				ProposerID: 2,
			},
			LastAcceptedValue: "oleg",
		},
		{
			LastAcceptedBallotNumber: &BallotNumber{
				Number:     2,
				ProposerID: 2,
			},
			LastAcceptedValue: "meow",
		},
		{
			LastAcceptedBallotNumber: &BallotNumber{
				Number:     3,
				ProposerID: 2,
			},
			LastAcceptedValue: "temlana",
		},
	}
}
func TestPromise(t *testing.T) {
	t.Run("Promise", func(t *testing.T) {
		acceptor := NewAcceptor(4)
		acceptor.Logs = setupLogs()
		log.Println(acceptor.Logs)
		testPrepareMessage := &PrepareMessage{
			ProposerBallotNumber: &BallotNumber{
				Number:     11,
				ProposerID: 1,
			},
		}

		p := acceptor.Promise(testPrepareMessage)

		assert.Equal(t, true, p.Promise)
		assert.Equal(t, int64(11), p.MaxPromissedBallotNumber.Number)
	})
	t.Run("False promise", func(t *testing.T) {
		acceptor := NewAcceptor(2)
		acceptor.Logs = setupLogs()
		acceptor.MaxPromissedBallotNumber = &BallotNumber{
			Number:     12,
			ProposerID: 2,
		}
		log.Println(acceptor.Logs)
		testPrepareMessage := &PrepareMessage{
			ProposerBallotNumber: &BallotNumber{
				Number:     3,
				ProposerID: 1,
			},
		}

		p := acceptor.Promise(testPrepareMessage)

		assert.Equal(t, false, p.Promise)
		assert.Equal(t, int64(12), p.MaxPromissedBallotNumber.Number)
		assert.Equal(t, 2, p.MaxPromissedBallotNumber.ProposerID)
	})
}

func TestAccept(t *testing.T) {
	t.Run("Accept", func(t *testing.T) {
		acceptor := NewAcceptor(4)
		resultBallot := &BallotNumber{
			Number:     12,
			ProposerID: 1,
		}
		// init logs
		logs := setupLogs()
		logs = append(logs, Log{
			LastAcceptedBallotNumber: resultBallot,
			LastAcceptedValue:        "murolando",
		})
		logs[0].LastAcceptedBallotNumber = resultBallot
		logs[1].LastAcceptedBallotNumber = resultBallot
		logs[2].LastAcceptedBallotNumber = resultBallot

		acceptor.Logs = setupLogs()

		testProposeMessage := &ProposeMessage{
			NewAcceptedBallotNumber: resultBallot,
			NewAcceptedLog:          logs,
		}

		am := acceptor.Accept(testProposeMessage)

		assert.Equal(t, resultBallot, am.NewAcceptedBallotNumber)
		assert.Equal(t, 4, len(acceptor.Logs))
		assert.Equal(t, logs[3].LastAcceptedBallotNumber.Number, acceptor.Logs[3].LastAcceptedBallotNumber.Number)
		assert.Equal(t, logs[3].LastAcceptedValue, acceptor.Logs[3].LastAcceptedValue)
		assert.Equal(t, int64(12), acceptor.Logs[1].LastAcceptedBallotNumber.Number)

	})

	t.Run("Failed to Accept", func(t *testing.T) {
		acceptor := NewAcceptor(4)

		// init logs
		logs := setupLogs()
		logs = append(logs, Log{
			LastAcceptedBallotNumber: &BallotNumber{
				Number:     12,
				ProposerID: 1,
			},
			LastAcceptedValue: "murolando",
		})

		acceptor.Logs = setupLogs()
		acceptor.MaxPromissedBallotNumber = &BallotNumber{
			Number:     24,
			ProposerID: 2,
		}
		acceptor.LastAcceptedBallotNumber = &BallotNumber{
			Number:     24,
			ProposerID: 2,
		}

		testProposeMessage := &ProposeMessage{
			NewAcceptedBallotNumber: &BallotNumber{
				Number:     12,
				ProposerID: 1,
			},
			NewAcceptedLog: logs,
		}

		am := acceptor.Accept(testProposeMessage)
		assert.Equal(t, int64(24), am.NewAcceptedBallotNumber.Number)
		assert.NotEqual(t, len(logs), len(acceptor.Logs))
	})
}
