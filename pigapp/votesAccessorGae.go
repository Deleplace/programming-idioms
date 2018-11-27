package main

import (
	. "github.com/Deleplace/programming-idioms/pig"

	"cloud.google.com/go/datastore"
	"golang.org/x/net/context"
)

// GaeVotesAccessor is a votesAccessor designed for the Google App Engine Datastore.
type GaeVotesAccessor struct {
}

func (va GaeVotesAccessor) idiomVote(c context.Context, vote IdiomVoteLog, nickname string) (newRating int, myVote int, err error) {
	// TODO: a transaction to make sure the counter is safely incremented

	delta, _, storedVote, errsave := va.saveIdiomVoteOrRemove(c, vote, nickname)
	err = errsave
	if storedVote != nil {
		myVote = storedVote.Value
	}

	if delta != 0 {
		_, idiom, errinc := dao.stealthIncrementIdiomRating(c, vote.IdiomId, delta)
		newRating, err = idiom.Rating, errinc
	}
	return
}

func (va GaeVotesAccessor) implVote(c context.Context, vote ImplVoteLog, nickname string) (newRating int, myVote int, err error) {
	// TODO a transaction for (vote save + idiom save).  Note that rating data is redundant as rating could be recomputed.

	_, idiom, errget := dao.getIdiomByImplID(c, vote.ImplId)
	if err != nil {
		err = errget
		return
	}

	vote.IdiomId = idiom.Id
	delta, _, storedVote, errsave := va.saveImplVoteOrRemove(c, vote, nickname)
	err = errsave
	if storedVote != nil {
		myVote = storedVote.Value
	}

	if delta != 0 {
		_, _, newRating, err = dao.stealthIncrementImplRating(c, vote.IdiomId, vote.ImplId, delta)
	}

	return
}

// This ancestor (voting booth) is specific to a nickname.
// Thus it should not lead to much contention.
func (va GaeVotesAccessor) createVoteIdiomLogNicknameAncestorKey(c context.Context, nickname string) (key *datastore.Key) {
	return datastore.NameKey("IdiomVoteLog", nickname, nil)
}

func (va GaeVotesAccessor) getIdiomVotes(c context.Context, nickname string) (keys []*datastore.Key, votes []*IdiomVoteLog, err error) {
	q := datastore.NewQuery("IdiomVoteLog").Ancestor(va.createVoteIdiomLogNicknameAncestorKey(c, nickname))
	votes = make([]*IdiomVoteLog, 0, 1)
	keys, err = ds.GetAll(c, q, &votes)
	if err != nil {
		return keys, votes, err
	}
	return keys, votes, nil
}

func (va GaeVotesAccessor) getIdiomVote(c context.Context, nickname string, idiomID int) (key *datastore.Key, vote *IdiomVoteLog, err error) {
	q := datastore.NewQuery("IdiomVoteLog").Ancestor(va.createVoteIdiomLogNicknameAncestorKey(c, nickname))
	q = q.Filter("IdiomId = ", idiomID)
	votes := make([]*IdiomVoteLog, 0, 1)
	keys, err := ds.GetAll(c, q, &votes)
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
		err = ds.Delete(c, key)
		if err == nil {
			delta = -existing.Value
			storedVote = nil
		} else {
			storedVote = existing
		}
	} else {
		key = datastore.IncompleteKey("IdiomVoteLog", va.createVoteIdiomLogNicknameAncestorKey(c, nickname))
		key, err = ds.Put(c, key, &vote)
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
	return datastore.NameKey("ImplVoteLog", nickname, nil)
}

func (va GaeVotesAccessor) getImplVotes(c context.Context, nickname string) (keys []*datastore.Key, votes []*ImplVoteLog, err error) {
	q := datastore.NewQuery("ImplVoteLog").Ancestor(va.createVoteImplLogNicknameAncestorKey(c, nickname))
	votes = make([]*ImplVoteLog, 0, 1)
	keys, err = ds.GetAll(c, q, &votes)
	if err != nil {
		return keys, votes, err
	}
	return keys, votes, nil
}

func (va GaeVotesAccessor) getImplVote(c context.Context, nickname string, implID int) (key *datastore.Key, vote *ImplVoteLog, err error) {
	q := datastore.NewQuery("ImplVoteLog").Ancestor(va.createVoteImplLogNicknameAncestorKey(c, nickname))
	q = q.Filter("ImplId = ", implID)
	votes := make([]*ImplVoteLog, 0, 1)
	keys, err := ds.GetAll(c, q, &votes)
	if err != nil {
		return key, vote, err
	}
	if len(votes) == 0 {
		return key, vote, nil
	}
	return keys[0], votes[0], nil
}

func (va GaeVotesAccessor) getImplVotesForIdiom(c context.Context, nickname string, idiomID int) (keys []*datastore.Key, votes []*ImplVoteLog, err error) {
	q := datastore.NewQuery("ImplVoteLog").Ancestor(va.createVoteImplLogNicknameAncestorKey(c, nickname))
	q = q.Filter("IdiomId = ", idiomID)
	keys, err = ds.GetAll(c, q, &votes)
	return
}

func (va GaeVotesAccessor) saveImplVoteOrRemove(c context.Context, vote ImplVoteLog, nickname string) (delta int, key *datastore.Key, storedVote *ImplVoteLog, err error) {
	key, existing, err := va.getImplVote(c, nickname, vote.ImplId)
	if err != nil {
		return
	}
	if existing != nil {
		// Well, this means the user has decided to click it again, it order to take back her vote
		err = ds.Delete(c, key)
		if err == nil {
			delta = -existing.Value
			storedVote = nil
		} else {
			storedVote = existing
		}
	} else {
		key = datastore.IncompleteKey("ImplVoteLog", va.createVoteImplLogNicknameAncestorKey(c, nickname))
		key, err = ds.Put(c, key, &vote)
		if err == nil {
			delta = vote.Value
			storedVote = &vote
		}
	}
	return
}

func (va GaeVotesAccessor) decorateIdiom(c context.Context, idiom *Idiom, username string) error {
	if username == "" {
		return nil
	}

	return ConcurrentWithAnyError(

		// Mark idiom already upvoted or downvoted by current user, if she did.
		func() error {
			_, vote, err := va.getIdiomVote(c, username, idiom.Id)
			if vote != nil {
				switch vote.Value {
				case -1:
					idiom.Deco.DownVoted = true
				case 1:
					idiom.Deco.UpVoted = true
				}
			}
			return err
		},

		// Mark each impl already upvoted or downvoted by current user, if she did.
		func() error {
			_, votes, err := va.getImplVotesForIdiom(c, username, idiom.Id)
			voteMap := make(map[int]int, len(votes))
			for _, vote := range votes {
				voteMap[vote.ImplId] = vote.Value
			}

			for i := range idiom.Implementations {
				impl := &idiom.Implementations[i]
				switch voteMap[impl.Id] {
				case -1:
					impl.Deco.DownVoted = true
				case 1:
					impl.Deco.UpVoted = true
				}
			}
			return err
		},
	)
}
