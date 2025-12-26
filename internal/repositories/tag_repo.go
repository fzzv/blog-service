package repositories

import (
	"strings"

	"blog-service/internal/models"

	"gorm.io/gorm"
)

type TagRepo struct {
	DB *gorm.DB
}

func NewTagRepo(db *gorm.DB) *TagRepo {
	return &TagRepo{DB: db}
}

func normalizeTag(s string) string {
	return strings.ToLower(strings.TrimSpace(s))
}

func (r *TagRepo) GetOrCreateByNames(names []string) ([]models.Tag, error) {
	out := make([]models.Tag, 0, len(names))
	seen := map[string]struct{}{}

	for _, n := range names {
		nn := normalizeTag(n)
		if nn == "" {
			continue
		}
		if _, ok := seen[nn]; ok {
			continue
		}
		seen[nn] = struct{}{}

		var t models.Tag
		err := r.DB.Where("name = ?", nn).First(&t).Error
		if err == nil {
			out = append(out, t)
			continue
		}
		if !IsNotFound(err) {
			return nil, err
		}

		t = models.Tag{Name: nn}
		if err := r.DB.Create(&t).Error; err != nil {
			// 并发创建时可能撞 unique，重查一次
			var t2 models.Tag
			if err2 := r.DB.Where("name = ?", nn).First(&t2).Error; err2 != nil {
				return nil, err
			}
			out = append(out, t2)
			continue
		}
		out = append(out, t)
	}
	return out, nil
}
