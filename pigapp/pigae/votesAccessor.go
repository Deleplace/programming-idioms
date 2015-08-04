package pigae

import (
	. "github.com/Deleplace/programming-idioms/pig"

	"appengine"
	"appengine/datastore"
)

type votesAccessor interface {
	votesGetter
	votesSaver
}

type votesGetter interface {
	getIdiomVotes(c appengine.Context, nickname string) (keys []*datastore.Key, votes []*IdiomVoteLog, err error)
	getIdiomVote(c appengine.Context, nickname string, idiomID int) (key *datastore.Key, vote *IdiomVoteLog, err error)
	getImplVotes(c appengine.Context, nickname string) (keys []*datastore.Key, votes []*ImplVoteLog, err error)
	getImplVote(c appengine.Context, nickname string, implID int) (key *datastore.Key, vote *ImplVoteLog, err error)
	decorateIdiom(c appengine.Context, idiom *Idiom, username string)
	decorateImpl(c appengine.Context, impl *Impl, username string)
}

type votesSaver interface {
	idiomVote(c appengine.Context, vote IdiomVoteLog, nickname string) (newRating int, myVote int, err error)
	implVote(c appengine.Context, vote ImplVoteLog, nickname string) (newRating int, myVote int, err error)
	createVoteIdiomLogNicknameAncestorKey(c appengine.Context, nickname string) (key *datastore.Key)
	saveIdiomVoteOrRemove(c appengine.Context, vote IdiomVoteLog, nickname string) (delta int, key *datastore.Key, storedVote *IdiomVoteLog, err error)
	createVoteImplLogNicknameAncestorKey(c appengine.Context, nickname string) (key *datastore.Key)
	saveImplVoteOrRemove(c appengine.Context, vote ImplVoteLog, nickname string) (delta int, key *datastore.Key, storedVote *ImplVoteLog, err error)
}
