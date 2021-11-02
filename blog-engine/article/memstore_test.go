package article

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMemStore(t *testing.T) {
	memStore := CreateMemStore()

	article, err := CreateArticle(validTitle, validContent)

	require.NoError(t, err)

	ctx := context.Background()

	err = memStore.SaveArticle(ctx, Article{})

	assert.Error(t, ErrNilArticle, err)

	err = memStore.SaveArticle(ctx, article)

	if assert.NoError(t, err) {
		a, err := memStore.FindArticleByID(ctx, article.ID)

		if assert.NoError(t, err) {
			assert.Equal(t, article, a)
		}
	}

	nonExistentID, _ := uuid.NewRandom()

	a, err := memStore.FindArticleByID(ctx, nonExistentID)

	assert.Equal(t, Article{}, a)
	assert.ErrorIs(t, ErrArticleNotFound, err)

}
