package pigae

import (
	. "github.com/Deleplace/programming-idioms/pig"

	"golang.org/x/net/context"
	"google.golang.org/appengine/datastore"
)

// GaeVotesAccessor is a votesAccessor designed for the Google App Engine Datastore.
type GaeVotesAccessor struct {
}

func (va GaeVotesAccessor) idiomVote(c context.Context, vote IdiomVoteLog, nickname string) (newRating int, myVote int, err error) {
	// TODO: a transaction to make sure the counter is safely incremented

	key, idiom, err := dao.getIdiom(c, vote.IdiomId)
	if err != nil {
		return
	}

	delta, _, storedVote, err := va.saveIdiomVoteOrRemove(c, vote, nickname)
	if storedVote != nil {
		myVote = storedVote.Value
	}

	if delta != 0 {
		idiom.Rating += delta
		err = dao.saveExistingIdiom(c, key, idiom)
		if err == nil {
			newRating = idiom.Rating
		}
	}
	return
}

func (va GaeVotesAccessor) implVote(c context.Context, vote ImplVoteLog, nickname string) (newRating int, myVote int, err error) {
	// TODO a transaction for (vote save + idiom save).  Note that rating data is redundant as rating could be recomputed.

	key, idiom, err := dao.getIdiomByImplID(c, vote.ImplId)
	if err != nil {
		return
	}

	vote.IdiomId = idiom.Id
	delta, _, storedVote, err := va.saveImplVoteOrRemove(c, vote, nickname)
	if storedVote != nil {
		myVote = storedVote.Value
	}

	if delta != 0 {
		// TODO: more efficient way than iterating?
		_, impl, _ := idiom.FindImplInIdiom(vote.ImplId)
		impl.Rating += delta
		err = dao.saveExistingIdiom(c, key, idiom)
		if err == nil {
			newRating = impl.Rating
		}
	}

	return
}

// This ancestor (voting booth) is specific to a nickname.
// Thus it should not lead to much contention.
func (va GaeVotesAccessor) createVoteIdiomLogNicknameAncestorKey(c context.Context, nickname string) (key *datastore.Key) {
	return datastore.NewKey(c, "IdiomVoteLog", nickname, 0, nil)
}

func (va GaeVotesAccessor) getIdiomVotes(c context.Context, nickname string) (keys []*datastore.Key, votes []*IdiomVoteLog, err error) {
	q := datastore.NewQuery("IdiomVoteLog").Ancestor(va.createVoteIdiomLogNicknameAncestorKey(c, nickname))
	votes = make([]*IdiomVoteLog, 0, 1)
	keys, err = q.GetAll(c, &votes)
	if err != nil {
		return keys, votes, err
	}
	return keys, votes, nil
}

func (va GaeVotesAccessor) getIdiomVote(c context.Context, nickname string, idiomID int) (key *datastore.Key, vote *IdiomVoteLog, err error) {
	q := datastore.NewQuery("IdiomVoteLog").Ancestor(va.createVoteIdiomLogNicknameAncestorKey(c, nickname))
	q = q.Filter("IdiomId = ", idiomID)
	votes := make([]*IdiomVoteLog, 0, 1)
	keys, err := q.GetAll(c, &votes)
	if err != nil {
		return key, vote, err
	}
	if len(votes) == 0 {
		return key, vote, nil
	}
	return keys[0], votes[0], nil
}

func (va GaeVotesAccessor) saveIdiomVoteOrRemove(c context.Context, vote IdiomVoteLog, nickname string) (delta int, key *datastore.Key, storedVote *IdiomVoteLog, err error) {
	key, existing, err := va.getIdiomVote(c, nickname, vote.IdiomId)
	if err != nil {
		return
	}
	if existing != nil {
		// Well, this means the user has decided to click it again, it order to take back her vote
		err = datastore.Delete(c, key)
		if err == nil {
			delta = -existing.Value
			storedVote = nil
		} else {
			storedVote = existing
		}
	} else {
		key = datastore.NewIncompleteKey(c, "IdiomVoteLog", va.createVoteIdiomLogNicknameAncestorKey(c, nickname))
		key, err = datastore.Put(c, key, &vote)
		if err == nil {
			delta = vote.Value
			storedVote = &vote
		}
	}
	return
}

// This ancestor (voting booth) is specific to a nickname.
// Thus it should not lead to much contention.
func (va GaeVotesAccessor) createVoteImplLogNicknameAncestorKey(c context.Context, nickname string) (key *datastore.Key) {
	return datastore.NewKey(c, "ImplVoteLog", nickname, 0, nil)
}

func (va GaeVotesAccessor) getImplVotes(c context.Context, nickname string) (keys []*datastore.Key, votes []*ImplVoteLog, err error) {
	q := datastore.NewQuery("ImplVoteLog").Ancestor(va.createVoteImplLogNicknameAncestorKey(c, nickname))
	votes = make([]*ImplVoteLog, 0, 1)
	keys, err = q.GetAll(c, &votes)
	if err != nil {
		return keys, votes, err
	}
	return keys, votes, nil
}

func (va GaeVotesAccessor) getImplVote(c context.Context, nickname string, implID int) (key *datastore.Key, vote *ImplVoteLog, err error) {
	q := datastore.NewQuery("ImplVoteLog").Ancestor(va.createVoteImplLogNicknameAncestorKey(c, nickname))
	q = q.Filter("ImplId = ", implID)
	votes := make([]*ImplVoteLog, 0, 1)
	keys, err := q.GetAll(c, &votes)
	if err != nil {
		return key, vote, err
	}
	if len(votes) == 0 {
		return key, vote, nil
	}
	return keys[0], votes[0], nil
}

func (va GaeVotesAccessor) saveImplVoteOrRemove(c context.Context, vote ImplVoteLog, nickname string) (delta int, key *datastore.Key, storedVote *ImplVoteLog, err error) {
	key, existing, err := va.getImplVote(c, nickname, vote.ImplId)
	if err != nil {
		return
	}
	if existing != nil {
		// Well, this means the user has decided to click it again, it order to take back her vote
		err = datastore.Delete(c, key)
		if err == nil {
			delta = -existing.Value
			storedVote = nil
		} else {
			storedVote = existing
		}
	} else {
		key = datastore.NewIncompleteKey(c, "ImplVoteLog", va.createVoteImplLogNicknameAncestorKey(c, nickname))
		key, err = datastore.Put(c, key, &vote)
		if err == nil {
			delta = vote.Value
			storedVote = &vote
		}
	}
	return
}

func (va GaeVotesAccessor) decorateIdiom(c context.Context, idiom *Idiom, username string) {
	if username == "" {
		return
	}
	if _, vote, _ := va.getIdiomVote(c, username, idiom.Id); vote != nil {
		switch vote.Value {
		case -1:
			idiom.Deco.DownVoted = true
		case 1:
			idiom.Deco.UpVoted = true
		}
	}
}

func (va GaeVotesAccessor) decorateImpl(c context.Context, impl *Impl, username string) {
	if username == "" {
		return
	}
	_, vote, _ := va.getImplVote(c, username, impl.Id)
	if vote != nil {
		switch vote.Value {
		case -1:
			impl.Deco.DownVoted = true
		case 1:
			impl.Deco.UpVoted = true
		}
	}
}
