package article

import (
	"context"
	"errors"
	"sync"

	"github.com/google/uuid"
)

type MemStore struct {
	articles map[uuid.UUID]Article

	mutex *sync.RWMutex
}

var (
	ErrNotImplemented  = errors.New("method not yet implemented")
	ErrArticleNotFound = errors.New("article not found")
)

// Create A MemStore

func CreateMemStore() *MemStore {
	return &MemStore{
		articles: make(map[uuid.UUID]Article),
		mutex:    &sync.RWMutex{},
	}
}

func (m *MemStore) FillArticle(articles ...Article) {
	if len(articles) == 0 {
		return
	}

	m.mutex.Lock()
	defer m.mutex.Unlock()

	for _, a := range articles {
		m.articles[a.ID] = a

	}
}

func (m *MemStore) FindArticleByID(ctx context.Context, id uuid.UUID) (Article, error) {
	if len(m.articles) == 0 {
		return Article{}, ErrArticleNotFound
	}
	m.mutex.RLock()
	defer m.mutex.RUnlock()
	article, ok := m.articles[id]

	if !ok {
		return Article{}, ErrArticleNotFound
	}

	return article, nil

}

func (m *MemStore) SaveArticle(ctx context.Context, article Article) error {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	if article.IsNil() {
		return ErrNilArticle
	}

	m.articles[article.ID] = article

	return nil

}
