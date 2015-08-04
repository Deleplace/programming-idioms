package pigae

import (
	. "github.com/Deleplace/programming-idioms/pig"

	"appengine"
	"appengine/datastore"
)

// VotesAccessorMock should be useful for unit testing. But it is dead code for now.
type VotesAccessorMock struct {
}

func (va VotesAccessorMock) idiomVote(c appengine.Context, vote IdiomVoteLog, nickname string) (int, error) {
	return -1, nil
}

func (va VotesAccessorMock) implVote(c appengine.Context, vote ImplVoteLog, nickname string) (int, error) {
	return -1, nil
}

func (va VotesAccessorMock) createVoteIdiomLogNicknameAncestorKey(c appengine.Context, nickname string) (key *datastore.Key) {
	return
}

func (va VotesAccessorMock) getIdiomVotes(c appengine.Context, nickname string) (keys []*datastore.Key, votes []*IdiomVoteLog, err error) {
	return nil, nil, nil
}

func (va VotesAccessorMock) getIdiomVote(c appengine.Context, nickname string, idiomID int) (key *datastore.Key, vote *IdiomVoteLog, err error) {
	return nil, nil, nil
}

func (va VotesAccessorMock) saveIdiomVoteOrRemove(c appengine.Context, vote IdiomVoteLog, nickname string) (*datastore.Key, error) {
	return nil, nil
}

func (va VotesAccessorMock) createVoteImplLogNicknameAncestorKey(c appengine.Context, nickname string) (key *datastore.Key) {
	return
}

func (va VotesAccessorMock) getImplVotes(c appengine.Context, nickname string) (keys []*datastore.Key, votes []*ImplVoteLog, err error) {
	return
}

func (va VotesAccessorMock) getImplVote(c appengine.Context, nickname string, implID int) (key *datastore.Key, vote *ImplVoteLog, err error) {
	return
}

func (va VotesAccessorMock) saveImplVoteOrRemove(c appengine.Context, vote ImplVoteLog, nickname string) (*datastore.Key, error) {
	return nil, nil
}

func (va VotesAccessorMock) decorateIdiom(c appengine.Context, idiom *Idiom, username string) {
}

func (va VotesAccessorMock) decorateImpl(c appengine.Context, impl *Impl, username string) {
}
