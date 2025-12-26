package handlers

import (
	"net/http"
	"strconv"
	"strings"

	"blog-service/internal/middleware"
	"blog-service/internal/models"
	"blog-service/internal/repositories"
	"blog-service/internal/services"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

type PostHandler struct {
	Posts    *services.PostService
	PostRepo *repositories.PostRepo // 用于只读/计数等
	V        *validator.Validate
}

type createPostReq struct {
	Title     string            `json:"title" validate:"required,max=200"`
	ContentMD string            `json:"content_md" validate:"required"`
	Status    models.PostStatus `json:"status" validate:"required,oneof=draft published"`
	Tags      []string          `json:"tags"`
}

func (h PostHandler) Create(c *gin.Context) {
	uid, _ := middleware.GetAuthUserID(c)

	var req createPostReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "bad_request"})
		return
	}
	if err := h.V.Struct(req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "validation_error"})
		return
	}

	p, err := h.Posts.Create(services.CreatePostInput{
		Title:     req.Title,
		ContentMD: req.ContentMD,
		Status:    req.Status,
		Tags:      req.Tags,
		AuthorID:  uid,
	})
	if err != nil {
		switch err {
		case services.ErrInvalidStatus, services.ErrTitleRequired, services.ErrContentRequired:
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": "internal_server_error"})
		}
		return
	}

	c.JSON(http.StatusCreated, postDetailDTO(p))
}

type updatePostReq struct {
	Title     *string            `json:"title" validate:"omitempty,max=200"`
	ContentMD *string            `json:"content_md"`
	Status    *models.PostStatus `json:"status" validate:"omitempty,oneof=draft published"`
	Tags      *[]string          `json:"tags"`
}

func (h PostHandler) Update(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil || id == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "bad_request"})
		return
	}

	var req updatePostReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "bad_request"})
		return
	}
	if err := h.V.Struct(req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "validation_error"})
		return
	}

	p, err := h.Posts.Update(uint(id), services.UpdatePostInput{
		Title:     req.Title,
		ContentMD: req.ContentMD,
		Status:    req.Status,
		Tags:      req.Tags,
	})
	if err != nil {
		switch err {
		case services.ErrPostNotFound:
			c.JSON(http.StatusNotFound, gin.H{"error": "not_found"})
		case services.ErrInvalidStatus, services.ErrTitleRequired, services.ErrContentRequired:
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": "internal_server_error"})
		}
		return
	}

	c.JSON(http.StatusOK, postDetailDTO(p))
}

func (h PostHandler) Delete(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil || id == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "bad_request"})
		return
	}
	if err := h.PostRepo.DeleteByID(uint(id)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal_server_error"})
		return
	}
	c.Status(http.StatusNoContent)
}

func (h PostHandler) List(c *gin.Context) {
	page, _ := strconv.Atoi(c.Query("page"))
	size, _ := strconv.Atoi(c.Query("size"))

	// admin 才允许带 status 参数查看草稿
	statusQ := strings.TrimSpace(c.Query("status"))
	role, _ := middleware.GetAuthRole(c)

	if role == "admin" && statusQ != "" {
		var st models.PostStatus
		if statusQ == "draft" {
			st = models.PostDraft
		} else if statusQ == "published" {
			st = models.PostPublished
		} else {
			c.JSON(http.StatusBadRequest, gin.H{"error": "bad_request"})
			return
		}
		items, total, err := h.PostRepo.ListAny(&st, page, size)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "internal_server_error"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"total": total, "items": postListDTO(items)})
		return
	}

	items, total, err := h.PostRepo.ListPublished(page, size)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal_server_error"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"total": total, "items": postListDTO(items)})
}

func (h PostHandler) GetBySlug(c *gin.Context) {
	slug := c.Param("slug")
	if slug == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "bad_request"})
		return
	}

	role, _ := middleware.GetAuthRole(c)

	var (
		p   *models.Post
		err error
	)
	// admin 若带 token，可以看草稿详情
	if role == "admin" {
		p, err = h.PostRepo.FindBySlugAny(slug)
	} else {
		p, err = h.PostRepo.FindBySlugPublished(slug)
	}
	if err != nil {
		if repositories.IsNotFound(err) {
			c.JSON(http.StatusNotFound, gin.H{"error": "not_found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal_server_error"})
		return
	}

	// 浏览量 +1（失败不阻断）
	_ = h.PostRepo.IncViewCount(p.ID)

	c.JSON(http.StatusOK, postDetailDTO(p))
}

type previewReq struct {
	ContentMD string `json:"content_md" validate:"required"`
}

func (h PostHandler) Preview(c *gin.Context) {
	var req previewReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "bad_request"})
		return
	}
	if err := h.V.Struct(req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "validation_error"})
		return
	}

	html, err := h.Posts.Preview(req.ContentMD)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal_server_error"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"content_html": html})
}

// ---- DTO helpers ----

func postListDTO(items []models.Post) []gin.H {
	out := make([]gin.H, 0, len(items))
	for _, p := range items {
		out = append(out, gin.H{
			"id":            p.ID,
			"title":         p.Title,
			"slug":          p.Slug,
			"status":        p.Status,
			"published_at":  p.PublishedAt,
			"author":        gin.H{"id": p.Author.ID, "username": p.Author.Username},
			"tags":          tagDTO(p.Tags),
			"view_count":    p.ViewCount,
			"like_count":    p.LikeCount,
			"comment_count": p.CommentCount,
			"created_at":    p.CreatedAt,
			"updated_at":    p.UpdatedAt,
		})
	}
	return out
}

func postDetailDTO(p *models.Post) gin.H {
	return gin.H{
		"id":            p.ID,
		"title":         p.Title,
		"slug":          p.Slug,
		"status":        p.Status,
		"published_at":  p.PublishedAt,
		"author":        gin.H{"id": p.Author.ID, "username": p.Author.Username},
		"tags":          tagDTO(p.Tags),
		"content_md":    p.ContentMD,
		"content_html":  p.ContentHTML,
		"view_count":    p.ViewCount,
		"like_count":    p.LikeCount,
		"comment_count": p.CommentCount,
		"created_at":    p.CreatedAt,
		"updated_at":    p.UpdatedAt,
	}
}

func tagDTO(tags []models.Tag) []string {
	out := make([]string, 0, len(tags))
	for _, t := range tags {
		out = append(out, t.Name)
	}
	return out
}
