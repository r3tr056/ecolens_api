package models

import (
	"context"
	"encoding/json"
	"fmt"
	"math/rand"
	"time"

	"github.com/go-redis/redis/v8"
	"gorm.io/gorm"
)

// For autocomplete
type MatchResult struct {
	Rank    int    `json:"rank"`
	Content string `json:"content"`
}

// for search results
type SearchResult struct {
	ResultID   int    `json:"id"`
	Rank       int    `json:"rank"`
	Title      string `json:"title"`
	Descripton string `json:"description"`
	SearchTerm string `json:"searchTerm"`
	Result     string `json:"result"`
}

type SearchResultPage struct {
	PageID  string         `json:"page_id"`
	Index   int            `json:"index"`
	Length  int            `json:"length"`
	Results []SearchResult `json:"results"`
}

func GenerateRandomPageID() string {
	randomNumber := rand.Intn(10000000)
	return fmt.Sprintf("page_%d", randomNumber)
}

func CreateSearchResult(rank int, title, description string) SearchResult {
	return SearchResult{
		Rank:       rank,
		Title:      title,
		Descripton: description,
	}
}

func SearchProducts(ctx context.Context, redisClient *redis.Client, db *gorm.DB, query string, page int, pageSize int) ([]Product, error) {
	// Redis caching
	cacheKey := fmt.Sprintf("marketproduct:%s:%d:%d", query, page, pageSize)
	cachedResult, err := redisClient.Get(ctx, cacheKey).Result()
	if err == nil {
		// Cache Hit : unmarshal the cache
		var products []Product
		if err := json.Unmarshal([]byte(cachedResult), &products); err == nil {
			return products, nil
		}
	}
	// Cache Miss : Perform the database search
	var products []Product

	offset := (page - 1) * pageSize
	limit := pageSize

	err = db.
		Table("products").
		Select("id, ts_rank_cd(to_tsvector('english', name || ' ' || description), to_tsquery('english', ?)) AS rank, name", query).
		Where("to_tsvector('english', name || ' ' || description) @@ to_tsquery('english', ?) OR name ILIKE ANY(ARRAY[?]) OR description ILIKE ANY(ARRAY[?])", query, query, []string{"%" + query + "%"}, []string{"%" + query + "%"}).
		Order("rank DESC").
		Offset(offset).
		Limit(limit).
		Find(&products).
		Error
	if err != nil {
		return nil, err
	}

	// Cache the results in Redis with a TTL (time-to-live)
	if marshalledResult, err := json.Marshal(products); err == nil {
		redisClient.Set(ctx, cacheKey, marshalledResult, 10*time.Minute)
	}

	return products, nil
}

func SearchMarketplaceProducts(ctx context.Context, redisClient *redis.Client, db *gorm.DB, query string, page int, pageSize int) ([]MarketPlaceProduct, error) {
	cacheKey := fmt.Sprintf("marketproduct:%s:%d:%d", query, page, pageSize)
	cachedResult, err := redisClient.Get(ctx, cacheKey).Result()
	if err == nil {
		var products []MarketPlaceProduct
		if err := json.Unmarshal([]byte(cachedResult), &products); err == nil {
			return products, nil
		}
	}

	var products []MarketPlaceProduct

	offset := (page - 1) * pageSize
	limit := pageSize

	err = db.
		Table("marketplace_products").
		Select("id, ts_rank_cd(to_tsvector('english', name || ' ' || description), to_tsquery('english', ?)) AS rank, name", query).
		Where("to_tsvector('english', name || ' ' || description) @@ to_tsquery('english', ?) OR name ILIKE ANY(ARRAY[?]) OR description ILIKE ANY(ARRAY[?])", query, query, []string{"%" + query + "%"}, []string{"%" + query + "%"}).
		Order("rank DESC").
		Offset(offset).
		Limit(limit).
		Find(&products).
		Error

	if err != nil {
		return nil, err
	}

	if marshalledResult, err := json.Marshal(products); err == nil {
		redisClient.Set(ctx, cacheKey, marshalledResult, 10*time.Minute)
	}

	return products, nil
}

func SearchReports(ctx context.Context, redisClient *redis.Client, db *gorm.DB, query string, page int, pageSize int) ([]Report, error) {
	cacheKey := fmt.Sprintf("reports:%s:%d:%d", query, page, pageSize)
	// Try to get results from the cache
	cachedResult, err := redisClient.Get(ctx, cacheKey).Result()
	if err == nil {
		var reports []Report
		if err := json.Unmarshal([]byte(cachedResult), &reports); err == nil {
			return reports, nil
		}
	}

	var reports []Report

	offset := (page - 1) * pageSize
	limit := pageSize

	err = db.
		Table("reports").
		Select("id, ts_rank_cd(to_tsvector('english', summary), to_tsquery('english', ?)) AS rank, name", query).
		Where("to_tsvector('english', summary) @@ to_tsquery('english', ?) OR summary ILIKE ANY(ARRAY[?])", query, []string{"%" + query + "%"}).
		Order("rank DESC").
		Offset(offset).
		Limit(limit).
		Find(&reports).
		Error

	if err != nil {
		return nil, err
	}

	if marshalledResult, err := json.Marshal(reports); err != nil {
		redisClient.Set(ctx, cacheKey, marshalledResult, 10*time.Minute)
	}

	return reports, nil
}
