package uploadmanifest

import (
	"encoding/json/jsontext"
	"encoding/json/v2"
	"errors"
	"fmt"
	"os"
)

var preserveJSONNumbers = json.WithUnmarshalers(
	json.UnmarshalFromFunc(func(
		dec *jsontext.Decoder,
		val *any,
	) error {
		if dec.PeekKind() == jsontext.KindNumber {
			*val = jsontext.Value(nil)
		}

		return errors.ErrUnsupported
	}),
)

func Load(path string) (File, error) {
	f, err := os.Open(path)
	if err != nil {
		return File{}, fmt.Errorf(
			"read manifest file %q: %w",
			path,
			err,
		)
	}

	defer func() {
		_ = f.Close()
	}()

	var mf File

	if err := json.UnmarshalRead(
		f,
		&mf,
		preserveJSONNumbers,
	); err != nil {
		return File{}, fmt.Errorf(
			"parse manifest file %q: %w",
			path,
			err,
		)
	}

	if len(mf.Items) == 0 {
		return File{}, fmt.Errorf(
			"manifest file %q contains no items",
			path,
		)
	}

	return mf, nil
}
