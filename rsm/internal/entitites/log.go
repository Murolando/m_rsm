package entitites

type Log struct {
	// (n_accepted,v_accepted)
	LastAcceptedBallotNumber *BallotNumber
	LastAcceptedValue        string
}

func NewEmptyLog() *Log {
	return &Log{
		LastAcceptedBallotNumber: &BallotNumber{},
		LastAcceptedValue:        DefaultValue,
	}
}
