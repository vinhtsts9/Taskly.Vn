package elasticsearch

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"Taskly.com/m/global"
	model "Taskly.com/m/internal/models"
	"github.com/gin-gonic/gin"
)

type SearchResponse struct {
	Hits struct {
		Hits []struct {
			Source model.Gig `json:"_source"`
			Sort   []any     `json:"sort"`
		} `json:"hits"`
	} `json:"hits"`
}

func SearchGigs(c *gin.Context) {
	// Parse query parameters
	keyword := c.Query("keyword")
	category := c.Query("category")
	minPrice := c.Query("minPrice")
	maxPrice := c.Query("maxPrice")
	sizeStr := c.DefaultQuery("size", "10")
	size, err := strconv.Atoi(sizeStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid 'size' parameter, must be an integer"})
		return
	}

	// Parse search_after
	var searchAfter []interface{}
	searchAfterStr := c.Query("search_after")
	if searchAfterStr != "" {
		parts := strings.Split(searchAfterStr, ",")
		if len(parts) == 2 {
			createdAt, err := strconv.ParseFloat(parts[0], 64)
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid 'search_after' parameter: first value must be a number"})
				return
			}
			searchAfter = append(searchAfter, createdAt)
			searchAfter = append(searchAfter, parts[1])
		} else {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid 'search_after' parameter, expected two values separated by a comma"})
			return
		}
	}

	// Build query
	queryString := "*"
	if keyword != "" {
		// Split keyword into terms and search for each term with wildcards
		// This allows searching for "Gig 17" to match "Sample Gig Title 17"
		terms := strings.Fields(keyword)
		for i, term := range terms {
			terms[i] = "*" + term + "*"
		}
		queryString = strings.Join(terms, " AND ")
	}

	query := map[string]any{
		"size": size,
		"sort": []any{
			map[string]any{"created_at": "desc"},
			map[string]any{"id": "desc"},
		},
		"query": map[string]any{
			"bool": map[string]any{
				"must": []any{
					map[string]any{
						"query_string": map[string]any{
							"query":            queryString,
							"fields":           []string{"title", "description"},
							"default_operator": "AND",
						},
					},
				},
				"filter": []any{},
			},
		},
	}

	if len(searchAfter) > 0 {
		query["search_after"] = searchAfter
	}

	if category != "" {
		query["query"].(map[string]any)["bool"].(map[string]any)["filter"] = append(
			query["query"].(map[string]any)["bool"].(map[string]any)["filter"].([]any),
			map[string]any{"term": map[string]any{"category": category}},
		)
	}

	if minPrice != "" || maxPrice != "" {
		priceRange := map[string]any{}
		if minPrice != "" {
			if _, err := strconv.ParseFloat(minPrice, 64); err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid 'minPrice' parameter, must be a number"})
				return
			}
			priceRange["gte"] = minPrice
		}
		if maxPrice != "" {
			if _, err := strconv.ParseFloat(maxPrice, 64); err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid 'maxPrice' parameter, must be a number"})
				return
			}
			priceRange["lte"] = maxPrice
		}
		query["query"].(map[string]any)["bool"].(map[string]any)["filter"] = append(
			query["query"].(map[string]any)["bool"].(map[string]any)["filter"].([]any),
			map[string]any{"range": map[string]any{"price": priceRange}},
		)
	}

	var buf bytes.Buffer
	if err := json.NewEncoder(&buf).Encode(query); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	res, err := global.Elasticsearch.Search(
		global.Elasticsearch.Search.WithContext(context.Background()),
		global.Elasticsearch.Search.WithIndex("dbserver1.public.gigs"),
		global.Elasticsearch.Search.WithBody(&buf),
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer res.Body.Close()

	var sr SearchResponse
	if err := json.NewDecoder(res.Body).Decode(&sr); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	results := make([]map[string]any, 0)
	var nextSearchAfter string
	if len(sr.Hits.Hits) > 0 {
		lastHit := sr.Hits.Hits[len(sr.Hits.Hits)-1]
		if len(lastHit.Sort) == 2 {
			nextSearchAfter = fmt.Sprintf("%v,%v", lastHit.Sort[0], lastHit.Sort[1])
		}
	}
	for _, hit := range sr.Hits.Hits {
		createdAtStr := ""
		// Type assertion để kiểm tra và chuyển đổi created_at
		if createdAtFloat, ok := hit.Source.CreatedAt.(float64); ok {
			microseconds := int64(createdAtFloat)
			t := time.Unix(0, microseconds*1000) // Chuyển micro giây thành nano giây
			createdAtStr = t.UTC().Format(time.RFC3339)
		}

		item := map[string]any{
			"gig": map[string]interface{}{
				"id":            hit.Source.ID,
				"user_id":       hit.Source.UserID,
				"title":         hit.Source.Title,
				"description":   hit.Source.Description,
				"category_id":   hit.Source.CategoryID,
				"price":         hit.Source.Price,
				"delivery_time": hit.Source.DeliveryTime,
				"image_url":     hit.Source.ImageURL,
				"status":        hit.Source.Status,
				"created_at":    createdAtStr, // Trả về chuỗi đã định dạng
			},
		}
		results = append(results, item)
	}

	response := gin.H{
		"data":              results,
		"next_search_after": nextSearchAfter,
	}
	c.JSON(http.StatusOK, response)
}
