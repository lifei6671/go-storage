package fileblob

import (
	"encoding/json"
	"os"
)

const attrsExt = ".attrs"

// xattrs stores extended attributes for an object. The format is like
// filesystem extended attributes, see
// https://www.freedesktop.org/wiki/CommonExtendedAttributes.
type xattrs struct {
	ContentType string `json:"user.content_type"`
}

// setAttrs creates a "path.attrs" file along with blob to store the attributes,
// it uses JSON format.
func setAttrs(path string, xa xattrs) error {
	f, err := os.Create(path + attrsExt)
	if err != nil {
		return err
	}
	if err := json.NewEncoder(f).Encode(xa); err != nil {
		f.Close()
		return err
	}
	return f.Close()
}

// getAttrs looks at the "path.attrs" file to retrieve the attributes and
// decodes them into a xattrs struct. It doesn't return error when there is no
// such .attrs file.
func getAttrs(path string) (xattrs, error) {
	f, err := os.Open(path + attrsExt)
	if err != nil {
		if os.IsNotExist(err) {
			// Handle gracefully for non-existing .attr files.
			return xattrs{
				ContentType: "application/octet-stream",
			}, nil
		}
		return xattrs{}, err
	}
	xa := new(xattrs)
	if err := json.NewDecoder(f).Decode(xa); err != nil {
		f.Close()
		return xattrs{}, err
	}
	return *xa, f.Close()
}
