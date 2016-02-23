package pigae

import (
	. "github.com/Deleplace/programming-idioms/pig"

	"golang.org/x/net/context"
	"google.golang.org/appengine/datastore"
)

type votesAccessor interface {
	votesGetter
	votesSaver
}

type votesGetter interface {
	getIdiomVotes(c context.Context, nickname string) (keys []*datastore.Key, votes []*IdiomVoteLog, err error)
	getIdiomVote(c context.Context, nickname string, idiomID int) (key *datastore.Key, vote *IdiomVoteLog, err error)
	getImplVotes(c context.Context, nickname string) (keys []*datastore.Key, votes []*ImplVoteLog, err error)
	getImplVote(c context.Context, nickname string, implID int) (key *datastore.Key, vote *ImplVoteLog, err error)
	decorateIdiom(c context.Context, idiom *Idiom, username string)
	decorateImpl(c context.Context, impl *Impl, username string)
}

type votesSaver interface {
	idiomVote(c context.Context, vote IdiomVoteLog, nickname string) (newRating int, myVote int, err error)
	implVote(c context.Context, vote ImplVoteLog, nickname string) (newRating int, myVote int, err error)
	createVoteIdiomLogNicknameAncestorKey(c context.Context, nickname string) (key *datastore.Key)
	saveIdiomVoteOrRemove(c context.Context, vote IdiomVoteLog, nickname string) (delta int, key *datastore.Key, storedVote *IdiomVoteLog, err error)
	createVoteImplLogNicknameAncestorKey(c context.Context, nickname string) (key *datastore.Key)
	saveImplVoteOrRemove(c context.Context, vote ImplVoteLog, nickname string) (delta int, key *datastore.Key, storedVote *ImplVoteLog, err error)
}
