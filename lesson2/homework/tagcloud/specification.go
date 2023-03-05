package tagcloud

// TagCloud aggregates statistics about used tags
type TagCloud struct {
	tagMap   map[string]uint
	tagStats []TagStat
}

// TagStat represents statistics regarding single tag
type TagStat struct {
	Tag             string
	OccurrenceCount int
}

// New should create a valid TagCloud instance
func New() *TagCloud {
	return &TagCloud{tagMap: make(map[string]uint), tagStats: make([]TagStat, 0)}
}

// AddTag should add a tag to the cloud if it wasn't present and increase tag occurrence count
// thread-safety is not needed
func (tc *TagCloud) AddTag(tag string) {
	index, hasValue := tc.tagMap[tag]
	if !hasValue {
		index = uint(len(tc.tagStats))
		tc.tagStats = append(tc.tagStats, TagStat{tag, 0})
	}
	tc.tagStats[index].OccurrenceCount++
	for index > 0 && tc.tagStats[index].OccurrenceCount > tc.tagStats[index-1].OccurrenceCount {
		tc.tagMap[tc.tagStats[index-1].Tag] = index
		tc.tagStats[index], tc.tagStats[index-1] = tc.tagStats[index-1], tc.tagStats[index]
		index--
	}
	tc.tagMap[tag] = index
}

// TopN should return top N most frequent tags ordered in descending order by occurrence count
// if there are multiple tags with the same occurrence count then the order is defined by implementation
// if n is greater that TagCloud size then all elements should be returned
// thread-safety is not needed
// there are no restrictions on time complexity
func (tc *TagCloud) TopN(n int) []TagStat {
	if n > len(tc.tagStats) {
		return tc.tagStats[:]
	}
	return tc.tagStats[:n]
}
