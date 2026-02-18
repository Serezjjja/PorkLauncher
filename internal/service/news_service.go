package service

// NewsArticle represents a news article
type NewsArticle struct {
	Title       string `json:"title"`
	DestURL     string `json:"dest_url"`
	Description string `json:"description"`
	ImageURL    string `json:"image_url"`
}

// NewsService handles news - disabled (returns empty data)
type NewsService struct{}

// NewNewsService creates a new news service
func NewNewsService() *NewsService {
	return &NewsService{}
}

// FetchNews returns empty news list (API disabled)
func (s *NewsService) FetchNews() ([]NewsArticle, error) {
	// Return empty list - no external API calls
	return []NewsArticle{}, nil
}

// FetchLatestNews returns nil (API disabled)
func (s *NewsService) FetchLatestNews() (*NewsArticle, error) {
	// Return nil - no external API calls
	return nil, nil
}
