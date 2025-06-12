package result

import "fmt"

const (
	responseErr = "unexpected google images response format\n%v"
)

// Parser parsers query responses, every query implementation must be aware of its expected input and enforce it
// In case of invalid input an error should be returned
type Parser interface {
	Parse(interface{}) (interface{}, error)
}

// GoogleImagesResultParser is the specific parser to parse Google Images results
type GoogleImagesResultParser struct {
}

// Parse takes the first item content and extracts the link field from it
func (g GoogleImagesResultParser) Parse(raw interface{}) (interface{}, error) {
	r, ok := raw.(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf(responseErr, r)
	}

	items, ok := r["items"].([]interface{})
	if !ok {
		return nil, fmt.Errorf(responseErr, r)
	}

	if len(items) < 1 {
		return nil, fmt.Errorf(responseErr, r)
	}

	item, ok := items[0].(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf(responseErr, r)
	}

	link, ok := item["link"].(string)
	if !ok {
		return nil, fmt.Errorf(responseErr, r)
	}

	return link, nil
}
