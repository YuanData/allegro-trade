package db

import "context"

type CreateMemberTxParams struct {
	CreateMemberParams
	AfterCreate func(member Member) error
}

type CreateMemberTxResult struct {
	Member Member
}

func (store *SQLStore) CreateMemberTx(ctx context.Context, arg CreateMemberTxParams) (CreateMemberTxResult, error) {
	var result CreateMemberTxResult

	err := store.execTx(ctx, func(q *Queries) error {
		var err error

		result.Member, err = q.CreateMember(ctx, arg.CreateMemberParams)
		if err != nil {
			return err
		}

		return arg.AfterCreate(result.Member)
	})

	return result, err
}
