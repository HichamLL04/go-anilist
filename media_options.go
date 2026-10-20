package anilist

type searchParams struct {
	format     MediaFormat
	seasonYear int
	limit      int
}

// SearchOption configures parameters for MediaService.Search.
type SearchOption func(*searchParams)

// WithFormat filters search results by media format (TV, Movie, OVA, etc.).
func WithFormat(f MediaFormat) SearchOption {
	return func(p *searchParams) {
		p.format = f
	}
}

// WithSeasonYear filters search results by release year.
func WithSeasonYear(year int) SearchOption {
	return func(p *searchParams) {
		p.seasonYear = year
	}
}

// WithLimit sets the maximum number of results to return (max 50).
func WithLimit(limit int) SearchOption {
	return func(p *searchParams) {
		p.limit = limit
	}
}