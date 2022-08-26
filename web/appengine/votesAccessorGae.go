package main

import (
	. "github.com/Deleplace/programming-idioms/idioms"

	"context"

	"google.golang.org/appengine/datastore"
)

// GaeVotesAccessor is a votesAccessor designed for the Google App Engine Datastore.
type GaeVotesAccessor struct {
}

func (va GaeVotesAccessor) idiomVote(ctx context.Context, vote IdiomVoteLog, nickname string) (newRating int, myVote int, err error) {
	// TODO: a transaction to make sure the counter is safely incremented

	delta, _, storedVote, errsave := va.saveIdiomVoteOrRemove(ctx, vote, nickname)
	err = errsave
	if storedVote != nil {
		myVote = storedVote.Value
	}

	if delta != 0 {
		_, idiom, errinc := dao.stealthIncrementIdiomRating(ctx, vote.IdiomId, delta)
		newRating, err = idiom.Rating, errinc
	}
	return
}

func (va GaeVotesAccessor) implVote(ctx context.Context, vote ImplVoteLog, nickname string) (newRating int, myVote int, err error) {
	// TODO a transaction for (vote save + idiom save).  Note that rating data is redundant as rating could be recomputed.

	_, idiom, errget := dao.getIdiomByImplID(ctx, vote.ImplId)
	if err != nil {
		err = errget
		return
	}

	vote.IdiomId = idiom.Id
	delta, _, storedVote, errsave := va.saveImplVoteOrRemove(ctx, vote, nickname)
	err = errsave
	if storedVote != nil {
		myVote = storedVote.Value
	}

	if delta != 0 {
		_, _, newRating, err = dao.stealthIncrementImplRating(ctx, vote.IdiomId, vote.ImplId, delta)
	}

	return
}

// This ancestor (voting booth) is specific to a nickname.
// Thus it should not lead to much contention.
func (va GaeVotesAccessor) createVoteIdiomLogNicknameAncestorKey(ctx context.Context, nickname string) (key *datastore.Key) {
	return datastore.NewKey(ctx, "IdiomVoteLog", nickname, 0, nil)
}

func (va GaeVotesAccessor) getIdiomVotes(ctx context.Context, nickname string) (keys []*datastore.Key, votes []*IdiomVoteLog, err error) {
	q := datastore.NewQuery("IdiomVoteLog").Ancestor(va.createVoteIdiomLogNicknameAncestorKey(ctx, nickname))
	votes = make([]*IdiomVoteLog, 0, 1)
	keys, err = q.GetAll(ctx, &votes)
	if err != nil {
		return keys, votes, err
	}
	return keys, votes, nil
}

func (va GaeVotesAccessor) getIdiomVote(ctx context.Context, nickname string, idiomID int) (key *datastore.Key, vote *IdiomVoteLog, err error) {
	q := datastore.NewQuery("IdiomVoteLog").Ancestor(va.createVoteIdiomLogNicknameAncestorKey(ctx, nickname))
	q = q.Filter("IdiomId = ", idiomID)
	votes := make([]*IdiomVoteLog, 0, 1)
	keys, err := q.GetAll(ctx, &votes)
	if err != nil {
		return key, vote, err
	}
	if len(votes) == 0 {
		return key, vote, nil
	}
	return keys[0], votes[0], nil
}

func (va GaeVotesAccessor) saveIdiomVoteOrRemove(ctx context.Context, vote IdiomVoteLog, nickname string) (delta int, key *datastore.Key, storedVote *IdiomVoteLog, err error) {
	key, existing, err := va.getIdiomVote(ctx, nickname, vote.IdiomId)
	if err != nil {
		return
	}
	if existing != nil {
		// Well, this means the user has decided to click it again, it order to take back her vote
		err = datastore.Delete(ctx, key)
		if err == nil {
			delta = -existing.Value
			storedVote = nil
		} else {
			storedVote = existing
		}
	} else {
		key = datastore.NewIncompleteKey(ctx, "IdiomVoteLog", va.createVoteIdiomLogNicknameAncestorKey(ctx, nickname))
		key, err = datastore.Put(ctx, key, &vote)
		if err == nil {
			delta = vote.Value
			storedVote = &vote
		}
	}
	return
}

// This ancestor (voting booth) is specific to a nickname.
// Thus it should not lead to much contention.
func (va GaeVotesAccessor) createVoteImplLogNicknameAncestorKey(ctx context.Context, nickname string) (key *datastore.Key) {
	return datastore.NewKey(ctx, "ImplVoteLog", nickname, 0, nil)
}

func (va GaeVotesAccessor) getImplVotes(ctx context.Context, nickname string) (keys []*datastore.Key, votes []*ImplVoteLog, err error) {
	q := datastore.NewQuery("ImplVoteLog").Ancestor(va.createVoteImplLogNicknameAncestorKey(ctx, nickname))
	votes = make([]*ImplVoteLog, 0, 1)
	keys, err = q.GetAll(ctx, &votes)
	if err != nil {
		return keys, votes, err
	}
	return keys, votes, nil
}

func (va GaeVotesAccessor) getImplVote(ctx context.Context, nickname string, implID int) (key *datastore.Key, vote *ImplVoteLog, err error) {
	q := datastore.NewQuery("ImplVoteLog").Ancestor(va.createVoteImplLogNicknameAncestorKey(ctx, nickname))
	q = q.Filter("ImplId = ", implID)
	votes := make([]*ImplVoteLog, 0, 1)
	keys, err := q.GetAll(ctx, &votes)
	if err != nil {
		return key, vote, err
	}
	if len(votes) == 0 {
		return key, vote, nil
	}
	return keys[0], votes[0], nil
}

func (va GaeVotesAccessor) getImplVotesForIdiom(ctx context.Context, nickname string, idiomID int) (keys []*datastore.Key, votes []*ImplVoteLog, err error) {
	q := datastore.NewQuery("ImplVoteLog").Ancestor(va.createVoteImplLogNicknameAncestorKey(ctx, nickname))
	q = q.Filter("IdiomId = ", idiomID)
	keys, err = q.GetAll(ctx, &votes)
	return
}

func (va GaeVotesAccessor) saveImplVoteOrRemove(ctx context.Context, vote ImplVoteLog, nickname string) (delta int, key *datastore.Key, storedVote *ImplVoteLog, err error) {
	key, existing, err := va.getImplVote(ctx, nickname, vote.ImplId)
	if err != nil {
		return
	}
	if existing != nil {
		// Well, this means the user has decided to click it again, it order to take back her vote
		err = datastore.Delete(ctx, key)
		if err == nil {
			delta = -existing.Value
			storedVote = nil
		} else {
			storedVote = existing
		}
	} else {
		key = datastore.NewIncompleteKey(ctx, "ImplVoteLog", va.createVoteImplLogNicknameAncestorKey(ctx, nickname))
		key, err = datastore.Put(ctx, key, &vote)
		if err == nil {
			delta = vote.Value
			storedVote = &vote
		}
	}
	return
}

func (va GaeVotesAccessor) decorateIdiom(ctx context.Context, idiom *Idiom, username string) error {
	if username == "" {
		return nil
	}

	return ConcurrentWithAnyError(

		// Mark idiom already upvoted or downvoted by current user, if she did.
		func() error {
			_, vote, err := va.getIdiomVote(ctx, username, idiom.Id)
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
			_, votes, err := va.getImplVotesForIdiom(ctx, username, idiom.Id)
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
