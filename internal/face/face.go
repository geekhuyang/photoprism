/*

Package face provides face landmark detection.

Copyright (c) 2018 - 2021 Michael Mayer <hello@photoprism.org>

    This program is free software: you can redistribute it and/or modify
    it under the terms of the GNU Affero General Public License as published
    by the Free Software Foundation, either version 3 of the License, or
    (at your option) any later version.

    This program is distributed in the hope that it will be useful,
    but WITHOUT ANY WARRANTY; without even the implied warranty of
    MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
    GNU Affero General Public License for more details.

    You should have received a copy of the GNU Affero General Public License
    along with this program.  If not, see <https://www.gnu.org/licenses/>.

    PhotoPrism® is a registered trademark of Michael Mayer.  You may use it as required
    to describe our software, run your own server, for educational purposes, but not for
    offering commercial goods, products, or services without prior written permission.
    In other words, please ask.

Feel free to send an e-mail to hello@photoprism.org if you have questions,
want to support our work, or just want to say hello.

Additional information can be found in our Developer Guide:
https://docs.photoprism.org/developer-guide/

*/

package face

import (
	"encoding/json"

	"github.com/photoprism/photoprism/internal/event"
)

var log = event.Log

// Faces is a list of face detection results.
type Faces []Face

// Face represents a face detection result.
type Face struct {
	Rows      int    `json:"rows,omitempty"`
	Cols      int    `json:"cols,omitempty"`
	Score     int    `json:"score,omitempty"`
	Face      Point  `json:"face,omitempty"`
	Eyes      Points `json:"eyes,omitempty"`
	Landmarks Points `json:"landmarks,omitempty"`
}

// Dim returns the max number of rows and cols as float32 to calculate relative coordinates.
func (f *Face) Dim() float32 {
	if f.Rows > f.Cols {
		return float32(f.Rows)
	}

	if f.Cols > 0 {
		return float32(f.Cols)
	}

	return float32(1)
}

// Marker returns the relative position on the image.
func (f *Face) Marker() Marker {
	return f.Face.Marker(Point{}, f.Dim())
}

// EyesMidpoint returns the point in between the eyes.
func (f *Face) EyesMidpoint() Point {
	if len(f.Eyes) != 2 {
		return Point{
			Name:  "midpoint",
			Row:   f.Face.Row,
			Col:   f.Face.Col,
			Scale: 0,
		}
	}

	return Point{
		Name:  "midpoint",
		Row:   (f.Eyes[0].Row + f.Eyes[1].Row) / 2,
		Col:   (f.Eyes[0].Col + f.Eyes[1].Col) / 2,
		Scale: 0,
	}
}

// RelativeLandmarks returns detected relative marker positions.
func (f *Face) RelativeLandmarks() Markers {
	p := f.EyesMidpoint()

	m := f.Landmarks.Markers(p, f.Dim())
	m = append(m, f.Eyes.Markers(p, f.Dim())...)

	return m
}

// RelativeLandmarksJSON returns detected relative marker positions as JSON.
func (f *Face) RelativeLandmarksJSON() (b []byte) {
	b, err := json.Marshal(f.RelativeLandmarks())

	if err != nil {
		log.Errorf("faces: %s", err)
		return []byte("{}")
	}

	return b
}
