package pigae

import (
	. "github.com/Deleplace/programming-idioms/pig"

	"golang.org/x/net/context"
	"google.golang.org/appengine/datastore"
)

// VotesAccessorMock should be useful for unit testing. But it is dead code for now.
type VotesAccessorMock struct {
}

func (va VotesAccessorMock) idiomVote(c context.Context, vote IdiomVoteLog, nickname string) (int, error) {
	return -1, nil
}

func (va VotesAccessorMock) implVote(c context.Context, vote ImplVoteLog, nickname string) (int, error) {
	return -1, nil
}

func (va VotesAccessorMock) createVoteIdiomLogNicknameAncestorKey(c context.Context, nickname string) (key *datastore.Key) {
	return
}

func (va VotesAccessorMock) getIdiomVotes(c context.Context, nickname string) (keys []*datastore.Key, votes []*IdiomVoteLog, err error) {
	return nil, nil, nil
}

func (va VotesAccessorMock) getIdiomVote(c context.Context, nickname string, idiomID int) (key *datastore.Key, vote *IdiomVoteLog, err error) {
	return nil, nil, nil
}

func (va VotesAccessorMock) saveIdiomVoteOrRemove(c context.Context, vote IdiomVoteLog, nickname string) (*datastore.Key, error) {
	return nil, nil
}

func (va VotesAccessorMock) createVoteImplLogNicknameAncestorKey(c context.Context, nickname string) (key *datastore.Key) {
	return
}

func (va VotesAccessorMock) getImplVotes(c context.Context, nickname string) (keys []*datastore.Key, votes []*ImplVoteLog, err error) {
	return
}

func (va VotesAccessorMock) getImplVote(c context.Context, nickname string, implID int) (key *datastore.Key, vote *ImplVoteLog, err error) {
	return
}

func (va VotesAccessorMock) saveImplVoteOrRemove(c context.Context, vote ImplVoteLog, nickname string) (*datastore.Key, error) {
	return nil, nil
}

func (va VotesAccessorMock) decorateIdiom(c context.Context, idiom *Idiom, username string) {
}

func (va VotesAccessorMock) decorateImpl(c context.Context, impl *Impl, username string) {
}
